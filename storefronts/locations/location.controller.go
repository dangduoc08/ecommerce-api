package locations

import (
	"strconv"

	"github.com/dangduoc08/ecommerce-api/admins/locations"
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/core"
)

type LocationController struct {
	common.REST
	common.Logger
	LocationProvider locations.LocationProvider
}

func (instance LocationController) NewController() core.Controller {
	return instance
}

func (instance LocationController) READ_VERSION_1(
	ctx gogo.Context,
	query gogo.Query,
) []*locations.LocationModel {
	var locationID *uint

	if query.Get("location_id") != "" {
		u64LocationID, err := strconv.ParseUint(query.Get("location_id"), 10, 64)
		if err == nil {
			uintLocationID := uint(u64LocationID)
			locationID = &uintLocationID
		}
	}

	locationRecs, err := instance.LocationProvider.FindManyBy(&locations.Query{
		LocationID: locationID,
	})
	if err != nil {
		instance.Logger.Debug(
			"READ.FindManyBy",
			"error", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		return []*locations.LocationModel{}
	}

	return locationRecs
}
