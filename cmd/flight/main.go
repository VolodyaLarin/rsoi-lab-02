package main

import (
	handler "github.com/VolodyaLarin/rsoi-lab-02/internal/flight/handlers"
	"github.com/VolodyaLarin/rsoi-lab-02/internal/flight/usecase"
	"github.com/VolodyaLarin/rsoi-lab-02/internal/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	//dsn := os.Getenv("DB_DSN")
	//db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	//if err != nil {
	//	log.WithError(err).Fatal("can't start db con")
	//}
	//err = db.AutoMigrate(repo.TicketModel{})
	//if err != nil {
	//	log.WithError(err).Fatal("can't migrate")
	//	return
	//}

	r := gin.Default()
	apiV1R := r.Group("/api/v1/")
	r.Use(utils.JSONLogMiddleware())

	ticketUc := &usecase.FlightUsecase{}
	handler.NewFlightHandlerV1(ticketUc).RegisterRoutes(apiV1R)

	r.Run(":8080")
}
