package manufacturers

import (
	"fmt"

	"github.com/dangduoc08/ecommerce-api/constants"
	"github.com/dangduoc08/ecommerce-api/dbs"
	"github.com/dangduoc08/gogo/core"
	"gorm.io/gorm/clause"
)

type ManufacturerProvider struct {
	dbs.DBProvider
}

func (instance ManufacturerProvider) NewProvider() core.Provider {
	return instance
}

func (instance ManufacturerProvider) FindByID(id uint) (*ManufacturerModel, error) {
	manufacturerRec := &ManufacturerModel{
		ID: id,
	}

	if err := instance.
		First(manufacturerRec).
		Error; err != nil {
		return nil, err
	}

	return manufacturerRec, nil
}

func (instance ManufacturerProvider) FindOneBy(query *Query) (*ManufacturerModel, error) {
	manufacturerRec := &ManufacturerModel{}
	manufacturerQueries := map[string]any{}

	if query.ID != 0 {
		manufacturerQueries["id"] = query.ID
	}

	if query.StoreID != 0 {
		manufacturerQueries["store_id"] = query.StoreID
	}

	if err := instance.
		Where(manufacturerQueries).
		First(manufacturerRec).
		Error; err != nil {
		return nil, err
	}

	return manufacturerRec, nil
}

func (instance ManufacturerProvider) FindManyBy(query *Query) ([]*ManufacturerModel, error) {
	manufacturerRecs := []*ManufacturerModel{}
	manufacturerQueries := map[string]any{}
	tx := instance.DB

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
		return []*ManufacturerModel{}, err
	}

	return manufacturerRecs, nil
}

func (instance ManufacturerProvider) CreateOne(data *Creation) (*ManufacturerModel, error) {
	manufacturerRec := &ManufacturerModel{
		Name:    data.Name,
		Logo:    data.Logo,
		StoreID: data.StoreID,
		Slug:    data.Slug,
	}

	if err := instance.
		Create(&manufacturerRec).
		Error; err != nil {
		return nil, err
	}

	return manufacturerRec, nil
}

func (instance ManufacturerProvider) UpdateByID(id uint, data *Update) (*ManufacturerModel, error) {
	manufacturerRec := &ManufacturerModel{
		ID:   id,
		Name: data.Name,
		Logo: data.Logo,
		Slug: data.Slug,
	}

	if err := instance.
		Clauses(clause.Returning{}).
		Updates(&manufacturerRec).
		Error; err != nil {
		return nil, err
	}

	return manufacturerRec, nil
}
