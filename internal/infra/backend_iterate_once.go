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
	"errors"
	"fmt"

	"github.com/cgi-fr/silo/pkg/silo"
	"github.com/cockroachdb/pebble"
)

type BackendInterateOnce struct {
	Backend
}

func NewBackendInterateOnce(path string) (BackendInterateOnce, error) {
	backend, err := NewBackend(path)
	if err != nil {
		return BackendInterateOnce{backend}, err
	}

	return BackendInterateOnce{backend}, nil
}

func (b BackendInterateOnce) Snapshot() silo.Snapshot { //nolint:ireturn
	return NewSnapshotInterateOnce(b.db)
}

type SnapshotInterateOnce struct {
	db     *pebble.DB
	iter   *pebble.Iterator
	pulled map[string]bool
}

const DefaultPulledMapCap = 128

func NewSnapshotInterateOnce(db *pebble.DB) silo.Snapshot { //nolint:ireturn
	return SnapshotInterateOnce{
		db:     db,
		iter:   nil,
		pulled: make(map[string]bool, DefaultPulledMapCap),
	}
}

func (s SnapshotInterateOnce) Next() (silo.DataNode, bool, error) {
	if s.iter == nil { //nolint:nestif
		var err error
		if s.iter, err = s.db.NewIter(&pebble.IterOptions{}); err != nil { //nolint:exhaustruct
			return silo.DataNode{Key: "", Data: ""}, false, fmt.Errorf("%w", err)
		}

		if !s.iter.First() {
			return silo.DataNode{Key: "", Data: ""}, false, nil
		}

		if _, pulled := s.pulled[string(s.iter.Key())]; !pulled {
			node, err := decodeKey(s.iter.Key())
			if err != nil {
				return silo.DataNode{Key: "", Data: ""}, false, fmt.Errorf("%w", err)
			}

			return node, true, nil
		}
	}

	for {
		if !s.iter.Next() {
			return silo.DataNode{Key: "", Data: ""}, false, nil
		}

		if _, pulled := s.pulled[string(s.iter.Key())]; !pulled {
			node, err := decodeKey(s.iter.Key())
			if err != nil {
				return silo.DataNode{Key: "", Data: ""}, false, fmt.Errorf("%w", err)
			}

			return node, true, nil
		}
	}
}

func (s SnapshotInterateOnce) PullAll(node silo.DataNode) ([]silo.DataNode, error) {
	key, err := node.Binary()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if _, pulled := s.pulled[string(key)]; pulled {
		return []silo.DataNode{}, nil
	}

	s.pulled[string(key)] = true

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

	return set, nil
}

func (s SnapshotInterateOnce) Close() error {
	if s.iter == nil {
		return nil
	}

	if err := s.iter.Close(); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func decodeKey(rawKey []byte) (silo.DataNode, error) {
	key, err := silo.DecodeDataNode(rawKey)
	if err != nil {
		return silo.DataNode{Key: "", Data: ""}, fmt.Errorf("%w", err)
	}

	return key, nil
}
