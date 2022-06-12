package service

import (
	"time"
	"webapi/dao/form_req"
	"webapi/dao/form_resp"
	"webapi/dao/mongo"
	"webapi/internal/password"
	"webapi/internal/wrapper"
	"webapi/models"
	"webapi/support"
	"webapi/utils"

	"github.com/globalsign/mgo/bson"
)

// CreateUserHandler 创建用户
func CreateUserHandler(ctx *wrapper.Context, reqBody interface{}) (err error) {
	traceCtx := ctx.Request().Context()
	req := reqBody.(*form_req.CreateUserReq)
	resp := form_resp.StatusResp{Status: support.StatusOK}
	if !utils.String.Compare(req.Password, req.Confirm) {
		support.SendApiErrorResponse(ctx, support.PasswordNotConfirm, 0)
		return nil
	}
	existQuery := bson.M{"user_id": req.UserId}
	if mongo.User.IsExist(traceCtx, existQuery) {
		support.SendApiErrorResponse(ctx, support.UserIsExist, 0)
		return nil
	}
	passwordStrengthLevel := utils.Logic.GetPasswordStrength(req.Password)
	if passwordStrengthLevel == 0 {
		support.SendApiErrorResponse(ctx, support.PasswordStrengthFailed, 0)
		return nil
	}
	// 创建账户
	userDoc := models.User{
		Role:            req.Role,
		UserName:        req.Username,
		UserId:          req.UserId,
		Password:        password.MakePassword(req.Password),
		LastPwdChangeTm: time.Now(),
		LastLoginTm:     time.Now(),
	}
	if req.Role == 2 {
		userDoc.Grade = req.Grade
		userDoc.Class = req.Class
	}
	if err = mongo.User.Create(traceCtx, userDoc); err != nil {
		support.SendApiErrorResponse(ctx, support.CreateUserFailed, 0)
		return nil
	}
	support.SendApiResponse(ctx, resp, "")
	return
}

// UserInfoHandler 获取用户信息
func UserInfoHandler(ctx *wrapper.Context, reqBody interface{}) (err error) {
	traceCtx := ctx.Request().Context()
	var userDoc models.User
	userDoc, err = mongo.User.FindByUserId(traceCtx, ctx.UserToken.UserId)
	if err != nil {
		support.SendApiErrorResponse(ctx, support.UserNotExist, 0)
		return nil
	}
	resp := form_resp.UserInfoResp{
		UserId:        userDoc.UserId,
		Role:          userDoc.Role,
		UserName:      userDoc.UserName,
		Grade:         userDoc.Grade,
		Class:         userDoc.Class,
		LoginTime:     utils.Time2String(time.Now()),
		LastLoginTime: utils.Time2String(userDoc.LastLoginTm),
	}
	support.SendApiResponse(ctx, resp, "success")
	return nil
}

// UserPasswordHandler 忘记密码
func UserPasswordHandler(ctx *wrapper.Context, reqBody interface{}) (err error) {
	traceCtx := ctx.Request().Context()
	req := reqBody.(*form_req.UserPasswordReq)
	if !utils.String.Compare(req.Password, req.Confirm) {
		support.SendApiErrorResponse(ctx, support.PasswordNotConfirm, 0)
		return nil
	}
	query := bson.M{"user_id": req.UserId, "username": req.UserName, "role": req.Role}
	_, err = mongo.User.FindOne(traceCtx, query)
	if err != nil {
		support.SendApiErrorResponse(ctx, support.UserNotExist, 0)
		return nil
	}
	upset := bson.M{"password": password.MakePassword(req.Password)}
	err = mongo.User.Update(traceCtx, query, upset)
	if err != nil {
		support.SendApiErrorResponse(ctx, support.UpdatePasswordFailed, 0)
		return nil
	}

	resp := form_resp.UserPasswordResp{
		Password: req.Password,
	}
	support.SendApiResponse(ctx, resp, "success")
	return nil
}

// ChangePasswordHandler 修改账户密码
func ChangePasswordHandler(ctx *wrapper.Context, reqBody interface{}) (err error) {
	traceCtx := ctx.Request().Context()
	req := reqBody.(*form_req.ChangePasswordReq)
	var userDoc models.User
	userDoc, err = mongo.User.FindByUserId(traceCtx, ctx.UserToken.UserId)
	if err != nil {
		support.SendApiErrorResponse(ctx, support.UserNotExist, 0)
		return nil
	}
	if !password.CheckPassword(req.Password, userDoc.Password) {
		support.SendApiErrorResponse(ctx, support.PasswordWrong, 0)
		return nil
	}
	query := bson.M{"user_id": ctx.UserToken.UserId}
	newPwd := password.MakePassword(req.Password)
	upset := bson.M{"password": newPwd, "last_pwd_change_tm": time.Now()}
	err = mongo.User.Update(traceCtx, query, upset)
	if err != nil {
		support.SendApiErrorResponse(ctx, support.UpdatePasswordFailed, 0)
		return nil
	}
	resp := form_resp.StatusResp{Status: "ok"}
	support.SendApiResponse(ctx, resp, "success")
	return nil
}
