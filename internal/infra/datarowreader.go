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
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/cgi-fr/silo/pkg/silo"
	"github.com/goccy/go-json"
)

const linebreak byte = 10

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

type DataRowReaderWriterJSONLine struct {
	input  *bufio.Scanner
	output *bufio.Writer
}

func NewDataRowReaderWriterJSONLine(input io.Reader, output io.Writer) *DataRowReaderWriterJSONLine {
	return &DataRowReaderWriterJSONLine{input: bufio.NewScanner(input), output: bufio.NewWriter(output)}
}

func (drr *DataRowReaderWriterJSONLine) ReadDataRow() (silo.DataRow, error) {
	if drr.input.Scan() {
		if err := drr.writeLine(); err != nil {
			return nil, err
		}

		data := silo.DataRow{}
		if err := json.UnmarshalNoEscape(drr.input.Bytes(), &data); err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		return data, nil
	}

	if err := drr.input.Err(); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return nil, nil
}

func (drr *DataRowReaderWriterJSONLine) writeLine() error {
	if _, err := drr.output.Write(drr.input.Bytes()); err != nil {
		return fmt.Errorf("%w", err)
	}

	if err := drr.output.WriteByte(linebreak); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (drr *DataRowReaderWriterJSONLine) Close() error {
	if drr.output == nil {
		return nil
	}

	return fmt.Errorf("%w", drr.output.Flush())
}
