package service

import (
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
	"x-ui/database"
	"x-ui/database/model"
	"x-ui/logger"
	"x-ui/util/common"
	"x-ui/xray"
)

type InboundService struct {
}

func (s *InboundService) GetInbounds(userId int) ([]*model.Inbound, error) {
	db := database.GetDB()
	var inbounds []*model.Inbound
	err := db.Model(model.Inbound{}).Where("user_id = ?", userId).Find(&inbounds).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return inbounds, nil
}

func (s *InboundService) GetAllInbounds() ([]*model.Inbound, error) {
	db := database.GetDB()
	var inbounds []*model.Inbound
	err := db.Model(model.Inbound{}).Preload("Clients").Find(&inbounds).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return inbounds, nil
}

func (s *InboundService) GetAllEnabledInbounds() ([]*model.Inbound, error) {
	db := database.GetDB()
	var inbounds []*model.Inbound
	err := db.Model(model.Inbound{}).Preload("Clients").Find(&inbounds).Where("Enable = ?", true).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return inbounds, nil
}

func (s *InboundService) checkPortExist(port int, ignoreId uuid.UUID) (bool, error) {
	var err error
	db := database.GetDB()
	db = db.Model(model.Inbound{}).Where("port = ?", port)
	if ignoreId != uuid.Nil {
		uuid.Must(ignoreId, err)
		if err != nil {
			return false, err
		}
		db = db.Where("id != ?", ignoreId)
	}
	var count int64
	err = db.Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *InboundService) getClients(inbound *model.Inbound) ([]model.Client, error) {
	settings := map[string][]model.Client{}
	json.Unmarshal([]byte(inbound.Settings), &settings)
	if settings == nil {
		return nil, fmt.Errorf("Setting is null")
	}

	clients := settings["clients"]
	if clients == nil {
		return nil, nil
	}
	return clients, nil
}

//
//func (s *InboundService) checkEmailsExist(emails map[string]bool, ignoreId uuid.UUID) (string, error) {
//	db := database.GetDB()
//	var inbounds []*model.Inbound
//	db = db.Model(model.Inbound{}).Preload("Clients").Where("Protocol in ?", []model.Protocol{model.VMess, model.VLESS})
//	//if ignoreId > 0 {
//	db = db.Where("id != ?", ignoreId)
//	//}
//	db = db.Find(&inbounds)
//	if db.Error != nil {
//		return "", db.Error
//	}
//
//	for _, inbound := range inbounds {
//		clients, err := s.getClients(inbound)
//		if err != nil {
//			return "", err
//		}
//
//		for _, client := range clients {
//			if emails[client.Email] {
//				return client.Email, nil
//			}
//		}
//	}
//	return "", nil
//}
//
//func (s *InboundService) checkEmailExistForInbound(inbound *model.Inbound) (string, error) {
//	clients, err := s.getClients(inbound)
//	if err != nil {
//		return "", err
//	}
//	emails := make(map[string]bool)
//	for _, client := range clients {
//		if client.Email != "" {
//			if emails[client.Email] {
//				return client.Email, nil
//			}
//			emails[client.Email] = true
//		}
//	}
//	return s.checkEmailsExist(emails, inbound.ID)
//}

func (s *InboundService) AddInbound(inbound *model.Inbound) (*model.Inbound, error) {
	exist, err := s.checkPortExist(inbound.Port, uuid.Nil)
	if err != nil {
		return inbound, err
	}
	if exist {
		return inbound, common.NewError("Port already exists:", inbound.Port)
	}
	db := database.GetDB()

	err = db.Save(inbound).Error
	//if err == nil {
	//	s.UpdateClientStat(inbound.ID, inbound.Settings)
	//}
	return inbound, err
}

