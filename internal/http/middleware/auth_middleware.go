package middleware

import (
	"errors"

	"github.com/golang-jwt/jwt"
	"github.com/valyala/fasthttp"
)

const (
	AuthorizationHeader = "Authorization"
)

type CustomClaims struct {
	UserID            string `json:"user_id"`
	AccountOperatorID string `json:"account_operator_id"`
	jwt.StandardClaims
}

type Identity struct {
	UserID            string `json:"user_id"`
	AccountOperatorID string `json:"account_operator_id"`
}

// todo delete when rpc will works
// Deprecated
func Middleware(h fasthttp.RequestHandler, key string) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		tokenString := string(ctx.Request.Header.Peek(AuthorizationHeader))
		if tokenString == "" {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { // TODO ask actual Method
				return nil, errors.New("unexpected signing method")
			}
			return []byte(key), nil
		})
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			return
		}

		err = token.Claims.Valid()
		if !token.Valid || err != nil {
			ctx.SetStatusCode(fasthttp.StatusForbidden)
			return
		}

		claims := token.Claims.(*CustomClaims)
		_ = Identity{
			UserID:            claims.UserID,
			AccountOperatorID: claims.AccountOperatorID,
		}

		h(ctx)
	}
}

// Deprecated
func MiddlewareSetupResponse(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {

		ctx.Response.Header.Set(fasthttp.HeaderAccessControlAllowOrigin, "*")
		ctx.Response.Header.Set(fasthttp.HeaderContentType, "application/json")
		ctx.Response.Header.Set(fasthttp.HeaderAccessControlAllowCredentials, "true")
		ctx.Response.Header.Set(fasthttp.HeaderAccessControlAllowMethods, "GET,HEAD,OPTIONS,POST,PUT")
		ctx.Response.Header.Set(fasthttp.HeaderAccessControlAllowHeaders, "Access-Control-Allow-Headers, "+
			"Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, "+
			"Access-Control-Request-Headers, Authorization")

		h(ctx)
	}
}
