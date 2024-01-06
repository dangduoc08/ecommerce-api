package controllers

import (
	"strconv"

	"github.com/dangduoc08/ecommerce-api/locations/models"
	"github.com/dangduoc08/ecommerce-api/locations/providers"
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/core"
)

type REST struct {
	common.REST
	common.Logger
	providers.DBHandler
}

func (self REST) NewController() core.Controller {
	self.
		Prefix("v1").
		Prefix("locations")

	return self
}

func (self REST) READ(
	ctx gooh.Context,
	query gooh.Query,
) []*models.Location {
	var locationID *uint

	if query.Get("location_id") != "" {
		u64LocationID, err := strconv.ParseUint(query.Get("location_id"), 10, 64)
		if err == nil {
			uintLocationID := uint(u64LocationID)
			locationID = &uintLocationID
		}
	}

	locations, err := self.FindManyBy(&providers.Query{
		LocationID: locationID,
	})
	if err != nil {
		self.Logger.Debug(
			"READ.FindManyBy",
			"error", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		return []*models.Location{}
	}

	return locations
}
