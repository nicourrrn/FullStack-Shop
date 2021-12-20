package repo

import (
	"fullstack-shop/pkg/login/model/obj/order"
	"fullstack-shop/pkg/login/model/obj/product"
)

type BDKey struct {
	Key, Data string
}

// TODO Запушить с кооперативного компа обновленные User и Supplier
//type UserRepositoryInterface interface {
//	GetByKey(key BDKey) (*models.User, error)
//	Post(user *models.User) error
//}
//
//type SupplierRepositoryInterface interface {
//	GetByKey(key BDKey) (*models.Supplier, error)
//	Post(user *models.Supplier) error
//}

type ProductRepositoryInterface interface {
	GetByKey(key BDKey) (*product.Product, error)
	Post(user *product.Product) error
}

type OrderRepositoryInterface interface {
	GetByKey(key BDKey) (*order.Order, error)
	Post(user *order.Order) error
}
