package tests

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/util"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"testing"
)

func TestExampleRBACModel(t *testing.T) {
	e, _ := casbin.NewEnforcer("../static/model_test.conf", "../static/policy_test.csv")

	e.AddNamedMatchingFunc("g", "KeyMatch", util.KeyMatch)
	e.AddNamedMatchingFunc("g2", "KeyMatch5", util.KeyMatch5)

	testEnforce(t, e, "staff:marketing", "/api/v1/admin/me/track_activity", "POST", true)
	// Super Admin
	testEnforce(t, e, "super_admin", "/any-path", "/any-action", true)

	// Sales
	testEnforce(t, e, "leader:sales", "/any-path", "/any-action", false)
	testEnforce(t, e, "leader:sales", "/sales/create", "POST", true)
	testEnforce(t, e, "staff:sales", "/sales/create", "POST", false)
	testEnforce(t, e, "staff:sales", "/sales/1", "GET", true)
	testEnforce(t, e, "leader:sales", "/sales/write-but-for-staff", "PUT", true)
	testEnforce(t, e, "staff:sales", "/sales/write-but-for-staff", "PUT", true)
	testEnforce(t, e, "staff:sales", "/sales/write-but-not-for-staff", "PUT", false)

	// Dev
	testEnforce(t, e, "leader:dev", "/sales/any", "PUT", true)
	testEnforce(t, e, "leader:dev", "/any-path", "PUT", true)
	testEnforce(t, e, "staff:dev", "/any-path", "PUT", true)

	// Marketing
	testEnforce(t, e, "leader:marketing", "/blog/any", "PUT", true)
	testEnforce(t, e, "staff:marketing", "/blog/any", "GET", true)

	testEnforce(t, e, "staff:marketing", "/api/v1/admin/notifications", "GET", true)
	testEnforce(t, e, "staff:marketing", "/api/v1/admin/notifications?page=1&limit=12", "GET", true)
	testEnforce(t, e, "staff:marketing", "/api/v1/admin/notifications/any/any", "GET", true)
	testEnforce(t, e, "staff:marketing", "/api/v1/admin/notifications/1", "PUT", false)

	// Operator
	testEnforce(t, e, "leader:operator", "/settings", "GET", true)
	testEnforce(t, e, "leader:operator", "/settings/any", "PUT", true)
	testEnforce(t, e, "leader:operator", "/settings/any-1/any-2", "PUT", true)
	testEnforce(t, e, "staff:operator", "/settings", "GET", true)
	testEnforce(t, e, "staff:operator", "/settings/any", "PUT", false)
	testEnforce(t, e, "staff:operator", "/settings/any-1/any-2", "PUT", false)

	testEnforce(t, e, "super_admin", "/api/v1/admin/resources", "GET", true)
	testEnforce(t, e, "super_admin", "/api/v1/admin/posts/1/archive", "DELETE", true)
}

func TestEnforceList(t *testing.T) {
	e, _ := casbin.NewEnforcer("../static/model_test.conf", "../static/policy_test.csv")
	e.AddNamedMatchingFunc("g2", "KeyMatch2", util.KeyMatch2)
	all := e.GetAllNamedSubjects("p")
	res := models.ListResources(e, all...)
	//res := models.ListResources(e, "staff:sales")
	helper.PrintJSON(res)
}

func testEnforce(t *testing.T, e *casbin.Enforcer, sub interface{}, obj interface{}, act string, res bool) {
	t.Helper()
	if myRes, err := e.Enforce(sub, obj, act); err != nil {
		t.Errorf("Enforce Error: %s", err)
	} else if myRes != res {
		t.Errorf("%s, %v, %s: %t, supposed to be %t", sub, obj, act, myRes, res)
	}
}

//api/v1/seo/:root -> home, ...
//column -> string ,metadata -> []object{json}
