package rsa

import (
	"encoding/base64"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"github.com/rotisserie/eris"
)

type Client struct {
	cfg *config.Configuration
}

func New(cfg *config.Configuration) *Client {
	return &Client{
		cfg: cfg,
	}
}

func (c *Client) Validate(encrypted string) error {
	privateKey, err := os.ReadFile(c.cfg.RsaPrivatePemFile)
	if err != nil {
		return eris.Wrapf(err, "Failed to load private pem %s", c.cfg.RsaPrivatePemFile)
	}

	cipherText, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return eris.Wrapf(err, "Failed to decode string %s", encrypted)
	}

	originText, err := Decrypt(privateKey, []byte(cipherText))
	if err != nil {
		return eris.Wrapf(err, "Failed to descrypt %s", string(cipherText))
	}

	var parts = strings.Split(string(originText), "|")

	if len(parts) != 2 {
		return eris.Errorf("Invalid data")
	}

	if parts[0] != c.cfg.RsaSecret {
		return eris.Errorf("Invalid secret")
	}

	ts, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return eris.Wrapf(err, "Invalid timestamp")
	}

	if ts < time.Now().Unix() {
		return eris.Wrapf(err, "Invalid timestamp")
	}

	return nil
}
