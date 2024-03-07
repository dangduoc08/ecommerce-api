package providers

import (
	"fmt"

	"github.com/dangduoc08/ecommerce-api/constants"
	dbProviders "github.com/dangduoc08/ecommerce-api/dbs/providers"
	"github.com/dangduoc08/ecommerce-api/groups/models"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/core"
	"github.com/dangduoc08/gogo/modules/config"
	"gorm.io/gorm/clause"
)

type DBHandler struct {
	dbProviders.DB
	config.ConfigService
	common.Logger
}

func (instance DBHandler) NewProvider() core.Provider {
	return instance
}

func (instance DBHandler) FindOneBy(query *Query) (*models.Group, error) {
	groupRec := &models.Group{}
	groupQueries := map[string]any{}

	if query.ID != 0 {
		groupQueries["id"] = query.ID
	}

	if query.StoreID != 0 {
		groupQueries["store_id"] = query.StoreID
	}

	if err := instance.
		Where(groupQueries).
		First(groupRec).
		Error; err != nil {
		return nil, err
	}

	return groupRec, nil
}

func (instance DBHandler) FindManyBy(query *Query) ([]*models.Group, error) {
	groupRecs := []*models.Group{}
	groupQueries := map[string]any{}

	tx := instance.DB.DB
	if query != nil {
		if len(query.IDs) > 0 {
			tx = tx.Where("id IN ?", query.IDs)
		}
	}

	if query.StoreID != 0 {
		groupQueries["store_id"] = query.StoreID
	}

	if query.Order == "" {
		query.Order = constants.ASC
	}

	sort := ""
	if query.Sort != "" {
		sort = fmt.Sprintf("%v %v", query.Sort, query.Order)
	}

	if err := tx.
		Order(sort).
		Limit(query.Limit).
		Offset(query.Offset).
		Where(groupQueries).
		Find(&groupRecs).
		Error; err != nil {
		return []*models.Group{}, err
	}

	return groupRecs, nil
}

func (instance DBHandler) CreateOne(data *Creation) (*models.Group, error) {
	group := &models.Group{
		Name:        data.Name,
		Permissions: data.Permissions,
		StoreID:     data.StoreID,
	}

	if err := instance.Create(group).Error; err != nil {
		return nil, err
	}

	return group, nil
}

func (instance DBHandler) UpdateByID(id uint, data *Update) (*models.Group, error) {
	groupRec := &models.Group{
		ID:          id,
		StoreID:     data.StoreID,
		Name:        data.Name,
		Permissions: data.Permissions,
	}

	if err := instance.
		Clauses(clause.Returning{}).
		Updates(&groupRec).
		Error; err != nil {
		return nil, err
	}

	return groupRec, nil
}

func (instance DBHandler) DeleteByID(id uint) error {
	tx := instance.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Exec("DELETE FROM users_groups WHERE group_id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&models.Group{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
