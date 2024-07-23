package service

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/adamnasrudin03/go-skeleton-chi/app/configs"
	"github.com/adamnasrudin03/go-skeleton-chi/app/dto"
	"github.com/adamnasrudin03/go-skeleton-chi/app/models"
	"github.com/adamnasrudin03/go-skeleton-chi/app/repository/mocks"
	"github.com/adamnasrudin03/go-skeleton-chi/pkg/driver"

	response_mapper "github.com/adamnasrudin03/go-helpers/response-mapper/v1"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TeamMemberServiceTestSuite struct {
	suite.Suite
	repo        *mocks.TeamMemberRepository
	ctx         context.Context
	service     TeamMemberService
	teamMember  models.TeamMember
	teamMembers []models.TeamMember
}

func (srv *TeamMemberServiceTestSuite) SetupTest() {
	var (
		cfg    = &configs.Configs{}
		logger = driver.Logger(cfg)
	)
	srv.teamMember = models.TeamMember{
		ID:             1,
		Name:           "adam",
		UsernameGithub: "adamnasrudin.vercel.app",
		Email:          "adam@example.com",
	}
	srv.teamMembers = []models.TeamMember{
		srv.teamMember,
		{
			ID:             2,
			Name:           "adam nasrudin",
			UsernameGithub: "adamnasrudin03",
			Email:          "adamnasrudin@example.com",
		},
	}

	srv.repo = &mocks.TeamMemberRepository{}
	srv.ctx = context.Background()
	srv.service = NewTeamMemberService(srv.repo, cfg, logger)
}

func TestTeamMemberService(t *testing.T) {
	suite.Run(t, new(TeamMemberServiceTestSuite))
}

