package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	query "breadboard/.sqlc-auto-gen"

	"github.com/gorilla/mux"
)

type Club struct {
	Name        string                 `json:"name"`
	Alias       string                 `json:"alias"`
	Faculty     []Faculty              `json:"faculty"`
	Branch      []string               `json:"branch"`
	Kind        string                 `json:"kind"`
	Description string                 `json:"description"`
	Socials     map[string]interface{} `json:"socials"`
	Members     []int64                `json:"members"`
}

type Faculty struct {
	Name   string `json:"name"`
	Mobile string `json:"mobile"`
}

// CreateClubFaculty creates a new faculty incharge for a group.
func CreateClubFaculty(db *sql.DB) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(db)
	return func(w http.ResponseWriter, r *http.Request) {
		group_name := mux.Vars(r)["name"]
		vars := r.URL.Query()

		name := vars.Get("name")
		if name == "" {
			RespondError(w, 400, "Required query param, name, missing")
			return
		}

		mobileStr := vars.Get("mobile")
		if mobileStr == "" {
			RespondError(w, 400, "Required query param, mobile, missing")
			return
		}
		mobile, err := strconv.Atoi(mobileStr)
		if err != nil {
			RespondError(w, 400, "Mobile paramter must only contain digits")
			return
		}

		params := query.CreateClubFacultyParams{
			Name:  group_name,
			EmpID: int32(mobile),
		}
		err = queries.CreateClubFaculty(ctx, params)
		if err != nil {
			log.Println(err)
			RespondError(w, 500, "Something went wrong while inserting details to our database")
			return
		}

		RespondJSON(w, 200, "Added "+name+" as a faculty incharge of "+group_name+" successfully!")
	}
}

// CreateClubMember adds a new member to a club.
func CreateClubMember(db *sql.DB) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(db)

	type CreateClubMemberParams struct {
		ClubNameOrAlias string   `json:"club_name_or_alias"`
		RollNumber      string   `json:"roll_number"`
		Position        string   `json:"position"`
		ExtraGroups     []string `json:"extra_groups"`
		Comments        string   `json:"comments"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var clubMember CreateClubMemberParams
		json.NewDecoder(r.Body).Decode(&clubMember)

		clubMember.ClubNameOrAlias = mux.Vars(r)["name"]

		// TODO: add parameter validation middleware.

		params := query.CreateClubMemberParams{
			ClubNameOrAlias: clubMember.ClubNameOrAlias,
			RollNumber:      clubMember.RollNumber,
			Position:        clubMember.Position,
			ExtraGroups:     clubMember.ExtraGroups,
			Comments:        sql.NullString{String: clubMember.Comments, Valid: true},
		}
		err := queries.CreateClubMember(ctx, params)
		if err != nil {
			log.Println(err)
			RespondError(w, 500, "Something went wrong while inserting details to our database")
			return
		}

		RespondJSON(w, 200, "Successfully added '"+clubMember.RollNumber+"' to "+clubMember.ClubNameOrAlias+" as: "+clubMember.Position)
	}
}

// CreateClubSocial adds a new social media handle of a group.
func CreateClubSocial(db *sql.DB) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(db)
	return func(w http.ResponseWriter, r *http.Request) {
		group_name := mux.Vars(r)["name"]
		vars := r.URL.Query()

		platform_type := vars.Get("type")
		if platform_type == "" {
			RespondError(w, 400, "Required query param, type, missing")
			return
		}

		link := vars.Get("link")
		if link == "" {
			RespondError(w, 400, "Required query param, link, missing")
			return
		}

		params := query.CreateClubSocialParams{
			Name:         group_name,
			PlatformType: platform_type,
			Link:         link,
		}
		err := queries.CreateClubSocial(ctx, params)
		if err != nil {
			log.Println(err)
			RespondError(w, 500, "Something went wrong while inserting details to our database")
			return
		}

		RespondJSON(w, 200, "Added link of the "+platform_type+" handle for "+group_name+" successfully!")
	}
}

// DeleteClubFaculty deletes an existing faculty incharge of a group.
func DeleteClubFaculty(db *sql.DB) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(db)
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			RespondError(w, 400, "Id paramter must only contain digits")
			return
		}

		params := query.DeleteClubFacultyParams{
			Name:  vars["name"],
			EmpID: int32(id),
		}
		err = queries.DeleteClubFaculty(ctx, params)
		if err != nil {
			log.Println(err)
			RespondError(w, 500, "Something went wrong while deleting details from our database")
			return
		}

		RespondJSON(w, 200, "Removed "+vars["id"]+" as a faculty incharge of "+vars["name"]+" successfully!")
	}
}

// DeleteClubMember deletes an existing member of a club.
func DeleteClubMember(db *sql.DB) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(db)

	type DeleteClubMemberParams struct {
		ClubNameOrAlias string `json:"club_name_or_alias"`
		RollNumber      string `json:"roll_number"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		clubMember := DeleteClubMemberParams{
			ClubNameOrAlias: mux.Vars(r)["name"],
			RollNumber:      mux.Vars(r)["roll"],
		}

		// TODO: add parameter validation middleware.

		params := query.DeleteClubMemberParams{
			ClubNameOrAlias: clubMember.ClubNameOrAlias,
			RollNumber:      clubMember.RollNumber,
		}
		err := queries.DeleteClubMember(ctx, params)
		if err != nil {
			log.Println(err)
			RespondError(w, 500, "Something went wrong while deleting details from our database")
			return
		}

		RespondJSON(w, 200, "Successfully removed '"+clubMember.RollNumber+"' from "+clubMember.ClubNameOrAlias)
	}
}

