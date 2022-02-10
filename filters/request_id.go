package filters

import (
	"encoding/hex"
	"math/rand"
	"wego"
)

func genRequestID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	sid := hex.EncodeToString(b)
	return sid
}

const keyRequestID = "X-Request-ID"
func GenRequestID() wego.HandlerFunc {
	return func(c *wego.WebContext) {
		rid :=c.Input.GetHeader(keyRequestID)
		if rid == "" {
			rid = genRequestID()
			c.Input.Header.Add(keyRequestID, rid)
		}
		c.SetHeader(keyRequestID, rid)
		c.Next()
	}
}
