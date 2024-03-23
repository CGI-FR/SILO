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

package silo

import "github.com/cgi-fr/silo/pkg/multimap"

type BackendInMemory struct {
	links multimap.Multimap[DataNode, DataNode]
}

func NewBackendInMemory() *BackendInMemory {
	return &BackendInMemory{
		links: multimap.Multimap[DataNode, DataNode]{},
	}
}

func (b *BackendInMemory) Store(key DataNode, value DataNode) error {
	b.links.Add(key, value)

	return nil
}

func (b *BackendInMemory) Snapshot() Snapshot { //nolint:ireturn
	return &BackendInMemory{
		links: b.links.Copy(),
	}
}

func (b *BackendInMemory) Close() error {
	return nil
}

func (b *BackendInMemory) Next() (DataNode, bool, error) {
	key, present := b.links.RandomKey()

	return key, present, nil
}

func (b *BackendInMemory) PullAll(node DataNode) ([]DataNode, error) {
	return b.links.Delete(node), nil
}
