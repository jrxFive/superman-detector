package schemas

type Geo struct {
	Lat    float64 `json:"lat"`
	Lon    float64 `json:"lon"`
	Radius uint16  `json:"radius"`
}
