package middlewares

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/labstack/echo/v4"
)

type bodyDumpResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *bodyDumpResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

func (w *bodyDumpResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w *bodyDumpResponseWriter) Flush() {
	w.ResponseWriter.(http.Flusher).Flush()
}

func (w *bodyDumpResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

func (m *Middleware) Logger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if regexp.MustCompile("/ws").MatchString(c.Request().RequestURI) {
				return next(c)
			}

			var cc = c.(*models.CustomContext)
			var req = c.Request()
			var res = c.Response()
			var start = time.Now()
			var fields []zapcore.Field

			var method = cc.Request().Method
			var needPrintRequest = method != http.MethodGet
			var needPrintResponse = needPrintRequest
			var msg = ""

			if needPrintRequest {
				// Request
				reqBody := []byte{}
				if c.Request().Body != nil { // Read
					reqBody, _ = io.ReadAll(c.Request().Body)
				}
				c.Request().Body = io.NopCloser(bytes.NewBuffer(reqBody)) // Reset

				if len(reqBody) > 0 {
					if json.Valid(reqBody) {
						fields = append(fields, zap.Any(logger.ApiReqDataKey, json.RawMessage(reqBody)))
					} else {
						fields = append(fields, zap.Any(logger.ApiResDataKey, string(reqBody)))
					}
				}

			}

			if needPrintResponse {
				resBody := new(bytes.Buffer)
				mw := io.MultiWriter(c.Response().Writer, resBody)
				writer := &bodyDumpResponseWriter{Writer: mw, ResponseWriter: c.Response().Writer}
				c.Response().Writer = writer

				if resBody.Len() > 0 {
					var data = resBody.Bytes()
					if json.Valid(data) {
						fields = append(fields, zap.Any(logger.ApiResDataKey, json.RawMessage(data)))
					} else {
						fields = append(fields, zap.Any(logger.ApiResDataKey, string(data)))
					}

				}

			}

			if err = next(c); err != nil {
				errMsg, stack := errs.ParseErisJSON(err)
				msg = errMsg
				fields = append(fields, zap.Any("stack", stack))
				c.Error(err)
			}

			fields = append(fields, cc.GetHttpRequestInfo(req, res).ToZapFields(start)...)

			switch {
			case res.Status >= 500:
				if msg != "" {
					msg = fmt.Sprintf("Internal Server Error: %s", msg)
				} else {
					msg = "Internal Server Error"
				}

				cc.CustomLogger.Error(msg, fields...)
			case res.Status >= 400:
				if msg != "" {
					msg = fmt.Sprintf("Client Error: %s", msg)
				} else {
					msg = "Client Error"
				}
				cc.CustomLogger.Error(msg, fields...)
			case res.Status >= 300:
				cc.CustomLogger.Info("Redirection", fields...)
			default:
				cc.CustomLogger.Info("Success", fields...)
			}

			return nil
		}
	}
}
