package orders

import (
	"encoding/json"
	"net/http"
	"online-store-go/pkg/error_utils"
	"online-store-go/pkg/logger_utils"

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
		logger_utils.Error("Error when get raw data", jsonErr)
		c.JSON(http.StatusBadRequest, error_utils.NewBadRequestError("invalid json body"))
		return
	}

	var data map[string]interface{}
	jsonErr = json.Unmarshal(jsonData, &data)
	if jsonErr != nil {
		logger_utils.Error("Error when unmarshal jsonData", jsonErr)
		c.JSON(http.StatusBadRequest, error_utils.NewBadRequestError("invalid json body"))
		return
	}

	err := h.s.CreateOrder(data)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "order created successfully",
	})
}
