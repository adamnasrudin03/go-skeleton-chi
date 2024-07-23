package models

import "fmt"

func KeyCacheTeamMemberDetail(id uint64) string {
	return fmt.Sprintf("team_member_detail_%d", id)
}
