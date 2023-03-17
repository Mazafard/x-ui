package xray

import "x-ui/database/model"

type ClientTraffic struct {
	Id         int    `json:"id" form:"id" gorm:"primaryKey;autoIncrement"`
	InboundId  int    `json:"inboundId" form:"inboundId"`
	Enable     bool   `json:"enable" form:"enable"`
	Email      string `json:"email" form:"email" gorm:"unique"`
	Up         int64  `json:"up" form:"up"`
	Down       int64  `json:"down" form:"down"`
	ExpiryTime int64  `json:"expiryTime" form:"expiryTime"`
	Total      int64  `json:"total" form:"total"`
}

type NewTraffic struct {
	Id         int           `json:"id" form:"id" gorm:"primaryKey;autoIncrement"`
	Inbound    model.Inbound `json:"inboundId" form:"inboundId"`
	Enable     bool          `json:"enable" form:"enable"`
	Client     model.Client  `json:"email" form:"email" gorm:"unique"`
	Up         int64         `json:"up" form:"up"`
	Down       int64         `json:"down" form:"down"`
	ExpiryTime int64         `json:"expiryTime" form:"expiryTime"`
	Total      int64         `json:"total" form:"total"`
}
