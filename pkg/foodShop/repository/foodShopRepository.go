package repository

import "github.com/TewApirat/food-shop/pkg/foodShop/model"

type FoodShopRepository interface {
	ListMenuItems() ([]model.MenuItem, error)
	FindMenuItemByCode(code model.MenuItemCode) (model.MenuItem, error)
	ListPromotions() ([]model.Promotion, error)
}
