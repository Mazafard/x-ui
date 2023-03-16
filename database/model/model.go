package model

import (
	"bytes"
	"fmt"
	"github.com/xtls/xray-core/common/uuid"
	"x-ui/util/json_util"
)

type Protocol string

const (
	VMess       Protocol = "vmess"
	VLESS       Protocol = "vless"
	Dokodemo    Protocol = "Dokodemo-door"
	Http        Protocol = "http"
	Trojan      Protocol = "trojan"
	Shadowsocks Protocol = "shadowsocks"
)

type User struct {
	Id       int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Username string `json:"username"`
	Password string `json:"password"`
}
type ClientTraffic struct {
	Id        int    `json:"id" form:"id" gorm:"primaryKey;autoIncrement"`
	InboundId int    `json:"inboundId" form:"enable"`
	Enable    bool   `json:"enable" form:"enable"`
	Email     string `json:"email" form:"email" gorm:"unique"`
	Up        int64  `json:"up" form:"up"`
	Down      int64  `json:"down" form:"down"`
}
type Inbound struct {
	Id          int             `json:"id" form:"id" gorm:"primaryKey;autoIncrement"`
	UserId      int             `json:"-"`
	Up          int64           `json:"up" form:"up"`
	Down        int64           `json:"down" form:"down"`
	Total       int64           `json:"total" form:"total"`
	Remark      string          `json:"remark" form:"remark"`
	Enable      bool            `json:"enable" form:"enable"`
	ExpiryTime  int64           `json:"expiryTime" form:"expiryTime"`
	ClientStats []ClientTraffic `gorm:"foreignKey:InboundId;references:Id" json:"clientStats" form:"clientStats"`
	Clients     []*Client       `gorm:"many2many:inbound_client;"`

	// config part
	Listen         string   `json:"listen" form:"listen"`
	Port           int      `json:"port" form:"port" gorm:"unique"`
	Protocol       Protocol `json:"protocol" form:"protocol"`
	Settings       string   `json:"settings" form:"settings"`
	StreamSettings string   `json:"streamSettings" form:"streamSettings"`
	Tag            string   `json:"tag" form:"tag" gorm:"unique"`
	Sniffing       string   `json:"sniffing" form:"sniffing"`
}

type InboundConfig struct {
	Listen         json_util.RawMessage `json:"listen"` // listen 不能为空字符串
	Port           int                  `json:"port"`
	Protocol       string               `json:"protocol"`
	Settings       json_util.RawMessage `json:"settings"`
	StreamSettings json_util.RawMessage `json:"streamSettings"`
	Tag            string               `json:"tag"`
	Sniffing       json_util.RawMessage `json:"sniffing"`
}

func (c *InboundConfig) Equals(other *InboundConfig) bool {
	if !bytes.Equal(c.Listen, other.Listen) {
		return false
	}
	if c.Port != other.Port {
		return false
	}
	if c.Protocol != other.Protocol {
		return false
	}
	if !bytes.Equal(c.Settings, other.Settings) {
		return false
	}
	if !bytes.Equal(c.StreamSettings, other.StreamSettings) {
		return false
	}
	if c.Tag != other.Tag {
		return false
	}
	if !bytes.Equal(c.Sniffing, other.Sniffing) {
		return false
	}
	return true
}

func (i *Inbound) GenXrayInboundConfig() *InboundConfig {
	listen := i.Listen
	if listen != "" {
		listen = fmt.Sprintf("\"%v\"", listen)
	}
	return &InboundConfig{
		Listen:         json_util.RawMessage(listen),
		Port:           i.Port,
		Protocol:       string(i.Protocol),
		Settings:       json_util.RawMessage(i.Settings),
		StreamSettings: json_util.RawMessage(i.StreamSettings),
		Tag:            i.Tag,
		Sniffing:       json_util.RawMessage(i.Sniffing),
	}
}

type Setting struct {
	Id    int    `json:"id" form:"id" gorm:"primaryKey;autoIncrement"`
	Key   string `json:"key" form:"key"`
	Value string `json:"value" form:"value"`
}
type Client struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4" json:"id" gorm:"primaryKey"`
	Creator    int        `json:"-"`
	AlterIds   uint16     `json:"alterId"`
	Enable     bool       `json:"enable" form:"enable"`
	Email      string     `json:"email"`
	LimitIP    int        `json:"limitIp"`
	Security   string     `json:"security"`
	TotalGB    int64      `json:"totalGB" form:"totalGB"`
	ExpiryTime int64      `json:"expiryTime" form:"expiryTime"`
	Inbound    []*Inbound `gorm:"many2many:inbound_client;"`
}

type InboundClient struct {
	ClientID      int `gorm:"primaryKey"`
	InboundID     int `gorm:"primaryKey"`
	ClientTraffic []ClientTraffic
}

type InboundClientIps struct {
	Id          int    `json:"id" gorm:"primaryKey;autoIncrement"`
	ClientEmail string `json:"clientEmail" form:"clientEmail" gorm:"unique"`
	Ips         string `json:"ips" form:"ips"`
}
