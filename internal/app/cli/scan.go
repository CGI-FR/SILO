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

func NewScanCommand(parent string, stderr *os.File, stdout *os.File, stdin *os.File) *cobra.Command {
	var (
		passthrough bool
		only        []string
		aliases     map[string]string
	)

	cmd := &cobra.Command{ //nolint:exhaustruct
		Use:     "scan path",
		Short:   "Ingest data from stdin and update silo database stored in given path",
		Example: "  lino pull database --table client | " + parent + " scan clients",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := scan(cmd, args[0], passthrough, only, aliases); err != nil {
				log.Fatal().Err(err).Int("return", 1).Msg("end SILO")
			}
		},
	}

	cmd.Flags().BoolVarP(&passthrough, "passthrough", "p", false, "pass stdin to stdout")
	cmd.Flags().StringSliceVarP(&only, "only", "o", []string{}, "only scan these columns, exclude all others")
	cmd.Flags().StringToStringVarP(&aliases, "alias", "a", map[string]string{}, "use given aliases for each columns")

	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetIn(stdin)

	return cmd
}

func scan(cmd *cobra.Command, path string, passthrough bool, only []string, aliases map[string]string) error {
	backend, err := infra.NewBackend(path)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	defer backend.Close()

	driver := silo.NewDriver(backend, nil, silo.WithKeys(only), silo.WithAliases(aliases))

	var reader silo.DataRowReader

	if passthrough {
		reader = infra.NewDataRowReaderWriterJSONLine(cmd.InOrStdin(), cmd.OutOrStdout())
	} else {
		reader, err = infra.NewDataRowReaderJSONLine()
	}

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if !passthrough {
		observer := infra.NewScanObserver()
		defer observer.Close()

		if err := driver.Scan(reader, observer); err != nil {
			return fmt.Errorf("%w", err)
		}
	} else if err := driver.Scan(reader); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
