/*
 *
 * Copyright 2023 puzzlewidgetserver authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
package puzzlewidgetserver

import (
	"errors"
	"strconv"
)

var errNotInt = errors.New("value is not an int")
var errNotMap = errors.New("value is not a map")
var errFilesType = errors.New("field Files is not of the expected type")
var errNotSlice = errors.New("value is not a slice")
var errNotString = errors.New("value is not a string")

func AsMap(value any) (Data, error) {
	if value == nil {
		return nil, nil
	}
	m, ok := value.(Data)
	if !ok {
		return nil, errNotMap
	}
	return m, nil
}

func AsSlice(value any) ([]any, error) {
	if value == nil {
		return nil, nil
	}
	s, ok := value.([]any)
	if !ok {
		return nil, errNotSlice
	}
	return s, nil
}

func AsString(value any) (string, error) {
	if value == nil {
		return "", nil
	}
	s, ok := value.(string)
	if !ok {
		return "", errNotString
	}
	return s, nil
}

func AsUint64(value any) (uint64, error) {
	if value == nil {
		return 0, nil
	}
	switch casted := value.(type) {
	case uint:
		return uint64(casted), nil
	case uint8:
		return uint64(casted), nil
	case uint16:
		return uint64(casted), nil
	case uint32:
		return uint64(casted), nil
	case uint64:
		return uint64(casted), nil
	case int:
		return uint64(casted), nil
	case int8:
		return uint64(casted), nil
	case int16:
		return uint64(casted), nil
	case int32:
		return uint64(casted), nil
	case int64:
		return uint64(casted), nil
	case float32:
		return uint64(casted), nil
	case float64:
		return uint64(casted), nil
	case string:
		i, err := strconv.ParseUint(casted, 10, 64)
		if err != nil {
			return 0, err
		}
		return i, nil
	}
	return 0, errNotInt
}

func GetFiles(data Data) (map[string][]byte, error) {
	value := data[filesKey]
	if value == nil {
		return nil, nil
	}
	m, ok := value.(map[string][]byte)
	if !ok {
		return nil, errFilesType
	}
	return m, nil
}