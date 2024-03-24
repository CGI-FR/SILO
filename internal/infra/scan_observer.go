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

package infra

import (
	"fmt"
	"time"

	"github.com/cgi-fr/silo/pkg/silo"
	"github.com/gosuri/uilive"
)

type ScanObserver struct {
	rowCount, linkCount int
	output              *uilive.Writer
}

func NewScanObserver() *ScanObserver {
	output := uilive.New()
	output.Start()
	output.RefreshInterval = time.Second

	return &ScanObserver{
		rowCount:  0,
		linkCount: 0,
		output:    output,
	}
}

func (o *ScanObserver) IngestedRow(_ silo.DataRow) {
	o.rowCount++
	fmt.Fprintf(o.output, "Ingested (links/rows) : %d / %d\n", o.linkCount, o.rowCount)
	o.output.Flush()
}

func (o *ScanObserver) IngestedLink(_ silo.DataLink) {
	o.linkCount++
	fmt.Fprintf(o.output, "Ingested (links/rows) : %d / %d\n", o.linkCount, o.rowCount)
	o.output.Flush()
}

func (o *ScanObserver) Close() {
	o.output.Stop()
}
