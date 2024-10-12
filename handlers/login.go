package handlers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

type LoginForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func (h *Handler) loginHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func (h *Handler) submitLoginForm(c *gin.Context) {
	var formData LoginForm
	if err := c.ShouldBind(&formData); err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": err.Error()})
		return
	}
	username := os.Getenv("USER_LOGIN")
	password := os.Getenv("USER_PASSWORD")
	if formData.Username != username || formData.Password != password {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"data": "неправильные логин или пароль",
		})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"foo": "bar",
		"nbf": time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, err := token.SignedString([]byte("hmacSampleSecret"))
	if err != nil {
		logrus.Error(err)
	}
	fmt.Println(tokenString)
	c.Request.Header.Add("Token", tokenString)
	c.Redirect(http.StatusFound, "/backups")
}
