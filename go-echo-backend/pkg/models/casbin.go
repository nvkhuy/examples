package models

import (
	"github.com/casbin/casbin/v2"
	"strings"
)

type casbinGroupAction struct {
	Path   string `json:"path"`
	Action string `json:"action"`
}
type casbinGroup struct {
	Name    string               `json:"name"`
	Actions []*casbinGroupAction `json:"actions"`
}
type CasbinResource struct {
	Role   string         `json:"role"`
	Team   string         `json:"team"`
	Groups []*casbinGroup `json:"groups"`
}
type CasbinResources []CasbinResource

func ListResources(e *casbin.Enforcer, roles ...string) (res CasbinResources) {
	res = make(CasbinResources, len(roles))
	for i, role := range roles {
		split := strings.Split(role, ":")
		if len(split) > 0 {
			res[i].Role = split[0]
		}
		if len(split) > 1 {
			res[i].Team = split[1]
		}
		policySlice := e.GetFilteredPolicy(0, role)
		groupIndex := make(map[string]int)
		for _, pl := range policySlice { // pl -> role, group, action
			groupName := pl[1]
			groupAction := pl[2]
			groups := e.GetFilteredNamedGroupingPolicy("g2", 1, groupName)
			for _, gr := range groups { // gr -> path, group
				grId, ok := groupIndex[pl[1]]
				pathName := gr[0]
				if ok {
					grv := res[i].Groups[grId]
					grv.Actions = append(grv.Actions, &casbinGroupAction{
						Path:   pathName,
						Action: groupAction,
					})
				} else {
					res[i].Groups = append(res[i].Groups, &casbinGroup{
						Name: groupName,
						Actions: []*casbinGroupAction{
							{
								Path:   pathName,
								Action: groupAction,
							},
						},
					})
					groupIndex[pl[1]] = len(res[i].Groups) - 1
				}
			}

		}
	}
	return
}
