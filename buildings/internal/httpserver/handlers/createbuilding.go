package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrvin/tasks-go/buildings/internal/storage"
)

type BuildingCreator interface {
	CreateBuilding(ctx context.Context, building *storage.Building) error
}

// NewCreateBuilding create new building.
//
//	@Summary      Create building
//	@Description  create new building
//	@Tags         buildings
//	@Accept       json
//	@Produce      json
//	@Param        input body storage.Building true "building data"
//	@Success      201  {string} string "OK"
//	@Router       /buildings [post]
func NewCreateBuilding(creator BuildingCreator) gin.HandlerFunc {
	return func(c *gin.Context) {
		var building storage.Building
		if err := c.ShouldBindJSON(&building); err != nil {
			log.Println("Error while unmarshalling create building request body: %w", err)
			c.AbortWithError(http.StatusBadRequest, err) //nolint:errcheck
		}
		if err := creator.CreateBuilding(c.Request.Context(), &building); err != nil {
			log.Println("Error while save building to storage: %w", err)
			c.AbortWithError(http.StatusInternalServerError, err) //nolint:errcheck
			return
		}

		c.JSON(http.StatusCreated, gin.H{"status": "OK"})
	}
}
