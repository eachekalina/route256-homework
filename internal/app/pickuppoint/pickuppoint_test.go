package pickuppoint

var SamplePickUpPoint = PickUpPoint{
	Id:      1,
	Name:    "Generic pick-up point",
	Address: "5, Test st., Moscow",
	Contact: "test@example.com",
}

var SamplePickUpPointSlice = []PickUpPoint{
	{
		Id:      1,
		Name:    "Generic pick-up point",
		Address: "5, Test st., Moscow",
		Contact: "test@example.com",
	},
	{
		Id:      2,
		Name:    "Another pick-up point",
		Address: "19, Sample st., Moscow",
		Contact: "sample@example.com",
	},
}