func (s *InboundService) AddInbounds(inbounds []*model.Inbound) error {
	for _, inbound := range inbounds {
		exist, err := s.checkPortExist(inbound.Port, uuid.Nil)
		if err != nil {
			return err
		}
		if exist {
			return common.NewError("Port already exists:", inbound.Port)
		}
	}

	db := database.GetDB()
	tx := db.Begin()
	var err error
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	for _, inbound := range inbounds {
		err = tx.Save(inbound).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *InboundService) DelInbound(id uuid.UUID) error {
	db := database.GetDB()
	return db.Delete(model.Inbound{}, id).Error
}

func (s *InboundService) GetInbound(id uuid.UUID) (*model.Inbound, error) {
	db := database.GetDB()
	inbound := &model.Inbound{}
	err := db.Debug().Model(model.Inbound{}).First(inbound, "id=?", id).Error
	if err != nil {
		return nil, err
	}
	return inbound, nil
}

func (s *InboundService) GetInboundClients(id uuid.UUID) (*[]model.Client, error) {
	db := database.GetDB()
	clients := &[]model.Client{}
	err := db.Debug().Model(model.Client{}).Where("Inbound_id=?", id).Find(clients).Error
	if err != nil {
		return nil, err
	}
	return clients, nil
}

func (s *InboundService) GetInboundsId(ids []uuid.UUID) (inbounds []*model.Inbound, err error) {
	db := database.GetDB()

	//if err := db.Model(model.Inbound{}).Where("id IN (?)", ids).Find(&inbounds).Error; err != nil {
	if err := db.Model(model.Inbound{}).Find(&inbounds, ids).Error; err != nil {
		return nil, err
	}

	return inbounds, nil
}

func (s *InboundService) UpdateInbound(inbound *model.Inbound) (*model.Inbound, error) {
	exist, err := s.checkPortExist(inbound.Port, inbound.ID)
	if err != nil {
		return inbound, err
	}
	if exist {
		return inbound, common.NewError("port already exists:", inbound.Port)
	}

	oldInbound, err := s.GetInbound(inbound.ID)
	if err != nil {
		return inbound, err
	}
	oldInbound.Up = inbound.Up
	oldInbound.Down = inbound.Down
	oldInbound.Total = inbound.Total
	oldInbound.Remark = inbound.Remark
	oldInbound.Enable = inbound.Enable
	oldInbound.ExpiryTime = inbound.ExpiryTime
	oldInbound.Listen = inbound.Listen
	oldInbound.Port = inbound.Port
	oldInbound.Protocol = inbound.Protocol
	oldInbound.Settings = inbound.Settings
	oldInbound.StreamSettings = inbound.StreamSettings
	oldInbound.Sniffing = inbound.Sniffing
	oldInbound.Tag = fmt.Sprintf("inbound-%v", inbound.Port)

	//s.UpdateClientStat(inbound.ID, inbound.Settings)
	db := database.GetDB()
	return inbound, db.Save(oldInbound).Error
}

func (s *InboundService) AddTraffic(traffics []*xray.Traffic) (err error) {
	if len(traffics) == 0 {
		return nil
	}
	db := database.GetDB()
	db = db.Model(model.Inbound{})
	tx := db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	for _, traffic := range traffics {
		if traffic.IsInbound {
			err = tx.Where("tag = ?", traffic.Tag).
				UpdateColumn("up", gorm.Expr("up + ?", traffic.Up)).
				UpdateColumn("down", gorm.Expr("down + ?", traffic.Down)).
				Error
			if err != nil {
				return
			}
		}
	}
	return
}
func (s *InboundService) AddClientTraffic(traffics []*xray.ClientTraffic) (err error) {
	if len(traffics) == 0 {
		return nil
	}
	db := database.GetDB()
	dbInbound := db.Model(model.Inbound{})

	db = db.Model(xray.NewTraffic{})
	tx := db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	txInbound := dbInbound.Begin()
	defer func() {
		if err != nil {
			txInbound.Rollback()
		} else {
			txInbound.Commit()
		}
	}()

	for _, traffic := range traffics {
		inbound := &model.Inbound{}

		err := txInbound.Preload("Clients", "Email = ?", "%"+traffic.Email+"%").First(inbound).Error
		//err := txInbound.Where("settings like ?", "%"+traffic.Email+"%").First(inbound).Error
		traffic.InboundId = inbound.ID
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// delete removed client record
				clientErr := s.DelClientStat(tx, traffic.Email)
				logger.Warning(err, traffic.Email, clientErr)

			}
			continue
		}
		// get settings clients
		//settings := map[string][]model.Client{}
		//json.Unmarshal([]byte(inbound.Settings), &settings)
		//clients := settings["clients"]
		//for _, client := range clients {
		//	if traffic.Email == client.Email {
		//		traffic.ExpiryTime = client.ExpiryTime
		//		traffic.Total = client.TotalGB
		//	}
		//}
		if tx.Where("inbound_id = ?", inbound.ID).Where("email = ?", traffic.Email).
			UpdateColumn("enable", true).
			UpdateColumn("expiry_time", traffic.ExpiryTime).
			UpdateColumn("total", traffic.Total).
			UpdateColumn("up", gorm.Expr("up + ?", traffic.Up)).
			UpdateColumn("down", gorm.Expr("down + ?", traffic.Down)).RowsAffected == 0 {
			err = tx.FirstOrCreate(traffic).Error
		}

		if err != nil {
			logger.Warning("AddClientTraffic update data ", err)
			continue
		}

	}
	return
}

