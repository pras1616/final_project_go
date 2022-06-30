package helpers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	defaultUsername = "hacktiv"
	defaultPassword = "12345678"
	defaultHasAuth  = "h12bi3pu5924bf924y02bf2"
)

func BasicAuth(c *gin.Context) {
	user, password, defaultHasAuth := c.Request.BasicAuth()
	if defaultHasAuth && user == defaultUsername && password == defaultPassword {
		log.WithFields(log.Fields{
			"user": user,
		}).Info("User authenticated")
		fmt.Println("User authenticated")
	} else {
		c.Abort()
		c.Writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
		c.JSON(401, gin.H{
			"code":    401,
			"data":    "-",
			"message": "Not Allow, Unauthorized",
		})
		return
	}
}
