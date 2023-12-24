package controllers

import (
	"strconv"

	"github.com/dangduoc08/ecommerce-api/locations/models"
	"github.com/dangduoc08/ecommerce-api/locations/providers"
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/core"
)

type LocationREST struct {
	common.REST
	Logger     common.Logger
	LocationDB providers.LocationDB
}

func (self LocationREST) NewController() core.Controller {
	self.
		Prefix("v1").
		Prefix("locations")

	return self
}

func (self LocationREST) READ(
	c gooh.Context,
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

	locations, err := self.LocationDB.FindManyBy(&providers.LocationQuery{
		LocationID: locationID,
	})
	if err != nil {
		self.Logger.Debug(
			"LocationREST.READ.LocationDB.FindManyBy",
			"error", err.Error(),
			"X-Request-ID", c.GetID(),
		)
		return []*models.Location{}
	}

	return locations
}
