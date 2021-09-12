package middleware

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/lonelypale/goutils/errors"
	"github.com/lonelypale/goutils/ref"
	"github.com/lonelypale/goutils/token"
)

const (
	TokenRequestLabel = "token_request_label"
)

func Token(tkn token.Token, opt *token.Options) gin.HandlerFunc {
	tokenType := ref.PrimitiveType(reflect.TypeOf(tkn))

	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			errorHandler(c, http.StatusUnauthorized, errors.New("no Authorization"))
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			errorHandler(c, http.StatusUnauthorized, errors.New("no Bearer"))
			return
		}

		reqToken, ok := reflect.New(tokenType).Interface().(token.Token)
		if !ok {
			errorHandler(c, http.StatusInternalServerError, errors.New("bad token type"))
			return
		}

		jwtToken, err := token.ParseToken(parts[1], opt.SecretKey, reqToken)
		if err != nil {
			errorHandler(c, http.StatusUnauthorized, err)
			return
		}

		if opt.Cache != nil {
			key := reqToken.CacheKey()
			cacheToken := reflect.New(tokenType).Interface().(token.Token)
			if err := opt.Cache.Get(key, cacheToken); err != nil {
				errorHandler(c, http.StatusUnauthorized, err)
				return
			}

			if cacheToken.GetSignature() != jwtToken.Signature {
				errorHandler(c, http.StatusUnauthorized, errors.New("token signature mismatch"))
				return
			}
		}

		c.Set(TokenRequestLabel, reqToken)
		c.Next()
	}
}

func errorHandler(c *gin.Context, code int, errs ...error) {
	c.String(code, http.StatusText(code))
	c.Abort()

	for _, err := range errs {
		log.Error(err)
	}
}
