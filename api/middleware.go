package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/parin/simplebank/token"
	"net/http"
	"strings"
)

const (
	authorizationHeaderKey = "authorization"
	authorizationTypeBearer = "Bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc  {
	return func(ctx *gin.Context) {
		// check if the authorization  header is provided
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeaderKey) == 0 {
			err := errors.New("authorization header is not  provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,errorResponse(err))
			return
		}
		// check if authorization header consists of 2 elements
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2  {
			err := errors.New("invalid authorization  header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,errorResponse(err))
			return
		}
		// check if field one is bearer
		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s",authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,errorResponse(err))
			return
		}
		// Verifies the token
		accessToken := fields[1]
		payload,err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,errorResponse(err))
			return
		}
		ctx.Set(authorizationPayloadKey,payload)
		ctx.Next()
	}
}
