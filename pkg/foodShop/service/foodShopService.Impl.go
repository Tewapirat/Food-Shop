package service

import (
	"fmt"
	"strings"

	"github.com/TewApirat/food-shop/pkg/foodShop/domain"
	"github.com/TewApirat/food-shop/pkg/foodShop/model"
	_foodShopException "github.com/TewApirat/food-shop/pkg/foodShop/exception"
	_foodShopRepository "github.com/TewApirat/food-shop/pkg/foodShop/repository"
)

var pairDiscountEligibleCodes = map[model.MenuItemCode]bool{
	"ORANGE": true,
	"PINK":   true,
	"GREEN":  true,
}

const (
	pairBundleSize          = 2
	pairDiscountPercent     = 5
	memberDiscountPercent   = 10
)

type foodShopServiceImpl struct {
	foodShopRepository _foodShopRepository.FoodShopRepository
}

func NewFoodShopServiceImpl(foodShopRepository _foodShopRepository.FoodShopRepository) FoodShopService {
	return &foodShopServiceImpl{foodShopRepository: foodShopRepository}
}

func (s *foodShopServiceImpl) GetMenuCatalog() ([]model.MenuItem, error) {
	return s.foodShopRepository.ListMenuItems()
}

func (s *foodShopServiceImpl) GetPromotions() ([]model.Promotion, error) {
	return s.foodShopRepository.ListPromotions()
}

// 1. check empty order
// 2. pepare state for calculation
// 3. loop items: validate → normalize → lookup → total price


func (s *foodShopServiceImpl) QuoteOrder(req model.PurchasingRequest) (model.OrderQuote, error) {
	if len(req.Items) == 0 {
		return model.OrderQuote{}, &_foodShopException.EmptyOrderError{}
	}

	qtyByCode := make(map[model.MenuItemCode]int)
	priceByCode := make(map[model.MenuItemCode]domain.Money)

	var subtotal domain.Money

	for rawCode, qty := range req.Items {
		if qty < 1 {
			return model.OrderQuote{}, &_foodShopException.InvalidQuantityError{Qty: qty}
		}

		code, err := normalizeItemCode(rawCode)
		if err != nil {
			return model.OrderQuote{}, err
		}

		menuItem, err := s.foodShopRepository.FindMenuItemByCode(code)
		if err != nil {
			// ถ้า repo มี error ชนิด not found ของตัวเอง แนะนำให้ map เป็น domain/usecase error ให้ชัด
			// ที่นี่ทำแบบห่อด้วย context เพื่อ debug ง่าย และปล่อยให้ controller ตัดสินใจ response
			return model.OrderQuote{}, fmt.Errorf("find menu item by code %s: %w", code, err)
		}

		// เก็บราคาต่อ code ไว้ (กัน lookup ซ้ำใน calculatePairDiscount)
		priceByCode[code] = menuItem.Price
		qtyByCode[code] += qty

		subtotal = subtotal.Add(menuItem.Price.MulInt(qty))
	}

	pairDiscount, err := calculatePairDiscount(qtyByCode, priceByCode)
	if err != nil {
		return model.OrderQuote{}, err
	}

	afterPairDiscount := subtotal.Sub(pairDiscount)

	var memberDiscount domain.Money
	if req.Member {
		memberDiscount = afterPairDiscount.Percent(memberDiscountPercent)
	}

	total := afterPairDiscount.Sub(memberDiscount)

	return model.OrderQuote{
		Subtotal:       subtotal,
		PairDiscount:   pairDiscount,
		MemberDiscount: memberDiscount,
		Total:          total,
	}, nil
}

// แยก normalize ออกมาเพื่อให้อ่านง่าย + เทสต์ง่าย
func normalizeItemCode(raw string) (model.MenuItemCode, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", &_foodShopException.InvalidItemCodeError{}
	}
	return model.MenuItemCode(strings.ToUpper(trimmed)), nil
}

// pure function: ไม่มี IO, ไม่มี repo call, ไม่กลืน error
func calculatePairDiscount(
	qtyByCode map[model.MenuItemCode]int,
	priceByCode map[model.MenuItemCode]domain.Money,
) (domain.Money, error) {

	totalDiscount := domain.Money(0)

	for code, qty := range qtyByCode {
		if !pairDiscountEligibleCodes[code] || qty < pairBundleSize {
			continue
		}

		unitPrice, ok := priceByCode[code]
		if !ok {
			// เป็น “สภาวะผิดปกติ” เพราะ code นี้ถูกคิด subtotal ไปแล้วแต่ไม่มีราคา
			return domain.Money(0), &_foodShopException.MenuItemPriceMissingError{Code: code}
		}

		pairCount := qty / pairBundleSize
		bundleValue := unitPrice.MulInt(pairBundleSize)
		discountPerBundle := bundleValue.Percent(pairDiscountPercent)

		totalDiscount = totalDiscount.Add(discountPerBundle.MulInt(pairCount))
	}

	return totalDiscount, nil
}
