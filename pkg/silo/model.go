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
	"bytes"
	"encoding/gob"
	"fmt"
	"strings"
)

type DataRow map[string]any

type DataNode struct {
	Key  string
	Data any
}

type DataLink struct {
	E1 DataNode
	E2 DataNode
}

func DecodeDataNode(data []byte) (DataNode, error) {
	var result DataNode

	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)

	if err := decoder.Decode(&result); err != nil {
		return DataNode{}, fmt.Errorf("%w", err)
	}

	return result, nil
}

func (n DataNode) Binary() ([]byte, error) {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)

	if err := encoder.Encode(n); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return buf.Bytes(), nil
}

func (n DataNode) String() string {
	result := &strings.Builder{}
	result.Grow(512) //nolint:gomnd

	result.WriteString(n.Key)
	result.WriteRune('=')

	toStringRepresentationBuffered(n.Data, result)

	return result.String()
}
