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
		Percent int8 `json:"percent"`
	}
)