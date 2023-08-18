package event

import "time"

type OrderDeleted struct {
	Name    string
	Payload interface{}
}

func NewOrderDeleted() *OrderDeleted {
	return &OrderDeleted{
		Name: "OrderDeleted",
	}
}

func (e *OrderDeleted) GetName() string {
	return e.Name
}

func (e *OrderDeleted) GetPayload() interface{} {
	return e.Payload
}

func (e *OrderDeleted) SetPayload(payload interface{}) {
	e.Payload = payload
}

func (e *OrderDeleted) GetDateTime() time.Time {
	return time.Now()
}
