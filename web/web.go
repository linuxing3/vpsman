package web

import (
	"fmt"
	"net/http"
	"strconv"
	"github.com/linuxing3/vpsman/util"
	"github.com/linuxing3/vpsman/web/controller"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr/v2"
)

func userRouter(router *gin.Engine) {
	user := router.Group("/xray/user")
	{
		user.GET("", func(c *gin.Context) {
			requestUser := RequestUsername(c)
			if requestUser == "admin" {
				c.JSON(200, controller.UserList(""))
			} else {
				c.JSON(200, controller.UserList(requestUser))
			}
		})
		user.GET("/page", func(c *gin.Context) {
			curPageStr := c.DefaultQuery("curPage", "1")
			pageSizeStr := c.DefaultQuery("pageSize", "10")
			curPage, _ := strconv.Atoi(curPageStr)
			pageSize, _ := strconv.Atoi(pageSizeStr)
			c.JSON(200, controller.PageUserList(curPage, pageSize))
		})
		user.POST("", func(c *gin.Context) {
			username := c.PostForm("username")
			password := c.PostForm("password")
			c.JSON(200, controller.CreateUser(username, password))
		})
		user.POST("/update", func(c *gin.Context) {
			sid := c.PostForm("id")
			username := c.PostForm("username")
			password := c.PostForm("password")
			c.JSON(200, controller.UpdateUser(sid, username, password))
		})
		user.POST("/expire", func(c *gin.Context) {
			sid := c.PostForm("id")
			sDays := c.PostForm("useDays")
			useDays, _ := strconv.Atoi(sDays)
			c.JSON(200, controller.SetExpire(sid, uint(useDays)))
		})
		user.DELETE("/expire", func(c *gin.Context) {
			sid := c.Query("id")
			c.JSON(200, controller.CancelExpire(sid))
		})
		user.DELETE("", func(c *gin.Context) {
			sid := c.Query("id")
			c.JSON(200, controller.DelUser(sid))
		})
	}
}

func xrayRouter(router *gin.Engine) {
	router.POST("/xray/start", func(c *gin.Context) {
		c.JSON(200, controller.Start())
	})
	router.POST("/xray/stop", func(c *gin.Context) {
		c.JSON(200, controller.Stop())
	})
	router.POST("/xray/restart", func(c *gin.Context) {
		c.JSON(200, controller.Restart())
	})
	router.GET("/xray/loglevel", func(c *gin.Context) {
		c.JSON(200, controller.GetLogLevel())
	})
	router.POST("/xray/update", func(c *gin.Context) {
		c.JSON(200, controller.Update())
	})
	// router.POST("/xray/switch", func(c *gin.Context) {
	// 	tType := c.DefaultPostForm("type", "xray")
	// 	c.JSON(200, controller.SetTrojanType(tType))
	// })
	router.POST("/xray/loglevel", func(c *gin.Context) {
		slevel := c.DefaultPostForm("level", "1")
		// level, _ := strconv.Atoi(slevel)
		c.JSON(200, controller.SetLogLevel(slevel))
	})
	router.POST("/xray/domain", func(c *gin.Context) {
		c.JSON(200, controller.SetDomain(c.PostForm("domain")))
	})
	router.GET("/xray/log", func(c *gin.Context) {
		controller.Log(c)
	})
}

func staticRouter(router *gin.Engine) {
	box := packr.New("trojanBox", "./templates")
	router.Use(func(c *gin.Context) {
		requestUrl := c.Request.URL.Path
		if box.Has(requestUrl) || requestUrl == "/" {
			http.FileServer(box).ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	})
}

// Start web启动入口
func Start(host string, port int, isSSL bool) {
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	// 静态路由
	staticRouter(router)
	// 中间件
	router.Use(Auth(router).MiddlewareFunc())
	// 动态路由
	xrayRouter(router)
	userRouter(router)
	// 定时任务
	controller.SheduleTask()
	controller.CollectTask()
	// 打开端口
	util.OpenPort(port)
	// 启动服务器
  router.Run(fmt.Sprintf("%s:%d", host, port))
}
