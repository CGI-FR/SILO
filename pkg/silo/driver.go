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
	snapshot := d.backend.Snapshot()

	for count := 0; ; count++ {
		entryNode, present := snapshot.Next()
		if !present {
			break
		}

		_ = d.writer.Write(entryNode, strconv.Itoa(count))

		done := map[string]any{entryNode: nil}

		if err := d.dump(snapshot, entryNode, done, count); err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	return nil
}

func (d *Driver) dump(snapshot Snapshot, node string, done map[string]any, id int) error {
	connectedNodes, err := snapshot.PullAll(node)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	for _, connectedNode := range connectedNodes {
		if _, ok := done[connectedNode]; !ok {
			_ = d.writer.Write(connectedNode, strconv.Itoa(id))
			done[connectedNode] = nil

			if err := d.dump(snapshot, connectedNode, done, id); err != nil {
				return fmt.Errorf("%w", err)
			}
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

		links := Scan(datarow)

		for _, link := range links {
			if err := d.backend.Store(link.E1.String(), link.E2.String()); err != nil {
				return fmt.Errorf("%w: %w", ErrPersistingData, err)
			}

			if err := d.backend.Store(link.E2.String(), link.E1.String()); err != nil {
				return fmt.Errorf("%w: %w", ErrPersistingData, err)
			}
		}
	}

	return nil
}

func Scan(datarow DataRow) []DataLink {
	nodes := []DataNode{}
	links := []DataLink{}

	for key, value := range datarow {
		if value != nil {
			nodes = append(nodes, DataNode{Key: key, Data: value})
		}
	}

	if len(nodes) == 1 {
		links = append(links, DataLink{E1: nodes[0], E2: nodes[0]})
	}

	// find all pairs in nodes
	for i := 0; i < len(nodes); i++ {
		for j := i + 1; j < len(nodes); j++ {
			links = append(links, DataLink{E1: nodes[i], E2: nodes[j]})
		}
	}

	return links
}
