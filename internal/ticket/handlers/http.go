package handlers

import (
	"github.com/VolodyaLarin/rsoi-lab-02/internal/ticket"
	"github.com/VolodyaLarin/rsoi-lab-02/internal/ticket/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type TicketHandlerV1 struct {
	uc *usecase.TicketUsecase
}

func NewTicketHandlerV1(uc *usecase.TicketUsecase) *TicketHandlerV1 {
	return &TicketHandlerV1{uc: uc}
}

func (p TicketHandlerV1) RegisterRoutes(router gin.IRouter) {

	router.GET("/tickets/", p.list)
	router.POST("/tickets/", p.create)
	router.DELETE("/tickets/:uid", p.delete)
}

func (h TicketHandlerV1) list(ctx *gin.Context) {
	username := ctx.GetHeader("X-User-Name")
	if username == "" {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	uidQ := ctx.QueryArray("uids[]")
	var uids []uuid.UUID
	for _, i := range uidQ {
		uid, err := uuid.Parse(i)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, nil)
			return
		}
		uids = append(uids, uid)
	}
	log.Warning(uidQ)

	err, tickets := h.uc.List(ctx, &bonus.TicketFilter{
		Username: username,
		Uids:     uids,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	ctx.JSON(http.StatusOK, tickets)
}

type TicketCreate struct {
	FlightNumber string `json:"flightNumber" binding:"required"`
	Price        int64  `json:"price" binding:"required"`
}

func (h TicketHandlerV1) create(ctx *gin.Context) {
	username := ctx.GetHeader("X-User-Name")
	if username == "" {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	buyData := TicketCreate{}

	err := ctx.ShouldBindJSON(&buyData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err, ticket := h.uc.Buy(ctx, username, buyData.FlightNumber, buyData.Price)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	ctx.JSON(http.StatusCreated, ticket)
}

func (h TicketHandlerV1) delete(ctx *gin.Context) {
	username := ctx.GetHeader("X-User-Name")
	uid, err := uuid.Parse(ctx.Param("uid"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, nil)
		return
	}

	err, _ = h.uc.Cancel(ctx, username, uid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	ctx.Status(http.StatusNoContent)
}
