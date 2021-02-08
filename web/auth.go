package web

import (
	"fmt"
	"time"

	"github.com/linuxing3/vpsman/core"
	"github.com/linuxing3/vpsman/web/controller"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

var (
	identityKey    = "id"
	authMiddleware *jwt.GinJWTMiddleware
	err            error
)

// Login auth用户验证结构体
type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func init() {
	authMiddleware, err = jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "k8s-manager",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		SendCookie:  true,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*Login); ok {
				return jwt.MapClaims{
					identityKey: v.Username,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &Login{
				Username: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var (
				loginVals Login
			)
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			if err != nil {
				return nil, err
			}
			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(*Login); ok {
				return true
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		fmt.Println("JWT Error:" + err.Error())
	}
}

func updateUser(c *gin.Context) {
	responseBody := controller.ResponseBody{Msg: "success"}
	defer controller.TimeCost(time.Now(), &responseBody)
	username := c.DefaultPostForm("username", "admin")
	pass := c.PostForm("password")
	s := core.Sqlite{ Path: "" }
	s.UpdateUserORM("", username, pass, "")
	// set value in leveldb for sessions
	c.JSON(200, responseBody)
}

// RequestUsername 获取请求接口的用户名
func RequestUsername(c *gin.Context) string {
	claims := jwt.ExtractClaims(c)
	return claims[identityKey].(string)
}

// Auth 权限router
func Auth(r *gin.Engine) *jwt.GinJWTMiddleware {
	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		fmt.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": 404, "message": "Page not found"})
	})
	r.GET("/auth/check", func(c *gin.Context) {
		result := "admin"
		if result == "" {
			c.JSON(201, gin.H{"code": 201, "message": "No administrator account found inside the database", "data": nil})
		} else {
      title := "xray 管理平台"
			c.JSON(200, gin.H{
				"code":    200,
				"message": "success",
				"data": map[string]string{
					"title": title,
				},
			})
		}
	})
	r.POST("/auth/login", authMiddleware.LoginHandler)
	r.POST("/auth/register", updateUser)
	authO := r.Group("/auth")
	authO.Use(authMiddleware.MiddlewareFunc())
	{
		authO.GET("/loginUser", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"code":    200,
				"message": "success",
				"data": map[string]string{
					"username": RequestUsername(c),
				},
			})
		})
		authO.POST("/reset_pass", updateUser)
		authO.POST("/logout", authMiddleware.LogoutHandler)
		authO.POST("/refresh_token", authMiddleware.RefreshHandler)
	}
	return authMiddleware
}