func (s *InboundService) DisableInvalidInbounds() (int64, error) {
	db := database.GetDB()
	now := time.Now().Unix() * 1000
	result := db.Model(model.Inbound{}).
		Where("((total > 0 and up + down >= total) or (expiry_time > 0 and expiry_time <= ?)) and enable = ?", now, true).
		Update("enable", false)
	err := result.Error
	count := result.RowsAffected
	return count, err
}
func (s *InboundService) DisableInvalidClients() (int64, error) {
	db := database.GetDB()
	now := time.Now().Unix() * 1000
	result := db.Debug().Model(model.Client{}).Where("Enable = ?", true).Preload("Inbound", "Total_up > 0 ").
		Where("((total_up > 0 and total_up + total_down >= total_gb) or (expiry_time > 0 and expiry_time <= ?))", now).
		Update("enable", false)
	err := result.Error
	count := result.RowsAffected
	return count, err
}

//func (s *InboundService) UpdateClientStat(inboundId uuid.UUID, inboundSettings string) error {
//	db := database.GetDB()
//
//	// get settings clients
//	settings := map[string][]model.Client{}
//	json.Unmarshal([]byte(inboundSettings), &settings)
//	//var clients []model.Client
//	//db.Model(model.Client{}).Preload("Inbounds").Preload("Traffics").Update()
//	clients := settings["clients"]
//	for _, client := range clients {
//		result := db.Model(xray.ClientTraffic{}).
//			Where("inbound_id = ? and email = ?", inboundId, client.InboundClients).
//			Updates(map[string]interface{}{"enable": true, "total": client.TotalGB, "expiry_time": client.ExpiryTime})
//		if result.RowsAffected == 0 {
//			clientTraffic := xray.ClientTraffic{}
//			clientTraffic.InboundId = inboundId
//			clientTraffic.Email = client.Email
//			clientTraffic.ExpiryTime = client.ExpiryTime
//			clientTraffic.Enable = true
//			clientTraffic.Up = 0
//			clientTraffic.Down = 0
//			db.Create(&clientTraffic)
//		}
//		err := result.Error
//		if err != nil {
//			return err
//		}
//
//	}
//	return nil
//}

func (s *InboundService) DelClientStat(tx *gorm.DB, email string) error {
	return tx.Where("email = ?", email).Delete(xray.ClientTraffic{}).Error
}

func (s *InboundService) GetInboundClientIps(clientEmail string) (string, error) {
	db := database.GetDB()
	InboundClientIps := &model.InboundClientIps{}
	err := db.Model(model.InboundClientIps{}).Where("client_email = ?", clientEmail).First(InboundClientIps).Error
	if err != nil {
		return "", err
	}
	return InboundClientIps.Ips, nil
}
func (s *InboundService) ClearClientIps(clientEmail string) error {
	db := database.GetDB()

	result := db.Model(model.InboundClientIps{}).
		Where("client_email = ?", clientEmail).
		Update("ips", "")
	err := result.Error

	if err != nil {
		return err
	}
	return nil
}
func (s *InboundService) ResetClientTraffic(clientEmail string) error {
	db := database.GetDB()

	result := db.Model(xray.ClientTraffic{}).
		Where("email = ?", clientEmail).
		Update("up", 0).
		Update("down", 0)

	err := result.Error

	if err != nil {
		return err
	}
	return nil
}

//func (s *InboundService) GetClientTrafficById(id uuid.UUID) (traffic *xray.ClientTraffic, err error) {
//	db := database.GetDB()
//	inbound := &model.Inbound{}
//	traffic = &xray.ClientTraffic{}
//
//	err = db.Model(model.Inbound{}).Where("settings like ?", id).First(inbound).Error
//	if err != nil {
//		if err == gorm.ErrRecordNotFound {
//			logger.Warning(err)
//			return nil, err
//		}
//	}
//	traffic.InboundId = inbound.ID
//
//	// get settings clients
//	settings := map[string][]model.Client{}
//	json.Unmarshal([]byte(inbound.Settings), &settings)
//	clients := settings["clients"]
//	for _, client := range clients {
//		if id == client.ID {
//			traffic.Email = client.Email
//		}
//	}
//	err = db.Model(xray.ClientTraffic{}).Where("email = ?", traffic.Email).First(traffic).Error
//	if err != nil {
//		logger.Warning(err)
//		return nil, err
//	}
//	return traffic, err
//}
