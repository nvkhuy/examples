package models

import "sort"

func (pc *UserProductClasses) Uniq() UserProductClasses {
	sort.Slice(*pc, func(i, j int) bool {
		x, y := (*pc)[i], (*pc)[j]
		return x.Conf > y.Conf
	})
	m := make(map[string]UserProductClass)
	for _, p := range *pc {
		if _, ok := m[p.Class]; !ok {
			m[p.Class] = p
		}
	}
	var classes UserProductClasses
	for _, v := range m {
		classes = append(classes, v)
	}
	sort.Slice(classes, func(i, j int) bool {
		return classes[i].Conf > classes[j].Conf
	})
	return classes
}
