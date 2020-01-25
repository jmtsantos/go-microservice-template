package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// IndexGET returns the index page
func (h *Handlers) IndexGET(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}
