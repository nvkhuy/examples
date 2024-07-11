package helper

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/brianvoe/sjwt"
	"github.com/rotisserie/eris"
)

type GenerateRemoteURLParams struct {
	BaseURL  string
	Issuer   string
	Secret   string
	Expiry   time.Duration
	UserID   string
	Metadata map[string]string
}

func GenerateRemoteURL(params GenerateRemoteURLParams) (string, error) {
	var claims = sjwt.New()
	claims.SetIssuer("share")
	claims.SetExpiresAt(time.Now().Add(params.Expiry))
	claims.Set("user_id", params.UserID)

	var token = claims.Generate([]byte(params.Secret))

	redirectLink, err := url.Parse(params.BaseURL)
	if err != nil || redirectLink == nil {
		return "", eris.Wrap(err, "")
	}

	var values = redirectLink.Query()
	values.Add("token", token)

	for key, v := range params.Metadata {
		values.Add(key, v)
	}
	redirectLink.RawQuery = values.Encode()

	return redirectLink.String(), nil
}

func AddURLQuery(link string, values map[string]string) string {
	url, err := url.Parse(link)
	if err != nil {
		return link
	}

	var query = url.Query()
	for k, v := range values {
		query.Add(k, v)
	}

	url.RawQuery = query.Encode()

	return url.String()
}

func FixProductImages(urls []string) (results []string) {
	const base = "https://"
	for _, urlPath := range urls {
		if len(urlPath) > len(base) && urlPath[:len(base)] != base {
			urlPath = fmt.Sprintf("%s%s", base, urlPath)
		}
		validExts := []string{".jpeg", "jpg", ".png", ".webp"}
		if _, err := url.ParseRequestURI(urlPath); err == nil {
			ext := filepath.Ext(urlPath)
			for _, validExt := range validExts {
				if ok := strings.Contains(ext, validExt); ok {
					results = append(results, urlPath)
					break
				}
			}
		}
	}
	return
}
