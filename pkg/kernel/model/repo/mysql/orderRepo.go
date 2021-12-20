package databases

import (
	"fullstack-shop/pkg/kernel/model/obj/order"
	"fullstack-shop/pkg/kernel/model/repo"
)

type OrderDBRepo struct {
}

func (u OrderDBRepo) GetByKey(key repo.BDKey) (*order.Order, error) {
	return nil, nil
}

func (u OrderDBRepo) Post(order *order.Order) error {
	return nil
}
