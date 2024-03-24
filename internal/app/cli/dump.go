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

package cli

import (
	"fmt"
	"os"

	"github.com/cgi-fr/silo/internal/infra"
	"github.com/cgi-fr/silo/pkg/silo"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func NewDumpCommand(parent string, stderr *os.File, stdout *os.File, stdin *os.File) *cobra.Command {
	cmd := &cobra.Command{ //nolint:exhaustruct
		Use:     "dump path",
		Short:   "Dump silo database stored in given path into stdout",
		Example: "  " + parent + " dump clients",
		Args:    cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			if err := dump(args[0]); err != nil {
				log.Fatal().Err(err).Int("return", 1).Msg("end SILO")
			}
		},
	}

	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetIn(stdin)

	return cmd
}

func dump(path string) error {
	backend, err := infra.NewBackend(path)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	defer backend.Close()

	driver := silo.NewDriver(backend, infra.NewDumpJSONLine())

	if err := driver.Dump(); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
