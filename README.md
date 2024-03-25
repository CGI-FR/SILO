# SILO

SILO (Sparse Input Linked Output) is an open-source command-line interface (CLI) tool designed for processing data in JSONLine format. It provides functionality to ingest data from standard input (stdin) and isolate entites (which are groups of related values) into a file, allowing users to create a referential of all entities discovered within the JSONLine data.

SILO can be used in addition to LINO and PIMO tools to generate consistency sources.

Here is an short example where SILO can be useful :

- TableA contains data about clients, and have these direct identifiers as columns : ID_CLIENT, EMAIL_CLIENT, ACCOUNT_NUMBER
- TableB contains data about clients, and have these direct identifiers as columns : ID_CLIENT, EMAIL_CLIENT

Unfortunately, available dataset contains a lot of duplication and null values.

**`TableA`**

| ID_CLIENT | EMAIL_CLIENT        | ACCOUNT_NUMBER |
| --------- | ------------------- | -------------- |
| 0001      | jonh.doe@domain.com |                |
|           |                     | C01            |

**`TableB`**

| ACCOUNT_NUMBER | EMAIL_CLIENT        |
| -------------- | ------------------- |
| C01            | jonh.doe@domain.com |

SILO will be able to generate the following referential.

**`Output of SILO`**

| UUID                                 | ID_CLIENT | EMAIL_CLIENT        | ACCOUNT_NUMBER |
| ------------------------------------ | --------- | ------------------- | -------------- |
| 79cc287b-3640-49c1-9e6a-86cff87cce41 | 0001      | jonh.doe@domain.com | C01            |

By leveraging SILO's capabilities, users can efficiently identify and link related records across disparate datasets, even in cases where direct identifiers are missing or duplicated.

## Installation

To install SILO, follow these steps:

1. Download the released tar.gz corresponding to your operating system
2. Extract the tar.gz
3. Optionnaly, move the `silo` binary to a shared path like `/usr/bin/silo`

## Usage

SILO provides two main commands:

### silo scan

The silo scan command is used to ingest data from stdin in JSONLine format, persisted on disk for future reference. Here's how to use it:

```console
$ silo scan my-silo < input.jsonl
⣾ Scanned 5 rows, found 15 links (4084 row/s) [0s]
```

Analysis data is persisted on disk on the `my-silo` path relative to the current directory.

### silo dump

The silo dump command is used to dump each connected entity into a file. This allows users to create a referential of all entities discovered within the JSONLine data. Here's how to use it:

```console
$ silo dump my-silo
{"uuid":"19bef352-ed87-4de8-a4ea-65f1d7db9ced","id":"ID1","key":2}
{"uuid":"19bef352-ed87-4de8-a4ea-65f1d7db9ced","id":"ID2","key":"2"}
{"uuid":"19bef352-ed87-4de8-a4ea-65f1d7db9ced","id":"ID3","key":2.2}
{"uuid":"19bef352-ed87-4de8-a4ea-65f1d7db9ced","id":"ID4","key":"00002"}
{"uuid":"60d7e970-ca56-410f-86f3-a6c1e67f032a","id":"ID2","key":"1"}
{"uuid":"60d7e970-ca56-410f-86f3-a6c1e67f032a","id":"ID4","key":"00001"}
{"uuid":"60d7e970-ca56-410f-86f3-a6c1e67f032a","id":"ID3","key":1.1}
{"uuid":"60d7e970-ca56-410f-86f3-a6c1e67f032a","id":"ID1","key":1}
{"uuid":"a628e8b5-69a7-4707-8f81-da2200ae1e1f","id":"ID2","key":"3"}
{"uuid":"a628e8b5-69a7-4707-8f81-da2200ae1e1f","id":"ID3","key":3.3}
{"uuid":"a628e8b5-69a7-4707-8f81-da2200ae1e1f","id":"ID4","key":"00003"}
{"uuid":"a628e8b5-69a7-4707-8f81-da2200ae1e1f","id":"ID1","key":3}
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

Copyright (C) 2024 CGI France

SILO is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

SILO is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
 along with SILO.  If not, see <http://www.gnu.org/licenses/>.
