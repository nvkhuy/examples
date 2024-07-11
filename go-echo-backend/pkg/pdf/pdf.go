package pdf

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/brianvoe/sjwt"
	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"github.com/engineeringinflow/inflow-backend/pkg/locker"
)

type GetPDFParams struct {
	URL      string `json:"url"`
	Selector string `json:"selector"`

	Landscape         bool `json:"landscape"`
	DisableLocker     bool `json:"disable_locker"`
	PrintBackground   bool `json:"print_background"`
	PreferCssPageSize bool `json:"prefer_css_page_size"`
}

type Client struct {
	locker *locker.Locker
	config *config.Configuration
}

func New(cfg *config.Configuration) *Client {
	return &Client{
		locker: locker.New(cfg),
		config: cfg,
	}
}

func (c *Client) GetPDF(params GetPDFParams) (result []byte, err error) {
	if c.locker != nil && !params.DisableLocker {
		cancel, err := c.locker.AcquireLock(fmt.Sprintf("processing_pdf_%s", params.URL), time.Minute*5)
		if err != nil {
			return nil, err
		}
		defer cancel()
	}

	var claims = sjwt.New()
	claims.SetIssuer(fmt.Sprintf("ghost|%s", c.config.SuperAdminUserID))
	claims.SetExpiresAt(time.Now().Add(time.Minute * 10))
	claims.Set("id", c.config.SuperAdminUserID)
	claims.Set("aud", "super_admin")
	var token = claims.Generate([]byte(c.config.JWTSecret))

	u, err := url.Parse(c.config.LambdaAPIPDFURL)
	if err != nil {
		err = fmt.Errorf("error parsing base URL: %v", err)
		return
	}

	q := u.Query()
	q.Set("url", params.URL)
	q.Set("selector", params.Selector)
	q.Set("jwt_token", token)
	q.Set("landscape", fmt.Sprintf("%t", params.Landscape))
	q.Set("print_background", fmt.Sprintf("%t", params.PrintBackground))
	q.Set("prefer_css_page_size", fmt.Sprintf("%t", params.PreferCssPageSize))

	u.RawQuery = q.Encode()

	link := u.String()
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		err = fmt.Errorf("error creating request: %v", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("error making request: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("error reading response body: %v", err)
		return
	}

	if json.Valid(body) {
		var apiErr APIError
		if e := json.Unmarshal(body, &apiErr); e == nil && apiErr.Message != "" {
			return nil, &apiErr
		}
	}
	result = body
	return
}
