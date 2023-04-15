package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
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

type Inbound struct {
	ID         uuid.UUID `gorm:"type:uuid" json:"id" gorm:"primaryKey"`
	UserId     int       `json:"-"`
	Up         int64     `json:"up" form:"up"`
	Down       int64     `json:"down" form:"down"`
	Total      int64     `json:"total" form:"total"`
	Remark     string    `json:"remark" form:"remark"`
	Enable     bool      `json:"enable" form:"enable"`
	ExpiryTime int64     `json:"expiryTime" form:"expiryTime"`
	//ClientStats []ClientTraffic `gorm:"foreignKey:InboundId;references:Id" json:"clientStats" form:"clientStats"`
	Clients []Client `gorm:"foreignKey:InboundID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	// config part
	Listen         string   `json:"listen" form:"listen"`
	Port           int      `json:"port" form:"port" gorm:"unique"`
	Protocol       Protocol `json:"protocol" form:"protocol"`
	Settings       string   `json:"settings" form:"settings"`
	StreamSettings string   `json:"streamSettings" form:"streamSettings"`
	Tag            string   `json:"tag" form:"tag" gorm:"unique"`
	Sniffing       string   `json:"sniffing" form:"sniffing"`
}

func (i *Inbound) BeforeCreate(tx *gorm.DB) (err error) {
	// UUID version 4
	i.ID = uuid.NewV4()
	return
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

type xrayClient struct {
	AlterIds   int       `json:"AlterIds"`
	Enable     bool      `json:"Enable"`
	ID         uuid.UUID `gorm:"type:uuid;" json:"id" gorm:"primaryKey"`
	Email      string    `json:"Email"`
	LimitIP    int       `json:"LimitIP"`
	Security   string    `json:"Security"`
	TotalGB    int       `json:"TotalGB"`
	ExpiryTime int64     `json:"ExpiryTime"`
}

func (i *Inbound) GenXrayInboundConfig() *InboundConfig {
	listen := i.Listen

	var smallClient []xrayClient
	if listen != "" {
		listen = fmt.Sprintf("\"%v\"", listen)
	}
	junkClient, err := json.Marshal(i.Clients)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(junkClient, &smallClient)
	if err != nil {
		return nil
	}
	settings := map[string]interface{}{}
	err = json.Unmarshal([]byte(i.Settings), &settings)
	if err != nil {
		return nil
	}
	settings["clients"] = smallClient
	modifiedSettings, err := json.Marshal(settings)
	if err != nil {
		return nil
	}
	return &InboundConfig{
		Listen:         json_util.RawMessage(listen),
		Port:           i.Port,
		Protocol:       string(i.Protocol),
		Settings:       modifiedSettings,
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
	ID         uuid.UUID `gorm:"type:uuid;" json:"id" gorm:"primaryKey"`
	InboundID  uuid.UUID `gorm:"type:uuid;" json:"inbound_id"`
	Inbound    Inbound   `json:"-"`
	Creator    int       `json:"-"`
	AlterIds   uint16    `json:"alterId"`
	Enable     bool      `json:"enable" form:"enable"`
	LimitIP    int       `json:"limitIp" form:"limitIp"`
	Email      string    `json:"email" form:"email"`
	Security   string    `json:"security" form:"security"`
	TotalGB    int64     `json:"totalGB" form:"totalGB"`
	TotalUp    int64     `json:"totalUp" form:"totalUp"`
	TotalDown  int64     `json:"totalDown" form:"totalDown"`
	ExpiryTime int64     `json:"expiryTime" form:"expiryTime"`
}

func (c *Client) BeforeCreate(tx *gorm.DB) (err error) {
	// UUID version 4
	c.ID = uuid.NewV4()
	c.Email = fmt.Sprintf("%s@%s", c.ID, c.InboundID)
	return
}

type Traffic struct {
	ID        uuid.UUID `gorm:"type:uuid;" json:"id" gorm:"primaryKey"`
	ClientsID uuid.UUID `json:"inboundClientID" form:"inboundClientID"`
	Up        int64     `json:"up" form:"up"`
	Down      int64     `json:"down" form:"down"`
	CreatedAt time.Time
}

type InboundClientIps struct {
	Id          int    `json:"id" gorm:"primaryKey;autoIncrement"`
	ClientEmail string `json:"clientEmail" form:"clientEmail" gorm:"unique"`
	Ips         string `json:"ips" form:"ips"`
}
