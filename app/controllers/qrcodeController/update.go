package qrcodeController

import (
	"errors"

	"4u-go/app/apiException"
	"4u-go/app/services/qrcodeService"
	"4u-go/app/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type updateQrcodeData struct {
	ID           uint   `json:"id" binding:"required"`
	College      uint   `json:"college" binding:"required"`
	Department   string `json:"department" binding:"required"`
	Description  string `json:"description"`
	FeedbackType uint   `json:"feedback_type" binding:"required"`
	Location     string `json:"location" binding:"required"`
	Status       bool   `json:"status"`
}

// UpdateQrcode 更新权益码信息
func UpdateQrcode(c *gin.Context) {
	var data updateQrcodeData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	qrcode, err := qrcodeService.GetQrcodeById(data.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			apiException.AbortWithException(c, apiException.ResourceNotFound, err)
		} else {
			apiException.AbortWithException(c, apiException.ServerError, err)
		}
		return
	}

	{ // 更新权益码信息
		qrcode.College = data.College
		qrcode.Department = data.Department
		qrcode.Description = data.Description
		qrcode.FeedbackType = data.FeedbackType
		qrcode.Location = data.Location
		qrcode.Status = data.Status
	}

	err = qrcodeService.SaveQrcode(qrcode)

	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}
