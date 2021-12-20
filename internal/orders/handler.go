package orders

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type handler struct {
	s Service
}

type Handler interface {
	CreateOrder(*gin.Context)
}

func NewHandler(service Service) Handler {
	return &handler{
		s: service,
	}
}

func (h *handler) CreateOrder(c *gin.Context) {
	jsonData, jsonErr := c.GetRawData()
	if jsonErr != nil {
		c.JSON(http.StatusBadRequest, jsonErr.Error())
		return
	}

	var data map[string]interface{}
	json.Unmarshal(jsonData, &data)

	err := h.s.CreateOrder(data)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "order created successfully",
	})
}
