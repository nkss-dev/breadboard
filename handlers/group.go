package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

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

// CreateGroupAdmin creates a new admin for a group.
func CreateGroupAdmin(db *sql.DB) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(db)
	return func(w http.ResponseWriter, r *http.Request) {
		group_name := mux.Vars(r)["name"]
		vars := r.URL.Query()

		position := vars.Get("position")
		if position == "" {
			RespondError(w, 400, "Required query param, position, missing")
			return
		}

		rollStr := vars.Get("roll")
		if rollStr == "" {
			RespondError(w, 400, "Required query param, roll, missing")
			return
		}
		roll, err := strconv.Atoi(rollStr)
		if err != nil {
			RespondError(w, 400, "Roll paramter must only contain digits")
			return
		}

		params := query.CreateGroupAdminParams{
			Name:       group_name,
			Position:   position,
			RollNumber: int32(roll),
		}
		err = queries.CreateGroupAdmin(ctx, params)
		if err != nil {
			log.Println(err)
			RespondError(w, 500, "Something went wrong while inserting details to our database")
			return
		}

		RespondJSON(w, 200, "Added '"+fmt.Sprint(roll)+"' as the "+position+" of "+group_name+" successfully!")
	}
}

// CreateGroupFaculty creates a new faculty incharge for a group.
func CreateGroupFaculty(db *sql.DB) http.HandlerFunc {
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

		params := query.CreateGroupFacultyParams{
			Name_2: name,
			Mobile: int64(mobile),
			Name:   group_name,
		}
		err = queries.CreateGroupFaculty(ctx, params)
		if err != nil {
			log.Println(err)
			RespondError(w, 500, "Something went wrong while inserting details to our database")
			return
		}

		RespondJSON(w, 200, "Added "+name+" as a faculty incharge of "+group_name+" successfully!")
	}
}

// CreateGroupMember adds a new member to a group.
func CreateGroupMember(db *sql.DB) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(db)
	return func(w http.ResponseWriter, r *http.Request) {
		group_name := mux.Vars(r)["name"]

		rollStr := r.URL.Query().Get("roll")
		if rollStr == "" {
			RespondError(w, 400, "Required query param, roll, missing")
			return
		}
		roll, err := strconv.Atoi(rollStr)
		if err != nil {
			RespondError(w, 400, "Roll paramter must only contain digits")
			return
		}

		params := query.CreateGroupMemberParams{
			Name:       group_name,
			RollNumber: int32(roll),
		}
		err = queries.CreateGroupMember(ctx, params)
		if err != nil {
			log.Println(err)
			RespondError(w, 500, "Something went wrong while inserting details to our database")
			return
		}

		RespondJSON(w, 200, "Added '"+fmt.Sprint(roll)+"' as a member of "+group_name+" successfully!")
	}
}

// CreateGroupSocial adds a new social media handle of a group.
func CreateGroupSocial(db *sql.DB) http.HandlerFunc {
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

		params := query.CreateGroupSocialParams{
			Name:         group_name,
			PlatformType: platform_type,
			Link:         link,
		}
		err := queries.CreateGroupSocial(ctx, params)
		if err != nil {
			log.Println(err)
			RespondError(w, 500, "Something went wrong while inserting details to our database")
			return
		}

		RespondJSON(w, 200, "Added link of the "+platform_type+" handle for "+group_name+" successfully!")
	}
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

// UpdateGroupFaculty updates the mobile number of a group incharge.
func UpdateGroupFaculty(db *sql.DB) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(db)
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		mobileStr := r.URL.Query().Get("mobile")
		if mobileStr == "" {
			RespondError(w, 400, "Required query param, mobile, missing")
			return
		}
		mobile, err := strconv.Atoi(mobileStr)
		if err != nil {
			RespondError(w, 400, "Mobile paramter must only contain digits")
			return
		}

		params := query.UpdateGroupFacultyParams{
			Name:      vars["fname"],
			Mobile:    int64(mobile),
			GroupName: vars["name"],
		}
		err = queries.UpdateGroupFaculty(ctx, params)
		if err != nil {
			RespondError(w, 500, "Something went wrong while updating details to our database")
			return
		}

		RespondJSON(w, 200, "Faculty mobile number updated successfully")
	}
}

// UpdateGroupSocials updates the link of a social media handle for a group.
func UpdateGroupSocials(db *sql.DB) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(db)
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		link := r.URL.Query().Get("link")
		if link == "" {
			RespondError(w, 400, "Required query param, link, missing")
			return
		}

		params := query.UpdateGroupSocialsParams{
			PlatformType: vars["type"],
			Link:         link,
			GroupName:    vars["name"],
		}
		err := queries.UpdateGroupSocials(ctx, params)
		if err != nil {
			RespondError(w, 500, "Something went wrong while updating details to our database")
			return
		}

		RespondJSON(w, 200, "Social media link updated successfully")
	}
}
