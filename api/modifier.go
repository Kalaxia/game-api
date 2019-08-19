package api

const(
	ModifierTypeResource = "resource"
	ModifierTypeMoney = "credits"
	ModifierTypePoints = "points"
)

type(
	Modifier struct{
		Type string `json:"type"`
		Resource string `json:"resource"`
		Percent uint8 `json:"percent"`
	}
)