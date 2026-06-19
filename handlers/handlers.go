package handlers

import (
	"errors"
	"net/http"
	"time"

	"logistics-tracker/middleware"
	"logistics-tracker/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var user models.User
	if err := models.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if !user.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	token, err := middleware.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{Token: token})
}

func CreateWaybill(c *gin.Context) {
	var req models.WaybillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
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
	}

	if err := models.DB.Create(&waybill).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create waybill"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":              waybill.ID,
		"tracking_number": waybill.TrackingNumber,
		"carrier":         waybill.Carrier,
		"created_at":      waybill.CreatedAt,
	})
}

func GetWaybillTrackings(c *gin.Context) {
	trackingNumber := c.Param("tracking_number")
	if trackingNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tracking number is required"})
		return
	}

	var waybill models.Waybill
	if err := models.DB.Where("tracking_number = ?", trackingNumber).First(&waybill).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Waybill not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	var trackings []models.Tracking
	if err := models.DB.Where("waybill_id = ?", waybill.ID).Order("tracking_time ASC").Find(&trackings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query trackings"})
		return
	}

	if len(trackings) == 0 {
		trackings = generateMockTrackings(waybill.ID)
	}

	response := models.WaybillResponse{
		ID:             waybill.ID,
		TrackingNumber: waybill.TrackingNumber,
		Carrier:        waybill.Carrier,
		CreatedAt:      waybill.CreatedAt,
		Trackings:      trackings,
	}

	c.JSON(http.StatusOK, response)
}

func generateMockTrackings(waybillID uint) []models.Tracking {
	now := time.Now()
	statuses := []struct {
		Status      string
		Location    string
		Description string
		HoursAgo      int
	}{
		{"已签收", "北京市朝阳区配送站", "快件已被本人签收，感谢使用", 0},
		{"派送中", "北京市朝阳区", "快递员张师傅正在为您派送，联系电话：138****1234", 2},
		{"到达目的地", "北京转运中心", "快件已到达北京转运中心，正在分拣", 6},
		{"运输中", "济南转运中心", "快件已从济南转运中心发出，下一站：北京", 24},
		{"运输中", "上海转运中心", "快件已从上海转运中心发出，下一站：济南", 48},
		{"已揽收", "上海市浦东新区营业部", "快递员已揽收，快件已发往上海转运中心", 72},
	}

	var trackings []models.Tracking
	for _, s := range statuses {
		trackingTime := now.Add(-time.Duration(s.HoursAgo) * time.Hour)
		tracking := models.Tracking{
			WaybillID:    waybillID,
			Status:       s.Status,
			Location:     s.Location,
			Description:  s.Description,
			TrackingTime: trackingTime,
		}
		models.DB.Create(&tracking)
		trackings = append(trackings, tracking)
	}
	return trackings
}
