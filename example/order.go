package main

import (
	"edago"
	"fmt"

	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(userId int) string
}

type orderService struct {
	eventBus edago.EventBus
}

func NewOrderService(eventBus edago.EventBus) OrderService {
	res := &orderService{
		eventBus: eventBus,
	}
	res.eventBus.Subscribe(payFinished, res.payFinishedHandler)

	return res
}

var _ OrderService = (*orderService)(nil)

func (o *orderService) payFinishedHandler(userId int, orderId string) {
	fmt.Println("persist order")
}

func (o *orderService) CreateOrder(userId int) string {
	return uuid.New().String()
}
