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
	"os"
	"time"

	"github.com/cgi-fr/silo/pkg/silo"
	"github.com/schollz/progressbar/v3"
)

type DumpObserver struct {
	countTotal        int
	countComplete     int
	countConsistent   int
	countInconsistent int
	countEmpty        int
	bar               *progressbar.ProgressBar
}

func NewDumpObserver() *DumpObserver {
	//nolint:gomnd
	pgb := progressbar.NewOptions(-1,
		progressbar.OptionSetDescription("Dumping ... "),
		progressbar.OptionSetItsString("entity"),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowIts(),
		progressbar.OptionSpinnerType(11),
		progressbar.OptionThrottle(time.Millisecond*10),
		progressbar.OptionOnCompletion(func() { fmt.Fprintln(os.Stderr) }),
		// progressbar.OptionShowDescriptionAtLineEnd(),
	)

	return &DumpObserver{
		countTotal:        0,
		countComplete:     0,
		countConsistent:   0,
		countInconsistent: 0,
		countEmpty:        0,
		bar:               pgb,
	}
}

func (o *DumpObserver) Entity(status silo.Status, _ map[string]int) {
	o.countTotal++

	switch status {
	case silo.StatusEntityComplete:
		o.countComplete++
	case silo.StatusEntityConsistent:
		o.countConsistent++
	case silo.StatusEntityInconsistent:
		o.countInconsistent++
	case silo.StatusEntityEmpty:
		o.countEmpty++
	}

	_ = o.bar.Add(1)

	o.bar.Describe(fmt.Sprintf("Dumped %d entities / complete=%d / incomplete=%d / inconsistent=%d / empty=%d",
		o.countTotal,
		o.countComplete,
		o.countConsistent,
		o.countInconsistent,
		o.countEmpty))
}

func (o *DumpObserver) Close() {
	_ = o.bar.Close()
}
