package controller

import (
	"github.com/gin-gonic/gin"
	"ops-api/dao"
	"ops-api/service"
	"strconv"
)

var Account account

type account struct{}

// AddAccount 新增账号
// @Summary 新增账号
// @Description 账号相关接口
// @Tags 账号管理
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param user body service.AccountCreate true "账号信息"
// @Success 200 {string} json "{"code": 0, "msg": "创建成功", "data": nil}"
// @Router /api/v1/account [post]
func (a *account) AddAccount(c *gin.Context) {
	var account = &service.AccountCreate{}

	if err := c.ShouldBind(account); err != nil {
		Response(c, 90400, err.Error())
		return
	}

	userID := c.GetUint("id")
	data, err := service.Account.AddAccount(account, userID)
	if err != nil {
		Response(c, 90500, err.Error())
		return
	}

	CreateOrUpdateResponse(c, 0, "创建成功", data)
}

// AddAccounts 批量新增账号
// @Summary 批量新增账号
// @Description 账号相关接口
// @Tags 账号管理
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param user body service.BatchAccountCreate true "账号信息"
// @Success 200 {string} json "{"code": 0, "msg": "创建成功", "data": nil}"
// @Router /api/v1/accounts [post]
func (a *account) AddAccounts(c *gin.Context) {
	var account = &service.BatchAccountCreate{}

	if err := c.ShouldBind(account); err != nil {
		Response(c, 90400, err.Error())
		return
	}

	userID := c.GetUint("id")
	accounts, err := service.Account.AddAccounts(account.Accounts, userID)
	if err != nil {
		Response(c, 90500, err.Error())
		return
	}

	CreateOrUpdateResponse(c, 0, "创建成功", accounts)
}

// DeleteAccount 删除账号
// @Summary 删除账号
// @Description 账号相关接口
// @Tags 账号管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path int true "账号ID"
// @Success 200 {string} json "{"code": 0, "msg": "删除成功"}"
// @Router /api/v1/account/{id} [delete]
func (a *account) DeleteAccount(c *gin.Context) {

	// 获取账号ID
	accountId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		Response(c, 90500, err.Error())
		return
	}

	// 获取用户ID
	userID := c.GetUint("id")
	if err := service.Account.DeleteAccount(accountId, int(userID)); err != nil {
		Response(c, 90500, err.Error())
		return
	}

	Response(c, 0, "删除成功")
}

// UpdateAccount 更新账号信息
// @Summary 更新账号信息
// @Description 账号相关接口
// @Tags 账号管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param task body dao.AccountUpdate true "账号信息"
// @Success 200 {string} json "{"code": 0, "msg": "更新成功", "data": nil}"
// @Router /api/v1/account [put]
func (a *account) UpdateAccount(c *gin.Context) {
	var data = &dao.AccountUpdate{}

	if err := c.ShouldBind(&data); err != nil {
		Response(c, 90400, err.Error())
		return
	}

	userID := c.GetUint("id")
	account, err := service.Account.UpdateAccount(data, userID)
	if err != nil {
		Response(c, 90500, err.Error())
		return
	}

	CreateOrUpdateResponse(c, 0, "更新成功", account)
}

