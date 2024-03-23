// Copyright (C) 2023 CGI France
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
	"fmt"

	"github.com/cgi-fr/silo/pkg/silo"
)

type DumpJSONLine struct{}

func NewDumpJSONLine() *DumpJSONLine {
	return &DumpJSONLine{}
}

func (d *DumpJSONLine) Write(node silo.DataNode, uuid string) error {
	line := struct {
		UUID string `json:"uuid"`
		ID   string `json:"id"`
		Key  any    `json:"key"`
	}{
		UUID: uuid,
		ID:   node.Key,
		Key:  node.Data,
	}

	bytes, err := json.Marshal(line)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	println(string(bytes))

	return nil
}

func (d *DumpJSONLine) Close() error {
	return nil
}
