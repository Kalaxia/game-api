package api

import(
	"fmt"
)

func (f *ShipFrame) generateShipModelName(p *Player, tonnage uint8, shipType string) string {
	return f.Identifier + f.generateShipTypeIdentifier(shipType) + culturesData[f.Culture].Identifier + "-" + fmt.Sprintf("%d", p.Id % 10) + fmt.Sprintf("%02d", tonnage % 100)
}

func (f *ShipFrame) generateShipTypeIdentifier(shipType string) string {
	return map[string]string{
		ShipTypeFighter: "I",
		ShipTypeBomber: "B",
		ShipTypeFreighter: "R",
		ShipTypeCorvette: "O",
		ShipTypeFrigate: "F",
		ShipTypeCruiser: "U",
	}[shipType]
}