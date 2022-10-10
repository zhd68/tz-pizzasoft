package memory

import (
	"fmt"

	"github.com/zhd68/tz-pizzasoft/internal/model"
	"github.com/zhd68/tz-pizzasoft/internal/storage"
	"github.com/zhd68/tz-pizzasoft/pkg/generatorid"
)

type Storage struct {
	data           map[string]*model.Order
	generateIDFunc func() string
}

func New() *Storage {
	return &Storage{
		data:           make(map[string]*model.Order),
		generateIDFunc: generatorid.GenerateId(),
	}
}

func (s *Storage) SaveOrder(o *model.Order) (*model.Order, error) {
	id := s.generateIDFunc()
	o.OrderID = id
	s.data[o.OrderID] = o
	return o, nil
}

func (s *Storage) UpdateOrder(id string, items model.Items) (*model.Order, error) {
	order, ok := s.data[id]
	if !ok {
		return nil, fmt.Errorf("order <%s> not exist", id)
	}
	if order.Done {
		return nil, fmt.Errorf("order <%s> is done", id)
	}
	order.Items = append(order.Items, items...)
	s.data[id] = order
	return order, nil
}

func (s *Storage) GetOrder(id string) (*model.Order, error) {
	order, ok := s.data[id]
	if !ok {
		return nil, fmt.Errorf("order <%s> not exist", id)
	}
	return order, nil
}

func (s *Storage) DoneOrder(id string) (*model.Order, error) {
	order, ok := s.data[id]
	if !ok {
		return nil, fmt.Errorf("order <%s> not exist", id)
	}
	if order.Done {
		return nil, fmt.Errorf("order <%s> is done", id)
	}
	order.Done = true
	s.data[id] = order
	return order, nil
}

func (s *Storage) GetAllOrders(done *bool) ([]*model.Order, error) {
	if len(s.data) == 0 {
		return nil, storage.ErrNoSavedOrders
	}
	var orders []*model.Order
	if done == nil {
		for _, order := range s.data {
			orders = append(orders, order)
		}
	} else {
		for _, order := range s.data {
			if order.Done == *done {
				orders = append(orders, order)
			}
		}
	}

	return orders, nil
}
