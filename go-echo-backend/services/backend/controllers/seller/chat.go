package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/services/consumer/tasks"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// @Tags Seller-Chat
// @Summary Create chat message
// @Description This API allows seller to create chat message
// @Accept  json
// @Produce  json
// @Param data body models.CreateChatMessageRequest true
// @Success 200 {object} models.ChatMessage
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/chat/create_chat_message [post]
func SellerCreateChatMessage(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var params models.CreateChatMessageRequest
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	params.JwtClaimsInfo = claims
	params.SenderID = claims.GetUserID()

	msg, err := repo.NewChatRepo(cc.App.DB).CreateChatMessage(&params)

	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	_, err = tasks.SendChatMessageTask{
		ChatMessage: msg,
	}.Dispatch(cc.Request().Context())
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(msg)
}

// @Tags Seller-Chat
// @Summary Get chat message list
// @Description This API allows seller to list message with pagination
// @Accept  json
// @Produce  json
// @Param page query int false
// @Param limit query int false
// @Param inquiry_id query string false
// @Param purchase_order_id query string false
// @Param bulk_purchase_order_id query string false
// @Success 200 {object} query.Pagination{Records=[]models.ChatMessage}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/chat/get_message_list [get]
func SellerGetChatMessageList(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var params models.GetMessageListRequest
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	params.JwtClaimsInfo = claims

	messageList, err := repo.NewChatRepo(cc.App.DB).GetMessageList(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(messageList)
}

// @Tags Seller-Chat
// @Summary Get chat user relevant stage
// @Description This API allows seller to retrieve chat user relevant stage
// @Accept  json
// @Produce  json
// @Param inquiry_id query string false
// @Param purchase_order_id query string false
// @Param bulk_purchase_order_id query string false
// @Success 200 {object} []models.ChatRoomStatus
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/chat/get_chat_relevant_stage [get]
func SellerGetChatUserRelevantStage(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var params models.GetChatUserRelevantStageRequest
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	params.JwtClaimsInfo = claims

	chatRoomStatues, err := repo.NewChatRepo(cc.App.DB).GetChatRelevantStage(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(chatRoomStatues)
}
