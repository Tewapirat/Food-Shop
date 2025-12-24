package tests

import _foodShopModel "github.com/TewApirat/food-shop/pkg/foodShop/model"

func lineQtyMap(lines []_foodShopModel.OrderLine) map[_foodShopModel.MenuItemCode]int {
	m := make(map[_foodShopModel.MenuItemCode]int, len(lines))
	for _, ln := range lines {
		m[ln.Code] += ln.Qty
	}
	return m
}
