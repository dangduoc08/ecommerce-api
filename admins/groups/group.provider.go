package groups

import (
	"fmt"

	"github.com/dangduoc08/ecommerce-api/constants"
	"github.com/dangduoc08/ecommerce-api/dbs"
	"github.com/dangduoc08/gogo/core"
	"gorm.io/gorm/clause"
)

type GroupProvider struct {
	dbs.DBProvider
}

func (instance GroupProvider) NewProvider() core.Provider {
	return instance
}

func (instance GroupProvider) FindOneBy(query *Query) (*GroupModel, error) {
	groupRec := &GroupModel{}
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

func (instance GroupProvider) FindManyBy(query *Query) ([]*GroupModel, error) {
	groupRecs := []*GroupModel{}
	groupQueries := map[string]any{}

	tx := instance.DB
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
		return []*GroupModel{}, err
	}

	return groupRecs, nil
}

func (instance GroupProvider) CreateOne(data *Creation) (*GroupModel, error) {
	group := &GroupModel{
		Name:        data.Name,
		Permissions: data.Permissions,
		StoreID:     data.StoreID,
	}

	if err := instance.Create(group).Error; err != nil {
		return nil, err
	}

	return group, nil
}

func (instance GroupProvider) UpdateByID(id uint, data *Update) (*GroupModel, error) {
	groupRec := &GroupModel{
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

func (instance GroupProvider) DeleteByID(id uint) error {
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

	if err := tx.Delete(&GroupModel{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
