package handlers

import (
	"context"
	"database/sql"
	"net/http"

	"NKSS-backend/pkg/query"
)

type Group struct {
	Name        string                 `json:"name"`
	Alias       string                 `json:"alias"`
	Faculty     []Faculty              `json:"faculty"`
	Branch      string                 `json:"branch"`
	Kind        string                 `json:"kind"`
	Description string                 `json:"description"`
	Socials     map[string]interface{} `json:"socials"`
	Admins      []Admin                `json:"admins"`
	Members     []int64                `json:"members"`
}

type Faculty struct {
	Name   string `json:"name"`
	Mobile int64  `json:"mobile"`
}

type Admin struct {
	Position   string `json:"position"`
	RollNumber int64  `json:"roll_number"`
}

func GetGroups(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	queries := query.New(db)
	group_rows, err := queries.GetAllGroups(ctx)
	if err == sql.ErrNoRows {
		RespondError(w, 404, "No groups found!")
		return
	}
	if err != nil {
		RespondError(w, 500, "Something went wrong while fetching details from our database")
		return
	}

	var groups []Group
	for _, group_row := range group_rows {
		var group Group

		group.Name = group_row.Name
		group.Alias = group_row.Alias.String
		group.Branch = group_row.Branch.String
		group.Kind = group_row.Kind
		group.Description = group_row.Description.String
		group.Members = group_row.Members

		var faculties []Faculty
		for i, name := range group_row.FacultyNames {
			faculties = append(faculties, Faculty{Name: name, Mobile: group_row.FacultyMobiles[i]})
		}
		group.Faculty = faculties

		socials := make(map[string]interface{})
		for i, social_type := range group_row.SocialTypes {
			socials[social_type] = group_row.SocialLinks[i]
		}
		discord := make(map[string]interface{})
		discord["id"] = group_row.ServerID
		discord["invite"] = group_row.ServerInvite
		discord["roles"] = map[string]int64{
			"freshman":  group_row.FresherRole.Int64,
			"sophomore": group_row.SophomoreRole.Int64,
			"junior":    group_row.JuniorRole.Int64,
			"senior":    group_row.SeniorRole.Int64,
		}
		socials["discord"] = discord
		group.Socials = socials

		var admins []Admin
		for i, position := range group_row.AdminPositions {
			admins = append(admins, Admin{Position: position, RollNumber: group_row.AdminRolls[i]})
		}
		group.Admins = admins

		groups = append(groups, group)
	}

	RespondJSON(w, 200, groups)
}
