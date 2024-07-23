package dto

import (
	help "github.com/adamnasrudin03/go-helpers"
	response_mapper "github.com/adamnasrudin03/go-helpers/response-mapper/v1"
	"github.com/adamnasrudin03/go-skeleton-chi/app/models"
)

type TeamMemberDetailReq struct {
	ID             uint64 `json:"id"`
	Name           string `json:"name"`
	UsernameGithub string `json:"username_github"`
	Email          string `json:"email"`
	CustomColumn   string `json:"custom_column"`
	NotID          uint64 `json:"not_id"`
}

type TeamMemberCreateReq struct {
	Name           string `json:"name" validate:"required"`
	UsernameGithub string `json:"username_github" validate:"required"`
	Email          string `json:"email" validate:"required,email"`
}

type TeamMemberUpdateReq struct {
	ID             uint64 `json:"id" validate:"min=1"`
	Name           string `json:"name" validate:"required"`
	UsernameGithub string `json:"username_github" validate:"required"`
	Email          string `json:"email" validate:"required,email"`
}

type TeamMemberListReq struct {
	Search            string `json:"search" query:"search"`
	Limit             int    `json:"limit" query:"limit"`
	Offset            int    `json:"offset" query:"offset"`
	Page              int    `json:"page" query:"page"`
	OrderBy           string `json:"order_by" query:"order_by"`
	SortBy            string `json:"sort_by" query:"sort_by"`
	IsNoLimit         bool   `json:"is_no_limit" query:"is_no_limit"`
	IsNotDefaultQuery bool   `json:"is_not_default_query" query:"is_not_default_query"`
	CustomColumns     string `json:"custom_columns" query:"custom_columns"`
}

func (m *TeamMemberListReq) Validate() error {
	if m.Page <= 0 {
		m.Page = 1
	}

	if m.Limit <= 0 {
		m.Limit = 10
	}

	m.Search = help.ToLower(m.Search)

	m.OrderBy = help.ToUpper(m.OrderBy)
	if !models.IsValidOrderBy[m.OrderBy] && m.OrderBy != "" {
		return response_mapper.ErrInvalidFormat("order_by", "order_by")
	}

	m.SortBy = help.ToLower(m.SortBy)
	if m.OrderBy != "" && m.SortBy == "" {
		return response_mapper.ErrIsRequired("sort_by", "sort_by")
	}

	return nil
}

func (c *TeamMemberListReq) DefaultQuery() TeamMemberListReq {
	if c.Limit <= 0 {
		c.Limit = 10
	}

	if c.Page <= 0 {
		c.Page = 1
	}

	if c.Page > 0 {
		c.Offset = (c.Page - 1) * c.Limit
	}

	return *c
}
