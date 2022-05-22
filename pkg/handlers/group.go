package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"NKSS-backend/pkg/query"

	"github.com/gorilla/mux"
	"golang.org/x/exp/slices"
)

func GetAllMemberInfo(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	w.Header().set("Content-Type", "application/json")

	ctx := context.Background()
	queries := query.New(db)
	discord, err := queries.getAllGroups(ctx)
	if err == sql.ErrNoRows {
		respondError(w, 404, "No groups found!")
		return
	}
	faculty , err := queries.getAllFaculty(ctx)
	if err == sql.ErrNoRows {
		respondError(w, 404, "No faculty found!")
		return
	}
	social, err := queries.getAllGroupSocials(ctx)
	if err == sql.ErrNoRows {
		respondError(w, 404, "No socials found!")
		return
	}
	admin, err := queries.getAllGroupAdmins(ctx)
	if err == sql.ErrNoRows {
		respondError(w, 404, "No admins found!")
		return
	}
	member, err := queries.getAllGroupMembers(ctx)
	if err == sql.ErrNoRows {
		respondError(w, 404, "No members found!")
		return
	}
	
	var groups [](map[string]interface{})
	
	for g:=0; g < len(discord); g++ {
		
		var group map[string]interface{}

		group["name"] := discord[g]["name"]
		group["alias"] := discord[g]["alias"]
		group["branch"] := discord[g]["branch"]
		group["kind"] := discord[g]["kind"]
		group["description"] := discord[g]["description"]


		var group["faculty"] []map[string]string
		x := 
		for f:=0; f < len(faculty); f++ {
			thisfaculty = faculty[f]
			if thisfaculty["group_name"] == group["name"] {
				group["faculty"][x]["name"] = thisfaculty["name"]
				group["faculty"][x]["mobile"] = thisfaculty["mobile"]
				x += 1
			}
		}
		var group["social"] []map[string]string
		x := 0
		for s:=0; s < len(social); s++ {
			thissocial = social[s]	
			if thissocial["name"] == group["name"] {
				group["social"][x]["type"] = thissocial["type"]
				group["social"][x]["link"] = thissocial["link"]
				x += 1
			}
		}
		var group["admin"] []map[string]string
		x := 0
		for a:=0; a < len(admin); a++ {
			thisadmin = admin[a]
			if thisadmin["group_name"] == group["name"] {
				group["admin"][x]["roll_number"] = thisadmin["roll_number"]
				group["admin"][x]["position"] = thisadmin["position"]
				x += 1
			}
		}
		var group["member"] []map[string]string
		x := 0
		for m:=0; m < len(member); m++ {
			thismember = member[m]
			if thismember["group_name"] == group["name"] {
				group["member"][x]["roll_number"] = thisadmin["roll_number"]
				x += 1
			}
		}

		groups[g] = group

	}

	jsonResp, _ := json.Marshal(groups)
	w.Write(jsonResp)
	respondJSON(w, 200, groups)

}