// DeleteClubSocial deletes an existing social media handle of a group.
func DeleteClubSocial(db *sql.DB) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(db)
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		params := query.DeleteClubSocialParams{
			Name:         vars["name"],
			PlatformType: vars["type"],
		}
		err := queries.DeleteClubSocial(ctx, params)
		if err != nil {
			log.Println(err)
			RespondError(w, 500, "Something went wrong while deleting details from our database")
			return
		}

		RespondJSON(w, 200, "Removed '"+vars["type"]+"' as one of the social media handles of "+vars["name"]+" successfully!")
	}
}

// GetClub returns a handler to return a group's details
// based on the unique parameter passed.
//
// This handler takes in a name argument which is first
// checked as an alias and then as the name of a group.
func GetClub(db *sql.DB) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(db)
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		group, err := queries.GetClub(ctx, vars["name"])
		if err == sql.ErrNoRows {
			RespondError(w, 404, "No groups found!")
			return
		}
		if err != nil {
			fmt.Println(err)
			RespondError(w, 500, "Something went wrong while fetching details from our database")
			return
		}

		RespondJSON(w, 200, group)
	}
}

// GetClubs retrieves the group details from the database.
func GetClubs(db *sql.DB) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(db)
	return func(w http.ResponseWriter, r *http.Request) {
		groups, err := queries.GetClubs(ctx)
		if err == sql.ErrNoRows {
			RespondError(w, 404, "No groups found!")
			return
		}
		if err != nil {
			RespondError(w, 500, "Something went wrong while fetching details from our database")
			return
		}

		RespondJSON(w, 200, groups)
	}
}

// GetClubFaculty retrieves the management faculty of a group from the database.
func GetClubFaculty(db *sql.DB) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(db)
	return func(w http.ResponseWriter, r *http.Request) {
		group_name := mux.Vars(r)["name"]
		faculty, err := queries.GetClubFaculty(ctx, group_name)
		if err != nil {
			RespondError(w, 500, "Something went wrong while fetching details from our database")
			return
		}

		RespondJSON(w, 200, faculty)
	}
}

// GetClubSocials retrieves the social media links of a group from the database.
func GetClubSocials(db *sql.DB) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(db)
	return func(w http.ResponseWriter, r *http.Request) {
		group_name := mux.Vars(r)["name"]
		socials, err := queries.GetClubSocials(ctx, group_name)
		if err != nil {
			RespondError(w, 500, "Something went wrong while fetching details from our database")
			return
		}

		RespondJSON(w, 200, socials)
	}
}

// ReadClubMembers retrieves the members of a club.
func ReadClubMembers(db *sql.DB) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(db)

	type ReadClubMemberParams struct {
		ClubNameOrAlias string `json:"club_name_or_alias"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		clubMember := ReadClubMemberParams{
			ClubNameOrAlias: mux.Vars(r)["name"],
		}

		// TODO: add parameter validation middleware.

		members, err := queries.ReadClubMembers(ctx, clubMember.ClubNameOrAlias)
		if err != nil {
			RespondError(w, 500, "Something went wrong while fetching details from our database")
			return
		}

		RespondJSON(w, 200, members)
	}
}

// UpdateClubMember updates a club member's details.
func UpdateClubMember(db *sql.DB) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(db)

	type UpdateClubMemberParams struct {
		ClubNameOrAlias string   `json:"club_name_or_alias"`
		RollNumber      string   `json:"roll_number"`
		Position        string   `json:"position"`
		ExtraGroups     []string `json:"extra_groups"`
		Comments        string   `json:"comments"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var clubMember UpdateClubMemberParams
		json.NewDecoder(r.Body).Decode(&clubMember)

		clubMember.ClubNameOrAlias = mux.Vars(r)["name"]

		// TODO: add parameter validation middleware.

		params := query.UpdateClubMemberParams{
			ClubNameOrAlias: clubMember.ClubNameOrAlias,
			RollNumber:      clubMember.RollNumber,
			Position:        clubMember.Position,
			ExtraGroups:     clubMember.ExtraGroups,
			Comments:        sql.NullString{String: clubMember.Comments, Valid: true},
		}
		err := queries.UpdateClubMember(ctx, params)
		if err != nil {
			log.Println(err)
			RespondError(w, 500, "Something went wrong while updating details to our database")
			return
		}

		RespondJSON(w, 200, "Successfully updated '"+clubMember.RollNumber+"' to "+clubMember.ClubNameOrAlias+" as: "+clubMember.Position)
	}
}

// UpdateClubSocials updates the link of a social media handle for a group.
func UpdateClubSocials(db *sql.DB) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(db)
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		link := r.URL.Query().Get("link")
		if link == "" {
			RespondError(w, 400, "Required query param, link, missing")
			return
		}

		params := query.UpdateClubSocialsParams{
			PlatformType: vars["type"],
			Link:         link,
			Name:         vars["name"],
		}
		err := queries.UpdateClubSocials(ctx, params)
		if err != nil {
			RespondError(w, 500, "Something went wrong while updating details to our database")
			return
		}

		RespondJSON(w, 200, "Social media link updated successfully")
	}
}
