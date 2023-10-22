package main

import (
	handler "github.com/VolodyaLarin/rsoi-lab-02/internal/bonus/handlers"
	"github.com/VolodyaLarin/rsoi-lab-02/internal/bonus/repo"
	"github.com/VolodyaLarin/rsoi-lab-02/internal/bonus/usecase"
	"github.com/VolodyaLarin/rsoi-lab-02/internal/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.WithError(err).Fatal("can't start db con")
	}
	err = db.AutoMigrate(repo.PrivilegeModel{}, repo.PrivilegeHistoryModel{})
	if err != nil {
		log.WithError(err).Fatal("can't migrate")
		return
	}

	r := gin.Default()
	apiV1R := r.Group("/api/v1/")
	r.Use(utils.JSONLogMiddleware())

	repo := repo.NewBonusRepo(db)

	uc := usecase.NewBonusUsecase(repo)
	handler.NewBonusHandlerV1(uc).RegisterRoutes(apiV1R)

	r.Run(":8080")
}
