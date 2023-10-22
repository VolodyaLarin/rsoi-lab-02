package handlers

import (
	bonus2 "github.com/VolodyaLarin/rsoi-lab-02/internal/bonus"
	"github.com/VolodyaLarin/rsoi-lab-02/internal/bonus/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type BonusHandlerV1 struct {
	uc *usecase.BonusUsecase
}

func NewBonusHandlerV1(uc *usecase.BonusUsecase) *BonusHandlerV1 {
	return &BonusHandlerV1{uc: uc}
}

func (p BonusHandlerV1) RegisterRoutes(router gin.IRouter) {

	router.GET("/bonus/", p.list)
	router.POST("/bonus/", p.create)
	router.DELETE("/bonus/:uid", p.delete)
}

func (h BonusHandlerV1) list(ctx *gin.Context) {
	username := ctx.GetHeader("X-User-Name")
	if username == "" {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	err, bonusDetails := h.uc.List(ctx, username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	ctx.JSON(http.StatusOK, bonusDetails)
}

type BonusDetails struct {
	TicketUid       string `json:"ticketUid" binding:"required"`
	PaidFromBalance *bool  `json:"paidFromBalance" binding:"required"`
	FlightNumber    string `json:"flightNumber" binding:"required"`
	Price           int64  `json:"price" binding:"required"`
}

func (h BonusHandlerV1) create(ctx *gin.Context) {
	username := ctx.GetHeader("X-User-Name")
	if username == "" {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	buyData := BonusDetails{}

	err := ctx.ShouldBindJSON(&buyData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ticketUid, err := uuid.Parse(buyData.TicketUid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	data := &bonus2.BonusHistoryDto{}
	bonus := &bonus2.BonusDto{}

	if *buyData.PaidFromBalance {
		err, data, bonus = h.uc.MinusFromPayment(ctx, username, ticketUid, buyData.Price, false)
	} else {
		err, data, bonus = h.uc.PlusFromPayment(ctx, username, ticketUid, buyData.Price, false)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"item":      data,
		"privelege": bonus,
	})
}

func (h BonusHandlerV1) delete(ctx *gin.Context) {
	username := ctx.GetHeader("X-User-Name")
	uid, err := uuid.Parse(ctx.Param("uid"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, nil)
		return
	}

	err = h.uc.RevertFromPayment(ctx, username, uid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	ctx.Status(http.StatusNoContent)
}
