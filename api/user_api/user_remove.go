package user_api

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"Blog_gin/res"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserRemoveView 删除用户接口
// @Summary 删除用户
// @Description 根据传入的用户ID列表删除对应的用户和相关信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param token header string true "token"
// @Param request body models.RemoveRequest true "用户ID列表"
// @Success 200 {object} res.Response{}
// @Router /api/users [delete]
func (UserApi) UserRemoveView(c *gin.Context) {
	var cr models.RemoveRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithCode(res.ArgumentError, c)
		return
	}

	var userList []models.UserModel
	count := global.DB.Find(&userList, cr.IDList).RowsAffected
	if count == 0 {
		res.FailWithMessage("用户不存在", c)
		return
	}

	//事务
	err = global.DB.Transaction(func(tx *gorm.DB) error {
		//TODO:删除用户，信息表，评论表，用户收藏的文章，用户发布的文章
		err = tx.Delete(&userList).Error
		if err != nil {
			global.Log.Errorf("UserRemoveView Failed err:%+v\n", err.Error())
			tx.Rollback()
			return err
		}
		return err
	})

	if err != nil {
		global.Log.Errorf("删除用户失败%+v\n", err.Error())
		res.FailWithMessage("删除用户失败", c)
		return
	}
	res.OkWithMessage(fmt.Sprintf("共删除 %d 个用户", count), c)
}
