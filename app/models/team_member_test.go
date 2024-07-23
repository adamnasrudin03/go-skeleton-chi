package models

import "testing"

func TestTeamMember_TableName(t *testing.T) {
	tests := []struct {
		name string
		tr   TeamMember
		want string
	}{
		{
			name: "success",
			tr:   TeamMember{},
			want: "team_members",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.TableName(); got != tt.want {
				t.Errorf("TeamMember.TableName() = %v, want %v", got, tt.want)
			}
		})
	}
}
