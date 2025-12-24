package repository

import (
	"github.com/stretchr/testify/mock"

	"github.com/TewApirat/food-shop/pkg/foodShop/model"
)

type FoodShopRepositoryMock struct {
	mock.Mock
}

func (m *FoodShopRepositoryMock) FindMenuItemByCode(code model.MenuItemCode) (model.MenuItem, error) {
	args := m.Called(code)
	return args.Get(0).(model.MenuItem), args.Error(1)
}

func (m *FoodShopRepositoryMock) ListMenuItems() ([]model.MenuItem, error) {
	args := m.Called()
	return args.Get(0).([]model.MenuItem), args.Error(1)
}

func (m *FoodShopRepositoryMock) ListPromotions() ([]model.Promotion, error) {
	args := m.Called()
	return args.Get(0).([]model.Promotion), args.Error(1)
}
