package service

import (
	"errors"
	"testing"

	"github.com/adamnasrudin03/go-skeleton-chi/app/dto"
	"github.com/adamnasrudin03/go-skeleton-chi/app/models"
	"github.com/stretchr/testify/mock"
)

func (srv *TeamMemberServiceTestSuite) TestTeamMemberSrv_checkDuplicate() {
	params := dto.TeamMemberDetailReq{
		UsernameGithub: srv.teamMember.UsernameGithub,
		Email:          srv.teamMember.Email,
		NotID:          srv.teamMember.ID,
	}

	tests := []struct {
		name     string
		req      dto.TeamMemberDetailReq
		mockFunc func(input dto.TeamMemberDetailReq)
		wantErr  bool
	}{
		{
			name: "failed check email",
			req:  params,
			mockFunc: func(input dto.TeamMemberDetailReq) {
				srv.repo.On("GetDetail", mock.Anything, dto.TeamMemberDetailReq{
					CustomColumn: "id",
					Email:        input.Email,
					NotID:        input.NotID,
				}).Return(nil, errors.New("invalid check email")).Once()
			},
			wantErr: true,
		},
		{
			name: "duplicate email",
			req:  params,
			mockFunc: func(input dto.TeamMemberDetailReq) {
				srv.repo.On("GetDetail", mock.Anything, dto.TeamMemberDetailReq{
					CustomColumn: "id",
					Email:        input.Email,
					NotID:        input.NotID,
				}).Return(&models.TeamMember{ID: 101}, nil).Once()
			},
			wantErr: true,
		},
		{
			name: "failed check username_github",
			req:  params,
			mockFunc: func(input dto.TeamMemberDetailReq) {
				srv.repo.On("GetDetail", mock.Anything, dto.TeamMemberDetailReq{
					CustomColumn: "id",
					Email:        input.Email,
					NotID:        input.NotID,
				}).Return(nil, nil).Once()
				srv.repo.On("GetDetail", mock.Anything, dto.TeamMemberDetailReq{
					CustomColumn:   "id",
					UsernameGithub: input.UsernameGithub,
					NotID:          input.NotID,
				}).Return(nil, errors.New("invalid check username_github")).Once()
			},
			wantErr: true,
		},
		{
			name: "duplicate username_github",
			req:  params,
			mockFunc: func(input dto.TeamMemberDetailReq) {
				srv.repo.On("GetDetail", mock.Anything, dto.TeamMemberDetailReq{
					CustomColumn: "id",
					Email:        input.Email,
					NotID:        input.NotID,
				}).Return(nil, nil).Once()
				srv.repo.On("GetDetail", mock.Anything, dto.TeamMemberDetailReq{
					CustomColumn:   "id",
					UsernameGithub: input.UsernameGithub,
					NotID:          input.NotID,
				}).Return(&models.TeamMember{ID: 101}, nil).Once()
			},
			wantErr: true,
		},
		{
			name: "success",
			req:  params,
			mockFunc: func(input dto.TeamMemberDetailReq) {
				srv.repo.On("GetDetail", mock.Anything, dto.TeamMemberDetailReq{
					CustomColumn: "id",
					Email:        input.Email,
					NotID:        input.NotID,
				}).Return(nil, nil).Once()
				srv.repo.On("GetDetail", mock.Anything, dto.TeamMemberDetailReq{
					CustomColumn:   "id",
					UsernameGithub: input.UsernameGithub,
					NotID:          input.NotID,
				}).Return(nil, nil).Once()
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		srv.T().Run(tt.name, func(t *testing.T) {
			if tt.mockFunc != nil {
				tt.mockFunc(tt.req)
			}

			if err := srv.service.checkDuplicate(srv.ctx, tt.req); (err != nil) != tt.wantErr {
				t.Errorf("TeamMemberSrv.checkDuplicate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
