package handlers

import (
	"github.com/VolodyaLarin/rsoi-lab-02/internal/flight/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type FlightHandlerV1 struct {
	uc *usecase.FlightUsecase
}

type FlightAnswer struct {
	Page          int `json:"page"`
	PageSize      int `json:"pageSize"`
	TotalElements int `json:"totalElements"`

	Items []usecase.FlightDto `json:"items"`
}

func NewFlightHandlerV1(uc *usecase.FlightUsecase) *FlightHandlerV1 {
	return &FlightHandlerV1{uc: uc}
}

func (p FlightHandlerV1) RegisterRoutes(router gin.IRouter) {
	router.GET("/flights/", p.list)
}

func (h FlightHandlerV1) list(ctx *gin.Context) {
	uids := ctx.QueryArray("uid")
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		page = 0
	}
	size, err := strconv.Atoi(ctx.Query("size"))
	if err != nil {
		size = 1000
	}

	err, flights, count := h.uc.List(ctx, &usecase.FlightFilter{
		Limit:   int64(size),
		Offset:  int64(page * size),
		Numbers: uids,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	ctx.JSON(http.StatusOK, FlightAnswer{
		Page:          page,
		PageSize:      size,
		TotalElements: int(count),
		Items:         flights,
	})
}
