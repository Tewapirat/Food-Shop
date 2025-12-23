package service

import "github.com/TewApirat/food-shop/pkg/foodShop/model"

type FoodShopService interface {
	GetMenuCatalog() ([]model.MenuItem, error)
	GetPromotions() ([]model.Promotion, error)
	QuoteOrder(req model.PurchasingRequest) (model.OrderQuote, error)
}
