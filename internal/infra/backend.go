// Copyright (C) 2024 CGI France
//
// This file is part of SILO.
//
// SILO is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// SILO is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with SILO.  If not, see <http://www.gnu.org/licenses/>.

package infra

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/cgi-fr/silo/pkg/silo"
	"github.com/cockroachdb/pebble"
	"github.com/rs/zerolog/log"
)

func decode(value []byte) ([]silo.DataNode, error) {
	var set map[silo.DataNode]any

	decoder := gob.NewDecoder(bytes.NewBuffer(value))

	err := decoder.Decode(&set)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	items := make([]silo.DataNode, 0, len(set))

	for item := range set {
		items = append(items, item)
	}

	return items, nil
}

func encode(items []silo.DataNode) ([]byte, error) {
	result := new(bytes.Buffer)

	set := make(map[silo.DataNode]any, len(items))

	for _, item := range items {
		set[item] = nil
	}

	encoder := gob.NewEncoder(result)
	if err := encoder.Encode(set); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return result.Bytes(), nil
}

type Snapshot struct {
	db *pebble.Batch
}

func (s Snapshot) Next() (silo.DataNode, bool, error) {
	iter, err := s.db.NewIter(&pebble.IterOptions{}) //nolint:exhaustruct
	if errors.Is(err, pebble.ErrNotFound) {
		return silo.DataNode{Key: "", Data: ""}, false, nil
	} else if err != nil {
		return silo.DataNode{Key: "", Data: ""}, false, fmt.Errorf("%w", err)
	}

	defer iter.Close()

	if !iter.First() {
		return silo.DataNode{Key: "", Data: ""}, false, nil
	}

	key, err := silo.DecodeDataNode(iter.Key())
	if err != nil {
		return silo.DataNode{Key: "", Data: ""}, false, fmt.Errorf("%w", err)
	}

	return key, true, nil
}

func (s Snapshot) PullAll(node silo.DataNode) ([]silo.DataNode, error) {
	key, err := node.Binary()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	item, closer, err := s.db.Get(key)
	if errors.Is(err, pebble.ErrNotFound) {
		return []silo.DataNode{}, nil
	} else if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	defer closer.Close()

	set, err := decode(item)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if err := s.db.Delete(key, pebble.Sync); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return set, nil
}

func (s Snapshot) Close() error {
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

type Backend struct {
	db *pebble.DB
}

func (b Backend) Get(node silo.DataNode) ([]silo.DataNode, error) {
	key, err := node.Binary()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	item, closer, err := b.db.Get(key)
	if errors.Is(err, pebble.ErrNotFound) {
		return []silo.DataNode{}, nil
	} else if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	defer closer.Close()

	return decode(item)
}

func (b Backend) Store(key silo.DataNode, value silo.DataNode) error {
	nodes, err := b.Get(key)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	nodes = append(nodes, value)

	rawNodes, err := encode(nodes)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	rawKey, err := key.Binary()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if err := b.db.Set(rawKey, rawNodes, pebble.NoSync); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (b Backend) Snapshot() silo.Snapshot { //nolint:ireturn
	return Snapshot{b.db.NewIndexedBatch()}
}

func (b Backend) Close() error {
	if err := b.db.Flush(); err != nil {
		return fmt.Errorf("%w", err)
	}

	if err := b.db.Close(); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func NewBackend(path string) (Backend, error) {
	if err := checkDirectory(path); err != nil {
		return Backend{}, fmt.Errorf("unable to open database %v : %w", path, err)
	}

	database, err := pebble.Open(path, &pebble.Options{Logger: BackendLogger{}}) //nolint:exhaustruct
	if err != nil {
		return Backend{}, fmt.Errorf("unable to open database %v : %w", path, err)
	}

	return Backend{db: database}, nil
}

type BackendLogger struct{}

func (l BackendLogger) Infof(format string, args ...interface{}) {
	log.Info().Msgf(format, args...)
}

func (l BackendLogger) Fatalf(format string, args ...interface{}) {
	log.Error().Msgf(format, args...)
}

var ErrPathIsNotValid = errors.New("path is not valid")

func checkDirectory(path string) error {
	// if path does not exists => ok
	if stats, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	} else if !stats.IsDir() {
		// if path exists but is file => ko
		return fmt.Errorf("%w '%s'", ErrPathIsNotValid, path)
	}

	dir, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer dir.Close()

	// if path is empty directory => ok
	if _, err := dir.Readdirnames(1); errors.Is(err, io.EOF) {
		return nil
	} else if err != nil {
		return fmt.Errorf("%w", err)
	}

	// if path is not empty but CURRENT file does not exists => ko
	if stats, err := os.Stat(filepath.Join(path, "CURRENT")); os.IsNotExist(err) {
		return fmt.Errorf("%w '%s'", ErrPathIsNotValid, path)
	} else if stats.IsDir() {
		// if CURRENT exists but is a directory => ko
		return fmt.Errorf("%w '%s'", ErrPathIsNotValid, path)
	}

	return nil
}
