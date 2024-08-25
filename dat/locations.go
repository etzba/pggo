package dat

import "github.com/etzba/pggo/wire"

func GetLocations() []wire.Location {
	locations := []wire.Location{
		{
			Id:         1,
			Name:       "",
			Address:    "",
			Longtitude: 56.342155413,
			Latitude:   18.123876123,
		},
		{
			Id:         2,
			Name:       "",
			Address:    "",
			Longtitude: 56.342155413,
			Latitude:   18.123876123,
		},
		{
			Id:         3,
			Name:       "",
			Address:    "",
			Longtitude: 56.342155413,
			Latitude:   18.123876123,
		},
	}
	return locations
}
