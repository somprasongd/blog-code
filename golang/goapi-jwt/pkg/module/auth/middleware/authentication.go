package middleware

import (
	"goapi/pkg/common"
	"goapi/pkg/common/logger"
	"goapi/pkg/util"
	"strings"
)

var (
	ErrNoToken       = common.NewUnauthorizedError("the token is required")
	ErrValidateToken = common.NewUnexpectedError("error occurred while validating token")
	ErrInvalidToken  = common.NewUnauthorizedError("the token is invalid")
)

func Authentication(secretKey string) common.HandleFunc {
	return func(c common.HContext) error {
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

		return c.Next()
	}
}
