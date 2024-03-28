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
	"fmt"

	"github.com/cgi-fr/silo/pkg/silo"
	"github.com/cockroachdb/pebble"
)

type BackendFull struct {
	Backend
}

func NewBackendFull(path string) (BackendFull, error) {
	backend, err := NewBackend(path)
	if err != nil {
		return BackendFull{backend}, err
	}

	return BackendFull{backend}, nil
}

func (b BackendFull) Snapshot() silo.Snapshot { //nolint:ireturn
	return NewSnapshotFull(b.db)
}

type SnapshotFull struct {
	db     *pebble.DB
	nodes  map[string][]byte
	loaded bool
}

const DefaultFullMapCap = 1024

func NewSnapshotFull(db *pebble.DB) silo.Snapshot { //nolint:ireturn
	return &SnapshotFull{
		db:     db,
		nodes:  make(map[string][]byte, DefaultFullMapCap),
		loaded: false,
	}
}

func (s *SnapshotFull) Load() error {
	iter, err := s.db.NewIter(&pebble.IterOptions{}) //nolint:exhaustruct
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	for iter.First(); iter.Valid(); iter.Next() {
		s.nodes[string(iter.Key())] = iter.Value()
	}

	s.loaded = true

	return nil
}

func (s *SnapshotFull) Next() (silo.DataNode, bool, error) {
	if !s.loaded {
		if err := s.Load(); err != nil {
			return silo.DataNode{Key: "", Data: ""}, false, err
		}
	}

	for key := range s.nodes {
		node, err := decodeKey([]byte(key))
		if err != nil {
			return silo.DataNode{Key: "", Data: ""}, false, fmt.Errorf("%w", err)
		}

		return node, true, nil
	}

	return silo.DataNode{Key: "", Data: ""}, false, nil
}

func (s *SnapshotFull) PullAll(node silo.DataNode) ([]silo.DataNode, error) {
	key, err := node.Binary()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	item, has := s.nodes[string(key)]
	if !has {
		return []silo.DataNode{}, nil
	}

	set, err := decode(item)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	delete(s.nodes, string(key))

	return set, nil
}

func (s *SnapshotFull) Close() error {
	return nil
}
