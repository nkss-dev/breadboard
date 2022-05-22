package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"NKSS-backend/pkg/query"
)

type Group struct {
	Name        string
	Alias       string
	Faculty     []query.GroupFaculty
	Branch      string
	Kind        string
	Description string
	Socials     interface{}
	Admins      []query.GroupAdmin
	Members     []int32
}

func GetGroups(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	queries := query.New(db)
	allGroups, err := queries.GetAllGroups(ctx)
	if err == sql.ErrNoRows {
		respondError(w, 404, "No groups found!")
		return
	}
	allFaculties, err := queries.GetAllFaculty(ctx)
	if err == sql.ErrNoRows {
		respondError(w, 404, "No faculty found!")
		return
	}
	allSocials, err := queries.GetAllGroupSocials(ctx)
	if err == sql.ErrNoRows {
		respondError(w, 404, "No socials found!")
		return
	}
	allAdmins, err := queries.GetAllGroupAdmins(ctx)
	if err == sql.ErrNoRows {
		respondError(w, 404, "No admins found!")
		return
	}
	allMembers, err := queries.GetAllGroupMembers(ctx)
	if err == sql.ErrNoRows {
		respondError(w, 404, "No members found!")
		return
	}

	var groups []Group

	for _, group := range allGroups {
		name := group.Name
		alias := group.Alias
		branch := group.Branch
		kind := group.Kind
		description := group.Description

		var faculty []query.GroupFaculty
		for _, thisFaculty := range allFaculties {
			if thisFaculty.GroupName == name {
				faculty = append(faculty, thisFaculty)
			}
		}
		var socials []query.GroupSocial
		for _, social := range allSocials {
			if social.Name == name {
				socials = append(socials, social)
			}
		}
		var admins []query.GroupAdmin
		for _, admin := range allAdmins {
			if admin.GroupName == name {
				admins = append(admins, admin)
			}
		}
		var members []int32
		for _, member := range allMembers {
			if member.GroupName == name {
				members = append(members, member.RollNumber)
			}
		}

		fmt.Println(group)
		groups = append(
			groups,
			Group{
				Name:        name,
				Alias:       alias.String,
				Faculty:     faculty,
				Branch:      branch.String,
				Kind:        kind,
				Description: description.String,
				Socials:     socials,
				Admins:      admins,
				Members:     members,
			})
	}

	respondJSON(w, 200, groups)
}
