package handlers

import (
	"errors"
	"net/http"

	"logistics-tracker/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateWaybill(c *gin.Context) {
	var req models.WaybillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if ok, msg := ValidateTrackingNumber(req.TrackingNumber); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}
	if ok, msg := ValidateCarrier(req.Carrier); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	var existing models.Waybill
	err := models.DB.Where("tracking_number = ?", req.TrackingNumber).First(&existing).Error
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Waybill already exists"})
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	waybill := models.Waybill{
		TrackingNumber: req.TrackingNumber,
		Carrier:        req.Carrier,
		Status:         models.WaybillStatusPending,
	}

	if err := models.DB.Create(&waybill).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create waybill"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":              waybill.ID,
		"tracking_number": waybill.TrackingNumber,
		"carrier":         waybill.Carrier,
		"status":          waybill.Status,
		"created_at":      waybill.CreatedAt,
	})
}

func ListWaybills(c *gin.Context) {
	status := c.Query("status")

	query := models.DB.Model(&models.Waybill{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var waybills []models.Waybill
	if err := query.Order("created_at DESC").Find(&waybills).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query waybills"})
		return
	}

	results := make([]gin.H, 0, len(waybills))
	for _, w := range waybills {
		results = append(results, gin.H{
			"id":              w.ID,
			"tracking_number": w.TrackingNumber,
			"carrier":         w.Carrier,
			"status":          w.Status,
			"created_at":      w.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"total":  len(results),
		"status": status,
		"data":   results,
	})
}
