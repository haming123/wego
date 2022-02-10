package cache

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CookieStore struct {
	CookieName 	string
	HashKey 	string
}

func NewCookieStore(cookie_name string, hash_key string) *CookieStore {
	cs := &CookieStore{}
	cs.HashKey = hash_key
	cs.CookieName = cookie_name
	if cs.CookieName == "" {
		cs.CookieName = "sdata"
	}
	return cs
}

func (this *CookieStore) SaveData(ctx context.Context, sid string, data []byte, max_age uint) error {
	//Base64编码
	base64_encode := string(EncodeBase64(data))
	//添加过期时间
	tm_end_str := genExpiration(time.Now(), int64(max_age))
	time_encode := tm_end_str + "-" + base64_encode
	//添加签名
	sign := createSign(time_encode, this.HashKey)
	sign_encode := sign + "-" + time_encode

	out := ctx.Value("http.ResponseWriter").(http.ResponseWriter)
	cookieData := &http.Cookie{
		Name:     this.CookieName,
		Value:    sign_encode,
		Domain:   "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(max_age),
	}
	http.SetCookie(out, cookieData)
	return nil
}

func (this *CookieStore) ReadData(ctx context.Context, sid string) ([]byte, error) {
	req := ctx.Value("http.Request").(*http.Request)
	cookieData, err := req.Cookie(this.CookieName)
	if err != nil || cookieData.Value == "" {
		return nil, err
	}

	//验证签名
	sign_encode := cookieData.Value
	index := strings.Index(sign_encode, "-")
	if index < 0 {
		return nil, errors.New("invalid cookie value")
	}
	sign, time_encode := sign_encode[:index], sign_encode[index+1:]
	if sign == "" || time_encode == "" {
		return nil, errors.New("invalid cookie value")
	}
	if !checkSing(this.HashKey, time_encode, sign) {
		return nil, errors.New("invalid cookie value")
	}

	//验证过期时间
	index = strings.Index(time_encode, "-")
	if index < 0 {
		return nil, errors.New("invalid cookie value")
	}
	tm_end_str, base64_encode := time_encode[:index], time_encode[index+1:]
	if tm_end_str == "" || base64_encode == "" {
		return nil, errors.New("invalid cookie value")
	}
	tm_end := getExpiration(tm_end_str)
	if tm_end > 0 && tm_end < time.Now().Unix() {
		return nil, errors.New("invalid cookie value")
	}

	//Base64解码
	data, err := DecodeBase64([]byte(base64_encode))
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (this *CookieStore) Exist(ctx context.Context, key string) (bool, error) {
	data, err := this.ReadData(ctx, key)
	if err != nil {
		return false, err
	}
	return data != nil, nil
}

func (this *CookieStore) Delete(ctx context.Context, key string) error {
	out := ctx.Value("http.ResponseWriter").(http.ResponseWriter)
	cookieData := &http.Cookie{
		Name:     this.CookieName,
		Value:    "",
		Domain:   "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   0,
	}
	http.SetCookie(out, cookieData)
	return nil
}

func genExpiration(tm time.Time, max_age int64) string {
	return fmt.Sprintf("%d:%d",tm.Unix(), max_age)
}

func getExpiration(tm_end_str string) int64{
	index := strings.Index(tm_end_str, ":")
	if index < 1 {
		return 0
	}
	tm_str, age_str := tm_end_str[:index], tm_end_str[index+1:]
	tm, _ := strconv.ParseInt(tm_str, 10, 64)
	age, _ := strconv.ParseInt(age_str, 10, 64)
	return tm + age
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

func createSign(data string, hash_key string) string {
	if len(hash_key) == 0 {
		return ""
	}
	mac := hmac.New(sha1.New, []byte(hash_key))
	if _, err := io.WriteString(mac, data); err != nil {
		return ""
	}
	return hex.EncodeToString(mac.Sum(nil))
}

func checkSing(hash_key string, data string, sig string) bool {
	sign := createSign(data, hash_key)
	return hmac.Equal([]byte(sig), []byte(sign))
}
