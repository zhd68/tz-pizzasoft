package storage

import (
	"errors"

	"github.com/zhd68/tz-pizzasoft/internal/model"
)

type Storage interface {
	SaveOrder(o *model.Order) (*model.Order, error)
	UpdateOrder(id string, items model.Items) (*model.Order, error)
	GetOrder(id string) (*model.Order, error)
	DoneOrder(id string) (*model.Order, error)
	GetAllOrders(done *bool) ([]*model.Order, error)
}

var ErrNoSavedOrders = errors.New("no save orders")
