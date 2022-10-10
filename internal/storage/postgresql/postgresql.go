package postgresql

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/zhd68/tz-pizzasoft/internal/model"
	"github.com/zhd68/tz-pizzasoft/internal/storage"
	"github.com/zhd68/tz-pizzasoft/pkg/generatorid"
)

type Storage struct {
	db             *sql.DB
	generateIDFunc func() string
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("postgres", path)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}

	return &Storage{
		db:             db,
		generateIDFunc: generatorid.GenerateId(),
	}, nil
}

func (s *Storage) Init() error {
	q := `CREATE TABLE IF NOT EXISTS orders (
		order_id TEXT NOT NULL UNIQUE,
		items INTEGER[] NOT NULL,
		done BOOLEAN DEFAULT false
	);`

	_, err := s.db.Exec(q)
	if err != nil {
		return fmt.Errorf("can't create table: %w", err)
	}

	return nil
}

func (s *Storage) SaveOrder(o *model.Order) (*model.Order, error) {
	q := `INSERT INTO orders (order_id, items) VALUES ($1, $2);`

	id := s.generateIDFunc()
	o.OrderID = id

	if _, err := s.db.Exec(q, o.OrderID, pq.Array(o.Items)); err != nil {
		return nil, fmt.Errorf("can't save order: %w", err)
	}

	return o, nil
}

func (s *Storage) UpdateOrder(id string, items model.Items) (*model.Order, error) {
	q := `UPDATE orders SET items = $1 WHERE order_id = $2;`

	order, err := s.GetOrder(id)
	if err != nil {
		return nil, err
	}
	if order.Done {
		return nil, fmt.Errorf("order <%s> is done", id)
	}
	order.Items = append(order.Items, items...)

	if _, err := s.db.Exec(q, pq.Array(order.Items), order.OrderID); err != nil {
		return nil, fmt.Errorf("can't update order: %w", err)
	}

	return order, nil
}

func (s *Storage) GetOrder(id string) (*model.Order, error) {
	q := `SELECT order_id, items, done FROM orders WHERE order_id=$1;`

	order := &model.Order{}
	arr := []int32{}

	err := s.db.QueryRow(q, id).Scan(&order.OrderID, pq.Array(&arr), &order.Done)
	if err != nil {
		return nil, fmt.Errorf("order <%s> not exist %w", id, err)
	}

	order.Items = arr

	return order, nil
}

func (s *Storage) DoneOrder(id string) (*model.Order, error) {
	q := `UPDATE orders SET done = true WHERE order_id = $1;`

	order, err := s.GetOrder(id)
	if err != nil {
		return nil, err
	}
	if order.Done {
		return nil, fmt.Errorf("order <%s> is done", id)
	}

	order.Done = true

	if _, err := s.db.Exec(q, order.OrderID); err != nil {
		return nil, fmt.Errorf("can't update order: %w", err)
	}

	return order, nil
}

func (s *Storage) GetAllOrders(done *bool) ([]*model.Order, error) {
	qAll := `SELECT * FROM orders;`
	qDone := `SELECT * FROM orders WHERE done = $1;`

	var rows *sql.Rows
	var err error

	if done == nil {
		rows, err = s.db.Query(qAll)
		if err == sql.ErrNoRows {
			return nil, storage.ErrNoSavedOrders
		}
		if err != nil {
			return nil, fmt.Errorf("can't get orders: %w", err)
		}
	} else {
		rows, err = s.db.Query(qDone, *done)
		if err == sql.ErrNoRows {
			return nil, storage.ErrNoSavedOrders
		}
		if err != nil {
			return nil, fmt.Errorf("can't get orders: %w", err)
		}
	}
	defer rows.Close()

	var orders []*model.Order

	for rows.Next() {
		order := &model.Order{}
		arr := []int32{}
		err = rows.Scan(&order.OrderID, pq.Array(&arr), &order.Done)
		if err != nil {
			return nil, err
		}
		order.Items = arr
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
