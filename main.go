package main

import (
	"os"

	_foodShopRepository "github.com/TewApirat/food-shop/pkg/foodShop/repository"
	_foodShopService "github.com/TewApirat/food-shop/pkg/foodShop/service"
	_foodShopController "github.com/TewApirat/food-shop/pkg/foodShop/controller"
)

func main() {
	foodShopRepository := _foodShopRepository.NewFoodShopRepositoryDefault()
	foodShopService := _foodShopService.NewFoodShopServiceImpl(foodShopRepository)
	foodShopController := _foodShopController.NewFoodShopControllerImpl(
		os.Stdin, 
		os.Stdout, 
		foodShopService,
	)

	foodShopController.ServeCLI()
}
