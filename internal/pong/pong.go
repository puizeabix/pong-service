package pong

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	Pong(c *gin.Context)
}

type handler struct {
}

type PongMessage struct {
	Message string `json:"message"`
}

func NewHandler() Handler {
	return &handler{}
}

func (h *handler) Pong(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, PongMessage{
		Message: "hello",
	})
}
