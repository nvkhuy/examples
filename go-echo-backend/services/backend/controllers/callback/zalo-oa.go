package webhook

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
	"log"
)

// CallbackZaloOA Get settings taxes
// @Tags Common
// @Summary Get settings taxes
// @Description Get settings taxes
// @Accept  json
// @Produce  json
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/callback/zalo/permission [get]
func CallbackZaloOA(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	type Params struct {
		OAID string `json:"oa_id" param:"oa_id" query:"oa_id" form:"oa_id"`
		Code string `json:"code" param:"code" query:"code" form:"code"`
	}
	var params Params
	var err = cc.Bind(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}
	log.Println("zalo-oa")
	log.Println(params)

	return cc.Success(params)
}
