package breadboard

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	database "breadboard/database"
	h "breadboard/handlers"
	m "breadboard/middleware"

	"github.com/go-co-op/gocron"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/cors"
)

type server struct {
	conn   *pgxpool.Pool
	router *mux.Router
}

// NewServer returns a new app instance.
func NewServer() *server {
	conn, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln("Unable to connect to database:\n", err)
	}

	// !TODO: Make separate interface to initialise database
	database.Init(conn)

	// Initialize cronjob for fetching announcements
	cron := gocron.NewScheduler(time.UTC)
	cron.Every(1).Day().At("00:00").Do(func() {
		h.FetchAnnouncements(conn)
	})
	cron.StartAsync()

	s := server{conn: conn, router: mux.NewRouter().StrictSlash(true)}
	s.setRouters()
	return &s
}

// Run executes the app to listen on a given port and serve the routers.
func (s *server) Run() {
	c := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedOrigins:   []string{"http://localhost:3000", "https://nkss-website.up.railway.app", "https://nkss.getpsyched.dev"},
	})
	log.Fatalln(http.ListenAndServe(":"+os.Getenv("PORT"), c.Handler(s.router)))
}

// setRouters maps endpoints to the functions they must route to.
func (s *server) setRouters() {
	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to NKSS' API!"))
	})

	// Status
	s.router.Handle("/status/student/discord", h.GetDiscordLinkStatus(s.conn)).Methods("GET")

	// Announcements
	s.router.HandleFunc("/announcements", h.GetAnnouncements(s.conn)).Methods("GET")

	// Courses
	s.router.HandleFunc("/courses", h.GetCourses(s.conn)).Methods("GET")
	s.router.HandleFunc("/courses", m.Authenticator(h.CreateCourse(s.conn))).Methods("POST")
	s.router.HandleFunc("/courses/{code}", h.GetCourse(s.conn)).Methods("GET")

	// Clubs
	s.router.Handle("/clubs", h.GetClubs(s.conn)).Methods("GET")
	s.router.Handle("/clubs/{name}", h.GetClub(s.conn)).Methods("GET")

	s.router.Handle("/clubs/{name}/faculty", h.GetClubFaculty(s.conn)).Methods("GET")
	s.router.Handle("/clubs/{name}/faculty", m.Authenticator(h.CreateClubFaculty(s.conn))).Methods("POST")
	s.router.Handle("/clubs/{name}/faculty/{fname}", m.Authenticator(h.DeleteClubFaculty(s.conn))).Methods("DELETE")

	s.router.Handle("/clubs/{name}/members", m.Authenticator(h.ReadClubMembers(s.conn))).Methods("GET")
	s.router.Handle("/clubs/{name}/members", m.Authenticator(h.CreateClubMemberBulk(s.conn))).Methods("POST")
	s.router.Handle("/clubs/{name}/members/{roll}", m.Authenticator(h.CreateClubMember(s.conn))).Methods("POST")
	s.router.Handle("/clubs/{name}/members/{roll}", m.Authenticator(h.UpdateClubMember(s.conn))).Methods("PUT")
	s.router.Handle("/clubs/{name}/members", m.Authenticator(h.DeleteClubMemberBulk(s.conn))).Methods("DELETE")
	s.router.Handle("/clubs/{name}/members/{roll}", m.Authenticator(h.DeleteClubMember(s.conn))).Methods("DELETE")

	s.router.Handle("/clubs/{name}/socials", h.GetClubSocials(s.conn)).Methods("GET")
	s.router.Handle("/clubs/{name}/socials", m.Authenticator(h.CreateClubSocial(s.conn))).Methods("POST")
	s.router.Handle("/clubs/{name}/socials/{type}", m.Authenticator(h.UpdateClubSocials(s.conn))).Methods("PUT")
	s.router.Handle("/clubs/{name}/socials/{type}", m.Authenticator(h.DeleteClubSocial(s.conn))).Methods("DELETE")

	// Students
	s.router.Handle("/hostels", h.GetHostels(s.conn)).Methods("GET")
	s.router.Handle("/students/{id}", m.Authenticator(h.GetStudent(s.conn))).Methods("GET")
}
