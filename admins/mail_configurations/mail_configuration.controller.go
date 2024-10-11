package mailConfigurations

import (
	"github.com/dangduoc08/ecommerce-api/admins/mail_configurations/dtos"
	sharedLayers "github.com/dangduoc08/ecommerce-api/shared_layers"
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/core"
	"github.com/dangduoc08/gogo/exception"
)

type MailConfigurationController struct {
	common.REST
	common.Guard
	common.Interceptor
	common.Logger
	MailConfigurationProvider
}

func (instance MailConfigurationController) NewController() core.Controller {
	instance.
		BindGuard(
			sharedLayers.AuthGuard{},
		)

	return instance
}

func (instance MailConfigurationController) READ_VERSION_1(
	ctx gogo.Context,
	tokenClaimsDTO sharedLayers.TokenClaimsDTO,
) []*MailConfigurationModel {
	mailConfigurationRecs, err := instance.FindManyBy(1, 0, &Query{
		StoreID: tokenClaimsDTO.StoreID,
	})

	if err != nil {
		instance.Debug(
			"MailConfigurationController.READ_VERSION_1.FindManyBy",
			"error", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		return nil
	}

	return mailConfigurationRecs
}

func (instance MailConfigurationController) UPDATE_BY_id_VERSION_1(
	ctx gogo.Context,
	paramDTO dtos.UPDATE_BY_id_Param_DTO,
	bodyDTO dtos.UPDATE_BY_id_Body_DTO,
	tokenClaimsDTO sharedLayers.TokenClaimsDTO,
) *MailConfigurationModel {
	if _, err := instance.FindOneBy(&Query{
		ID:      paramDTO.ID,
		StoreID: tokenClaimsDTO.StoreID,
	}); err != nil {
		instance.Debug(
			"MailConfigurationController.UPDATE_BY_id_VERSION_1.FindOneBy",
			"error", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.UnprocessableEntityException(err.Error(), exception.ExceptionOptions{
			Cause: err,
		}))
	}

	mailConfigurationRec, err := instance.UpdateByID(paramDTO.ID, &Update{
		Host:     bodyDTO.Data.Host,
		Port:     bodyDTO.Data.Port,
		Username: bodyDTO.Data.Username,
		Password: bodyDTO.Data.Password,
	})

	if err != nil {
		instance.Debug(
			"MailConfigurationController.UPDATE_BY_id_VERSION_1.UpdateByID",
			"error", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.UnprocessableEntityException(err.Error(), exception.ExceptionOptions{
			Cause: err,
		}))
	}

	return mailConfigurationRec
}
