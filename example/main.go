package main

import (
	"edago"
	"time"
)

func main() {
	bus := edago.NewEventBus()
	payService := NewPayService(bus)
	orderService := NewOrderService(bus)

	userId := 1
	orderId := orderService.CreateOrder(userId)
	payService.Pay(userId, orderId)

	time.Sleep(1 * time.Second)
}
