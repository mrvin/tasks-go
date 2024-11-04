package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mrvin/tasks-go/buildings/internal/storage"
)

type BuildingLister interface {
	ListBuildings(ctx context.Context, city string, year int16, numFloors int16) ([]storage.Building, error)
}

// NewListBuildings lists all existing buildings.
//
//	@Summary      List buildings
//	@Description  get buildings
//	@Tags         buildings
//	@Produce      json
//	@Param        city query string false "equal city"
//	@Param        year query string false "equal year"
//	@Param        number_floors query string false  "equal number_floors"
//	@Success      200  {array} storage.Building
//	@Router       /buildings [get]
func NewListBuildings(lister BuildingLister) gin.HandlerFunc {
	return func(c *gin.Context) {
		city := c.Query("city")
		strYear := c.Query("year")
		strNumFloors := c.Query("number_floors")

		var err error
		year := int64(0)
		if strYear != "" {
			year, err = strconv.ParseInt(strYear, 10, 16)
			if err != nil {
				log.Println("Error while convert year to int: %w", err)
				c.AbortWithError(http.StatusBadRequest, err) //nolint:errcheck
				return
			}
		}
		numFloors := int64(0)
		if strNumFloors != "" {
			numFloors, err = strconv.ParseInt(strNumFloors, 10, 16)
			if err != nil {
				log.Println("Error while convert number floors to int: %w", err)
				c.AbortWithError(http.StatusBadRequest, err) //nolint:errcheck
				return
			}
		}

		buildings, err := lister.ListBuildings(c.Request.Context(), city, int16(year), int16(numFloors)) //nolint:gosec
		if err != nil {
			log.Println("Error while getting list of buildings from storage: %w", err)
			c.AbortWithError(http.StatusInternalServerError, err) //nolint:errcheck
			return
		}

		c.JSON(http.StatusOK, buildings)
	}
}
