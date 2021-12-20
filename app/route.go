package app

import "online-store-go/internal/orders"

var (
	orderHandler orders.Handler
)

func initHandler() {
	orderHandler = orders.NewHandler(orders.NewService(orders.NewRepo()))
}

func route() {
	initHandler()

	router.POST("/orders", orderHandler.CreateOrder)
}
