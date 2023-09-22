package user_api

import (
	"Blog_gin/models"
	"Blog_gin/models/ctype"
	"Blog_gin/res"
	"Blog_gin/service/common"
	"Blog_gin/utils/desens"
	"Blog_gin/utils/jwts"
	"github.com/gin-gonic/gin"
)

type UserResponse struct {
	models.UserModel
	RoleID int `json:"role_id"`
}

type UserListRequest struct {
	models.PageInfo
	Role int `json:"role" form:"role"`
}

// UserListView 用户列表
// @Tags 用户管理
// @Summary 用户列表
// @Description 用户列表
// @Param token header string true "token"
// @Param data query models.PageInfo   false  "查询参数"
// @Produce json
// @Success 200 {object} res.Response{data=res.ListResponse[models.UserModel]}
// @Router /api/users [get]
func (UserApi) UserListView(c *gin.Context) {
	_claims, _ := c.Get("claims")
	claims := _claims.(*jwts.CustomClaims)
	var page UserListRequest
	if err := c.ShouldBindJSON(&page); err != nil {
		res.FailWithCode(res.ArgumentError, c)
		return
	}

	var users []UserResponse
	list, count, _ := common.ComList(models.UserModel{Role: ctype.Role(page.Role)},
		common.Option{PageInfo: page.PageInfo,
			Likes: []string{"nick_name"},
		})

	for _, user := range list {
		if ctype.Role(claims.Role) != ctype.PermissionAdmin {
			user.UserName = ""
		}
		user.Tel = desens.DesensitizationTel(user.Tel)
		user.Email = desens.DesensitizationEmail(user.Email)
		//脱敏
		users = append(users, UserResponse{
			UserModel: user,
			RoleID:    int(user.Role),
		})

	}
	res.OkWithList(users, count, c)

}
