package dto

import (
	"reflect"
	"testing"

	"github.com/adamnasrudin03/go-skeleton-chi/app/models"
)

func TestTeamMemberListReq_Validate(t *testing.T) {
	tests := []struct {
		name    string
		m       *TeamMemberListReq
		wantErr bool
	}{
		{
			name: "invalid order by",
			m: &TeamMemberListReq{
				OrderBy: "invalid",
				SortBy:  "",
			},
			wantErr: true,
		},
		{
			name: "sort by required if order by provided",
			m: &TeamMemberListReq{
				OrderBy: models.OrderByASC,
				SortBy:  "",
			},
			wantErr: true,
		},
		{
			name:    "success",
			m:       &TeamMemberListReq{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("TeamMemberListReq.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTeamMemberListReq_DefaultQuery(t *testing.T) {
	tests := []struct {
		name string
		c    *TeamMemberListReq
		want TeamMemberListReq
	}{
		{
			name: "success",
			c:    &TeamMemberListReq{},
			want: TeamMemberListReq{
				Limit:  10,
				Offset: 0,
				Page:   1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.DefaultQuery(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TeamMemberListReq.DefaultQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
