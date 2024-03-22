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
	"fmt"
	"os"

	"github.com/cgi-fr/silo/pkg/silo"
	"github.com/goccy/go-json"
)

type DataRowReaderJSONLine struct {
	decoder *json.Decoder
}

func NewDataRowReaderJSONLine() (*DataRowReaderJSONLine, error) {
	return &DataRowReaderJSONLine{decoder: json.NewDecoder(os.Stdin)}, nil
}

func NewDataRowReaderJSONLineFromFile(filename string) (*DataRowReaderJSONLine, error) {
	source, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &DataRowReaderJSONLine{decoder: json.NewDecoder(source)}, nil
}

func (drr *DataRowReaderJSONLine) ReadDataRow() (silo.DataRow, error) {
	if drr.decoder.More() {
		data := silo.DataRow{}
		if err := drr.decoder.Decode(&data); err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		return data, nil
	}

	return nil, nil
}

func (drr *DataRowReaderJSONLine) Close() error {
	return nil
}
