package main

import (
	"os"

	_foodShopRepository "github.com/TewApirat/food-shop/pkg/foodShop/repository"
	_foodShopService "github.com/TewApirat/food-shop/pkg/foodShop/service"
	_foodShopController "github.com/TewApirat/food-shop/pkg/foodShop/controller"
	_orderHistoryReppsitory "github.com/TewApirat/food-shop/pkg/orderHistory/repository"
)

func main() {
	foodShopRepository := _foodShopRepository.NewFoodShopRepositoryDefault()
	orderHistoryRepository := _orderHistoryReppsitory.NewOrderHistoryRepositoryImpl()

	foodShopService := _foodShopService.NewFoodShopServiceImpl(
		foodShopRepository,
		orderHistoryRepository,	
	)
	foodShopController := _foodShopController.NewFoodShopControllerImpl(
		os.Stdin, 
		os.Stdout, 
		foodShopService,
	)

	foodShopController.ServeCLI()
}
