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

package silo_test

import (
	"testing"

	"github.com/cgi-fr/silo/pkg/silo"
	"github.com/stretchr/testify/require"
)

func TestNominal(t *testing.T) {
	t.Parallel()

	rows := []silo.DataRow{
		{"ID1": 1, "ID2": "1", "ID3": 1.10, "ID4": "00001"},
		{"ID1": 2, "ID2": "2", "ID3": 2.20, "ID4": "00002"},
	}
	input := silo.NewDataRowReaderInMemory(rows)

	backend := silo.NewBackendInMemory()
	writer := silo.NewDumpToStdout()
	driver := silo.NewDriver(backend, writer)

	err := driver.Scan(input)
	require.NoError(t, err)

	require.NoError(t, driver.Dump())
}

func TestPartialNull(t *testing.T) {
	t.Parallel()

	rows := []silo.DataRow{
		{"ID1": 1, "ID2": nil, "ID3": nil, "ID4": "00001"},
		{"ID1": nil, "ID2": "2", "ID3": 2.0, "ID4": nil},
	}
	input := silo.NewDataRowReaderInMemory(rows)

	backend := silo.NewBackendInMemory()
	writer := silo.NewDumpToStdout()
	driver := silo.NewDriver(backend, writer)

	err := driver.Scan(input)
	require.NoError(t, err)

	require.NoError(t, driver.Dump())
}

func TestPartialMissing(t *testing.T) {
	t.Parallel()

	rows := []silo.DataRow{
		{"ID1": 1, "ID4": "00001"},
		{"ID2": "2", "ID3": 2.0},
	}
	input := silo.NewDataRowReaderInMemory(rows)

	backend := silo.NewBackendInMemory()
	writer := silo.NewDumpToStdout()
	driver := silo.NewDriver(backend, writer)

	err := driver.Scan(input)
	require.NoError(t, err)

	require.NoError(t, driver.Dump())
}
