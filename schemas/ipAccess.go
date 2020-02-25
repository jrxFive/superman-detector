package schemas

type IPAccess struct {
	Ip        string `json:"ip"`
	Speed     int    `json:"speed"`
	Timestamp int64  `json:"timestamp"`
	Geo
}
