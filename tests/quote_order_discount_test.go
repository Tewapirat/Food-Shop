package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/TewApirat/food-shop/pkg/foodShop/domain"
	_foodShopModel "github.com/TewApirat/food-shop/pkg/foodShop/model"
	_foodShopRepository "github.com/TewApirat/food-shop/pkg/foodShop/repository"
	_foodShopService "github.com/TewApirat/food-shop/pkg/foodShop/service"

	_orderHistoryModel "github.com/TewApirat/food-shop/pkg/orderHistory/model"
	_orderHistoryRepository "github.com/TewApirat/food-shop/pkg/orderHistory/repository"
)

func TestQuoteOrder_DiscountPolicies(t *testing.T) {
	type tc struct {
		label         string
		in            _foodShopModel.PurchasingRequest
		setupMenuMock func(r *_foodShopRepository.FoodShopRepositoryMock)

		expectedSubtotal       domain.Money
		expectedPairDiscount   domain.Money
		expectedMemberDiscount domain.Money
		expectedTotal          domain.Money

		expectedQty map[_foodShopModel.MenuItemCode]int
	}

	cases := []tc{
		{
			label: "Pair: GREEN(4) => 2 bundles",
			in: _foodShopModel.PurchasingRequest{
				Items:  map[string]int{"GREEN": 4},
				Member: false,
			},
			setupMenuMock: func(r *_foodShopRepository.FoodShopRepositoryMock) {
				r.On("FindMenuItemByCode", _foodShopModel.MenuItemCode("GREEN")).
					Return(_foodShopModel.MenuItem{Code: "GREEN", Name: "Green set", Price: domain.THB(40)}, nil).
					Once()
			},
			expectedSubtotal:       domain.THB(160), // 40*4
			expectedPairDiscount:   domain.THB(8),   // 2 bundles * (5% of 80) = 8
			expectedMemberDiscount: domain.THB(0),
			expectedTotal:          domain.THB(152),

			expectedQty: map[_foodShopModel.MenuItemCode]int{"GREEN": 4},
		},
		{
			label: "Pair: GREEN(3) => 1 bundle + remainder",
			in: _foodShopModel.PurchasingRequest{
				Items:  map[string]int{"GREEN": 3},
				Member: false,
			},
			setupMenuMock: func(r *_foodShopRepository.FoodShopRepositoryMock) {
				r.On("FindMenuItemByCode", _foodShopModel.MenuItemCode("GREEN")).
					Return(_foodShopModel.MenuItem{Code: "GREEN", Name: "Green set", Price: domain.THB(40)}, nil).
					Once()
			},
			expectedSubtotal:       domain.THB(120),
			expectedPairDiscount:   domain.THB(4),
			expectedMemberDiscount: domain.THB(0),
			expectedTotal:          domain.THB(116),
			expectedQty:            map[_foodShopModel.MenuItemCode]int{"GREEN": 3},
		},
		{
			label: "Pair discount: PINK(2) => 1 bundle",
			in: _foodShopModel.PurchasingRequest{
				Items:  map[string]int{"PINK": 2},
				Member: false,
			},
			setupMenuMock: func(r *_foodShopRepository.FoodShopRepositoryMock) {
				r.On("FindMenuItemByCode", _foodShopModel.MenuItemCode("PINK")).
					Return(_foodShopModel.MenuItem{Code: "PINK", Name: "Pink set", Price: domain.THB(80)}, nil).
					Once()
			},
			expectedSubtotal:       domain.THB(160),
			expectedPairDiscount:   domain.THB(8),
			expectedMemberDiscount: domain.THB(0),
			expectedTotal:          domain.THB(152),
			expectedQty:            map[_foodShopModel.MenuItemCode]int{"PINK": 2},
		},
		{
			label: "Pair discount: PINK(4) => 2 bundles",
			in: _foodShopModel.PurchasingRequest{
				Items:  map[string]int{"PINK": 4},
				Member: false,
			},
			setupMenuMock: func(r *_foodShopRepository.FoodShopRepositoryMock) {
				r.On("FindMenuItemByCode", _foodShopModel.MenuItemCode("PINK")).
					Return(_foodShopModel.MenuItem{Code: "PINK", Name: "Pink set", Price: domain.THB(80)}, nil).
					Once()
			},
			expectedSubtotal:       domain.THB(320),
			expectedPairDiscount:   domain.THB(16),
			expectedMemberDiscount: domain.THB(0),
			expectedTotal:          domain.THB(304),
			expectedQty:            map[_foodShopModel.MenuItemCode]int{"PINK": 4},
		},
		{
			label: "Pair: GREEN(2)+ORANGE(2) => sum per code",
			in: _foodShopModel.PurchasingRequest{
				Items:  map[string]int{"GREEN": 2, "ORANGE": 2},
				Member: false,
			},
			setupMenuMock: func(r *_foodShopRepository.FoodShopRepositoryMock) {
				r.On("FindMenuItemByCode", _foodShopModel.MenuItemCode("GREEN")).
					Return(_foodShopModel.MenuItem{Code: "GREEN", Name: "Green set", Price: domain.THB(40)}, nil).
					Once()
				r.On("FindMenuItemByCode", _foodShopModel.MenuItemCode("ORANGE")).
					Return(_foodShopModel.MenuItem{Code: "ORANGE", Name: "Orange set", Price: domain.THB(120)}, nil).
					Once()
			},
			expectedSubtotal:       domain.THB(320), // 80 + 240
			expectedPairDiscount:   domain.THB(16),  // 4 + 12
			expectedMemberDiscount: domain.THB(0),
			expectedTotal:          domain.THB(304),
			expectedQty:            map[_foodShopModel.MenuItemCode]int{"GREEN": 2, "ORANGE": 2},
		},
		{
			label: "Pair: RED(2) not eligible => 0 pair discount",
			in: _foodShopModel.PurchasingRequest{
				Items:  map[string]int{"RED": 2},
				Member: false,
			},
			setupMenuMock: func(r *_foodShopRepository.FoodShopRepositoryMock) {
				r.On("FindMenuItemByCode", _foodShopModel.MenuItemCode("RED")).
					Return(_foodShopModel.MenuItem{Code: "RED", Name: "Red set", Price: domain.THB(50)}, nil).
					Once()
			},
			expectedSubtotal:       domain.THB(100),
			expectedPairDiscount:   domain.THB(0),
			expectedMemberDiscount: domain.THB(0),
			expectedTotal:          domain.THB(100),
			expectedQty:            map[_foodShopModel.MenuItemCode]int{"RED": 2},
		},
		{
			label: "No cross-code pairing: GREEN(1)+ORANGE(1) => 0 pair discount",
			in: _foodShopModel.PurchasingRequest{
				Items:  map[string]int{"GREEN": 1, "ORANGE": 1},
				Member: false,
			},
			setupMenuMock: func(r *_foodShopRepository.FoodShopRepositoryMock) {
				r.On("FindMenuItemByCode", _foodShopModel.MenuItemCode("GREEN")).
					Return(_foodShopModel.MenuItem{Code: "GREEN", Name: "Green set", Price: domain.THB(40)}, nil).
					Once()
				r.On("FindMenuItemByCode", _foodShopModel.MenuItemCode("ORANGE")).
					Return(_foodShopModel.MenuItem{Code: "ORANGE", Name: "Orange set", Price: domain.THB(120)}, nil).
					Once()
			},
			expectedSubtotal:       domain.THB(160), // 40 + 120
			expectedPairDiscount:   domain.THB(0),   // cross-code should NOT form a bundle
			expectedMemberDiscount: domain.THB(0),
			expectedTotal:          domain.THB(160),
			expectedQty: map[_foodShopModel.MenuItemCode]int{
				"GREEN":  1,
				"ORANGE": 1,
			},
		},
		{
			label: "Member without pair: RED(2), member=true => 10% off",
			in: _foodShopModel.PurchasingRequest{
				Items:  map[string]int{"RED": 2},
				Member: true,
			},
			setupMenuMock: func(r *_foodShopRepository.FoodShopRepositoryMock) {
				r.On("FindMenuItemByCode", _foodShopModel.MenuItemCode("RED")).
					Return(_foodShopModel.MenuItem{Code: "RED", Name: "Red set", Price: domain.THB(50)}, nil).
					Once()
			},
			expectedSubtotal:       domain.THB(100),
			expectedPairDiscount:   domain.THB(0),
			expectedMemberDiscount: domain.THB(10), // 10% of 100
			expectedTotal:          domain.THB(90),
			expectedQty: map[_foodShopModel.MenuItemCode]int{
				"RED": 2,
			},
		},
		{
			label: "Member + multi-code: GREEN(2)+RED(1), member=true => member stacks after pair on total",
			in: _foodShopModel.PurchasingRequest{
				Items:  map[string]int{"GREEN": 2, "RED": 1},
				Member: true,
			},
			setupMenuMock: func(r *_foodShopRepository.FoodShopRepositoryMock) {
				r.On("FindMenuItemByCode", _foodShopModel.MenuItemCode("GREEN")).
					Return(_foodShopModel.MenuItem{Code: "GREEN", Name: "Green set", Price: domain.THB(40)}, nil).
					Once()
				r.On("FindMenuItemByCode", _foodShopModel.MenuItemCode("RED")).
					Return(_foodShopModel.MenuItem{Code: "RED", Name: "Red set", Price: domain.THB(50)}, nil).
					Once()
			},
			expectedSubtotal:       domain.THB(130),
			expectedPairDiscount:   domain.THB(4),
			expectedMemberDiscount: satang(1260),  // 12.60 THB (10% of 126.00)
			expectedTotal:          satang(11340), // 113.40 THB

			expectedQty: map[_foodShopModel.MenuItemCode]int{"GREEN": 2, "RED": 1},
		},
	}

	for _, c := range cases {
		t.Run(c.label, func(t *testing.T) {
			foodShopRepositoryMock := new(_foodShopRepository.FoodShopRepositoryMock)
			orderHistoryRepositoryMock := new(_orderHistoryRepository.OrderHistoryRepositoryMock)

			foodShopRepositoryMock.Test(t)
			orderHistoryRepositoryMock.Test(t)

			c.setupMenuMock(foodShopRepositoryMock)

			orderHistoryRepositoryMock.
				On("Add", mock.MatchedBy(func(entry _orderHistoryModel.OrderHistoryEntry) bool {
					if entry.OrderNo != 1 {
						return false
					}
					if entry.CreatedAt.IsZero() {
						return false
					}
					if time.Since(entry.CreatedAt) > time.Minute {
						return false
					}

					if entry.Member != c.in.Member {
						return false
					}
					if entry.Subtotal != c.expectedSubtotal {
						return false
					}
					if entry.PairDiscount != c.expectedPairDiscount {
						return false
					}
					if entry.MemberDiscount != c.expectedMemberDiscount {
						return false
					}
					if entry.Total != c.expectedTotal {
						return false
					}

					gotQty := lineQtyMap(entry.Line)
					if len(gotQty) != len(c.expectedQty) {
						return false
					}
					for code, qty := range c.expectedQty {
						if gotQty[code] != qty {
							return false
						}
					}
					return true
				})).
				Return(nil).
				Once()

			foodShopService := _foodShopService.NewFoodShopServiceImpl(
				foodShopRepositoryMock,
				orderHistoryRepositoryMock,
			)

			res, err := foodShopService.QuoteOrder(c.in)
			assert.NoError(t, err)

			assert.Equal(t, c.expectedSubtotal, res.Subtotal)
			assert.Equal(t, c.expectedPairDiscount, res.PairDiscount)
			assert.Equal(t, c.expectedMemberDiscount, res.MemberDiscount)
			assert.Equal(t, c.expectedTotal, res.Total)

			foodShopRepositoryMock.AssertExpectations(t)
			orderHistoryRepositoryMock.AssertExpectations(t)
		})
	}
}
