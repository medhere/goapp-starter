package utils

import (
	"goapp/config"
	"time"
)

var JwTSigningKey = config.Env("JWT_SIGNING_KEY", "secret")
var RefreshTokenKey = config.Env("REFRESH_TOKEN_KEY", "refresh_secret")
var AccessTokenDuration = time.Hour * 12      // 12 hours
var RefreshTokenDuration = time.Hour * 24 * 7 // 7 days
