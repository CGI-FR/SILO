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
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

const defaultEntitySize = 10

const (
	statusEntityComplete     = "complete"
	statusEntityConsistent   = "consistent"
	statusEntityInconsistent = "inconsistent"
	statusEntityEmpty        = "empty"
)

type Entity struct {
	include []string
	nodes   map[DataNode]int
	counts  map[string]int
	uuid    string
}

func NewEntity(include []string, nodes ...DataNode) Entity {
	entity := Entity{
		include: include,
		nodes:   make(map[DataNode]int, defaultEntitySize),
		counts:  make(map[string]int, defaultEntitySize),
		uuid:    uuid.NewString(),
	}
	for _, node := range nodes {
		entity.Append(node)
	}

	return entity
}

func (s Entity) Append(node DataNode) bool {
	count, gotNode := s.nodes[node]
	if gotNode {
		s.nodes[node] = count + 1
	} else {
		s.nodes[node] = 1
	}

	count, gotKey := s.counts[node.Key]
	if !gotKey {
		s.counts[node.Key] = 1
	} else if !gotNode {
		s.counts[node.Key] = count + 1
	}

	return !gotNode
}

func (s Entity) UUID() string {
	return s.uuid
}

func (s Entity) Finalize() {
	msg := log.Info().Str("status", statusEntityConsistent)

	counts := s.counts

	if len(s.include) > 0 {
		counts = make(map[string]int, len(s.include))
		for _, key := range s.include {
			if s.counts[key] > 0 {
				counts[key] = s.counts[key]
			}
		}
	}

	if len(counts) == len(s.include) && len(s.include) > 0 {
		msg.Str("status", statusEntityComplete)
	} else if len(counts) == 0 {
		msg.Str("status", statusEntityEmpty)
	}

	for _, count := range counts {
		if count > 1 {
			msg = log.Warn().Str("status", statusEntityInconsistent)

			break
		}
	}

	msg.Str("uuid", s.UUID())

	for id, count := range counts {
		msg.Int("count-"+id, count)
	}

	msg.Msg("entity identified")
}
