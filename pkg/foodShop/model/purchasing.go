package model

import "github.com/TewApirat/food-shop/pkg/foodShop/domain"


type PurchasingRequest struct {
	Items  map[string]int `json:"items"`
	Member bool           `json:"member"`
}

type OrderLine struct {
	Code MenuItemCode
	Qty  int
}

type OrderQuote struct {
	Subtotal       domain.Money
	PairDiscount   domain.Money
	MemberDiscount domain.Money
	Total          domain.Money
}
