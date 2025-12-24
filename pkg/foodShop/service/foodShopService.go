package service

import (
	_foodShopModel "github.com/TewApirat/food-shop/pkg/foodShop/model"
	_orderHistoryModel "github.com/TewApirat/food-shop/pkg/orderHistory/model"
	

)
type FoodShopService interface {
	GetMenuCatalog() ([]_foodShopModel.MenuItem, error)
	GetPromotions() ([]_foodShopModel.Promotion, error)
	QuoteOrder(req _foodShopModel.PurchasingRequest) (_foodShopModel.OrderQuote, error)
	ListOrderHistory() ([]_orderHistoryModel.OrderHistoryEntry, error)
	CountOrderHistory() (int, error)

}
