package main

import (
	"fmt"
	orderModel "fullstack-shop/pkg/kernel/model/obj/order"
	"github.com/nicourrrn/littleLogger"
	"log"
	"os"
)

const DEBUG = true

func main() {
	logFile, _ := os.Create("logs/log.txt")
	logger, err := littleLogger.NewLogger(logFile, DEBUG)
	defer logger.Wait()
	if err != nil {
		log.Fatalln(err)
	}
	logger.Info("Hello")
	logger.SetFormatter(littleLogger.FormatterClassic)
	order := orderModel.NewOrder(10)
	logger.Debug("New order make")
	order.AddProduct(1)
	order.AddProduct(2)
	order.Cell(true)
	fmt.Printf("%#v", order)
}
