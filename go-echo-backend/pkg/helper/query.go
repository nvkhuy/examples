package helper

import (
	"fmt"
	"strings"
)

func ParseQueryStatus(status string) []string {
	if status == "" {
		return []string{}
	}

	var segments = strings.Split(status, ",")

	return segments
}

func ParseQuerySort(status string) [][]string {
	if status == "" {
		return [][]string{}
	}

	var segments = strings.Split(status, ",")
	var listStatus = [][]string{}

	for _, s := range segments {
		var segs = strings.Split(s, ":")
		if len(segs) == 2 {
			listStatus = append(listStatus, segs)
		}
	}

	return listStatus
}

func ParseQuerySQL(colName string, values []string) string {
	var where = ""
	for index, g := range values {
		if index == 0 {
			where = fmt.Sprintf("%s = '%s'", colName, g)
		} else {
			where = fmt.Sprintf("%s OR %s = '%s'", where, colName, g)
		}
	}

	return where
}

func ParseQueryWhere(columns []string) string {
	var newColumns []string
	for _, column := range columns {
		newColumns = append(newColumns, fmt.Sprintf("'%s'", column))
	}

	if len(columns) > 1 {
		return fmt.Sprintf("%v", strings.Join(newColumns, ","))
	}
	return strings.Join(newColumns, ",")
}

func ParseOrderBy(prefix string, matrix [][]string, excludePrefixForKeys ...string) []string {
	var orderBy = []string{}
	if len(matrix) > 0 {
		for _, s := range matrix {
			if len(s) == 2 {
				var key = s[0]
				var skipPrefix = false
				if len(excludePrefixForKeys) > 0 {
					skipPrefix = StringContains(excludePrefixForKeys, key)

				}

				if !skipPrefix && prefix != "" {
					key = fmt.Sprintf("%s.%s", prefix, key)
				}
				orderBy = append(orderBy, fmt.Sprintf("%s %s", key, s[1]))
			}
		}
	}
	return orderBy
}
