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
	"github.com/dangduoc08/gooh/core"
)

type DBHandler struct {
	dbProviders.DB
	authProviders.Cipher
	GroupDBHandler groupProviders.DBHandler
}

func (self DBHandler) NewProvider() core.Provider {
	return self
}

func (self DBHandler) FindByID(id uint) (*models.User, error) {
	userRec := &models.User{
		ID: id,
	}

	if err := self.DB.
		Preload("Groups").
		First(userRec).Error; err != nil {
		return nil, err
	}

	return userRec, nil
}

func (self DBHandler) FindOneBy(query *Query) (*models.User, error) {
	userRec := &models.User{}
	userQueries := map[string]any{}

	if query.Username != "" {
		userQueries["username"] = query.Username
	}

	if query.Email != "" {
		userQueries["email"] = query.Username
	}

	if err := self.DB.
		Where(userQueries).
		Preload("Groups").
		First(userRec).
		Error; err != nil {
		return nil, err
	}

	return userRec, nil
}

func (self DBHandler) FindManyBy(query *Query) ([]*models.User, error) {
	userRecs := []*models.User{}
	userQueries := map[string]any{}

	tx := self.DB.DB
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

func (self DBHandler) CreateOne(data *Creation) (*models.User, error) {

	hash, err := self.HashPassword(data.Password)
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
		groups, err := self.GroupDBHandler.FindManyBy(&groupProviders.Query{
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

	if err := self.Create(userRec).Error; err != nil {
		return nil, err
	}

	return userRec, nil
}

func (self DBHandler) ModifyOne(user *models.User) (*models.User, error) {
	return user, self.Save(user).Error
}
