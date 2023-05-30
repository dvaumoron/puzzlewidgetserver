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
var errNotFloat = errors.New("value is not an float")
var errNotMap = errors.New("value is not a map")
var errNotSlice = errors.New("value is not a slice")
var errNotString = errors.New("value is not a string")
var errFilesType = errors.New("field Files is not of the expected type")
var errEmptyUrl = errors.New("field CurrentUrl is empty")
var errNoUser = errors.New("field Id is 0")

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

func AsFloat64(value any) (float64, error) {
	if value == nil {
		return 0, nil
	}
	switch casted := value.(type) {
	case uint:
		return float64(casted), nil
	case uint8:
		return float64(casted), nil
	case uint16:
		return float64(casted), nil
	case uint32:
		return float64(casted), nil
	case uint64:
		return float64(casted), nil
	case int:
		return float64(casted), nil
	case int8:
		return float64(casted), nil
	case int16:
		return float64(casted), nil
	case int32:
		return float64(casted), nil
	case int64:
		return float64(casted), nil
	case float32:
		return float64(casted), nil
	case float64:
		return casted, nil
	case string:
		f, err := strconv.ParseFloat(casted, 64)
		if err != nil {
			return 0, err
		}
		return f, nil
	}
	return 0, errNotFloat
}

func GetFormData(data Data) (Data, error) {
	return AsMap(data[formKey])
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

func GetBaseUrl(levelToErase uint8, data Data) (string, error) {
	res, err := AsString(data[urlKey])
	if err != nil {
		return "", err
	}

	i := len(res) - 1
	if i == -1 {
		return "", errEmptyUrl
	}
	for count := uint8(0); count < levelToErase; {
		i--
		if res[i] == '/' {
			count++
		}
	}
	return res[:i+1], nil
}

func GetCurrentUserId(data Data) (uint64, error) {
	res, err := AsUint64(data[userKey])
	if err != nil {
		return 0, err
	}
	if res == 0 {
		return 0, errNoUser
	}
	return res, nil
}

func GetPaginationNames() []string {
	return []string{"pageNumber", "pageSize", "filter"}
}

func GetPagination(defaultPageSize uint64, data Data) (uint64, uint64, uint64, string) {
	pageNumber, _ := AsUint64(data["queryData/pageNumber"])
	if pageNumber == 0 {
		pageNumber = 1
	}
	pageSize, _ := AsUint64(data["queryData/pageSize"])
	if pageSize == 0 {
		pageSize = defaultPageSize
	}
	filter, _ := AsString(data["queryData/filter"])

	start := (pageNumber - 1) * pageSize
	end := start + pageSize

	return pageNumber, start, end, filter
}

func InitPagination(data Data, filter string, pageNumber uint64, end uint64, total uint64) {
	data["Filter"] = filter
	if pageNumber != 1 {
		data["PreviousPageNumber"] = pageNumber - 1
	}
	if end < total {
		data["NextPageNumber"] = pageNumber + 1
	}
	data["Total"] = total
}
