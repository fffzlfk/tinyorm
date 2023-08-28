package utils

import "strings"

func Difference(a []string, b []string, capsSensitive bool) (diff []string) {
	mapB := make(map[string]struct{})
	for _, v := range b {
		if !capsSensitive {
			v = strings.ToLower(v)
		}
		mapB[v] = struct{}{}
	}
	for _, v := range a {
		var lowerV string
		if !capsSensitive {
			lowerV = strings.ToLower(v)
		}
		if _, ok := mapB[lowerV]; !ok {
			diff = append(diff, v)
		}
	}
	return
}
