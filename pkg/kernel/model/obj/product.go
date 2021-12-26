package obj

import (
	"errors"
	"time"
)

type Product struct {
	ID          int
	CategoryID  int
	Prices      map[time.Time]float32
	Name        string
	Description string
}

func NewProduct(categoryID int, price float32, name, desc string) *Product {
	return &Product{CategoryID: categoryID, Prices: map[time.Time]float32{time.Now(): price}, Name: name, Description: desc}
}

func (p *Product) AddPrice(newPrice float32) {
	p.Prices[time.Now()] = newPrice
}

func (p Product) GetPriceByDate(prodDate time.Time) (float32, error) {
	foundedDate := time.Time{}
	for date, _ := range p.Prices {
		//TODO Переписать проверерку
		if prodDate.After(date) && date.Sub(prodDate) > foundedDate.Sub(date) {
			foundedDate = date
		}
	}
	if foundedDate.IsZero() {
		return -1, errors.New("date not found")
	}
	return p.Prices[foundedDate], nil

}
