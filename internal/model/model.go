package model

import "fmt"

type Order struct {
	OrderID string `json:"order_id"`
	Items   Items  `json:"items"`
	Done    bool   `json:"done"`
}

type Items []int32

//type Item int

func (i Items) ValidateItems() error {
	if i == nil {
		return fmt.Errorf("'items' is nil")
	}
	if len(i) == 0 {
		return fmt.Errorf("invalid length array 'items'")
	}
	for _, item := range i {
		if item <= 0 || item > 5000 {
			return fmt.Errorf("invalid items value")
		}
	}
	return nil
}
