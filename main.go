package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"plexcorp.tech/gosass/controllers"
	"plexcorp.tech/gosass/middleware"
)

func main() {
	location, err := time.LoadLocation(os.Getenv("TZ"))
	if err != nil {
		fmt.Println("Timezone entered is invalid:", err)
		return
	}

	time.Local = location

	go func() {
		for {
			RunJobs()
			time.Sleep(30 * time.Second)
		}
	}()

	fmt.Println("Running server on port: 8080")

	mysqlDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
	)

	os.Setenv("MYSQL_DSN", mysqlDSN)
	router := gin.Default()
	router.StaticFS("/static", http.Dir("./static"))

	allowedIps := os.Getenv("allowed_ips")
	router.SetTrustedProxies(strings.Split(allowedIps, ","))

	router.Use(middleware.DBMiddleware())
	router.Use(middleware.SetupSession())
	router.Use(middleware.AuthMiddleware())
	router.Use(func(c *gin.Context) {
		c.Next()
		db := c.MustGet("db").(*gorm.DB)
		sqlDB, err := db.DB()
		if err == nil {
			defer sqlDB.Close()
		}
	})

	controller := controllers.Controller{}

	router.GET("/users/logout", controller.Logout)
	router.GEt("/users/login", controller.LoginView)
	router.POST("/users/authenticate", controller.CheckLogin)
	router.POST("/users/2fa/authenticate", controller.TwoFactorAuthenticate)
	router.Any("/users/password/reset/:token", controller.ChangePassword)
	router.Any("/users/password/forgot", controller.ForgotPassword)
	router.POST("/users/setup/complete", controller.FinishInitialUserSetup)
	router.GET("/users/setup", controller.SetupIntialUser)

	router.GET("/user/list", controller.ListUsers)
	router.POST("/user/actions", controller.HandleUserActionsFormPost)
	router.POST("/user/profile/update", controller.UpdateProfile)
	router.GET("/user/profile", controller.MyProfile)
	router.GET("/user/create", controller.NewUser)

	router.GET("/denied", controller.AccessDenied)
	router.GET("/user/2factor/qrcode", controller.ShowQrCodePng)

	router.Run(":3001")

}
