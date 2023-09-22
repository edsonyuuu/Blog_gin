package message_api

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"Blog_gin/res"
	"github.com/gin-gonic/gin"
)

type MessageRequest struct {
	SendUserID uint   `json:"send_user_id" binding:"required"` // 发送人id
	RevUserID  uint   `json:"rev_user_id" binding:"required"`  // 接收人id
	Content    string `json:"content" binding:"required"`      // 消息内容
}

// MessageCreateView 是创建消息的处理函数
// @Summary 创建消息
// @Description 创建新的消息并发送给接收人
// @Tags 消息管理
// @Accept json
// @Produce json
// @Param token header string true "token"
// @Param messageRequest body MessageRequest true "消息请求体"
// @Success 200 {object} res.Response{}
// @Router /api/messages [post]
func (MessageApi) MessageCreateView(c *gin.Context) {
	//当前用户发送消息
	var cr MessageRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, &cr, c)
		return
	}
	var sendUser, RecvUser models.UserModel

	err = global.DB.Take(&sendUser, cr.SendUserID).Error
	if err != nil {
		res.FailWithMessage("发送人不存在", c)
		return
	}
	err = global.DB.Take(&RecvUser, cr.RevUserID).Error
	if err != nil {
		res.FailWithMessage("接收人不存在", c)
		return
	}

	err = global.DB.Create(&models.MessageModel{
		SendUserID:       cr.SendUserID,
		SendUserNickName: sendUser.NickName,
		SendUserAvatar:   sendUser.Avatar,
		RevUserID:        cr.RevUserID,
		RevUserNickName:  RecvUser.NickName,
		RevUserAvatar:    RecvUser.Avatar,
		IsRead:           false,
		Content:          cr.Content,
	}).Error
	if err != nil {
		global.Log.Error(err)
		res.FailWithMessage("消息发送失败", c)
		return
	}

	res.OkWithMessage("消息发送成功", c)

}
