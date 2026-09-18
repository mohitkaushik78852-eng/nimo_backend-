package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type claimDeviceRequest struct {
	DeviceID    string `json:"device_id" binding:"required"`
	HouseholdID string `json:"household_id" binding:"required"`
}

// ClaimDevice binds a physical device to a household once Wi-Fi/Bluetooth
// setup succeeds on the app side. Persist this to a real Device table once
// devices move off the in-memory mock data in this file.
func ClaimDevice(c *gin.Context) {
	var req claimDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       "claimed",
		"device_id":    req.DeviceID,
		"household_id": req.HouseholdID,
	})
}
