package service

import (
	"context"

	response_mapper "github.com/adamnasrudin03/go-helpers/response-mapper/v1"
	"github.com/adamnasrudin03/go-skeleton-chi/app/dto"
)

func (s TeamMemberSrv) checkDuplicate(ctx context.Context, req dto.TeamMemberDetailReq) error {
	var (
		opName = "TeamMemberService-checkDuplicate"
		err    error
	)
	detail, err := s.Repo.GetDetail(ctx, dto.TeamMemberDetailReq{
		CustomColumn: "id",
		Email:        req.Email,
		NotID:        req.NotID,
	})
	if err != nil {
		s.Logger.Errorf("%s, failed check duplicate email: %v", opName, err)
		return response_mapper.ErrDB()
	}

	if detail != nil && detail.ID > 0 {
		return response_mapper.ErrIsDuplicate("email", "email")
	}

	detail, err = s.Repo.GetDetail(ctx, dto.TeamMemberDetailReq{
		CustomColumn:   "id",
		UsernameGithub: req.UsernameGithub,
		NotID:          req.NotID,
	})
	if err != nil {
		s.Logger.Errorf("%s, failed check duplicate username_github: %v", opName, err)
		return response_mapper.ErrDB()
	}

	if detail != nil && detail.ID > 0 {
		return response_mapper.ErrIsDuplicate("username_github", "username_github")
	}

	return nil

}
