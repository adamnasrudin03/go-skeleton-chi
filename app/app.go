package app

import (
	"github.com/adamnasrudin03/go-skeleton-chi/app/configs"
	"github.com/adamnasrudin03/go-skeleton-chi/app/controller"
	"github.com/adamnasrudin03/go-skeleton-chi/app/repository"
	"github.com/adamnasrudin03/go-skeleton-chi/app/service"
	"github.com/adamnasrudin03/go-skeleton-chi/pkg/driver"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func WiringRepository(db *gorm.DB, cache *driver.RedisClient, cfg *configs.Configs, logger *logrus.Logger) *repository.Repositories {
	return &repository.Repositories{
		TeamMember: repository.NewTeamMemberRepository(db, *cache, cfg, logger),
	}
}

func WiringService(repo *repository.Repositories, cfg *configs.Configs, logger *logrus.Logger) *service.Services {
	return &service.Services{
		TeamMember: service.NewTeamMemberService(repo.TeamMember, cfg, logger),
	}
}

func WiringController(srv *service.Services, cfg *configs.Configs, logger *logrus.Logger, validator *validator.Validate) *controller.Controllers {
	return &controller.Controllers{
		TeamMember: controller.NewTeamMemberDelivery(srv.TeamMember, logger, validator),
	}
}
