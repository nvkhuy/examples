package casbin

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/util"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

func New(config *config.Configuration, db *gorm.DB) *casbin.Enforcer {
	var a *gormadapter.Adapter
	var err error
	a, err = gormadapter.NewAdapterByDB(db)
	if err != nil {
		panic(eris.Wrap(err, "Config casbin - NewAdapter error"))
	}
	var e *casbin.Enforcer
	e, err = casbin.NewEnforcer(config.CasbinModelConfURL, a)
	e.AddNamedMatchingFunc("g", "KeyMatch", util.KeyMatch)
	e.AddNamedMatchingFunc("g2", "KeyMatch5", util.KeyMatch5)
	if err != nil {
		panic(eris.Wrap(err, "Config casbin error"))
	}
	_ = e.LoadPolicy()
	return e
}
