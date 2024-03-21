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

type BackendInMemory struct {
	links []DataLink
	nodes []DataNode
}

func NewBackendInMemory() *BackendInMemory {
	return &BackendInMemory{
		links: []DataLink{},
		nodes: []DataNode{},
	}
}

func (b *BackendInMemory) StoreLink(link DataLink) error {
	b.links = append(b.links, link)

	return nil
}

func (b *BackendInMemory) StoreNode(node DataNode) error {
	b.nodes = append(b.nodes, node)

	return nil
}

func (b *BackendInMemory) ReadLinks() DataLinkReader { //nolint:ireturn
	return NewDataLinkReaderInMemory(b.links)
}

func (b *BackendInMemory) ReadNodes() DataNodeReader { //nolint:ireturn
	return NewDataNodeReaderInMemory(b.nodes)
}