func (srv *TeamMemberServiceTestSuite) TestTeamMemberSrv_GetByID() {
	tests := []struct {
		name     string
		id       uint64
		mockFunc func(input uint64)
		want     *models.TeamMember
		wantErr  bool
	}{
		{
			name: "Success with cache",
			id:   srv.teamMember.ID,
			mockFunc: func(input uint64) {
				key := models.KeyCacheTeamMemberDetail(input)
				res := srv.teamMember
				srv.repo.On("GetCache", mock.Anything, key, &models.TeamMember{
					ID: 0,
				}).Return(true).Run(func(args mock.Arguments) {
					target := args.Get(2).(*models.TeamMember)
					*target = res
				}).Once()
			},
			want:    &srv.teamMember,
			wantErr: false,
		},
		{
			name: "failed ge db",
			id:   srv.teamMember.ID,
			mockFunc: func(input uint64) {
				key := models.KeyCacheTeamMemberDetail(input)
				srv.repo.On("GetCache", mock.Anything, key, &models.TeamMember{
					ID: 0,
				}).Return(false).Once()

				srv.repo.On("GetDetail", mock.Anything, dto.TeamMemberDetailReq{
					ID: input,
				}).Return(nil, errors.New("invalid")).Once()
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "not found",
			id:   srv.teamMember.ID,
			mockFunc: func(input uint64) {
				key := models.KeyCacheTeamMemberDetail(input)
				srv.repo.On("GetCache", mock.Anything, key, &models.TeamMember{
					ID: 0,
				}).Return(false).Once()

				srv.repo.On("GetDetail", mock.Anything, dto.TeamMemberDetailReq{
					ID: input,
				}).Return(nil, nil).Once()
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success",
			id:   srv.teamMember.ID,
			mockFunc: func(input uint64) {
				key := models.KeyCacheTeamMemberDetail(input)
				srv.repo.On("GetCache", mock.Anything, key, &models.TeamMember{
					ID: 0,
				}).Return(false).Once()

				srv.repo.On("GetDetail", mock.Anything, dto.TeamMemberDetailReq{
					ID: input,
				}).Return(&srv.teamMember, nil).Once()

				srv.repo.On("CreateCache", mock.Anything, key, &srv.teamMember, time.Minute).Return().Once()

			},
			want:    &srv.teamMember,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		srv.T().Run(tt.name, func(t *testing.T) {
			if tt.mockFunc != nil {
				tt.mockFunc(tt.id)
			}

			got, err := srv.service.GetByID(srv.ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("TeamMemberSrv.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TeamMemberSrv.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (srv *TeamMemberServiceTestSuite) TestTeamMemberSrv_Create() {
	params := dto.TeamMemberCreateReq{
		Name:           srv.teamMember.Name,
		UsernameGithub: srv.teamMember.UsernameGithub,
		Email:          srv.teamMember.Email,
	}
	tests := []struct {
		name     string
		req      dto.TeamMemberCreateReq
		mockFunc func(input dto.TeamMemberCreateReq)
		want     *models.TeamMember
		wantErr  bool
	}{
		{
			name: "duplicate",
			req:  params,
			mockFunc: func(input dto.TeamMemberCreateReq) {
				srv.repo.On("GetDetail", mock.Anything, dto.TeamMemberDetailReq{
					CustomColumn: "id",
					Email:        input.Email,
				}).Return(nil, nil).Once()

				srv.repo.On("GetDetail", mock.Anything, dto.TeamMemberDetailReq{
					CustomColumn:   "id",
					UsernameGithub: input.UsernameGithub,
				}).Return(&models.TeamMember{ID: srv.teamMember.ID}, nil).Once()

			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "failed create record",
			req:  params,
			mockFunc: func(input dto.TeamMemberCreateReq) {
				srv.repo.On("GetDetail", mock.Anything, dto.TeamMemberDetailReq{
					CustomColumn: "id",
					Email:        input.Email,
				}).Return(nil, nil).Once()

				srv.repo.On("GetDetail", mock.Anything, dto.TeamMemberDetailReq{
					CustomColumn:   "id",
					UsernameGithub: input.UsernameGithub,
				}).Return(nil, nil).Once()

				record := &models.TeamMember{
					Name:           input.Name,
					Email:          input.Email,
					UsernameGithub: input.UsernameGithub,
				}
				srv.repo.On("Create", mock.Anything, record).Return(nil, errors.New("invalid")).Once()

			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success",
			req:  params,
			mockFunc: func(input dto.TeamMemberCreateReq) {
				srv.repo.On("GetDetail", mock.Anything, dto.TeamMemberDetailReq{
					CustomColumn: "id",
					Email:        input.Email,
				}).Return(nil, nil).Once()

				srv.repo.On("GetDetail", mock.Anything, dto.TeamMemberDetailReq{
					CustomColumn:   "id",
					UsernameGithub: input.UsernameGithub,
				}).Return(nil, nil).Once()

				record := &models.TeamMember{
					Name:           input.Name,
					Email:          input.Email,
					UsernameGithub: input.UsernameGithub,
				}
				srv.repo.On("Create", mock.Anything, record).Return(&srv.teamMember, nil).Once()

			},
			want:    &srv.teamMember,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		srv.T().Run(tt.name, func(t *testing.T) {
			if tt.mockFunc != nil {
				tt.mockFunc(tt.req)
			}

			got, err := srv.service.Create(srv.ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("TeamMemberSrv.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TeamMemberSrv.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (srv *TeamMemberServiceTestSuite) TestTeamMemberSrv_DeleteByID() {
	tests := []struct {
		name     string
		id       uint64
		mockFunc func(input uint64)
		wantErr  bool
	}{
		{
			name: "not found",
			id:   srv.teamMember.ID,
			mockFunc: func(input uint64) {
				key := models.KeyCacheTeamMemberDetail(input)
				srv.repo.On("GetCache", mock.Anything, key, &models.TeamMember{
					ID: 0,
				}).Return(false).Once()

				srv.repo.On("GetDetail", mock.Anything, dto.TeamMemberDetailReq{
					ID: input,
				}).Return(nil, nil).Once()
			},
			wantErr: true,
		},
		{
			name: "failed delete",
			id:   srv.teamMember.ID,
			mockFunc: func(input uint64) {
				key := models.KeyCacheTeamMemberDetail(input)
				srv.repo.On("GetCache", mock.Anything, key, &models.TeamMember{
					ID: 0,
				}).Return(false).Once()

				record := &models.TeamMember{ID: input}
				srv.repo.On("GetDetail", mock.Anything, dto.TeamMemberDetailReq{ID: input}).Return(record, nil).Once()
				srv.repo.On("CreateCache", mock.Anything, key, record, time.Minute).Return().Once()

				srv.repo.On("Delete", mock.Anything, &models.TeamMember{
					ID: input,
				}).Return(errors.New("invalid")).Once()
			},
			wantErr: true,
		},
		{
			name: "success",
			id:   srv.teamMember.ID,
			mockFunc: func(input uint64) {
				key := models.KeyCacheTeamMemberDetail(input)
				srv.repo.On("GetCache", mock.Anything, key, &models.TeamMember{
					ID: 0,
				}).Return(false).Once()

				srv.repo.On("GetDetail", mock.Anything, dto.TeamMemberDetailReq{ID: input}).Return(&srv.teamMember, nil).Once()
				srv.repo.On("CreateCache", mock.Anything, key, &srv.teamMember, time.Minute).Return().Once()

				srv.repo.On("Delete", mock.Anything, &models.TeamMember{ID: input}).Return(nil).Once()
				srv.repo.On("DeleteCache", mock.Anything, key).Return().Once()
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		srv.T().Run(tt.name, func(t *testing.T) {
			if tt.mockFunc != nil {
				tt.mockFunc(tt.id)
			}

			if err := srv.service.DeleteByID(srv.ctx, tt.id); (err != nil) != tt.wantErr {
				t.Errorf("TeamMemberSrv.DeleteByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (srv *TeamMemberServiceTestSuite) TestTeamMemberSrv_Update() {
	params := dto.TeamMemberUpdateReq{
		ID:             srv.teamMember.ID,
		Name:           srv.teamMember.Name,
		UsernameGithub: srv.teamMember.UsernameGithub,
		Email:          srv.teamMember.Email,
	}

	tests := []struct {
		name     string
		req      dto.TeamMemberUpdateReq
		mockFunc func(input dto.TeamMemberUpdateReq)
		wantErr  bool
	}{
		{
			name: "not found",
			req:  params,
			mockFunc: func(input dto.TeamMemberUpdateReq) {
				key := models.KeyCacheTeamMemberDetail(input.ID)
				srv.repo.On("GetCache", mock.Anything, key, &models.TeamMember{ID: 0}).Return(false).Once()
				srv.repo.On("GetDetail", mock.Anything, dto.TeamMemberDetailReq{ID: input.ID}).Return(nil, nil).Once()
			},
			wantErr: true,
		},
		{
			name: "email duplicate",
			req:  params,
			mockFunc: func(input dto.TeamMemberUpdateReq) {
				key := models.KeyCacheTeamMemberDetail(input.ID)
				srv.repo.On("GetCache", mock.Anything, key, &models.TeamMember{ID: 0}).Return(false).Once()
				srv.repo.On("GetDetail", mock.Anything, dto.TeamMemberDetailReq{ID: input.ID}).Return(&srv.teamMember, nil).Once()
				srv.repo.On("CreateCache", mock.Anything, key, &srv.teamMember, time.Minute).Return().Once()
				// duplicate email
				srv.repo.On("GetDetail", mock.Anything, dto.TeamMemberDetailReq{
					CustomColumn: "id",
					Email:        input.Email,
					NotID:        input.ID,
				}).Return(&models.TeamMember{ID: 101}, nil).Once()
			},
			wantErr: true,
		},
		{
			name: "failed update record",
			req:  params,
			mockFunc: func(input dto.TeamMemberUpdateReq) {
				key := models.KeyCacheTeamMemberDetail(input.ID)
				srv.repo.On("GetCache", mock.Anything, key, &models.TeamMember{ID: 0}).Return(false).Once()
				srv.repo.On("GetDetail", mock.Anything, dto.TeamMemberDetailReq{ID: input.ID}).Return(&srv.teamMember, nil).Once()
				srv.repo.On("CreateCache", mock.Anything, key, &srv.teamMember, time.Minute).Return().Once()
				// Check duplicate
				srv.repo.On("GetDetail", mock.Anything, dto.TeamMemberDetailReq{
					CustomColumn: "id",
					Email:        input.Email,
					NotID:        input.ID,
				}).Return(nil, nil).Once()
				srv.repo.On("GetDetail", mock.Anything, dto.TeamMemberDetailReq{
					CustomColumn:   "id",
					UsernameGithub: input.UsernameGithub,
					NotID:          input.ID,
				}).Return(nil, nil).Once()

				srv.repo.On("Update", mock.Anything, &models.TeamMember{
					ID:             input.ID,
					Name:           input.Name,
					Email:          input.Email,
					UsernameGithub: input.UsernameGithub,
				}).Return(errors.New("invalid")).Once()
			},
			wantErr: true,
		},
		{
			name: "success",
			req:  params,
			mockFunc: func(input dto.TeamMemberUpdateReq) {
				key := models.KeyCacheTeamMemberDetail(input.ID)
				srv.repo.On("GetCache", mock.Anything, key, &models.TeamMember{ID: 0}).Return(false).Once()
				srv.repo.On("GetDetail", mock.Anything, dto.TeamMemberDetailReq{ID: input.ID}).Return(&srv.teamMember, nil).Once()
				srv.repo.On("CreateCache", mock.Anything, key, &srv.teamMember, time.Minute).Return().Once()
				// Check duplicate
				srv.repo.On("GetDetail", mock.Anything, dto.TeamMemberDetailReq{
					CustomColumn: "id",
					Email:        input.Email,
					NotID:        input.ID,
				}).Return(nil, nil).Once()
				srv.repo.On("GetDetail", mock.Anything, dto.TeamMemberDetailReq{
					CustomColumn:   "id",
					UsernameGithub: input.UsernameGithub,
					NotID:          input.ID,
				}).Return(nil, nil).Once()

				srv.repo.On("Update", mock.Anything, &models.TeamMember{
					ID:             input.ID,
					Name:           input.Name,
					Email:          input.Email,
					UsernameGithub: input.UsernameGithub,
				}).Return(nil).Once()
				srv.repo.On("DeleteCache", mock.Anything, key).Return().Once()
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		srv.T().Run(tt.name, func(t *testing.T) {
			if tt.mockFunc != nil {
				tt.mockFunc(tt.req)
			}

			if err := srv.service.Update(srv.ctx, tt.req); (err != nil) != tt.wantErr {
				t.Errorf("TeamMemberSrv.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (srv *TeamMemberServiceTestSuite) TestTeamMemberSrv_GetList() {
	params := dto.TeamMemberListReq{
		Limit: 1,
		Page:  1,
	}

	tests := []struct {
		name     string
		req      dto.TeamMemberListReq
		mockFunc func(input dto.TeamMemberListReq)
		want     *response_mapper.Pagination
		wantErr  bool
	}{
		{
			name: "invalid params",
			req: dto.TeamMemberListReq{
				OrderBy: "invalid",
			},
			mockFunc: func(input dto.TeamMemberListReq) {
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "failed get records",
			req:  params,
			mockFunc: func(input dto.TeamMemberListReq) {
				srv.repo.On("GetList", mock.Anything, input).Return(nil, errors.New("invalid")).Once()
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success total records less than limit",
			req: dto.TeamMemberListReq{
				Limit: 10,
				Page:  1,
			},
			mockFunc: func(input dto.TeamMemberListReq) {
				srv.repo.On("GetList", mock.Anything, input).Return(srv.teamMembers, nil).Once()
			},
			want: &response_mapper.Pagination{
				Meta: response_mapper.Meta{
					Page:         1,
					Limit:        10,
					TotalRecords: len(srv.teamMembers),
				},
				Data: srv.teamMembers,
			},
			wantErr: false,
		},
		{
			name: "failed get total records",
			req:  params,
			mockFunc: func(input dto.TeamMemberListReq) {
				srv.repo.On("GetList", mock.Anything, input).Return([]models.TeamMember{srv.teamMembers[0]}, nil).Once()

				input.CustomColumns = "id"
				input.IsNotDefaultQuery = true
				input.Offset = (input.Page - 1) * input.Limit
				input.Limit = models.DefaultLimitIsTotalDataTrue * input.Limit
				srv.repo.On("GetList", mock.Anything, input).Return(nil, errors.New("invalid")).Once()
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success total records more than limit",
			req:  params,
			mockFunc: func(input dto.TeamMemberListReq) {
				srv.repo.On("GetList", mock.Anything, input).Return([]models.TeamMember{srv.teamMembers[0]}, nil).Once()

				input.CustomColumns = "id"
				input.IsNotDefaultQuery = true
				input.Offset = (input.Page - 1) * input.Limit
				input.Limit = models.DefaultLimitIsTotalDataTrue * input.Limit

				total := []models.TeamMember{}
				for i := 0; i < len(srv.teamMembers); i++ {
					val := models.TeamMember{ID: srv.teamMembers[i].ID}
					total = append(total, val)
				}
				srv.repo.On("GetList", mock.Anything, input).Return(total, nil).Once()
			},
			want: &response_mapper.Pagination{
				Meta: response_mapper.Meta{
					Page:         params.Page,
					Limit:        params.Limit,
					TotalRecords: len(srv.teamMembers),
				},
				Data: []models.TeamMember{srv.teamMembers[0]},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		srv.T().Run(tt.name, func(t *testing.T) {
			if tt.mockFunc != nil {
				tt.mockFunc(tt.req)
			}

			got, err := srv.service.GetList(srv.ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("TeamMemberSrv.GetList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TeamMemberSrv.GetList() = %v, want %v", got, tt.want)
			}
		})
	}
}
