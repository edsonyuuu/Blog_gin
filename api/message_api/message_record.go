package message_api

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"Blog_gin/res"
	"Blog_gin/utils/jwts"
	"fmt"
	"github.com/gin-gonic/gin"
)

type MessageRecordRequest struct {
	UserID uint `json:"user_id" form:"user_id" binding:"required"  msg:"请输入查询的用户id"`
}

// MessageRecordView 用户的消息记录
// @Tags 消息管理
// @Summary 用户的消息记录
// @Description 用户的消息记录
// @Param token header string true "token"
// @Param data query MessageRecordRequest    true  "查询参数"
// @Produce json
// @Success 200 {object} res.Response{data=[]models.MessageModel}
// @Router /api/messages_record [get]
func (MessageApi) MessageRecordView(c *gin.Context) {

	var cr MessageRecordRequest
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		fmt.Printf("错误  err:%+v\n", err.Error())
		res.FailWithError(err, &cr, c)
		return
	}

	_claims, _ := c.Get("claims")
	claims := _claims.(*jwts.CustomClaims)
	var _messageList []models.MessageModel
	var messageList = make([]models.MessageModel, 0)
	global.DB.Order("created_at asc").Find(&_messageList, "send_user_id = ? or rev_user_id = ?", claims.UserID, claims.UserID)

	for _, model := range _messageList {

		//判断是否为一个组
		if model.RevUserID == cr.UserID || model.SendUserID == cr.UserID {
			messageList = append(messageList, model)
		}
	}

	//这里设置是否为已读未读，之后再说
	res.OkWithData(messageList, c)
}
