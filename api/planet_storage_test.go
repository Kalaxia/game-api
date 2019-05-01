package api

func getStorageMock() *Storage {
	return &Storage{
		Id: 1,
		Capacity: 5000,
		Resources: map[string]uint16{
			"cristal": 2500,
			"geode": 1450,
		},
	}
}