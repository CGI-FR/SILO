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
	*config
	backend Backend
	writer  DumpWriter
}

func NewDriver(backend Backend, writer DumpWriter, options ...Option) *Driver {
	errs := []error{}
	config := DefaultConfig()

	for _, option := range options {
		if err := option.apply(config); err != nil {
			errs = append(errs, err)
		}
	}

	if err := config.validate(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		panic(errs)
	}

	return &Driver{
		backend: backend,
		writer:  writer,
		config:  config,
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

		entity := NewEntity(d.config.includeList, d.writer, entryNode)

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

func (d *Driver) Scan(input DataRowReader, observers ...ScanObserver) error {
	defer input.Close()

	for {
		datarow, err := input.ReadDataRow()
		if err != nil && !errors.Is(err, io.EOF) {
			return fmt.Errorf("%w: %w", ErrReadingNextInput, err)
		}

		if errors.Is(err, io.EOF) || datarow == nil {
			break
		}

		nodes, links := d.scan(datarow)

		log.Info().Int("links", len(links)).Interface("row", datarow).Msg("datarow scanned")

		if err := d.ingest(datarow, nodes, links, observers...); err != nil {
			return err
		}
	}

	return nil
}

func (d *Driver) ingest(datarow DataRow, nodes []DataNode, links []DataLink, observers ...ScanObserver) error {
	for _, link := range links {
		if err := d.backend.Store(link.E1, link.E2); err != nil {
			return fmt.Errorf("%w: %w", ErrPersistingData, err)
		}

		if err := d.backend.Store(link.E2, link.E1); err != nil {
			return fmt.Errorf("%w: %w", ErrPersistingData, err)
		}

		for _, observer := range observers {
			observer.IngestedLink(link)
		}
	}

	// optimization : self reference is useful only if no link has been found, and nodes will contain a single node
	if len(links) == 0 && len(nodes) > 0 {
		if err := d.backend.Store(nodes[0], nodes[0]); err != nil {
			return fmt.Errorf("%w: %w", ErrPersistingData, err)
		}
	}

	for _, observer := range observers {
		observer.IngestedRow(datarow)
	}

	return nil
}

func (d *Driver) scan(datarow DataRow) ([]DataNode, []DataLink) {
	nodes := []DataNode{}
	links := []DataLink{}

	for key, value := range datarow {
		if _, included := d.config.include[key]; value != nil && (included || len(d.config.include) == 0) {
			if alias, exist := d.config.aliases[key]; exist {
				key = alias
			}

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
