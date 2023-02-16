package nkssbackend

import (
	"database/sql"
	"log"
	"net/http"
	"os"
    "time"

	h "nkssbackend/handlers"
	"nkssbackend/internal/database"
	m "nkssbackend/middleware"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
    "github.com/go-co-op/gocron"
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

    // Initialize cronjob for fetching announcements
    cron := gocron.NewScheduler(time.UTC)
    cron.Every(1).Day().At("00:00").Do(func() {
       h.FetchAnnouncements(db)
    })
    cron.StartAsync()

	s := server{db: db, router: mux.NewRouter().StrictSlash(true)}
	s.setRouters()
	return &s
}

// Run executes the app to listen on a given port and serve the routers.
func (s *server) Run() {
	c := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedOrigins:   []string{"http://localhost:3000", "http://nksss.live", "https://nksss.live"},
	})
	log.Fatalln(http.ListenAndServe(":"+os.Getenv("PORT"), c.Handler(s.router)))
}

// setRouters maps endpoints to the functions they must route to.
func (s *server) setRouters() {
	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to NKSSS' API!"))
	})

	// Status
	s.router.Handle("/status/student/discord", h.GetDiscordLinkStatus(s.db)).Methods("GET")

	// Announcements
	s.router.HandleFunc("/announcements", h.GetAnnouncements(s.db)).Methods("GET")

	// Courses
	s.router.HandleFunc("/courses", h.GetCourses(s.db)).Methods("GET")
	s.router.HandleFunc("/courses", m.Authenticator(h.CreateCourse(s.db))).Methods("POST")
	s.router.HandleFunc("/courses/{code}", h.GetCourse(s.db)).Methods("GET")

	// Clubs
	s.router.Handle("/clubs", h.GetClubs(s.db)).Methods("GET")
	s.router.Handle("/clubs/{name}", h.GetClub(s.db)).Methods("GET")

	s.router.Handle("/clubs/{name}/admins", h.GetClubAdmins(s.db)).Methods("GET")
	s.router.Handle("/clubs/{name}/admins", m.Authenticator(h.CreateClubAdmin(s.db))).Methods("POST")
	s.router.Handle("/clubs/{name}/admins/{roll}", m.Authenticator(h.DeleteClubAdmin(s.db))).Methods("DELETE")

	s.router.Handle("/clubs/{name}/faculty", h.GetClubFaculty(s.db)).Methods("GET")
	s.router.Handle("/clubs/{name}/faculty", m.Authenticator(h.CreateClubFaculty(s.db))).Methods("POST")
	s.router.Handle("/clubs/{name}/faculty/{fname}", m.Authenticator(h.DeleteClubFaculty(s.db))).Methods("DELETE")

	s.router.Handle("/clubs/{name}/members", m.Authenticator(h.GetClubMembers(s.db))).Methods("GET")
	s.router.Handle("/clubs/{name}/members", m.Authenticator(h.CreateClubMember(s.db))).Methods("POST")
	s.router.Handle("/clubs/{name}/members/{roll}", m.Authenticator(h.DeleteClubMember(s.db))).Methods("DELETE")

	s.router.Handle("/clubs/{name}/socials", h.GetClubSocials(s.db)).Methods("GET")
	s.router.Handle("/clubs/{name}/socials", m.Authenticator(h.CreateClubSocial(s.db))).Methods("POST")
	s.router.Handle("/clubs/{name}/socials/{type}", m.Authenticator(h.UpdateClubSocials(s.db))).Methods("PUT")
	s.router.Handle("/clubs/{name}/socials/{type}", m.Authenticator(h.DeleteClubSocial(s.db))).Methods("DELETE")

	// Students
	s.router.Handle("/hostels", h.GetHostels(s.db)).Methods("GET")
	s.router.Handle("/students/{roll}", m.Authenticator(h.GetStudent(s.db))).Methods("GET")

	// Discord
	s.router.Handle("/discord/users/{discord_id}", m.Authenticator(h.GetStudentByDiscordID(s.db))).Methods("GET")
}
