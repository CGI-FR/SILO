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
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func toStringRepresentationBuffered(value any, stringbuffer *strings.Builder) {
	switch tvalue := value.(type) {
	case string:
		stringbuffer.WriteString("string(")
		stringbuffer.WriteString(tvalue)
		stringbuffer.WriteByte(')')
	case bool:
		stringbuffer.WriteString("bool(")
		stringbuffer.WriteString(strconv.FormatBool(tvalue))
		stringbuffer.WriteByte(')')
	case float64:
		stringbuffer.WriteString("number(")
		stringbuffer.WriteString(strconv.FormatFloat(tvalue, 'g', -1, 64))
		stringbuffer.WriteByte(')')
	case int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8:
		stringbuffer.WriteString("number(")
		stringbuffer.WriteString(fmt.Sprint(tvalue))
		stringbuffer.WriteByte(')')
	case json.Number:
		stringbuffer.WriteString("number(")
		stringbuffer.WriteString(string(tvalue))
		stringbuffer.WriteByte(')')
	case []any:
		stringbuffer.WriteString("slice(")

		for _, value := range tvalue {
			toStringRepresentationBuffered(value, stringbuffer)
		}

		stringbuffer.WriteByte(')')
	case nil:
		stringbuffer.WriteString("nil(nil)")
	}
}
