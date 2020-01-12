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