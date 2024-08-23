package wire

type Location struct {
	Id         int     `json:"id"`
	Name       string  `json:"name"`
	Address    string  `json:"address"`
	Longtitude float64 `json:"longtitude"`
	Latitude   float64 `json:"latitude"`
}
