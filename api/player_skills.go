package api

import(
	"math"
)

type(
	PlayerSkills struct {
		tableName struct{} `pg:"player__skills"`

		Id uint16 `json:"id"`
		Combat uint16 `json:"combat"`
		Industry uint16 `json:"industry"`
		Politics uint16 `json:"politics"`
		Production uint16 `json:"production"`
		Trade uint16 `json:"trade"`
	}
)

const(
	PlayerClassPolitician = "politician"
	PlayerClassIndustrialist = "industrialist"
	PlayerClassStrategist = "strategist"
	PlayerClassTrader = "trader"
	PlayerSkillCombat = "combat"
	PlayerSkillIndustry = "industry"
	PlayerSkillPolitics = "politics"
	PlayerSkillProduction = "production"
	PlayerSkillTrade = "trade"
)

func (p *Player) createSkills(class string) {
	skillPoints := getClassSkillPoints(class)

	skills := &PlayerSkills{
		Combat: skillPoints[PlayerSkillCombat],
		Industry: skillPoints[PlayerSkillIndustry],
		Politics: skillPoints[PlayerSkillPolitics],
		Production: skillPoints[PlayerSkillProduction],
		Trade: skillPoints[PlayerSkillTrade],
	}
	if err := Database.Insert(skills); err != nil {
		panic(NewException("Could not create player skills", err))
	}
	p.Skills = skills
}

func (p *Player) increaseSkill(skill string, points uint16) {
	playerSkillAffinities := map[string]map[string]float64{
		PlayerSkillCombat: map[string]float64{
			PlayerSkillCombat: 1,
			PlayerSkillPolitics: 0.75,
			PlayerSkillTrade: 0.75,
			PlayerSkillProduction: -0.75,
			PlayerSkillIndustry: -0.75,
		},
		PlayerSkillPolitics: map[string]float64{
			PlayerSkillCombat: -0.75,
			PlayerSkillPolitics: 1,
			PlayerSkillTrade: 0.75,
			PlayerSkillProduction: -0.75,
			PlayerSkillIndustry: -0.75,
		},
		PlayerSkillTrade: map[string]float64{
			PlayerSkillCombat: 0.75,
			PlayerSkillPolitics: 0.75,
			PlayerSkillTrade: 1,
			PlayerSkillProduction: -0.75,
			PlayerSkillIndustry: -0.75,
		},
		PlayerSkillProduction: map[string]float64{
			PlayerSkillCombat: -0.75,
			PlayerSkillPolitics: -0.75,
			PlayerSkillTrade: 0.75,
			PlayerSkillProduction: 1,
			PlayerSkillIndustry: 0.75,
		},
		PlayerSkillIndustry: map[string]float64{
			PlayerSkillCombat: -0.75,
			PlayerSkillPolitics: -0.75,
			PlayerSkillTrade: 0.75,
			PlayerSkillProduction: 0.75,
			PlayerSkillIndustry: 1,
		},
	}
	p.Skills.Combat += uint16(math.Floor(float64(points) * playerSkillAffinities[skill][PlayerSkillCombat]))
	p.Skills.Industry += uint16(math.Floor(float64(points) * playerSkillAffinities[skill][PlayerSkillIndustry]))
	p.Skills.Politics += uint16(math.Floor(float64(points) * playerSkillAffinities[skill][PlayerSkillPolitics]))
	p.Skills.Production += uint16(math.Floor(float64(points) * playerSkillAffinities[skill][PlayerSkillProduction]))
	p.Skills.Trade += uint16(math.Floor(float64(points) * playerSkillAffinities[skill][PlayerSkillTrade]))
}

func getClassSkillPoints(class string) map[string]uint16 {
	skills := map[string]map[string]uint16{
		PlayerClassIndustrialist: map[string]uint16{
			PlayerSkillCombat: 0,
			PlayerSkillIndustry: 100,
			PlayerSkillPolitics: 0,
			PlayerSkillProduction: 75,
			PlayerSkillTrade: 75,
		},
		PlayerClassPolitician: map[string]uint16{
			PlayerSkillCombat: 50,
			PlayerSkillIndustry: 0,
			PlayerSkillPolitics: 100,
			PlayerSkillProduction: 25,
			PlayerSkillTrade: 75,
		},
		PlayerClassStrategist: map[string]uint16{
			PlayerSkillCombat: 100,
			PlayerSkillIndustry: 25,
			PlayerSkillPolitics: 75,
			PlayerSkillProduction: 0,
			PlayerSkillTrade: 50,
		},
		PlayerClassTrader: map[string]uint16{
			PlayerSkillCombat: 0,
			PlayerSkillIndustry: 75,
			PlayerSkillPolitics: 75,
			PlayerSkillProduction: 0,
			PlayerSkillTrade: 100,
		},
	}
	return skills[class]
}