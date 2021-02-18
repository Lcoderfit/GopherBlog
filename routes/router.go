package routes

import (
	"GopherBlog/constant"
	"GopherBlog/controller"
	"GopherBlog/middleware"
	"GopherBlog/utils"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"net/http"
)

// createMyRender 自定义渲染器
func createMyRender() multitemplate.Renderer {
	p := multitemplate.NewRenderer()
	p.AddFromFiles("admin", "static/admin/index.html")
	p.AddFromFiles("front", "static/admin/index.html")
	return p
}

// InitRouter 初始化路由
func InitRouter() {
	// 设置项目启动模式, debug表示调试模式，test表示测试模式，release表示发布模式（用于生产环境）
	gin.SetMode(utils.ServerCfg.AppMode)
	r := gin.New()
	//r.HTMLRender = createMyRender()
	//
	r.Use(middleware.Log())
	// 当出现panic时会导致程序崩溃退出，该中间件会恢复panic导致的崩溃并返回http code 500
	r.Use(gin.Recovery())
	// 支持跨域资源共享
	r.Use(middleware.Cors())
	// TODO: 第一个参数是URL路径，第二个参数是项目路径, 几个函数的区别是什么？？
	//r.Static("/static", "/static/front/static")
	//r.Static("/admin", "/static/admin")
	//r.StaticFile("/favicon.ico", "/static/front/favicon.ico")

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
		auth.PUT("/admin/reset-password/:id", controller.ChangeUserPassword)
		auth.PUT("/users/:id", controller.EditUserInfo)
		auth.DELETE("/users/:id", controller.DeleteUser)

		// 个人信息模块
		auth.GET("/admin/profiles/:id", controller.GetProfileInfo)
		auth.PUT("/profiles/:id", controller.UpdateProfileInfo)

		// 分类功能
		auth.GET("/admin/categories", controller.GetCategoryList)
		auth.POST("/categories", controller.AddCategory)
		auth.PUT("/categories/:id", controller.EditCategoryInfo)
		auth.DELETE("/categories/:id", controller.DeleteCategory)

		// 文章模块
		auth.GET("/admin/articles/:id", controller.GetArticleInfo)
		auth.GET("/admin/articles", controller.GetArticleList)
		auth.POST("/articles", controller.AddArticle)
		auth.PUT("/articles/:id", controller.EditArticleInfo)
		auth.DELETE("/articles/:id", controller.DeleteArticle)

		// 评论模块
		auth.GET("/admin/comments/:id", controller.GetCommentList)
		auth.PUT("/comments-check/:id", controller.ApproveComment)
		auth.PUT("/comments-uncheck/:id", controller.TakeDownComment)
		auth.DELETE("/comments/:id", controller.DeleteComment)
	}

	// 设置路由组，定义无需鉴权的接口
	router := r.Group("/api/v1")
	{

		// 验证token
		router.POST("/admin/token-check", controller.CheckToken)
		// 登录模块
		router.POST("/login", controller.Login)
		//router.POST("/login_front", controller.LoginFront)

		// 用户信息模块
		router.POST("/users", controller.AddUser)
		router.GET("/users/:id", controller.GetUserInfo)
		router.GET("/users", controller.GetUserList)

		// 获取个人信息
		router.GET("/profiles/:id", controller.GetProfileInfo)

		// 文章分类模块
		router.GET("/categories/:id", controller.GetCategoryInfo)
		router.GET("/categories", controller.GetCategoryList)

		// 文章模块
		router.GET("/articles/:id", controller.GetArticleInfo)
		router.GET("/articles", controller.GetArticleList)
		router.GET("/categories/:id/articles", controller.GetArticleListByCategoryId) // 获取同一分类的所有文章

		// 评论模块
		router.POST("/comments", controller.AddComment)
		router.GET("/comments/:id", controller.GetCommentInfo)
		router.GET("/articles/:id/comments-count", controller.GetCommentCount)
		router.GET("/articles/:id/comments", controller.GetCommentList)

		router.POST("/test", controller.Test)
	}

	// 运行项目
	err := r.Run(utils.ServerCfg.HttpPort)
	if err != nil {
		utils.Logger.Panic(constant.ConvertForLog(constant.StartProjectError), err)
	}
}
