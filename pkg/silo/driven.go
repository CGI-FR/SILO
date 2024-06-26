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

type DataRowReader interface {
	ReadDataRow() (DataRow, error)
	Close() error
}

type Backend interface {
	Store(key DataNode, value DataNode) error
	Snapshot() Snapshot
	Close() error
}

type Snapshot interface {
	Next() (DataNode, bool, error)
	PullAll(node DataNode) ([]DataNode, error)
	Close() error
}

type DumpWriter interface {
	Write(node DataNode, uuid string) error
	Close() error
}

type ScanObserver interface {
	IngestedRow(row DataRow)
	IngestedLink(link DataLink)
}

type DumpObserver interface {
	Entity(status Status, counts map[string]int)
}
