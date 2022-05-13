package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Faculty struct {
	Group_Name string
	Name       string
	Mobile     string
}

type Discord struct {
	Name           string
	Id             string
	Invite         string
	Fresher_Role   string
	Sophomore_Role string
	Junior_Role    string
	Senior_Role    string
	Guest_Role     string
}

type Group struct {
	Name           string
	Id             string
	Invite         string
	Fresher_Role   string
	Sophomore_Role string
	Junior_Role    string
	Senior_Role    string
	Guest_Role     string
	Alias	string
	Branch	string
	Kind 	string
	Description	string
}

type Social struct {
	Name string
	Type string
	link string
}

type Admin struct {
	Group_Name  string
	Position    string
	Roll_Number string
}

type Member struct {
	Group_Name  string
	Roll_Number string
}

func execQuery(query string, object stringerface{}, w http.ResponseWriter, db *gorm.DB) {

	err = db.Raw(query).Scan(&object).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		resp := make(map[string]string)
		resp["message"] = "Error! The query went wrong somewhere!"
		jsonResp, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResp)
		return
	}

}

func GetAllMemberInfo(w http.ResponseWriter, r *http.Request) {

	dsn := "host=localhost user=gorm password=gorm dbname=group port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")

	var faculty []Faculty
	var discord []Discord
	var social []Social
	var admin []Admin
	var member []Member

	execQuery("SELECT * FROM groups NATURAL JOIN group_discord", &discord, w, db)

	execQuery("SELECT * FROM group_faculty", &faculty, w, db)
	execQuery("SELECT * FROM group_social", &social, w, db)
	execQuery("SELECT * FROM group_admin", &admin, w, db)
	execQuery("SELECT * FROM group_member", &member, w, db)


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
}
