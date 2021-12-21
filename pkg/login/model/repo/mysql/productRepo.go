package databases

import (
	"fullstack-shop/pkg/login/model/obj/product"
	"fullstack-shop/pkg/login/model/repo"
)

type ProductDBRepo struct {
}

func (u ProductDBRepo) GetByKey(key repo.BDKey) (*product.Product, error) {
	return nil, nil
}

func (u ProductDBRepo) Post(product *product.Product) error {
	return nil
}
