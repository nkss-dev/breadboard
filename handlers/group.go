package handlers

import (
	"context"
	"database/sql"
	"net/http"

	"nkssbackend/query"

	"github.com/gorilla/mux"
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

// ConstructGroup translates the row returned by sqlc into
// the struct Group for a better strucutre
func ConstructGroup(raw_group query.GetGroupRow) (group Group) {
	group.Name = raw_group.Name
	group.Alias = raw_group.Alias.String
	group.Branch = raw_group.Branch.String
	group.Kind = raw_group.Kind
	group.Description = raw_group.Description.String
	group.Members = raw_group.Members

	var faculties []Faculty
	for i, name := range raw_group.FacultyNames {
		faculties = append(faculties, Faculty{Name: name, Mobile: raw_group.FacultyMobiles[i]})
	}
	group.Faculty = faculties

	socials := make(map[string]interface{})
	for i, social_type := range raw_group.SocialTypes {
		socials[social_type] = raw_group.SocialLinks[i]
	}
	discord := make(map[string]interface{})
	discord["id"] = raw_group.ServerID
	discord["invite"] = raw_group.ServerInvite
	discord["roles"] = map[string]int64{
		"freshman":  raw_group.FresherRole.Int64,
		"sophomore": raw_group.SophomoreRole.Int64,
		"junior":    raw_group.JuniorRole.Int64,
		"senior":    raw_group.SeniorRole.Int64,
	}
	socials["discord"] = discord
	group.Socials = socials

	var admins []Admin
	for i, position := range raw_group.AdminPositions {
		admins = append(admins, Admin{Position: position, RollNumber: raw_group.AdminRolls[i]})
	}
	group.Admins = admins

	return group
}

// GetGroup returns a handler to return a group's details
// based on the unique parameter passed.
//
// This handler takes in a name argument which is first
// checked as an alias and then as the name of a group.
func GetGroup(db *sql.DB) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(db)
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		group_row, err := queries.GetGroup(ctx, vars["name"])
		if err == sql.ErrNoRows {
			RespondError(w, 404, "No groups found!")
			return
		}
		if err != nil {
			RespondError(w, 500, "Something went wrong while fetching details from our database")
			return
		}

		group := ConstructGroup(group_row)
		RespondJSON(w, 200, group)
	}
}

// GetGroups retrieves the group details from the database.
func GetGroups(db *sql.DB) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(db)
	return func(w http.ResponseWriter, r *http.Request) {
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
			groups = append(groups, ConstructGroup(query.GetGroupRow(group_row)))
		}

		RespondJSON(w, 200, groups)
	}
}

// GetGroupAdmins retrieves the admins of a group from the database.
func GetGroupAdmins(db *sql.DB) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(db)
	return func(w http.ResponseWriter, r *http.Request) {
		group_name := mux.Vars(r)["name"]
		admins, err := queries.GetGroupAdmins(ctx, group_name)
		if err != nil {
			RespondError(w, 500, "Something went wrong while fetching details from our database")
			return
		}

		RespondJSON(w, 200, admins)
	}
}

// GetGroupFaculty retrieves the management faculty of a group from the database.
func GetGroupFaculty(db *sql.DB) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(db)
	return func(w http.ResponseWriter, r *http.Request) {
		group_name := mux.Vars(r)["name"]
		faculty, err := queries.GetGroupFaculty(ctx, group_name)
		if err != nil {
			RespondError(w, 500, "Something went wrong while fetching details from our database")
			return
		}

		RespondJSON(w, 200, faculty)
	}
}

// GetGroupMembers retrieves the members of a group from the database.
func GetGroupMembers(db *sql.DB) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(db)
	return func(w http.ResponseWriter, r *http.Request) {
		group_name := mux.Vars(r)["name"]
		members, err := queries.GetGroupMembers(ctx, group_name)
		if err != nil {
			RespondError(w, 500, "Something went wrong while fetching details from our database")
			return
		}

		RespondJSON(w, 200, members)
	}
}

// GetGroupSocials retrieves the social media links of a group from the database.
func GetGroupSocials(db *sql.DB) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(db)
	return func(w http.ResponseWriter, r *http.Request) {
		group_name := mux.Vars(r)["name"]
		socials, err := queries.GetGroupSocials(ctx, group_name)
		if err != nil {
			RespondError(w, 500, "Something went wrong while fetching details from our database")
			return
		}

		RespondJSON(w, 200, socials)
	}
}
