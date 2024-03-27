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

type Option interface {
	applier
}

type option func(*config) error

func (f option) apply(cfg *config) error {
	return f(cfg)
}

type applier interface {
	apply(cfg *config) error
}

func Alias(key, alias string) Option { //nolint:ireturn
	applier := func(cfg *config) error {
		cfg.aliases[key] = alias

		return nil
	}

	return option(applier)
}

func Include(key string) Option { //nolint:ireturn
	applier := func(cfg *config) error {
		if _, exist := cfg.include[key]; !exist {
			cfg.includeList = append(cfg.includeList, key)
		}

		cfg.include[key] = true

		return nil
	}

	return option(applier)
}

func WithAliases(aliases map[string]string) Option { //nolint:ireturn
	applier := func(cfg *config) error {
		for key, alias := range aliases {
			cfg.aliases[key] = alias
		}

		return nil
	}

	return option(applier)
}

func WithKeys(keys []string) Option { //nolint:ireturn
	applier := func(cfg *config) error {
		for _, key := range keys {
			if _, exist := cfg.include[key]; !exist {
				cfg.includeList = append(cfg.includeList, key)
			}

			cfg.include[key] = true
		}

		return nil
	}

	return option(applier)
}
