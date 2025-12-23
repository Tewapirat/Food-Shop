package model

import "github.com/TewApirat/food-shop/pkg/foodShop/domain"

type MenuItemCode string

type MenuItem struct {
	Code  MenuItemCode
	Name  string
	Price domain.Money
}
