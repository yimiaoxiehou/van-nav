package main

import (
	"fmt"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

func main() {
	initDB()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	// 嵌入文件夹

	// public,_ := fs.ReadDir("./public")
	// router.StaticFS("/",http.FS(fs))

	router.GET("/manifest.json", ManifastHanlder)
	router.Use(static.Serve("/", BinaryFileSystem(fs, "public")))
	// router.Use(static.Serve("/", static.LocalFile("./public", true)))
	api := router.Group("/api")
	{
		// 获取数据的路由
		api.GET("/", GetAllHandler)
		// 获取用户信息

		api.POST("/login", LoginHandler)
		api.GET("/logout", LogoutHandler)
		api.GET("/img", getLogoImgHandler)
		// 管理员用的
		admin := api.Group("/admin")
		admin.Use(JWTMiddleware())
		{
			admin.POST("/apiToken", AddApiTokenHandler)
			admin.DELETE("/apiToken/:id", DeleteApiTokenHandler)
			admin.GET("/all", GetAdminAllDataHandler)

			admin.GET("/exportTools", ExportToolsHandler)

			admin.POST("/importTools", ImportToolsHandler)

			admin.PUT("/user", UpdateUserHandler)

			admin.PUT("/setting", UpdateSettingHandler)

			admin.POST("/tool", AddToolHandler)
			admin.DELETE("/tool/:id", DeleteToolHandler)
			admin.PUT("/tool/:id", UpdateToolHandler)

			admin.POST("/catelog", AddCatelogHandler)
			admin.DELETE("/catelog/:id", DeleteCatelogHandler)
			admin.PUT("/catelog/:id", UpdateCatelogHandler)
		}
	}
	fmt.Println("应用启动成功，网址:   http://localhost:6412")
	router.Run(":6412")
}
