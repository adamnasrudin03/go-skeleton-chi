package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/adamnasrudin03/go-skeleton-chi/app/configs"
	"github.com/adamnasrudin03/go-skeleton-chi/app/dto"
	"github.com/adamnasrudin03/go-skeleton-chi/app/models"
	"github.com/adamnasrudin03/go-skeleton-chi/pkg/driver"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TeamMemberRepository interface {
	CreateCache(ctx context.Context, key string, data interface{}, ttl time.Duration)
	DeleteCache(ctx context.Context, key string)
	GetCache(ctx context.Context, key string, res interface{}) bool
	GetDetail(ctx context.Context, req dto.TeamMemberDetailReq) (*models.TeamMember, error)
	Create(ctx context.Context, req *models.TeamMember) (*models.TeamMember, error)
	Update(ctx context.Context, req *models.TeamMember) error
	Delete(ctx context.Context, req *models.TeamMember) error
	GetList(ctx context.Context, req dto.TeamMemberListReq) ([]models.TeamMember, error)
}

type TeamMemberRepo struct {
	DB     *gorm.DB
	Cache  driver.RedisClient
	Cfg    *configs.Configs
	Logger *logrus.Logger
}

func NewTeamMemberRepository(
	db *gorm.DB,
	redis driver.RedisClient,
	cfg *configs.Configs,
	logger *logrus.Logger,
) TeamMemberRepository {
	return &TeamMemberRepo{
		DB:     db,
		Cache:  redis,
		Cfg:    cfg,
		Logger: logger,
	}
}

func (r *TeamMemberRepo) CreateCache(ctx context.Context, key string, data interface{}, ttl time.Duration) {
	var (
		opName = "TeamMemberRepository-CreateCache"
		err    error
	)
	if ttl == 0 {
		ttl = r.Cfg.Redis.DefaultCacheTimeOut
	}

	err = r.Cache.Set(key, data, ttl)
	if err != nil {
		r.Logger.Errorf("%v error: %v ", opName, err)
		return
	}
}

func (r *TeamMemberRepo) DeleteCache(ctx context.Context, key string) {
	var (
		opName = "TeamMemberRepository-DeleteCache"
		err    error
	)
	err = r.Cache.Del(key)
	if err != nil {
		r.Logger.Errorf("%v error: %v ", opName, err)
		return
	}
}

func (r *TeamMemberRepo) GetCache(ctx context.Context, key string, res interface{}) bool {
	var (
		opName = "TeamMemberRepository-GetCache"
		err    error
	)

	data, err := r.Cache.Get(key)
	if err != nil {
		r.Logger.Errorf("%v error: %v ", opName, err)
		return false
	}

	err = json.Unmarshal([]byte(data), &res)
	if err != nil {
		r.Logger.Errorf("%v Unmarshal error: %v ", opName, err)
		return false
	}

	return true
}

func (r *TeamMemberRepo) GetDetail(ctx context.Context, req dto.TeamMemberDetailReq) (*models.TeamMember, error) {
	var (
		opName = "TeamMemberRepository-GetDetail"
		err    error
		resp   *models.TeamMember
		column = "*"
	)
	if req.CustomColumn != "" {
		column = req.CustomColumn
	}

	db := r.DB.WithContext(ctx).Model(&models.TeamMember{}).Select(column)

	if req.ID > 0 {
		db = db.Where("id = ?", req.ID)
	}
	if req.NotID > 0 {
		db = db.Where("id != ?", req.NotID)
	}
	if req.Email != "" {
		db = db.Where("email = ?", req.Email)
	}
	if req.UsernameGithub != "" {
		db = db.Where("username_github = ?", req.UsernameGithub)
	}

	err = db.First(&resp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		r.Logger.Errorf("%v error: %v ", opName, err)
		return nil, err
	}

	return resp, nil
}

func (r *TeamMemberRepo) Create(ctx context.Context, req *models.TeamMember) (*models.TeamMember, error) {
	var (
		opName = "TeamMemberRepository-Create"
		err    error
	)
	err = r.DB.WithContext(ctx).Create(req).Error
	if err != nil {
		r.Logger.Errorf("%v error: %v ", opName, err)
		return nil, err
	}

	return req, nil
}

func (r *TeamMemberRepo) Update(ctx context.Context, req *models.TeamMember) error {
	var (
		opName = "TeamMemberRepository-Update"
		err    error
	)
	err = r.DB.WithContext(ctx).Model(&models.TeamMember{}).Where("id = ?", req.ID).Updates(req).Error
	if err != nil {
		r.Logger.Errorf("%v error: %v ", opName, err)
		return err
	}

	return nil
}

func (r *TeamMemberRepo) Delete(ctx context.Context, req *models.TeamMember) error {
	var (
		opName = "TeamMemberRepository-Delete"
		err    error
	)

	err = r.DB.WithContext(ctx).Where("id = ?", req.ID).Delete(&models.TeamMember{}).Error
	if err != nil {
		r.Logger.Errorf("%v error: %v ", opName, err)
		return err
	}

	return nil
}

func (r *TeamMemberRepo) GetList(ctx context.Context, req dto.TeamMemberListReq) ([]models.TeamMember, error) {
	var (
		opName = "TeamMemberRepository-GetList"
		err    error
		resp   []models.TeamMember
		column = "*"
	)
	if req.CustomColumns != "" {
		column = req.CustomColumns
	}

	db := r.DB.WithContext(ctx).Model(&models.TeamMember{}).Select(column)
	if req.Search != "" {
		db = db.Where("email LIKE ?", "%"+req.Search+"%")
	}

	if !req.IsNotDefaultQuery {
		req = req.DefaultQuery()
	}
	if !req.IsNoLimit {
		db = db.Offset(int(req.Offset)).Limit(int(req.Limit))
	}

	if models.IsValidOrderBy[req.OrderBy] && req.SortBy != "" {
		db = db.Order(req.SortBy + " " + req.OrderBy)
	}

	err = db.Find(&resp).Error
	if err != nil {
		r.Logger.Errorf("%v error: %v ", opName, err)
		return nil, err
	}

	return resp, nil
}
