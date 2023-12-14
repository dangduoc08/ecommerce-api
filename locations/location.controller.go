package locations

import (
	"strconv"

	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/core"
)

type LocationController struct {
	common.REST
	Logger           common.Logger
	LocationProvider LocationProvider
}

func (self LocationController) NewController() core.Controller {
	self.Prefix("v1").Prefix("locations")

	return self
}

func (self LocationController) READ(query gooh.Query) []*Location {
	var locationID *uint

	if query.Get("location_id") != "" {
		u64LocationID, err := strconv.ParseUint(query.Get("location_id"), 10, 64)
		if err == nil {
			uintLocationID := uint(u64LocationID)
			locationID = &uintLocationID
		}
	}

	locations, err := self.LocationProvider.FindManyBy(locationID)
	if err != nil {
		self.Logger.Debug(
			"LocationProvider.FindManyBy",
			"error", err.Error(),
			"locations", locations,
		)
		return []*Location{}
	}

	return locations
}
