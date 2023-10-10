package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	query "breadboard/.sqlc-auto-gen"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
func CreateClubFaculty(conn *pgxpool.Pool) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(conn)
	return func(w http.ResponseWriter, r *http.Request) {
		clubName := mux.Vars(r)["name"]
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
			ClubName: clubName,
			EmpID:    int32(mobile),
		}
		err = queries.CreateClubFaculty(ctx, params)
		if err != nil {
			log.Println(err)
			RespondError(w, 500, "Something went wrong while inserting details to our database")
			return
		}

		RespondJSON(w, 200, "Added "+name+" as a faculty incharge of "+clubName+" successfully!")
	}
}

// CreateClubMember adds a new member to a club.
func CreateClubMember(conn *pgxpool.Pool) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(conn)

	return func(w http.ResponseWriter, r *http.Request) {
		var clubMember query.CreateClubMemberParams
		json.NewDecoder(r.Body).Decode(&clubMember)

		clubMember.ClubName = mux.Vars(r)["name"]
		clubMember.RollNumber = mux.Vars(r)["roll"]

		// TODO: add parameter validation middleware.

		err := queries.CreateClubMember(ctx, clubMember)
		if err != nil {
			log.Println(err)
			RespondError(w, 500, "Something went wrong while inserting details to our database")
			return
		}

		RespondJSON(w, 200, "Successfully added '"+clubMember.RollNumber+"' to "+clubMember.ClubName+" as: "+clubMember.Position)
	}
}

// CreateClubMemberBulk adds many new members to a club.
func CreateClubMemberBulk(conn *pgxpool.Pool) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(conn)

	return func(w http.ResponseWriter, r *http.Request) {
		var clubMembers []query.CreateClubMemberBulkParams
		json.NewDecoder(r.Body).Decode(&clubMembers)

		var rollNumbers []string
		for index, clubMember := range clubMembers {
			clubMembers[index].ClubName = mux.Vars(r)["name"]
			rollNumbers = append(rollNumbers, clubMember.RollNumber)
		}

		// TODO: add parameter validation middleware.

		_, err := queries.CreateClubMemberBulk(ctx, clubMembers)
		if err != nil {
			log.Println(err)
			RespondError(w, 500, "Something went wrong while inserting details to our database")
			return
		}

		RespondJSON(w, 200, "Sucessfully added the following students to "+mux.Vars(r)["name"]+": "+strings.Join(rollNumbers, "\n"))
	}
}

// CreateClubSocial adds a new social media handle of a group.
func CreateClubSocial(conn *pgxpool.Pool) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(conn)
	return func(w http.ResponseWriter, r *http.Request) {
		clubName := mux.Vars(r)["name"]
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
			ClubName:     clubName,
			PlatformType: platform_type,
			Link:         link,
		}
		err := queries.CreateClubSocial(ctx, params)
		if err != nil {
			log.Println(err)
			RespondError(w, 500, "Something went wrong while inserting details to our database")
			return
		}

		RespondJSON(w, 200, "Added link of the "+platform_type+" handle for "+clubName+" successfully!")
	}
}

// DeleteClubFaculty deletes an existing faculty incharge of a group.
func DeleteClubFaculty(conn *pgxpool.Pool) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(conn)
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			RespondError(w, 400, "Id paramter must only contain digits")
			return
		}

		params := query.DeleteClubFacultyParams{
			ClubName: vars["name"],
			EmpID:    int32(id),
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
func DeleteClubMember(conn *pgxpool.Pool) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(conn)

	type DeleteClubMemberParams struct {
		ClubName   string `json:"club_name"`
		RollNumber string `json:"roll_number"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		clubMember := DeleteClubMemberParams{
			ClubName:   mux.Vars(r)["name"],
			RollNumber: mux.Vars(r)["roll"],
		}

		// TODO: add parameter validation middleware.

		params := query.DeleteClubMemberParams{
			ClubName:   clubMember.ClubName,
			RollNumber: clubMember.RollNumber,
		}
		err := queries.DeleteClubMember(ctx, params)
		if err != nil {
			log.Println(err)
			RespondError(w, 500, "Something went wrong while deleting details from our database")
			return
		}

		RespondJSON(w, 200, "Successfully removed '"+clubMember.RollNumber+"' from "+clubMember.ClubName)
	}
}

// DeleteClubMemberBulk deletes existing members of a club.
func DeleteClubMemberBulk(conn *pgxpool.Pool) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(conn)

	return func(w http.ResponseWriter, r *http.Request) {
		var rollNumbers []string
		json.NewDecoder(r.Body).Decode(&rollNumbers)

		// TODO: add parameter validation middleware.

		params := query.DeleteClubMemberBulkParams{
			ClubName:    mux.Vars(r)["name"],
			RollNumbers: rollNumbers,
		}
		err := queries.DeleteClubMemberBulk(ctx, params)
		if err != nil {
			log.Println(err)
			RespondError(w, 500, "Something went wrong while deleting details from our database")
			return
		}

		RespondJSON(w, 200, "Successfully removed the following roll numbers from "+params.ClubName+": "+strings.Join(rollNumbers, "\n"))
	}
}

// DeleteClubSocial deletes an existing social media handle of a group.
func DeleteClubSocial(conn *pgxpool.Pool) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(conn)
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		params := query.DeleteClubSocialParams{
			ClubName:     vars["name"],
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
func GetClub(conn *pgxpool.Pool) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(conn)
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		group, err := queries.GetClub(ctx, vars["name"])
		if err == pgx.ErrNoRows {
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
func GetClubs(conn *pgxpool.Pool) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(conn)
	return func(w http.ResponseWriter, r *http.Request) {
		groups, err := queries.GetClubs(ctx)
		if err == pgx.ErrNoRows {
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
func GetClubFaculty(conn *pgxpool.Pool) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(conn)
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
func GetClubSocials(conn *pgxpool.Pool) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(conn)
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
func ReadClubMembers(conn *pgxpool.Pool) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(conn)

	type ReadClubMemberParams struct {
		ClubName string `json:"club_name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		clubMember := ReadClubMemberParams{
			ClubName: mux.Vars(r)["name"],
		}

		// TODO: add parameter validation middleware.

		members, err := queries.ReadClubMembers(ctx, clubMember.ClubName)
		if err != nil {
			RespondError(w, 500, "Something went wrong while fetching details from our database")
			return
		}

		RespondJSON(w, 200, members)
	}
}

// UpdateClubMember updates a club member's details.
func UpdateClubMember(conn *pgxpool.Pool) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(conn)

	return func(w http.ResponseWriter, r *http.Request) {
		var clubMember query.UpdateClubMemberParams
		json.NewDecoder(r.Body).Decode(&clubMember)

		clubMember.ClubName = mux.Vars(r)["name"]
		clubMember.RollNumber = mux.Vars(r)["roll"]

		// TODO: add parameter validation middleware.

		err := queries.UpdateClubMember(ctx, clubMember)
		if err != nil {
			log.Println(err)
			RespondError(w, 500, "Something went wrong while updating details to our database")
			return
		}

		RespondJSON(w, 200, "Successfully updated '"+clubMember.RollNumber+"' to "+clubMember.ClubName+" as: "+clubMember.Position)
	}
}

// UpdateClubSocials updates the link of a social media handle for a group.
func UpdateClubSocials(conn *pgxpool.Pool) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(conn)
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
			ClubName:     vars["name"],
		}
		err := queries.UpdateClubSocials(ctx, params)
		if err != nil {
			RespondError(w, 500, "Something went wrong while updating details to our database")
			return
		}

		RespondJSON(w, 200, "Social media link updated successfully")
	}
}
