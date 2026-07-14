package backup

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Ccccraz/cogmoteGO/internal/commonTypes"
	"github.com/Ccccraz/cogmoteGO/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createRequest struct {
	Source      Source      `json:"source" binding:"required"`
	Destination Destination `json:"destination" binding:"required"`
}

func RegisterRoutes(r gin.IRouter, sourceRoots, sambaRoots []config.BackupRoot) error {
	service, err := newService(sourceRoots, sambaRoots)
	if err != nil {
		return fmt.Errorf("create backup service: %w", err)
	}
	r.GET("/backups", func(c *gin.Context) { currentHandler(c, service) })
	r.POST("/backups", func(c *gin.Context) { createHandler(c, service) })
	return nil
}

func createHandler(c *gin.Context, service BackupService) {
	var request createRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commonTypes.APIError{Error: "invalid backup request", Detail: err.Error()})
		return
	}
	task, err := service.Create(uuid.NewString(), request.Source, request.Destination)
	if err != nil {
		if errors.Is(err, errBackupRunning) {
			c.JSON(http.StatusConflict, commonTypes.APIError{Error: "backup already running", Detail: err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, commonTypes.APIError{Error: "failed to create backup task", Detail: err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, task)
}

func currentHandler(c *gin.Context, service BackupService) {
	task, ok := service.Current()
	if !ok {
		c.JSON(http.StatusNotFound, commonTypes.APIError{Error: "backup task not found", Detail: ""})
		return
	}
	c.JSON(http.StatusOK, task)
}
