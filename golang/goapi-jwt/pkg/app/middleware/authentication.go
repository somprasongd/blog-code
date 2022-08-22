package middleware

import (
	"goapi/pkg/common"
	"goapi/pkg/common/logger"
	"goapi/pkg/util"
	"strings"

	"golang.org/x/exp/slices"
)

var (
	ErrNoToken       = common.NewUnauthorizedError("the token is required")
	ErrValidateToken = common.NewUnexpectedError("error occurred while validating token")
	ErrInvalidToken  = common.NewUnauthorizedError("the token is invalid or expired")
)

func Authentication(secretKey string, excludeList map[string][]string) common.HandleFunc {
	return func(c common.HContext) error {
		public := false

		if methods, ok := excludeList[c.Path()]; ok {
			public = slices.Contains(methods, c.Method())
		}

		if !public && strings.Contains(c.Path(), "/healthz") {
			public = true
		}

		if !public && strings.Contains(c.Path(), "/swagger/") {
			public = true
		}

		if !public && strings.Contains(c.Path(), "/thirdpartySwagger/") {
			public = true
		}

		if !public {
			auth := c.Authorization()
			// validate token
			if auth == "" {
				return common.ResponseError(c, ErrNoToken)
			}

			token := strings.TrimPrefix(auth, "Bearer ")
			valid, claims, err := util.ValidateToken(token, secretKey)

			if err != nil {
				logger.ErrorWithReqId(err.Error(), c.RequestId())
				return common.ResponseError(c, ErrValidateToken)
			}

			if !valid {
				return common.ResponseError(c, ErrInvalidToken)
			}

			c.Locals("user", claims)
		}

		return c.Next()
	}
}
