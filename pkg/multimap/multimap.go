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

package multimap

type Multimap[K, V comparable] map[K]map[V]int

// Add a key/value pair to the multimap.
func (m Multimap[K, V]) Add(key K, value V) {
	set, ok := m[key]
	if !ok {
		set = make(map[V]int)
	}

	set[value]++

	m[key] = set
}

// Delete values associated to key.
func (m Multimap[K, V]) Delete(key K) []V {
	set, ok := m[key]
	if !ok {
		return []V{}
	}

	values := make([]V, 0, len(set))
	for value := range set {
		values = append(values, value)
	}

	delete(m, key)

	return values
}

// Get values associated to key.
func (m Multimap[K, V]) Get(key K) []V {
	set, ok := m[key]
	if !ok {
		return []V{}
	}

	values := make([]V, 0, len(set))
	for value := range set {
		values = append(values, value)
	}

	return values
}

// Get a random key in the multimap.
func (m Multimap[K, V]) RandomKey() (K, bool) { //nolint:ireturn
	for key := range m {
		return key, true
	}

	return *new(K), false
}

// Count the number of values associated to key.
func (m Multimap[K, V]) Count(key K) int {
	return len(m[key])
}

// Copy the multimap.
func (m Multimap[K, V]) Copy() Multimap[K, V] {
	multimapCopy := make(Multimap[K, V], len(m))

	for key, values := range m {
		valuesCopy := make(map[V]int, len(values))
		for value, count := range values {
			valuesCopy[value] = count
		}

		multimapCopy[key] = valuesCopy
	}

	return multimapCopy
}
