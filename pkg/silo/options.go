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

type option func(*Config) error

func (f option) apply(cfg *Config) error {
	return f(cfg)
}

type applier interface {
	apply(cfg *Config) error
}

func Alias(key, alias string) Option { //nolint:ireturn
	applier := func(cfg *Config) error {
		cfg.Aliases[key] = alias

		return nil
	}

	return option(applier)
}

func Include(key string) Option { //nolint:ireturn
	applier := func(cfg *Config) error {
		cfg.Include[key] = true

		return nil
	}

	return option(applier)
}

func WithAliases(aliases map[string]string) Option { //nolint:ireturn
	applier := func(cfg *Config) error {
		for key, alias := range aliases {
			cfg.Aliases[key] = alias
		}

		return nil
	}

	return option(applier)
}

func WithKeys(keys []string) Option { //nolint:ireturn
	applier := func(cfg *Config) error {
		for _, key := range keys {
			cfg.Include[key] = true
		}

		return nil
	}

	return option(applier)
}
