package models

type LoginEvent struct {
	Model
	Username  string `gorm:"index:idx_username" json:"username"`
	Timestamp int64  `gorm:"index:idx_timestamp,idx_username" json:"timestamp"`
	EventUuid string
	IPAddress string
	Latitude  float64
	Longitude float64
	Radius    uint16
}
