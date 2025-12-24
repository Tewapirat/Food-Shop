package model

type PairDiscountPolicy struct {
	EligibleCodes map[MenuItemCode]bool 
	DiscountPercent int64                
	BundleSize int                    
}

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
