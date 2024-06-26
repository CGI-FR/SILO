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

import "errors"

type config struct {
	include     map[string]bool
	includeList []string
	aliases     map[string]string
}

func newConfig() *config {
	config := config{
		include:     map[string]bool{},
		includeList: []string{},
		aliases:     map[string]string{},
	}

	return &config
}

func (cfg *config) validate() error {
	var errs []error

	for key := range cfg.aliases {
		if _, ok := cfg.include[key]; !ok && len(cfg.include) > 0 {
			errs = append(errs, &ConfigScanAliasIsNotIncludedError{alias: key})
		}
	}

	if len(errs) != 0 {
		return errors.Join(errs...)
	}

	return nil
}
