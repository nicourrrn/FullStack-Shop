package main

import (
	"fmt"
	orderModel "fullstack-shop/pkg/kernel/model/obj/order"
	"github.com/nicourrrn/littleLogger"
	"log"
	"os"
)

func main() {
	logFile, _ := os.Create("logs/log.txt")
	logger, err := littleLogger.NewLogger(logFile, 1)
	if err != nil {
		log.Fatalln(err)
	}
	logger.Info("Hello")
	logger.SetFormatter(func() string {
		return "$msg\n"
	})
	order := orderModel.NewOrder(10)
	logger.Debug("New order make")
	order.AddProduct(1)
	order.AddProduct(2)
	order.Cell(true)
	fmt.Printf("%#v", order)
}
