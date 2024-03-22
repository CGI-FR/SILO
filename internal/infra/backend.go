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
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cockroachdb/pebble"
)

func decode(value []byte) ([]string, error) {
	var set map[string]any

	err := json.Unmarshal(value, &set)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	items := make([]string, 0, len(set))

	for item := range set {
		items = append(items, item)
	}

	return items, nil
}

func encode(items []string) ([]byte, error) {
	set := make(map[string]any, len(items))

	for _, item := range items {
		set[item] = nil
	}

	rawNodes, err := json.Marshal(set)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return rawNodes, nil
}

type Snapshot struct {
	db *pebble.Batch
}

func (s Snapshot) Next() (string, bool, error) {
	iter, err := s.db.NewIter(&pebble.IterOptions{}) //nolint:exhaustruct
	if errors.Is(err, pebble.ErrNotFound) {
		return "", false, nil
	} else if err != nil {
		return "", false, fmt.Errorf("%w", err)
	}

	defer iter.Close()

	if !iter.First() {
		return "", false, nil
	}

	return string(iter.Key()), true, nil
}

func (s Snapshot) PullAll(node string) ([]string, error) {
	var (
		set []string
		err error
	)

	item, closer, err := s.db.Get([]byte(node))
	if errors.Is(err, pebble.ErrNotFound) {
		return []string{}, nil
	} else if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	defer closer.Close()

	set, err = decode(item)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if err := s.db.Delete([]byte(node), pebble.Sync); err != nil {
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

func (b Backend) Get(node string) ([]string, error) {
	item, closer, err := b.db.Get([]byte(node))
	if errors.Is(err, pebble.ErrNotFound) {
		return []string{}, nil
	} else if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	defer closer.Close()

	return decode(item)
}

func (b Backend) Store(key string, value string) error {
	nodes, err := b.Get(key)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	nodes = append(nodes, value)

	rawNodes, err := encode(nodes)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if err := b.db.Set([]byte(key), rawNodes, pebble.NoSync); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (b Backend) Snapshot() Snapshot {
	return Snapshot{b.db.NewIndexedBatch()}
}

func (b Backend) Close() error {
	if err := b.db.Close(); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func NewBackend(path string) (Backend, error) {
	database, err := pebble.Open(path, &pebble.Options{})
	if err != nil {
		return Backend{}, fmt.Errorf("unable to open database %v : %w", path, err)
	}

	return Backend{db: database}, nil
}
