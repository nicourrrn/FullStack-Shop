package order

import (
	"errors"
	"time"
)

type Order struct {
	ID         int
	UserID     int
	ProductIDs []int
	CellDate   time.Time
	Paid       bool
	// start -> adding prod -> send -> finish -> dispute -> returned or finish
	Status string
}

func NewOrder(userID int) *Order {
	return &Order{UserID: userID, ProductIDs: make([]int, 0), Status: "adding prod"}
}

func (o *Order) AddProduct(productID int) error {
	o.ProductIDs = append(o.ProductIDs, productID)
	return nil
}

func (o *Order) RemoveProduct(productID int) error {
	for i, pId := range o.ProductIDs {
		if pId == productID {
			o.ProductIDs = append(o.ProductIDs[:i], o.ProductIDs[i+1:]...)
			return nil
		}
	}
	return errors.New("product not found")
}

func (o *Order) Cell(byCard bool) {
	o.CellDate = time.Now()
	o.Status = "send"
	o.Paid = byCard
}
