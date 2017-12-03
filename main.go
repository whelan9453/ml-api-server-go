package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func authCheck(c *gin.Context) {
	a := c.GetHeader(HeaderAuth)
	k := c.GetHeader(HeaderRestKey)
	pID := c.GetHeader(HeaderProjectID)
	pname := c.GetHeader(HeaderProjectName)

	if a == "" || k == "" || pID == "" || pname == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    ErrUnauthorized,
			"message": "Insufficient information.",
		})
		return
	}

	if k != os.Getenv(EnvRestKey) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code":    ErrUnauthorized,
			"message": "Invalid key.",
		})
		return
	}

	url := fmt.Sprintf("%s%s", os.Getenv(EnvGenAPIHost), os.Getenv(EnvGenEndPointAuth))
	req, errReq := http.NewRequest("GET", url, nil)
	req.Header.Set(HeaderAuth, a)
	req.Header.Set(HeaderProjectName, pname)
	if errReq != nil {
		log.Panicln("Making auth request failed.", errReq)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    ErrInternal,
			"message": "Something is wrong with the ml server.",
		})
		return
	}

	client := &http.Client{}
	resp, errResp := client.Do(req)
	if errResp != nil {
		log.Panicln("Auth failed.", errResp)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    ErrInternal,
			"message": "Something is wrong with the general server.",
		})
		return
	}

	//TODO
	io.Copy(os.Stdout, resp.Body)

	c.Next()
}

func setUpTestingEnv() {
	gin.SetMode(gin.DebugMode)
	os.Setenv(EnvServerPort, "5566")
	os.Setenv(EnvRestKey, "amllMDRzdTNzdTs2cnUgMTkg")
	os.Setenv(EnvAPIVer, "1")
	os.Setenv(EnvWhereAmI, "testing")
	os.Setenv(EnvGenAPIHost, "http://10.78.24.114:30000")
	os.Setenv(EnvGenEndPointAuth, "auth")
}

func main() {
	if os.Getenv(EnvWhereAmI) != "production" {
		fmt.Printf("Current env is %s\n", os.Getenv(EnvWhereAmI))
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
