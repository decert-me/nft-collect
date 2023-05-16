package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"nft-collect/internal/app/model/response"
	"nft-collect/internal/app/utils"
	"strings"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// jwt鉴权取头部信息： x-token
		token := c.Request.Header.Get("x-token")
		if token == "" {
			response.FailWithDetailed(gin.H{"reload": true}, "授权已过期或非法访问1", c)
			c.Abort()
			return
		}
		j := utils.NewJWT()
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == utils.TokenExpired {
				response.FailWithDetailed(gin.H{"reload": true}, "授权已过期", c)
				c.Abort()
				return
			}
			response.FailWithDetailed(gin.H{"reload": true}, err.Error(), c)
			c.Abort()
			return
		}
		c.Set("address", strings.ToLower(claims.Address))
	}
}

func Addr() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 鉴权头部信息： x-token
		token := c.Request.Header.Get("x-token")
		fmt.Println("token", token)
		if token != "" {
			j := utils.NewJWT()
			// 解析token包含的信息
			claims, err := j.ParseTokenWithNoAuth(token)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("claims", claims)
			c.Set("address", strings.ToLower(claims.Address))
		}
	}
}
