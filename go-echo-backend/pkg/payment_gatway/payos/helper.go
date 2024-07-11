package payos

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/samber/lo"
)

func (c *Client) generateHMAC(data string) string {
	hmac := hmac.New(sha256.New, []byte(c.cfg.PayosChecksumKey))

	// compute the HMAC
	hmac.Write([]byte(data))
	dataHmac := hmac.Sum(nil)

	return hex.EncodeToString(dataHmac)
}

func (c *Client) sortObjDataByAlphabet(obj map[string]interface{}) map[string]interface{} {
	var keys = lo.Keys(obj)
	sort.Strings(keys)

	var sortedMap map[string]interface{} = make(map[string]interface{})
	for _, key := range keys {
		sortedMap[key] = obj[key]
	}

	return sortedMap
}

func (c *Client) isValidData(obj map[string]interface{}, sig string) bool {
	var params = lo.Map(lo.Keys(obj), func(key string, index int) string {
		return fmt.Sprintf("%s=%v", key, obj[key])
	})

	var data = strings.Join(params, "&")

	return c.generateHMAC(data) == sig
}
