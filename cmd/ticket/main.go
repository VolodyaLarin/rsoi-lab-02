package main

import (
	handler "github.com/VolodyaLarin/rsoi-lab-02/internal/ticket/handlers"
	"github.com/VolodyaLarin/rsoi-lab-02/internal/ticket/repo"
	"github.com/VolodyaLarin/rsoi-lab-02/internal/ticket/usecase"
	"github.com/VolodyaLarin/rsoi-lab-02/internal/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"os"
)

// set gin write to logrus debug.
func main() {
	dsn := os.Getenv("DB_DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.WithError(err).Fatal("can't start db con")
	}
	err = db.AutoMigrate(repo.TicketModel{})
	if err != nil {
		log.WithError(err).Fatal("can't migrate")
		return
	}

	r := gin.Default()
	r.Use(utils.JSONLogMiddleware())

	apiV1R := r.Group("/api/v1/")

	rep := repo.NewTicketRepo(db)
	ticketUc := usecase.NewTicketUsecase(rep)
	handler.NewTicketHandlerV1(ticketUc).RegisterRoutes(apiV1R)

	r.GET("/manage/health", func(context *gin.Context) {
		context.Status(http.StatusOK)
	})

	r.Run(":8080")
}
