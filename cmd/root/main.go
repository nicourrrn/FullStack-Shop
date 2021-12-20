package main

import (
	"fmt"
	orderModel "fullstack-shop/pkg/kernel/model/obj/order"
)

func main() {
	order := orderModel.NewOrder(10)
	order.AddProduct(1)
	order.AddProduct(2)
	order.Cell(true)
	fmt.Printf("%#v", order)
}
