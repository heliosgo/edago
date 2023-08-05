package main

import "edago"

type PayService interface {
	Pay(userId int, orderId string)
}

type payService struct {
	eventBus edago.EventBus
}

func NewPayService(eventBus edago.EventBus) PayService {
	return &payService{
		eventBus: eventBus,
	}
}

var _ PayService = (*payService)(nil)

func (p *payService) Pay(userId int, orderId string) {
	// pay finished
	p.eventBus.Publish(payFinished, userId, orderId)
}
