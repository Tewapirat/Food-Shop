package model

// PairDiscountPolicy เป็นโครงสร้าง "กฎโปรคู่" (data/config เท่านั้น)
type PairDiscountPolicy struct {
	EligibleCodes map[MenuItemCode]bool // เช่น ORANGE/PINK/GREEN
	DiscountPercent int64                // เช่น 5 (หมายถึง 5%)
	BundleSize int                      // เช่น 2 (หมายถึง 2 ชิ้นต่อ 1 คู่)
}

// DefaultPairDiscountPolicy เป็น default config ของโปรคู่
func DefaultPairDiscountPolicy() PairDiscountPolicy {
	return PairDiscountPolicy{
		EligibleCodes: map[MenuItemCode]bool{
			"ORANGE": true,
			"PINK":   true,
			"GREEN":  true,
		},
		DiscountPercent: 5,
		BundleSize:      2,
	}
}
