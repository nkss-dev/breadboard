package nkssbackend

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	h "nkssbackend/handlers"
	"nkssbackend/internal/database"
	m "nkssbackend/middleware"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type server struct {
	db     *sql.DB
	router *mux.Router
}

// NewServer returns a new app instance.
func NewServer() *server {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln(err)
	}
	// !TODO: Make separate interface to initialise database
	database.Init(db)

	s := server{db: db, router: mux.NewRouter().StrictSlash(true)}
	s.setRouters()
	return &s
}

// Run executes the app to listen on a given port and serve the routers.
func (s *server) Run() {
	log.Fatalln(http.ListenAndServe(":"+os.Getenv("PORT"), s.router))
}

// setRouters maps endpoints to the functions they must route to.
func (s *server) setRouters() {
	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to NKSSS' API!"))
	})

	s.router.HandleFunc("/announcements", h.GetAnnouncements()).Methods("GET")

	s.router.HandleFunc("/courses", h.GetCourses(s.db)).Methods("GET")
	s.router.HandleFunc("/courses/{code}", h.GetCourse(s.db)).Methods("GET")

	// Groups
	s.router.Handle("/groups", m.Authenticator(h.GetGroups(s.db))).Methods("GET")
	s.router.Handle("/groups/{name}", m.Authenticator(h.GetGroup(s.db))).Methods("GET")

	s.router.Handle("/groups/{name}/admins", h.GetGroupAdmins(s.db)).Methods("GET")
	s.router.Handle("/groups/{name}/admins", m.Authenticator(h.CreateGroupAdmin(s.db))).Methods("POST")
	s.router.Handle("/groups/{name}/admins/{roll}", m.Authenticator(h.DeleteGroupAdmin(s.db))).Methods("DELETE")

	s.router.Handle("/groups/{name}/faculty", h.GetGroupFaculty(s.db)).Methods("GET")
	s.router.Handle("/groups/{name}/faculty", m.Authenticator(h.CreateGroupFaculty(s.db))).Methods("POST")
	s.router.Handle("/groups/{name}/faculty/{fname}", m.Authenticator(h.DeleteGroupFaculty(s.db))).Methods("DELETE")

	s.router.Handle("/groups/{name}/members", m.Authenticator(h.GetGroupMembers(s.db))).Methods("GET")
	s.router.Handle("/groups/{name}/members", m.Authenticator(h.CreateGroupMember(s.db))).Methods("POST")
	s.router.Handle("/groups/{name}/members/{roll}", m.Authenticator(h.DeleteGroupMember(s.db))).Methods("DELETE")

	s.router.Handle("/groups/{name}/socials", h.GetGroupSocials(s.db)).Methods("GET")
	s.router.Handle("/groups/{name}/socials", m.Authenticator(h.CreateGroupSocial(s.db))).Methods("POST")
	s.router.Handle("/groups/{name}/socials/{type}", m.Authenticator(h.UpdateGroupSocials(s.db))).Methods("PUT")
	s.router.Handle("/groups/{name}/socials/{type}", m.Authenticator(h.DeleteGroupSocial(s.db))).Methods("DELETE")

	s.router.Handle("/students/{roll}", m.Authenticator(h.GetStudentByRoll(s.db))).Methods("GET")
	s.router.Handle("/students/{roll}/member", m.Authenticator(h.GetStudentClubMemberships(s.db))).Methods("GET")
	s.router.Handle("/students/{roll}/admin", m.Authenticator(h.IsStudentAdmin(s.db))).Methods("GET")
}
