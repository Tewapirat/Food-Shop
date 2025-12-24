package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/TewApirat/food-shop/pkg/foodShop/domain"
	_foodShopException "github.com/TewApirat/food-shop/pkg/foodShop/exception"
	_foodShopModel "github.com/TewApirat/food-shop/pkg/foodShop/model"
	_foodShopRepository "github.com/TewApirat/food-shop/pkg/foodShop/repository"
	_foodShopService "github.com/TewApirat/food-shop/pkg/foodShop/service"
	_orderHistoryModel "github.com/TewApirat/food-shop/pkg/orderHistory/model"
	_orderHistoryRepository "github.com/TewApirat/food-shop/pkg/orderHistory/repository"
)

func satang(v int64) domain.Money { return domain.Money(v) }

func TestQuoteOrderSuccess(t *testing.T) {
	type tc struct {
		label        string
		in           _foodShopModel.PurchasingRequest
		expected     _foodShopModel.OrderQuote
		expectedQty  map[_foodShopModel.MenuItemCode]int 
		setupMenuMock func(r * _foodShopRepository.FoodShopRepositoryMock)
	}

	cases := []tc{
		{
			label: "Success: no member, pair discount applies to GREEN(2)",
			in: _foodShopModel.PurchasingRequest{
				Items:  map[string]int{"RED": 1, "GREEN": 2},
				Member: false,
			},
			expected: _foodShopModel.OrderQuote{
				Subtotal:       domain.THB(130),
				PairDiscount:   domain.THB(4),
				MemberDiscount: domain.THB(0),
				Total:          domain.THB(126),
			},
			expectedQty: map[_foodShopModel.MenuItemCode]int{
				"RED":   1,
				"GREEN": 2,
			},
			setupMenuMock: func(r * _foodShopRepository.FoodShopRepositoryMock) {
				r.On("FindMenuItemByCode", _foodShopModel.MenuItemCode("RED")).
					Return(_foodShopModel.MenuItem{Code: "RED", Name: "Red set", Price: domain.THB(50)}, nil).
					Once()
				r.On("FindMenuItemByCode", _foodShopModel.MenuItemCode("GREEN")).
					Return(_foodShopModel.MenuItem{Code: "GREEN", Name: "Green set", Price: domain.THB(40)}, nil).
					Once()
			},
		},
		{
			label: "Success: member stacks after pair discount (GREEN 2)",
			in: _foodShopModel.PurchasingRequest{
				Items:  map[string]int{"GREEN": 2},
				Member: true,
			},
			expected: _foodShopModel.OrderQuote{
				Subtotal:       domain.THB(80),
				PairDiscount:   domain.THB(4),
				MemberDiscount: satang(760),  // 7.60 THB
				Total:          satang(6840), // 68.40 THB
			},
			expectedQty: map[_foodShopModel.MenuItemCode]int{
				"GREEN": 2,
			},
			setupMenuMock: func(r * _foodShopRepository.FoodShopRepositoryMock) {
				r.On("FindMenuItemByCode", _foodShopModel.MenuItemCode("GREEN")).
					Return(_foodShopModel.MenuItem{Code: "GREEN", Name: "Green set", Price: domain.THB(40)}, nil).
					Once()
			},
		},
		{
			label: "Success: normalize code (\" green \") works",
			in: _foodShopModel.PurchasingRequest{
				Items:  map[string]int{" green ": 2},
				Member: false,
			},
			expected: _foodShopModel.OrderQuote{
				Subtotal:       domain.THB(80),
				PairDiscount:   domain.THB(4),
				MemberDiscount: domain.THB(0),
				Total:          domain.THB(76),
			},
			expectedQty: map[_foodShopModel.MenuItemCode]int{
				"GREEN": 2,
			},
			setupMenuMock: func(r * _foodShopRepository.FoodShopRepositoryMock) {
				// normalize แล้วจะเรียกด้วย "GREEN"
				r.On("FindMenuItemByCode", _foodShopModel.MenuItemCode("GREEN")).
					Return(_foodShopModel.MenuItem{Code: "GREEN", Name: "Green set", Price: domain.THB(40)}, nil).
					Once()
			},
		},
	}

	for _, c := range cases {
		t.Run(c.label, func(t *testing.T) {
			foodShopRepositoryMock := new(_foodShopRepository.FoodShopRepositoryMock)
			orderHistoryRepositoryMock := new(_orderHistoryRepository.OrderHistoryRepositoryMock)

			// เตรียม menu item mock ให้ตรงกับเคส
			c.setupMenuMock(foodShopRepositoryMock)

			// Expect: QuoteOrder สำเร็จต้อง Add 1 ครั้ง
			orderHistoryRepositoryMock.
				On("Add", mock.MatchedBy(func(entry _orderHistoryModel.OrderHistoryEntry) bool {
					// orderNo จะเริ่ม 1 เพราะเราสร้าง service ใหม่ต่อ subtest
					if entry.OrderNo != 1 {
						return false
					}
					// CreatedAt ควรถูกเซ็ต (ไม่ zero)
					if entry.CreatedAt.IsZero() {
						return false
					}
					// ป้องกันค่าหลุด ๆ แบบอนาคตไกลเกิน (optional)
					if time.Since(entry.CreatedAt) > time.Minute {
						return false
					}

					if entry.Member != c.in.Member {
						return false
					}
					if entry.Subtotal != c.expected.Subtotal {
						return false
					}
					if entry.PairDiscount != c.expected.PairDiscount {
						return false
					}
					if entry.MemberDiscount != c.expected.MemberDiscount {
						return false
					}
					if entry.Total != c.expected.Total {
						return false
					}

					// ตรวจ Lines แบบไม่พึ่งลำดับจาก map iteration
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

			result, err := foodShopService.QuoteOrder(c.in)
			assert.NoError(t, err)


			assert.Equal(t, c.expected.Subtotal, result.Subtotal)
			assert.Equal(t, c.expected.PairDiscount, result.PairDiscount)
			assert.Equal(t, c.expected.MemberDiscount, result.MemberDiscount)
			assert.Equal(t, c.expected.Total, result.Total)

			foodShopRepositoryMock.AssertExpectations(t)
			orderHistoryRepositoryMock.AssertExpectations(t)
		})
	}
}
func TestQuoteOrderFail(t *testing.T) {
	type tc struct {
		label string
		in    _foodShopModel.PurchasingRequest
		check func(t *testing.T, err error)
	}

	cases := []tc{
		{
			label: "Fail: empty order",
			in: _foodShopModel.PurchasingRequest{
				Items:  map[string]int{},
				Member: false,
			},
			check: func(t *testing.T, err error) {
				var target *_foodShopException.EmptyOrderError
				assert.ErrorAs(t, err, &target)
			},
		},
		{
			label: "Fail: invalid quantity",
			in: _foodShopModel.PurchasingRequest{
				Items:  map[string]int{"RED": 0},
				Member: false,
			},
			check: func(t *testing.T, err error) {
				var target *_foodShopException.InvalidQuantityError
				assert.ErrorAs(t, err, &target)
				assert.Equal(t, 0, target.Qty)
			},
		},
		{
			label: "Fail: invalid item code (blank)",
			in: _foodShopModel.PurchasingRequest{
				Items:  map[string]int{"   ": 1},
				Member: false,
			},
			check: func(t *testing.T, err error) {
				var target *_foodShopException.InvalidItemCodeError
				assert.ErrorAs(t, err, &target)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.label, func(t *testing.T) {
			foodShopRepositoryMock := new(_foodShopRepository.FoodShopRepositoryMock)
			orderHistoryRepositoryMock := new(_orderHistoryRepository.OrderHistoryRepositoryMock)

			foodShopService := _foodShopService.NewFoodShopServiceImpl(
				foodShopRepositoryMock,
				orderHistoryRepositoryMock,
			)

			result, err := foodShopService.QuoteOrder(c.in)

			assert.Equal(t, _foodShopModel.OrderQuote{}, result)
			assert.Error(t, err)
			c.check(t, err)

			// Fail เคสเหล่านี้ return ก่อนเรียก repo/add จึงไม่ต้อง set expectation เพิ่ม
			foodShopRepositoryMock.AssertExpectations(t)
			orderHistoryRepositoryMock.AssertExpectations(t)
		})
	}
}
