package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const sessionCookie = "sprinto_session"

func sessionSecret() []byte {
	if s := os.Getenv("SESSION_SECRET"); s != "" {
		return []byte(s)
	}
	return []byte("sprinto-dev-secret-change-in-prod")
}

func setSession(c *gin.Context, userID uint) {
	payload := strconv.FormatUint(uint64(userID), 10)
	mac := hmac.New(sha256.New, sessionSecret())
	mac.Write([]byte(payload))
	sig := base64.URLEncoding.EncodeToString(mac.Sum(nil))
	val := base64.URLEncoding.EncodeToString([]byte(payload)) + "." + sig
	c.SetCookie(sessionCookie, val, 86400*30, "/", "", false, true)
}

func getSessionUserID(c *gin.Context) (uint, bool) {
	val, err := c.Cookie(sessionCookie)
	if err != nil || val == "" {
		return 0, false
	}
	parts := strings.SplitN(val, ".", 2)
	if len(parts) != 2 {
		return 0, false
	}
	payloadBytes, err := base64.URLEncoding.DecodeString(parts[0])
	if err != nil {
		return 0, false
	}
	sigBytes, err := base64.URLEncoding.DecodeString(parts[1])
	if err != nil {
		return 0, false
	}
	mac := hmac.New(sha256.New, sessionSecret())
	mac.Write(payloadBytes)
	if !hmac.Equal(sigBytes, mac.Sum(nil)) {
		return 0, false
	}
	id, err := strconv.ParseUint(string(payloadBytes), 10, 64)
	if err != nil {
		return 0, false
	}
	return uint(id), true
}

func clearSession(c *gin.Context) {
	c.SetCookie(sessionCookie, "", -1, "/", "", false, true)
}
