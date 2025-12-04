package controllers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Render(c *gin.Context, name string, data gin.H) {
	if data == nil {
		data = gin.H{}
	}

	session := sessions.Default(c)
	if u := session.Get("user"); u != nil {
		data["User"] = u
	}

	c.HTML(http.StatusOK, name, data)
}
