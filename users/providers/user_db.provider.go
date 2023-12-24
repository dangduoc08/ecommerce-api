package providers

import (
	"fmt"

	authProviders "github.com/dangduoc08/ecommerce-api/auths/providers"
	"github.com/dangduoc08/ecommerce-api/constants"
	dbProviders "github.com/dangduoc08/ecommerce-api/db/providers"
	groupModels "github.com/dangduoc08/ecommerce-api/groups/models"
	groupProviders "github.com/dangduoc08/ecommerce-api/groups/providers"
	"github.com/dangduoc08/ecommerce-api/users/models"
	"github.com/dangduoc08/ecommerce-api/utils"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/core"
	"github.com/dangduoc08/gooh/modules/config"
)

type UserDB struct {
	DBProvider    dbProviders.DB
	ConfigService config.ConfigService
	Logger        common.Logger
	GroupDB       groupProviders.GroupDB
	AuthHandler   authProviders.AuthHandler
}

func (self UserDB) NewProvider() core.Provider {
	return self
}

func (self UserDB) IsDuplicated(data []map[string]string) bool {
	for _, kv := range data {
		for k, v := range kv {
			var userRec models.User
			self.DBProvider.DB.Where(fmt.Sprintf("%v = ?", k), fmt.Sprintf("%v", v)).First(&userRec)
			if userRec.ID != 0 {
				return true
			}
		}
	}

	return false
}

func (self UserDB) FindByID(id uint) (*models.User, error) {
	userRec := &models.User{
		ID: id,
	}

	if err := self.DBProvider.DB.
		Preload("Groups").
		First(userRec).Error; err != nil {
		return nil, err
	}

	return userRec, nil
}

func (self UserDB) FindOneBy(query *UserQuery) (*models.User, error) {
	userRec := &models.User{}
	userQueries := map[string]any{}

	if query.Username != "" {
		userQueries["username"] = query.Username
	}

	if query.Email != "" {
		userQueries["email"] = query.Username
	}

	if err := self.DBProvider.DB.
		Where(userQueries).
		Preload("Groups").
		First(userRec).
		Error; err != nil {
		return nil, err
	}

	return userRec, nil
}

func (self UserDB) FindManyBy(query *UserQuery) ([]*models.User, error) {
	userRecs := []*models.User{}
	userQueries := map[string]any{}

	tx := self.DBProvider.DB
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
		return []*models.User{}, err
	}

	return userRecs, nil
}

func (self UserDB) CreateOne(data *UserCreation) (*models.User, error) {

	hash, err := self.AuthHandler.HashPassword(data.Password)
	if err != nil {
		return nil, err
	}

	// check whether group ids
	// were existed
	groups := []*groupModels.Group{}
	userRec := &models.User{
		StoreID:   data.StoreID,
		Username:  data.Username,
		Email:     data.Email,
		Hash:      hash,
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Groups:    groups,
	}
	if len(data.GroupIDs) > 0 {
		groups, err := self.GroupDB.FindManyBy(&groupProviders.GroupQuery{
			IDs:     data.GroupIDs,
			StoreID: data.StoreID,
			Limit:   len(data.GroupIDs),
		})

		if err != nil {
			return nil, err
		}

		storedGroupIDs := utils.ArrMap(groups, func(gr *groupModels.Group, idx int) uint {
			return gr.ID
		})

		for index, id := range data.GroupIDs {
			if !utils.ArrIncludes(storedGroupIDs, id) {
				return nil, fmt.Errorf("%v[%v]=%v doesn't exist", "group_ids", index, id)
			}
		}

		userRec.Groups = groups
	}

	if err := self.DBProvider.DB.Create(userRec).Error; err != nil {
		return nil, err
	}

	return userRec, nil
}

func (self UserDB) ModifyOne(user *models.User) (*models.User, error) {
	return user, self.DBProvider.DB.Save(user).Error
}
