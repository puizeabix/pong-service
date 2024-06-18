package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	HealthCheck(c *gin.Context)
}

type handler struct {
}

type ServerStatus struct {
	Status string `json:"status"`
}

var (
	ServerStatusOK    = ServerStatus{Status: "OK"}
	ServerStatusError = ServerStatus{Status: "ERROR"}
)

func NewHandler() Handler {
	return &handler{}
}

func (h *handler) HealthCheck(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, ServerStatusOK)
}
