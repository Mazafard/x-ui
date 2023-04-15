package service

import (
	"errors"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"x-ui/database"
	"x-ui/database/model"
)

type ClientService struct {
}

func (s *ClientService) GetClients(userId int) ([]*model.Client, error) {
	db := database.GetDB()
	var clients []*model.Client
	err := db.Debug().Preload("Inbound").Find(&clients).Where("creator = ?", userId).First(&clients).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return clients, nil
}
func (s *ClientService) GetClient(userId int, clientId uuid.UUID) (client *model.Client, err error) {
	db := database.GetDB()
	err = db.Debug().Model(model.Client{}).Where("creator = ? AND id = ?", userId, clientId).First(&client).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return client, nil
}

func (s *ClientService) AddClient(client model.Client) (model.Client, error) {
	var err error
	db := database.GetDB()

	err = db.Save(&client).Error
	//if err == nil {
	//	s.UpdateClientStat(client.Id, inbound.Settings)
	//}
	return client, err
}

func (s *ClientService) ChangeStatusClient(client *model.Client, status bool) (*model.Client, error) {
	var err error
	db := database.GetDB()

	err = db.Model(&client).Update("enable", status).Error
	//if err == nil {
	//	s.UpdateClientStat(client.Id, inbound.Settings)
	//}
	return client, err
}

func (s *ClientService) DelClient(clientId uuid.UUID) (err error) {
	db := database.GetDB()
	return db.Delete(model.Client{}, clientId).Error
}
