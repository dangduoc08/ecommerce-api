package users

import (
	"fmt"

	"github.com/dangduoc08/ecommerce-api/constants"
	"github.com/dangduoc08/ecommerce-api/dbs"

	"github.com/dangduoc08/gogo/core"
)

type UserProvider struct {
	dbs.DBProvider
}

func (instance UserProvider) NewProvider() core.Provider {
	return instance
}

func (instance UserProvider) FindByID(id uint) (*UserModel, error) {
	userRec := &UserModel{
		ID: id,
	}

	if err := instance.DB.
		Preload("Groups").
		First(userRec).Error; err != nil {
		return nil, err
	}

	return userRec, nil
}

func (instance UserProvider) FindOneBy(query *Query) (*UserModel, error) {
	userRec := &UserModel{}
	userQueries := map[string]any{}

	if query.Username != "" {
		userQueries["username"] = query.Username
	}

	if query.Email != "" {
		userQueries["email"] = query.Username
	}

	if err := instance.DB.
		Where(userQueries).
		Preload("Groups").
		First(userRec).
		Error; err != nil {
		return nil, err
	}

	return userRec, nil
}

func (instance UserProvider) FindManyBy(query *Query) ([]*UserModel, error) {
	userRecs := []*UserModel{}
	userQueries := map[string]any{}

	tx := instance.DB
	if query != nil {
		if len(query.Statuses) > 0 {
			tx = tx.Where("status IN ?", query.Statuses)
		}
	}

	if query.Username != "" {
		userQueries["username"] = query.Username
	}

	if query.Email != "" {
		userQueries["email"] = query.Email
	}

	if query.StoreID != 0 {
		userQueries["store_id"] = query.StoreID
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
		Where(userQueries).
		Find(&userRecs).
		Error; err != nil {
		return []*UserModel{}, err
	}

	return userRecs, nil
}

func (instance UserProvider) CreateOne(data *Creation) (*UserModel, error) {
	return nil, nil
	// hash, err := instance.HashPassword(data.Password)
	// if err != nil {
	// 	return nil, err
	// }

	// // check whether group ids
	// // were existed
	// groups := []*groupModels.Group{}
	// userRec := &UserModel{
	// 	StoreID:   data.StoreID,
	// 	Username:  data.Username,
	// 	Email:     data.Email,
	// 	Hash:      hash,
	// 	FirstName: data.FirstName,
	// 	LastName:  data.LastName,
	// 	Groups:    groups,
	// }
	// if len(data.GroupIDs) > 0 {
	// 	groups, err := instance.GroupDBHandler.FindManyBy(&groupProviders.Query{
	// 		IDs:     data.GroupIDs,
	// 		StoreID: data.StoreID,
	// 		Limit:   len(data.GroupIDs),
	// 	})

	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	storedGroupIDs := utils.ArrMap(groups, func(gr *groupModels.Group, idx int) uint {
	// 		return gr.ID
	// 	})

	// 	for index, id := range data.GroupIDs {
	// 		if !utils.ArrIncludes(storedGroupIDs, id) {
	// 			return nil, fmt.Errorf("%v[%v]=%v doesn't exist", "group_ids", index, id)
	// 		}
	// 	}

	// 	userRec.Groups = groups
	// }

	// if err := instance.Create(userRec).Error; err != nil {
	// 	return nil, err
	// }

	// return userRec, nil
}

func (instance UserProvider) ModifyOne(user *UserModel) (*UserModel, error) {
	return user, instance.Save(user).Error
}

func (instance UserProvider) CheckDuplicated(data []map[string]string) bool {
	for _, kv := range data {
		for k, v := range kv {
			var userRec UserModel
			instance.Where(fmt.Sprintf("%v = ?", k), fmt.Sprintf("%v", v)).First(&userRec)
			if userRec.ID != 0 {
				return true
			}
		}
	}

	return false
}
