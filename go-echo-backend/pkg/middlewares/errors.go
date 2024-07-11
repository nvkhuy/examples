package middlewares

import (
	"encoding/json"
	"net/http"

	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/models"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgconn"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/rotisserie/eris"
	"github.com/stripe/stripe-go/v74"
)

// CustomErrorHandler error handler
func (m *Middleware) CustomErrorHandler(er error, c echo.Context) {
	var err = eris.Cause(er)
	var header = http.StatusBadRequest
	var message interface{} = "An error occured"
	var detail interface{} = err.Error()
	var code = header

	m.App.DB.CustomLogger.Errorf("*** error type=%T err=%+v", err, err)

	switch e := err.(type) {
	case *stripe.CardError:
		header = http.StatusBadRequest
		code = header
		detail = e
		message = e.Error()
		if e.Error() != "" {
			var stripeErr stripe.Error
			if e := json.Unmarshal([]byte(e.Error()), &stripeErr); e == nil && stripeErr.Msg != "" {
				message = stripeErr.Msg
			}
		}
	case *stripe.Error:
		header = http.StatusBadRequest
		code = e.HTTPStatusCode
		detail = e
		message = e.Msg
		if e.Msg != "" {
			message = e.Error()
		}

	case *validator.InvalidValidationError:
		header = http.StatusUnprocessableEntity
		code = header
		detail = e
		message = e.Error()
	case validator.ValidationErrors:
		header = http.StatusUnprocessableEntity
		code = header
		detail = e
		message = "Unprocessable Entity"
	case *echo.HTTPError:
		header = e.Code
		code = header
		detail = e
		message = e.Message
	case *pq.Error:
		header = http.StatusBadRequest
		code = header
		detail = e
		message = "Internal Server Error"
	case *pgconn.PgError:
		header = http.StatusBadRequest
		code = header
		detail = e
		message = "Internal Server Error"
	case *errs.Error:
		header = e.Header
		code = e.Code
		detail = e.Detail
		message = e.Error()

	}

	var response = models.M{
		"code":    code,
		"message": message,
	}

	if resp := c.Response(); resp != nil {
		response["request_id"] = resp.Header().Get(echo.HeaderXRequestID)
	}

	if !m.App.Config.IsProd() {
		response["detail"] = detail
		response["message"] = message

	}

	c.JSON(header, response)
}
