package exception

import (
	"fmt"
	"github.com/TewApirat/food-shop/pkg/foodShop/model"
)

type MenuItemPriceMissingError struct {
	Code model.MenuItemCode
}

func (e *MenuItemPriceMissingError) Error() string {
	return fmt.Sprintf("Error: menu item price missing for code: %s", e.Code)
}