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
	"errors"
	"fmt"
)

var (
	ErrReadingNextInput     = errors.New("error while reading next input")
	ErrPersistingData       = errors.New("error while persisting data")
	ErrReadingPersistedData = errors.New("error while reading persisted data")
)

type ConfigScanAliasIsNotIncludedError struct {
	alias string
}

func (e *ConfigScanAliasIsNotIncludedError) Error() string {
	return fmt.Sprintf("configuration error : alias [%s] is not included", e.alias)
}
