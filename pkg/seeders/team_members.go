package seeders

import (
	"github.com/adamnasrudin03/go-skeleton-chi/app/models"
	"gorm.io/gorm"
)

func InitTeamMembers(db *gorm.DB) {
	tx := db.Begin()
	var teamMembers = []models.TeamMember{}
	tx.Select("id").Limit(1).Find(&teamMembers)
	if len(teamMembers) == 0 {
		teamMembers = []models.TeamMember{
			{
				Name:           "Adam Nasrudin",
				UsernameGithub: "adamnasrudin03",
				Email:          "adamnasrudin@example.com",
			},
		}
		tx.Create(&teamMembers)
	}

	tx.Commit()
}