// BatchUpdateAccountOwner 批量更新账号所有者
// @Summary 批量更新账号所有者
// @Description 账号相关接口
// @Tags 账号管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param account_ids query []int true "账号ID列表"
// @Param new_owner_id query int true "新所有者ID"
// @Success 200 {string} json "{"code": 0, "msg": "更新成功", "data": nil}"
// @Router /api/v1/account/owners [put]
func (a *account) BatchUpdateAccountOwner(c *gin.Context) {
	var data struct {
		Accounts    []uint `json:"accounts" binding:"required"`
		OwnerUserID uint   `json:"owner_user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		Response(c, 90400, err.Error())
		return
	}

	oldOwnerID := c.GetUint("id")
	accounts, err := service.Account.BatchUpdateAccountOwner(data.Accounts, oldOwnerID, data.OwnerUserID)
	if err != nil {
		Response(c, 90500, err.Error())
		return
	}

	CreateOrUpdateResponse(c, 0, "更新成功", accounts)
}

// UpdateAccountUser 用户分享
// @Summary 用户分享
// @Description 账号相关接口
// @Tags 账号管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param users body dao.AccountUpdateUser true "用户ID列表"
// @Success 200 {string} json "{"code": 0, "msg": "更新成功", "data": nil}"
// @Router /api/v1/account/users [put]
func (a *account) UpdateAccountUser(c *gin.Context) {
	var data = &dao.AccountUpdateUser{}

	// 解析请求参数
	if err := c.ShouldBind(&data); err != nil {
		Response(c, 90400, err.Error())
		return
	}

	// 更新用户信息
	userId := c.GetUint("id")
	account, err := service.Account.UpdateAccountUser(data, userId)
	if err != nil {
		Response(c, 90500, err.Error())
		return
	}

	CreateOrUpdateResponse(c, 0, "更新成功", account)
}

// UpdatePassword 更改密码
// @Summary 更改密码
// @Description 账号相关接口
// @Tags 账号管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param user body dao.AccountUpdatePassword true "用户信息"
// @Success 200 {string} json "{"code": 0, "msg": "更新成功", "data": nil}"
// @Router /api/v1/account/password [put]
func (a *account) UpdatePassword(c *gin.Context) {
	var account = &dao.AccountUpdatePassword{}

	// 解析请求参数
	if err := c.ShouldBind(&account); err != nil {
		Response(c, 90400, err.Error())
		return
	}

	// 更新用户信息
	userId := c.GetUint("id")
	if err := service.Account.UpdatePassword(account, userId); err != nil {
		Response(c, 90500, err.Error())
		return
	}

	CreateOrUpdateResponse(c, 0, "更新成功", nil)
}

// GetAccountList 获取账号列表
// @Summary 获取账号列表
// @Description 账号相关接口
// @Tags 账号管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param page query int true "分页"
// @Param limit query int true "分页大小"
// @Param name query string false "账号信息"
// @Success 200 {string} json "{"code": 0, "data": []}"
// @Router /api/v1/accounts [get]
func (a *account) GetAccountList(c *gin.Context) {
	params := new(struct {
		Name  string `form:"name"`
		Page  int    `form:"page" binding:"required"`
		Limit int    `form:"limit" binding:"required"`
	})
	if err := c.Bind(params); err != nil {
		Response(c, 90400, err.Error())
		return
	}

	userID := c.GetUint("id")
	data, err := service.Account.GetAccountList(params.Name, userID, params.Page, params.Limit)
	if err != nil {
		Response(c, 90500, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": data,
	})
}

// GetAccountPassword 获取账号密码
// @Summary 获取账号密码
// @Description 账号相关接口
// @Tags 账号管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path int true "账号ID"
// @Success 200 {string} json "{"code": 0, "data": []}"
// @Router /api/v1/account/password/{id} [delete]
func (a *account) GetAccountPassword(c *gin.Context) {

	// 获取账号ID
	accountID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		Response(c, 90400, err.Error())
		return
	}

	username, _ := c.Get("username")
	userId := c.GetUint("id")
	password, err := service.Account.GetAccountPassword(uint(accountID), username.(string), userId)
	if err != nil {
		Response(c, 90500, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "获取成功，请妥善使用",
		"data": password,
	})
}

// GetSMSCode 获取验证码
// @Summary 获取验证码
// @Description 账号相关接口
// @Tags 账号管理
// @Accept application/json
// @Produce application/json
// @Param user body service.GetVerification true "用户信息"
// @Success 200 {string} json "{"code": 0, "msg": "校验码已发送，5分钟之内有效"}"
// @Router /api/v1/account/code [post]
func (a *account) GetSMSCode(c *gin.Context) {

	var data = &service.GetVerification{}

	// 解析请求参数
	if err := c.ShouldBind(&data); err != nil {
		Response(c, 90400, err.Error())
		return
	}

	// 获取短信验证码
	userID := c.GetUint("id")
	if err := service.Account.GetSMSCode(data, userID); err != nil {
		Response(c, 90500, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "校验码已发送，5分钟之内有效",
	})
}

// CodeVerification 校验验证码
// @Summary 校验验证码
// @Description 账号相关接口
// @Tags 账号管理
// @Param user body service.CodeVerification true "验证码信息"
// @Success 200 {string} json "{"code": 0, "msg": "验证成功，本次验证有效期为10分钟"}"
// @Router /api/v1/account/code_verification [post]
func (a *account) CodeVerification(c *gin.Context) {
	var data = &service.CodeVerification{}

	// 解析请求参数
	if err := c.ShouldBind(&data); err != nil {
		Response(c, 90400, err.Error())
		return
	}

	// 更新用户信息
	userID := c.GetUint("id")
	if err := service.Account.CodeVerification(userID, data); err != nil {
		Response(c, 90500, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "验证成功，本次验证有效期为10分钟",
	})
}
