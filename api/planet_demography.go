package api

import(
	"math"
)

const(
	populationPointRatio = 1000000
	planetPublicOrderRebellious = 20
	planetPublicOrderUnrest = 40
	planetPublicOrderStable = 60
	planetPublicOrderGood = 80
	planetPublicOrderPerfect = 95
	planetTaxRateVeryLow = 1
	planetTaxRateLow = 2
	planetTaxRateNormal = 3
	planetTaxRateHigh = 4
	planetTaxRateVeryHigh = 5
)

func (p *Planet) calculatePopulationGrowth() {
	p.Population = uint(math.Ceil(float64(p.Population) * (1.0 + p.calculatePopulationGrowthRate() - p.calculatePopulationDeclineRate())))
}

func (p *Planet) calculatePopulationDeclineRate() float64 {
	return 0.005
}

func (p *Planet) calculatePopulationGrowthRate() float64 {
	return 0.015
}

func (p *Planet) calculateTaxes() {
	wage := int32(math.Floor(float64(p.Population) * 0.0001 * float64(p.TaxRate)))
	p.Player.updateWallet(wage)
	p.Player.update()

	po := int8(p.PublicOrder) + p.processTaxesPublicOrderEffect()
	if po < 0 {
		po = 0
	}
	p.PublicOrder = uint8(po)
}

func (p *Planet) processTaxesPublicOrderEffect() int8 {
	return map[uint8]int8{
		planetTaxRateVeryLow: 2,
		planetTaxRateLow: 1,
		planetTaxRateNormal: 0,
		planetTaxRateHigh: -1,
		planetTaxRateVeryHigh: -2,
	}[p.TaxRate]
}