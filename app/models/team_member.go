package models

// TeamMember represents the model for an TeamMembers
type TeamMember struct {
	ID             uint64 `json:"id" gorm:"primaryKey"`
	Name           string `json:"name" gorm:"not null"`
	UsernameGithub string `json:"username_github" gorm:"not null;uniqueIndex"`
	Email          string `json:"email" gorm:"not null;uniqueIndex"`
	DefaultModel
}

func (TeamMember) TableName() string {
	return "team_members"
}
