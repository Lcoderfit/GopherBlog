package routes

import (
	"GopherBlog/controller"
	"GopherBlog/middleware"
	"GopherBlog/utils"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 自定义渲染器
func createMyRender() multitemplate.Renderer {
	p := multitemplate.NewRenderer()
	p.AddFromFiles("admin", "static/admin/index.html")
	p.AddFromFiles("front", "static/admin/index.html")
	return p
}

func InitRouter() {
	// 设置项目启动模式, debug表示调试模式，test表示测试模式，release表示发布模式（用于生产环境）
	gin.SetMode(utils.ServerCfg.AppMode)
	r := gin.New()
	r.HTMLRender = createMyRender()
	//r.Use(middleware.Log())
	// TODO: 第一个参数是URL路径，第二个参数是项目路径, 几个函数的区别是什么？？
	r.Static("/static", "/static/front/static")
	r.Static("/admin", "/static/admin")
	r.StaticFile("/favicon.ico", "/static/front/favicon.ico")

	// 用户主页路由
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "front", nil)
	})

	// 后台管理系统主页路由
	r.GET("/admin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin", nil)
	})

	// 需要鉴权的接口
	auth := r.Group("/api/v1")
	// 添加JWT中间件
	auth.Use(middleware.JwtToken())
	{
		// 用户模块
		auth.GET("/admin/users", controller.GetUserList)
		auth.PUT("/admin/change_password/:id", controller.ChangeUserPassword)
		auth.PUT("/user/:id", controller.EditUserInfo)
		auth.DELETE("/user/:id", controller.DeleteUser)

		//

	}

	// TODO:需要重构成RESTful API
	// 设置路由组，定义无需鉴权的接口
	router := r.Group("/api/v1")
	{

		// 验证token
		router.POST("/admin/check_token")
		// 登录模块
		router.POST("/login", controller.Login)
		//router.POST("/login_front", controller.LoginFront)

		// 用户信息模块
		router.POST("/user/add", controller.AddUser)
		router.GET("/user/:id", controller.GetUserInfo)
		router.GET("/users", controller.GetUserList)

		// 获取个人信息
		router.GET("/profile/:id", controller.GetProfile)

		// 文章分类模块
		router.GET("/categories", controller.GetCategoryList)
		router.GET("/category/:id", controller.GetCategoryInfo)

		// 文章模块
		//router.POST("/article/add", controller.CreateArticle)
		router.GET("/article/info/:id", controller.GetArticleInfo)
		router.GET("/articles", controller.GetArticleList)
		router.GET("/article/list/:id", controller.GetArticleListByCategoryId) // 获取同一分类的所有文章

		// 评论模块
		router.POST("/comment/add", controller.AddComment)
		router.GET("/comment/info/:id", controller.GetCommentInfo)
		router.GET("/comment_count", controller.GetCommentCount)
		router.GET("/comment/list/:id", controller.GetCommentList)
	}

	// 运行项目
	err := r.Run(utils.ServerCfg.HttpPort)
	if err != nil {
		utils.Logger.Error("项目启动失败: ", err)
	}
}
