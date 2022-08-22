package middleware

import (
	"goapi/pkg/common"
	"goapi/pkg/common/logger"

	"github.com/casbin/casbin/v2"
	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrUnauthorize = common.NewForbiddenError("you are not allowed to access this resource")
	ErrEnforce     = common.NewUnexpectedError("error occurred while enforce")
)

func Authorize(enforcer *casbin.Enforcer) common.HandleFunc {
	return func(c common.HContext) error {
		public := c.Locals("public").(bool)
		if public {
			return c.Next()
		}

		u := c.Locals("user").(jwt.MapClaims)
		role := u["role"].(string)

		ok, err := enforcer.Enforce(role, c.Path(), c.Method())

		if err != nil {
			logger.ErrorWithReqId(err.Error(), c.RequestId())
			return common.ResponseError(c, ErrEnforce)
		}

		if !ok {
			return common.ResponseError(c, ErrUnauthorize)
		}

		return c.Next()
	}
}
