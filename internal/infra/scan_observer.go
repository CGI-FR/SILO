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
	"github.com/schollz/progressbar/v3"
)

type ScanObserver struct {
	rowCount  int
	linkCount int
	bar       *progressbar.ProgressBar
}

func NewScanObserver() *ScanObserver {
	//nolint:gomnd
	pgb := progressbar.NewOptions(-1,
		progressbar.OptionSetDescription("Scanning ... "),
		progressbar.OptionSetItsString("row"),
		// progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionSpinnerType(11),
		progressbar.OptionThrottle(time.Millisecond*10),
	)

	return &ScanObserver{
		rowCount:  0,
		linkCount: 0,
		bar:       pgb,
	}
}

func (o *ScanObserver) IngestedRow(_ silo.DataRow) {
	o.rowCount++
	_ = o.bar.Add(1)
	o.bar.Describe(fmt.Sprintf("Scanned %d rows, found %d links", o.rowCount, o.linkCount))
}

func (o *ScanObserver) IngestedLink(_ silo.DataLink) {
	o.linkCount++
	_ = o.bar.Add(1)
	o.bar.Describe(fmt.Sprintf("Scanned %d rows, found %d links", o.rowCount, o.linkCount))
}

func (o *ScanObserver) Close() {
	_ = o.bar.Close()
}
