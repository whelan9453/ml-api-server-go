package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func authCheck(c *gin.Context) {
	k := c.GetHeader(HeaderRestKey)
	if k == "" || k != os.Getenv(EnvRestKey) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code":    ErrUnauthorized,
			"message": "Invalid token.",
		})
	}

	c.Next()
}

func setUpTestingEnv() {
	os.Setenv(EnvServerPort, "5566")
	os.Setenv(EnvRestKey, "amllMDRzdTNzdTs2cnUgMTkg")
}

func main() {

	if os.Getenv(EnvWhereAmI) != "production" {
		fmt.Printf("Current env is %s", os.Getenv(EnvWhereAmI))
		fmt.Println("Set up testing envionment variables")
		setUpTestingEnv()
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(authCheck)
	router.GET("/versionInfo", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"version": "0.0.1",
			"message": "Support version information request.",
		})
	})

	router.Run(fmt.Sprintf(":%v", os.Getenv(EnvServerPort)))
}
