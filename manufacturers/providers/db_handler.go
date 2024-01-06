package providers

import (
	"fmt"

	"github.com/dangduoc08/ecommerce-api/constants"
	dbProviders "github.com/dangduoc08/ecommerce-api/db/providers"
	"github.com/dangduoc08/ecommerce-api/manufacturers/models"
	"github.com/dangduoc08/gooh/core"
	"gorm.io/gorm/clause"
)

type DBHandler struct {
	dbProviders.DB
}

func (self DBHandler) NewProvider() core.Provider {
	return self
}

func (self DBHandler) FindByID(id uint) (*models.Manufacturer, error) {
	manufacturerRec := &models.Manufacturer{
		ID: id,
	}

	if err := self.
		First(manufacturerRec).
		Error; err != nil {
		return nil, err
	}

	return manufacturerRec, nil
}

func (self DBHandler) FindOneBy(query *Query) (*models.Manufacturer, error) {
	manufacturerRec := &models.Manufacturer{}
	manufacturerQueries := map[string]any{}

	if query.ID != 0 {
		manufacturerQueries["id"] = query.ID
	}

	if query.StoreID != 0 {
		manufacturerQueries["store_id"] = query.StoreID
	}

	if err := self.
		Where(manufacturerQueries).
		First(manufacturerRec).
		Error; err != nil {
		return nil, err
	}

	return manufacturerRec, nil
}

func (self DBHandler) FindManyBy(query *Query) ([]*models.Manufacturer, error) {
	manufacturerRecs := []*models.Manufacturer{}
	manufacturerQueries := map[string]any{}
	tx := self.DB.DB

	if query.StoreID != 0 {
		manufacturerQueries["store_id"] = query.StoreID
	}

	if query.Sort != "" {
		if query.Order == "" {
			query.Order = constants.ASC
		}
		tx = tx.Order(fmt.Sprintf("%v %v", query.Sort, query.Order))
	}

	if err := tx.
		Limit(query.Limit).
		Offset(query.Offset).
		Where(manufacturerQueries).
		Find(&manufacturerRecs).
		Error; err != nil {
		return []*models.Manufacturer{}, err
	}

	return manufacturerRecs, nil
}

func (self DBHandler) CreateOne(data *Creation) (*models.Manufacturer, error) {
	manufacturerRec := &models.Manufacturer{
		Name:    data.Name,
		Logo:    data.Logo,
		StoreID: data.StoreID,
		Slug:    data.Slug,
	}

	if err := self.
		Create(&manufacturerRec).
		Error; err != nil {
		return nil, err
	}

	return manufacturerRec, nil
}

func (self DBHandler) UpdateByID(id uint, data *Update) (*models.Manufacturer, error) {
	manufacturerRec := &models.Manufacturer{
		ID:   id,
		Name: data.Name,
		Logo: data.Logo,
		Slug: data.Slug,
	}

	if err := self.
		Clauses(clause.Returning{}).
		Updates(&manufacturerRec).
		Error; err != nil {
		return nil, err
	}

	return manufacturerRec, nil
}
