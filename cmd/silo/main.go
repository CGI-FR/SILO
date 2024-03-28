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

package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/cgi-fr/silo/internal/app/cli"
	"github.com/mattn/go-isatty"
	"github.com/pkg/profile"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// Provisioned by ldflags.
var (
	name      string //nolint: gochecknoglobals
	version   string //nolint: gochecknoglobals
	commit    string //nolint: gochecknoglobals
	buildDate string //nolint: gochecknoglobals
	builtBy   string //nolint: gochecknoglobals

	verbosity string //nolint: gochecknoglobals
	jsonlog   bool   //nolint: gochecknoglobals
	debug     bool   //nolint: gochecknoglobals
	colormode string //nolint: gochecknoglobals
	profiling string //nolint: gochecknoglobals
)

func main() {
	var profiler interface{ Stop() }

	cobra.OnInitialize(initLog)

	rootCmd := &cobra.Command{ //nolint:exhaustruct
		Use:     name,
		Short:   "Sparse Input Linked Output",
		Long:    `SILO ingest data from stdin and isolate entites (which are groups of related values) into a file.`,
		Example: "  silo scan my-silo < data.jsonl\n  silo dump my-silo -w > entities.jsonl",
		Version: fmt.Sprintf(`%v (commit=%v date=%v by=%v)
Copyright (C) 2024 CGI France
License GPLv3: GNU GPL version 3 <https://gnu.org/licenses/gpl.html>.
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.`, version, commit, buildDate, builtBy),
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
			log.Info().
				Str("verbosity", verbosity).
				Bool("log-json", jsonlog).
				Bool("debug", debug).
				Str("color", colormode).
				Msg("start SILO")

			if profiling == "cpu" {
				profiler = profile.Start(profile.CPUProfile, profile.ProfilePath("."), profile.Quiet)
			} else if profiling == "memory" || profiling == "mem" {
				profiler = profile.Start(profile.MemProfile, profile.ProfilePath("."), profile.Quiet)
			}
		},
		PersistentPostRun: func(_ *cobra.Command, _ []string) {
			log.Info().Int("return", 0).Msg("end SILO")

			if profiling == "cpu" || profiling == "memory" || profiling == "mem" {
				profiler.Stop()
			}
		},
	}

	rootCmd.PersistentFlags().StringVarP(&verbosity, "verbosity", "v", "warn",
		"set level of log verbosity : none (0), error (1), warn (2), info (3), debug (4), trace (5)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "add debug information to logs (very slow)")
	rootCmd.PersistentFlags().BoolVar(&jsonlog, "log-json", false, "output logs in JSON format")
	rootCmd.PersistentFlags().StringVar(&colormode, "color", "auto", "use colors in log outputs : yes, no or auto")
	rootCmd.PersistentFlags().StringVar(&profiling, "profiling", "",
		"create a pprof file - use 'cpu' to create a CPU pprof file or 'mem' to create an memory pprof file")

	scanCmd := cli.NewScanCommand(name, os.Stderr, os.Stdout, os.Stdin)
	dumpCmd := cli.NewDumpCommand(name, os.Stderr, os.Stdout, os.Stdin)

	rootCmd.AddGroup(&cobra.Group{ID: "main", Title: "Main Commands:"})

	scanCmd.GroupID = "main"
	dumpCmd.GroupID = "main"

	rootCmd.AddCommand(scanCmd, dumpCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Err(err).Msg("error when executing command")
		os.Exit(1)
	}
}

func initLog() {
	color := false

	switch strings.ToLower(colormode) {
	case "auto":
		if isatty.IsTerminal(os.Stdout.Fd()) && runtime.GOOS != "windows" {
			color = true
		}
	case "yes", "true", "1", "on", "enable":
		color = true
	}

	if jsonlog {
		log.Logger = zerolog.New(os.Stderr)
	} else {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: !color}) //nolint:exhaustruct
	}

	if debug {
		log.Logger = log.Logger.With().Caller().Logger()
	}

	setVerbosity()
}

func setVerbosity() {
	switch verbosity {
	case "trace", "5":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug", "4":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info", "3":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn", "2":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error", "1":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}
}
