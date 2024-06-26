# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Types of changes

- `Added` for new features.
- `Changed` for changes in existing functionality.
- `Deprecated` for soon-to-be removed features.
- `Removed` for now removed features.
- `Fixed` for any bug fixes.
- `Security` in case of vulnerabilities.

## [0.3.0]

- `Added` cpu and memory profiling with `--profiling mem|cpu` flag
- `Fixed` performance issues on dump in exchange for higher RAM consumption, using `--limited-ram` flag will fall back to the 0.2.0 dump version

## [0.2.0]

- `Added` flag `--include` (short `-i`) to only scan/dump a specific list of fields, this flag is repeatable
- `Added` flag `--alias` (short `-a`) to rename fields on the fly, this flag is repeatable
- `Added` flag `--watch` (short `-w`) to the dump command
- `Fixed` self reference link are no longer counted in the links counter while scanning

## [0.1.0]

- `Added` initial version
