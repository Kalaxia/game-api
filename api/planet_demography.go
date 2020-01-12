package api

import(
	"math"
)

const populationPointRatio = 1000000

func (p *Planet) calculatePopulationGrowth() {
	p.Population = uint(math.Ceil(float64(p.Population) * (1.0 + p.calculatePopulationGrowthRate() - p.calculatePopulationDeclineRate())))
}

func (p *Planet) calculatePopulationDeclineRate() float64 {
	return 0.005
}

func (p *Planet) calculatePopulationGrowthRate() float64 {
	return 0.015
}