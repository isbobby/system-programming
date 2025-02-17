package main

func main() {
	InspectStruct(Bob{
		Name:   "Bob",
		Age:    25,
		Height: 178.8,
		HousePet: Pet{
			Name:    "Ruffer",
			Species: "dog",
		},
	})
}
