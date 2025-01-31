package lostAndFoundController

import (
	"errors"

	"4u-go/app/apiException"
	"4u-go/app/services/lostAndFoundService"
	"4u-go/app/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type getLostAndFoundListData struct {
	Type   bool  `json:"type"`   // 1-失物 0-寻物
	Campus uint8 `json:"campus"` // 校区 0-其他 1-朝晖 2-屏峰 3-莫干山
	Kind   uint8 `json:"kind"`   // 物品种类 0-全部 1-其他 2-饭卡 3-电子 4-文体 5-衣包 6-证件
}
type getLostAndFoundListResponse struct {
	LostAndFoundList []lostAndFoundElement `json:"list"`
}
type lostAndFoundElement struct {
	ID           uint     `json:"id"`
	Imgs         []string `json:"imgs"`
	Name         string   `json:"name"`
	Place        string   `json:"place"`
	Time         string   `json:"time"`
	Introduction string   `json:"introduction"`
	Kind         uint8    `json:"kind"`
}

// GetLostAndFoundList 获取失物招领列表
func GetLostAndFoundList(c *gin.Context) {
	var data getLostAndFoundListData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	list, err := lostAndFoundService.GetLostAndFoundList(data.Type, data.Campus, data.Kind)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}

	lostAndFoundList := make([]lostAndFoundElement, 0)
	for _, record := range list {
		// 将string转为[]string
		imgs := utils.StringToStrings(record.Imgs)
		lostAndFoundList = append(lostAndFoundList, lostAndFoundElement{
			ID:           record.ID,
			Imgs:         imgs,
			Name:         record.Name,
			Place:        record.Place,
			Time:         record.Time,
			Introduction: record.Introduction,
		})
	}

	utils.JsonSuccessResponse(c, getLostAndFoundListResponse{
		LostAndFoundList: lostAndFoundList,
	})
}

type getLostAndFoundContentData struct {
	ID uint `json:"id" binding:"required"`
}

// GetLostAndFoundContact 获取失物招领联系方式
func GetLostAndFoundContact(c *gin.Context) {
	var data getLostAndFoundContentData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	contact, err := lostAndFoundService.GetLostAndFoundContact(data.ID, utils.GetUser(c).StudentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			apiException.AbortWithException(c, apiException.ResourceNotFound, err)
		} else {
			apiException.AbortWithException(c, apiException.ServerError, err)
		}
		return
	}

	utils.JsonSuccessResponse(c, contact)
}

type latestLostAndFoundResponse struct {
	Type         bool   `json:"type"`
	Imgs         string `json:"imgs"`
	Name         string `json:"name"`
	Place        string `json:"place"`
	Introduction string `json:"introduction"`
}

// GetLatestLostAndFound 获取最新失物招领
func GetLatestLostAndFound(c *gin.Context) {
	record, err := lostAndFoundService.GetLatestLostAndFound()
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}

	utils.JsonSuccessResponse(c, latestLostAndFoundResponse{
		Type:         record.Type,
		Imgs:         record.Imgs,
		Name:         record.Name,
		Place:        record.Place,
		Introduction: record.Introduction,
	})
}

type getLostAndFoundStatusData struct {
	Status uint8 `json:"status"` // 状态 0-已撤回 1-已审核 2-审核中
}
type getLostAndFoundStatusResponse struct {
	List []lostAndFoundStatusElement `json:"list"`
}
type lostAndFoundStatusElement struct {
	ID           uint     `json:"id"`
	Type         bool     `json:"type"`
	Imgs         []string `json:"imgs"`
	Name         string   `json:"name"`
	Kind         uint8    `json:"kind"`
	Place        string   `json:"place"`
	Time         string   `json:"time"`
	Introduction string   `json:"introduction"`
	IsApproved   uint8    `json:"is_approved"`
}

// GetUserLostAndFoundStatus 查看失物招领信息的状态
func GetUserLostAndFoundStatus(c *gin.Context) {
	var data getLostAndFoundStatusData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	list, err := lostAndFoundService.GetUserLostAndFoundStatus(utils.GetUser(c).StudentID, data.Status)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}

	lostAndFoundList := make([]lostAndFoundStatusElement, 0)
	for _, record := range list {
		// 将string转为[]string
		imgs := utils.StringToStrings(record.Imgs)
		lostAndFoundList = append(lostAndFoundList, lostAndFoundStatusElement{
			ID:           record.ID,
			Type:         record.Type,
			Imgs:         imgs,
			Name:         record.Name,
			Kind:         record.Kind,
			Place:        record.Place,
			Time:         record.Time,
			Introduction: record.Introduction,
			IsApproved:   record.IsApproved,
		})
	}
	utils.JsonSuccessResponse(c, getLostAndFoundStatusResponse{
		List: lostAndFoundList,
	})
}
