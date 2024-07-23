package models

import "testing"

func TestKeyCacheTeamMemberDetail(t *testing.T) {
	type args struct {
		id uint64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "success",
			args: args{
				id: 1,
			},
			want: "team_member_detail_1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KeyCacheTeamMemberDetail(tt.args.id); got != tt.want {
				t.Errorf("KeyCacheTeamMemberDetail() = %v, want %v", got, tt.want)
			}
		})
	}
}
