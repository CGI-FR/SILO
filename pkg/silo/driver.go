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

	"github.com/rs/zerolog/log"
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

	defer snapshot.Close()

	for count := 0; ; count++ {
		entryNode, hasNext, err := snapshot.Next()
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		if !hasNext {
			break
		}

		entity := NewEntity(entryNode)

		if err := d.writer.Write(entryNode, entity.UUID()); err != nil {
			return fmt.Errorf("%w", err)
		}

		if err := d.dump(snapshot, entryNode, entity); err != nil {
			return fmt.Errorf("%w", err)
		}

		entity.Finalize()
	}

	return nil
}

func (d *Driver) dump(snapshot Snapshot, node DataNode, entity Entity) error {
	connectedNodes, err := snapshot.PullAll(node)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	for _, connectedNode := range connectedNodes {
		if entity.Append(connectedNode) {
			if err := d.writer.Write(connectedNode, entity.UUID()); err != nil {
				return fmt.Errorf("%w", err)
			}

			if err := d.dump(snapshot, connectedNode, entity); err != nil {
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

		log.Info().Int("links", len(links)).Interface("row", datarow).Msg("datarow scanned")

		for _, link := range links {
			if err := d.backend.Store(link.E1, link.E2); err != nil {
				return fmt.Errorf("%w: %w", ErrPersistingData, err)
			}

			if err := d.backend.Store(link.E2, link.E1); err != nil {
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
