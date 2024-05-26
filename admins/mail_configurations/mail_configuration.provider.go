package mailConfigurations

import (
	"fmt"

	"github.com/dangduoc08/ecommerce-api/constants"
	"github.com/dangduoc08/ecommerce-api/dbs"
	"github.com/dangduoc08/gogo/core"
	"gorm.io/gorm/clause"
)

type MailConfigurationProvider struct {
	dbs.DBProvider
}

func (instance MailConfigurationProvider) NewProvider() core.Provider {

	return instance
}

func (instance MailConfigurationProvider) FindByID(
	id uint,
) (*MailConfigurationModel, error) {
	mailRec := &MailConfigurationModel{ID: id}

	if err := instance.Take(mailRec).Error; err != nil {
		return nil, err
	}

	return mailRec, nil
}

func (instance MailConfigurationProvider) FindOneBy(queries ...*Query) (*MailConfigurationModel, error) {
	mailConfigurationRec := &MailConfigurationModel{}
	tx := instance.DB

	for _, query := range queries {
		mailConfigurationQueries := map[string]any{}
		if query.StoreID != 0 {
			mailConfigurationQueries["store_id"] = query.StoreID
		}

		if query.ID != 0 {
			mailConfigurationQueries["id"] = query.ID
		}

		tx = tx.Or(mailConfigurationQueries)
	}

	if err := tx.
		First(mailConfigurationRec).
		Error; err != nil {
		return nil, err
	}

	return mailConfigurationRec, nil
}

func (instance MailConfigurationProvider) FindManyBy(limit int, offset int, queries ...*Query) ([]*MailConfigurationModel, error) {
	mailConfigurationRec := []*MailConfigurationModel{}
	tx := instance.DB

	for _, query := range queries {
		mailConfigurationQueries := map[string]any{}
		if query.StoreID != 0 {
			mailConfigurationQueries["store_id"] = query.StoreID
		}

		if query.ID != 0 {
			mailConfigurationQueries["id"] = query.ID
		}

		tx = tx.Or(mailConfigurationQueries)

		if query.Sort != "" {
			if query.Order == "" {
				query.Order = constants.ASC
			}
			tx = tx.Order(fmt.Sprintf("%v %v", query.Sort, query.Order))
		}
	}

	if err := tx.
		Limit(limit).
		Offset(offset).
		Find(&mailConfigurationRec).
		Error; err != nil {
		return nil, err
	}

	return mailConfigurationRec, nil
}

func (instance MailConfigurationProvider) UpdateByID(id uint, data *Update) (*MailConfigurationModel, error) {
	mailConfigurationRec := &MailConfigurationModel{
		ID:       id,
		Host:     data.Host,
		Port:     data.Port,
		Username: data.Username,
		Password: data.Password,
	}

	if err := instance.
		Clauses(clause.Returning{}).
		Updates(&mailConfigurationRec).
		Error; err != nil {
		return nil, err
	}

	return mailConfigurationRec, nil
}
