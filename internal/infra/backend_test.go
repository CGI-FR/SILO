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

package infra_test

import (
	"os"
	"testing"

	"github.com/cgi-fr/silo/internal/infra"
	"github.com/cgi-fr/silo/pkg/silo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNominal(t *testing.T) {
	t.Parallel()

	backend, err := infra.NewBackend("silo-pebble")

	require.NoError(t, err)

	defer os.RemoveAll("silo-pebble")
	defer backend.Close()

	require.NoError(t, backend.Store(silo.DataNode{Key: "ID1", Data: "1"}, silo.DataNode{Key: "ID2", Data: "1"}))
	require.NoError(t, backend.Store(silo.DataNode{Key: "ID2", Data: "1"}, silo.DataNode{Key: "ID1", Data: "1"}))
	require.NoError(t, backend.Store(silo.DataNode{Key: "ID1", Data: "1"}, silo.DataNode{Key: "ID3", Data: "1"}))
	require.NoError(t, backend.Store(silo.DataNode{Key: "ID3", Data: "1"}, silo.DataNode{Key: "ID1", Data: "1"}))
	require.NoError(t, backend.Store(silo.DataNode{Key: "ID1", Data: "1"}, silo.DataNode{Key: "ID4", Data: "1"}))
	require.NoError(t, backend.Store(silo.DataNode{Key: "ID4", Data: "1"}, silo.DataNode{Key: "ID1", Data: "1"}))
	require.NoError(t, backend.Store(silo.DataNode{Key: "ID2", Data: "1"}, silo.DataNode{Key: "ID3", Data: "1"}))
	require.NoError(t, backend.Store(silo.DataNode{Key: "ID3", Data: "1"}, silo.DataNode{Key: "ID2", Data: "1"}))
	require.NoError(t, backend.Store(silo.DataNode{Key: "ID2", Data: "1"}, silo.DataNode{Key: "ID4", Data: "1"}))
	require.NoError(t, backend.Store(silo.DataNode{Key: "ID4", Data: "1"}, silo.DataNode{Key: "ID2", Data: "1"}))
	require.NoError(t, backend.Store(silo.DataNode{Key: "ID3", Data: "1"}, silo.DataNode{Key: "ID4", Data: "1"}))
	require.NoError(t, backend.Store(silo.DataNode{Key: "ID4", Data: "1"}, silo.DataNode{Key: "ID3", Data: "1"}))

	require.NoError(t, backend.Store(silo.DataNode{Key: "ID1", Data: "2"}, silo.DataNode{Key: "ID2", Data: "2"}))
	require.NoError(t, backend.Store(silo.DataNode{Key: "ID2", Data: "2"}, silo.DataNode{Key: "ID1", Data: "2"}))
	require.NoError(t, backend.Store(silo.DataNode{Key: "ID1", Data: "2"}, silo.DataNode{Key: "ID3", Data: "2"}))
	require.NoError(t, backend.Store(silo.DataNode{Key: "ID3", Data: "2"}, silo.DataNode{Key: "ID1", Data: "2"}))
	require.NoError(t, backend.Store(silo.DataNode{Key: "ID1", Data: "2"}, silo.DataNode{Key: "ID4", Data: "2"}))
	require.NoError(t, backend.Store(silo.DataNode{Key: "ID4", Data: "2"}, silo.DataNode{Key: "ID1", Data: "2"}))
	require.NoError(t, backend.Store(silo.DataNode{Key: "ID2", Data: "2"}, silo.DataNode{Key: "ID3", Data: "2"}))
	require.NoError(t, backend.Store(silo.DataNode{Key: "ID3", Data: "2"}, silo.DataNode{Key: "ID2", Data: "2"}))
	require.NoError(t, backend.Store(silo.DataNode{Key: "ID2", Data: "2"}, silo.DataNode{Key: "ID4", Data: "2"}))
	require.NoError(t, backend.Store(silo.DataNode{Key: "ID4", Data: "2"}, silo.DataNode{Key: "ID2", Data: "2"}))
	require.NoError(t, backend.Store(silo.DataNode{Key: "ID3", Data: "2"}, silo.DataNode{Key: "ID4", Data: "2"}))
	require.NoError(t, backend.Store(silo.DataNode{Key: "ID4", Data: "2"}, silo.DataNode{Key: "ID3", Data: "2"}))

	id1, err := backend.Get(silo.DataNode{Key: "ID1", Data: "1"})

	require.NoError(t, err)
	assert.Len(t, id1, 3)

	snapshot := backend.Snapshot()

	defer snapshot.Close()

	next, ok, err := snapshot.Next()

	require.NoError(t, err)
	require.True(t, ok)

	idnext, err := snapshot.PullAll(next)

	require.NoError(t, err)
	assert.Len(t, idnext, 3)
}
