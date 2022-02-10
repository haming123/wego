package wego

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
)

type FlashStat int

const (
	FLASH_OK 		FlashStat = iota
	FLASH_ERROR
)

type WebFlash struct {
	Stat 	FlashStat
	Data 	map[string]string
}

func NewFlash() *WebFlash {
	return &WebFlash{
		Data: make(map[string]string),
	}
}

func getFlashCookieName(ctx *WebContext) string {
	cookieName := "sflag"
	cfg := ctx.engine.Config.SessionParam
	if cfg.CookieName != "" {
		cookieName = cfg.CookieName + "_flash"
	}
	return cookieName
}

func (this *WebFlash)SaveSuccessData(ctx *WebContext) error {
	this.Stat = FLASH_OK
	return this.SaveData(ctx)
}

func (this *WebFlash)SaveErrorData(ctx *WebContext) error {
	this.Stat = FLASH_ERROR
	return this.SaveData(ctx)
}

func (this *WebFlash)SaveData(ctx *WebContext) error {
	json_data, err := json.Marshal(*this)
	if err != nil {
		return err
	}
	base64_data := EncodeBase64(json_data)
	cookie_data := url.QueryEscape(string(base64_data))

	cookie := &http.Cookie {
		Name:     getFlashCookieName(ctx),
		Value:    cookie_data,
		Domain:   "",
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(ctx.Output.ResponseWriter, cookie)
	return nil
}

func (this *WebFlash)ReadData(ctx *WebContext) error {
	cookie_data, err := ctx.Input.Cookie(getFlashCookieName(ctx))
	if err != nil {
		return err
	}
	if cookie_data == "" {
		return nil
	}

	base64_data, err := url.QueryUnescape(cookie_data)
	if err != nil {
		return  err
	}

	//Base64解码
	json_data, err := DecodeBase64([]byte(base64_data))
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(json_data), this)
	if err != nil {
		return err
	}

	return nil
}

func EncodeBase64(value []byte) []byte {
	encoded := make([]byte, base64.URLEncoding.EncodedLen(len(value)))
	base64.URLEncoding.Encode(encoded, value)
	return encoded
}

func DecodeBase64(value []byte) ([]byte, error) {
	decoded := make([]byte, base64.URLEncoding.DecodedLen(len(value)))
	b, err := base64.URLEncoding.Decode(decoded, value)
	if err != nil {
		return nil, err
	}
	return decoded[:b], nil
}

func ReadFlash(ctx *WebContext) *WebFlash {
	flash := &WebFlash{}
	flash.ReadData(ctx)
	cookie := &http.Cookie {
		Name:     getFlashCookieName(ctx),
		Value:    "",
		Domain:   "",
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(ctx.Output.ResponseWriter, cookie)
	return flash
}
