package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GinTemplateUtil is a utility for handling HTML templates using Gin.
type GinTemplateUtil struct {
	router *gin.Engine
}

// NewGinTemplateUtil creates a new GinTemplateUtil instance.
func NewGinTemplateUtil(templateDir string) *GinTemplateUtil {
	r := gin.Default()
	r.LoadHTMLGlob(templateDir + "/*")
	return &GinTemplateUtil{router: r}
}

// Render renders a template with the given data.
func (gtu *GinTemplateUtil) Render(c *gin.Context, name string, data gin.H) {
	c.HTML(http.StatusOK, name, data)
}

// GetRouter returns the Gin router.
func (gtu *GinTemplateUtil) GetRouter() *gin.Engine {
	return gtu.router
}
