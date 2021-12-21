package databases

import (
	"fullstack-shop/pkg/login/model/obj/order"
	"fullstack-shop/pkg/login/model/repo"
)

type OrderDBRepo struct {
}

func (u OrderDBRepo) GetByKey(key repo.BDKey) (*order.Order, error) {
	return nil, nil
}

func (u OrderDBRepo) Post(order *order.Order) error {
	return nil
}
