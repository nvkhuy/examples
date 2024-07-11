package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/services/consumer/tasks"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// @Tags Buyer-ChatRoom
// @Summary Get chat chat room list
// @Description This API allows admin to list chat room with pagination
// @Accept  json
// @Produce  json
// @Param page query int false
// @Param limit query int false
// @Success 200 {object} query.Pagination{Records=[]models.ChatRoom}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/chat/get_chat_room_list [get]
func BuyerGetChatRoomList(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var params models.GetChatRoomListRequest
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	params.JwtClaimsInfo = claims

	chatRooms, err := repo.NewChatRoomRepo(cc.App.DB).GetChatRoomList(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(chatRooms)
}

// @Tags Buyer-ChatRoom
// @Summary Create chat room
// @Description This API allows buyer to create chat room
// @Accept  json
// @Produce  json
// @Param data body models.CreateChatRoomRequest true
// @Success 200 {object} models.ChatRoom
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/chat/create_chat_room [post]
func BuyerCreateChatRoom(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var params models.CreateChatRoomRequest
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	params.JwtClaimsInfo = claims
	if claims.GetRole().IsBuyer() {
		params.BuyerID = claims.GetUserID()
	}
	chatRoom, err := repo.NewChatRoomRepo(cc.App.DB).CreateChatRoom(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(chatRoom)
}

// @Tags Buyer-ChatRoom
// @Summary Mark seen message
// @Description This API allows buyer to mark seen message
// @Accept  json
// @Produce  json
// @Param room_id query string true
// @Success 200 {object} boolean
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/chat/mark_seen_chat_room_message [put]
func BuyerMarkSeenChatRoomMessage(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var params models.MarkSeenChatRoomMessageRequest
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	params.JwtClaimsInfo = claims

	err = repo.NewChatRoomRepo(cc.App.DB).MarkSeenChatRoomMessage(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	_, err = tasks.SeenChatRoomTask{
		RoomID:     params.RoomID,
		SeenUserID: params.GetUserID(),
	}.Dispatch(cc.Request().Context())
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(true)
}

// @Tags Buyer-ChatRoom
// @Summary Get unseen message count on specific room
// @Description This API allows buyer to get unseen message count on specific room
// @Accept  json
// @Produce  json
// @Param inquiry_id query string true
// @Param purchase_order_id query string true
// @Param bulk_purchase_order_id query string true
// @Success 200 {object} models.CountUnSeenChatMessageResponse
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/chat/count_unseen_message_on_room [get]
func BuyerCountUnseenChatMessageOnRoom(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var params models.CountUnSeenChatMessageOnRoomRequest
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	params.JwtClaimsInfo = claims
	if claims.GetRole().IsBuyer() {
		params.BuyerID = claims.GetUserID()
	}
	count, err := repo.NewChatRoomRepo(cc.App.DB).CountUnSeenChatMessageOnRoom(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(count)
}
