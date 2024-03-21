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

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/dominikbraun/graph"
)

type Driver struct {
	backend Backend
	writer  DumpWriter
}

func NewDriver(backend Backend, writer DumpWriter) *Driver {
	return &Driver{
		backend: backend,
		writer:  writer,
	}
}

func (d *Driver) Dump() error {
	nodes, _ := d.ReadAllNodes()
	links, _ := d.ReadAllLinks()
	nodemap := map[string]DataNode{}

	grph := graph.New[string, DataNode](func(n DataNode) string { return n.String() })

	for _, node := range nodes {
		nodemap[node.String()] = node

		if err := grph.AddVertex(node); err != nil && !errors.Is(err, graph.ErrVertexAlreadyExists) {
			return fmt.Errorf("%w", err)
		}
	}

	for _, link := range links {
		if err := grph.AddEdge(link.E1.String(), link.E2.String()); err != nil && !errors.Is(err, graph.ErrEdgeAlreadyExists) {
			return fmt.Errorf("%w", err)
		}
	}

	for count := 0; ; count++ {
		if len(nodemap) == 0 {
			break
		}

		for key := range nodemap {
			err := graph.BFS[string, DataNode](grph, key, func(hash string) bool {
				_ = d.writer.Write(nodemap[hash], strconv.Itoa(count))

				delete(nodemap, hash)

				return false
			})
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			break
		}
	}

	return nil
}

func (d *Driver) Scan(input DataRowReader) error {
	defer input.Close()

	for {
		datarow, err := input.ReadDataRow()
		if err != nil && !errors.Is(err, io.EOF) {
			return fmt.Errorf("%w: %w", ErrReadingNextInput, err)
		}

		if errors.Is(err, io.EOF) || datarow == nil {
			break
		}

		nodes, links := Scan(datarow)

		for _, node := range nodes {
			if err := d.backend.StoreNode(node); err != nil {
				return fmt.Errorf("%w: %w", ErrPersistingData, err)
			}
		}

		for _, link := range links {
			if err := d.backend.StoreLink(link); err != nil {
				return fmt.Errorf("%w: %w", ErrPersistingData, err)
			}
		}
	}

	return nil
}

func (d *Driver) ReadAllNodes() ([]DataNode, error) {
	nodes := []DataNode{}
	reader := d.backend.ReadNodes()

	for {
		node, err := reader.ReadDataNode()
		if err != nil && !errors.Is(err, io.EOF) {
			return nodes, fmt.Errorf("%w: %w", ErrReadingPersistedData, err)
		}

		if errors.Is(err, io.EOF) {
			break
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}

func (d *Driver) ReadAllLinks() ([]DataLink, error) {
	links := []DataLink{}
	reader := d.backend.ReadLinks()

	for {
		link, err := reader.ReadDataLink()
		if err != nil && !errors.Is(err, io.EOF) {
			return links, fmt.Errorf("%w: %w", ErrReadingPersistedData, err)
		}

		if errors.Is(err, io.EOF) {
			break
		}

		links = append(links, link)
	}

	return links, nil
}

func Scan(datarow DataRow) ([]DataNode, []DataLink) {
	nodes := []DataNode{}
	links := []DataLink{}

	for key, value := range datarow {
		if value != nil {
			nodes = append(nodes, DataNode{Key: key, Data: value})
		}
	}

	// find all pairs in nodes
	for i := 0; i < len(nodes); i++ {
		for j := i + 1; j < len(nodes); j++ {
			links = append(links, DataLink{E1: nodes[i], E2: nodes[j]})
		}
	}

	return nodes, links
}
