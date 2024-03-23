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

	"github.com/cgi-fr/silo/pkg/silo"
	"github.com/cockroachdb/pebble"
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
		return silo.DataNode{}, false, nil
	} else if err != nil {
		return silo.DataNode{}, false, fmt.Errorf("%w", err)
	}

	defer iter.Close()

	if !iter.First() {
		return silo.DataNode{}, false, nil
	}

	key, err := silo.DecodeDataNode(iter.Key())
	if err != nil {
		return silo.DataNode{}, false, fmt.Errorf("%w", err)
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
	database, err := pebble.Open(path, &pebble.Options{}) //nolint:exhaustruct
	if err != nil {
		return Backend{}, fmt.Errorf("unable to open database %v : %w", path, err)
	}

	return Backend{db: database}, nil
}
