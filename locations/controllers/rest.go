package controllers

import (
	"strconv"

	"github.com/dangduoc08/ecommerce-api/locations/models"
	"github.com/dangduoc08/ecommerce-api/locations/providers"
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/core"
)

type REST struct {
	common.REST
	common.Logger
	providers.DBHandler
}

func (instance REST) NewController() core.Controller {
	instance.
		Prefix("v1").
		Prefix("locations")

	return instance
}

func (instance REST) READ(
	ctx gogo.Context,
	query gogo.Query,
) []*models.Location {
	var locationID *uint

	if query.Get("location_id") != "" {
		u64LocationID, err := strconv.ParseUint(query.Get("location_id"), 10, 64)
		if err == nil {
			uintLocationID := uint(u64LocationID)
			locationID = &uintLocationID
		}
	}

	locations, err := instance.FindManyBy(&providers.Query{
		LocationID: locationID,
	})
	if err != nil {
		instance.Logger.Debug(
			"READ.FindManyBy",
			"error", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		return []*models.Location{}
	}

	return locations
}
