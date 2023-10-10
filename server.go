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
	pool   *pgxpool.Pool
	router *mux.Router
}

// NewServer returns a new app instance.
func NewServer() *server {
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln("Unable to connect to database:\n", err)
	}

	// !TODO: Make separate interface to initialise database
	database.Init(pool)

	// Initialize cronjob for fetching announcements
	cron := gocron.NewScheduler(time.UTC)
	cron.Every(1).Day().At("00:00").Do(func() {
		h.FetchAnnouncements(pool)
	})
	cron.StartAsync()

	s := server{pool: pool, router: mux.NewRouter().StrictSlash(true)}
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
	s.router.Handle("/status/student/discord", h.GetDiscordLinkStatus(s.pool)).Methods("GET")

	// Announcements
	s.router.HandleFunc("/announcements", h.GetAnnouncements(s.pool)).Methods("GET")

	// Courses
	s.router.HandleFunc("/courses", h.GetCourses(s.pool)).Methods("GET")
	s.router.HandleFunc("/courses", m.Authenticator(h.CreateCourse(s.pool))).Methods("POST")
	s.router.HandleFunc("/courses/{code}", h.GetCourse(s.pool)).Methods("GET")

	// Clubs
	s.router.Handle("/clubs", h.GetClubs(s.pool)).Methods("GET")
	s.router.Handle("/clubs/{name}", h.GetClub(s.pool)).Methods("GET")

	s.router.Handle("/clubs/{name}/faculty", h.GetClubFaculty(s.pool)).Methods("GET")
	s.router.Handle("/clubs/{name}/faculty", m.Authenticator(h.CreateClubFaculty(s.pool))).Methods("POST")
	s.router.Handle("/clubs/{name}/faculty/{fname}", m.Authenticator(h.DeleteClubFaculty(s.pool))).Methods("DELETE")

	s.router.Handle("/clubs/{name}/members", m.Authenticator(h.ReadClubMembers(s.pool))).Methods("GET")
	s.router.Handle("/clubs/{name}/members", m.Authenticator(h.CreateClubMemberBulk(s.pool))).Methods("POST")
	s.router.Handle("/clubs/{name}/members/{roll}", m.Authenticator(h.CreateClubMember(s.pool))).Methods("POST")
	s.router.Handle("/clubs/{name}/members/{roll}", m.Authenticator(h.UpdateClubMember(s.pool))).Methods("PUT")
	s.router.Handle("/clubs/{name}/members", m.Authenticator(h.DeleteClubMemberBulk(s.pool))).Methods("DELETE")
	s.router.Handle("/clubs/{name}/members/{roll}", m.Authenticator(h.DeleteClubMember(s.pool))).Methods("DELETE")

	s.router.Handle("/clubs/{name}/socials", h.GetClubSocials(s.pool)).Methods("GET")
	s.router.Handle("/clubs/{name}/socials", m.Authenticator(h.CreateClubSocial(s.pool))).Methods("POST")
	s.router.Handle("/clubs/{name}/socials/{type}", m.Authenticator(h.UpdateClubSocials(s.pool))).Methods("PUT")
	s.router.Handle("/clubs/{name}/socials/{type}", m.Authenticator(h.DeleteClubSocial(s.pool))).Methods("DELETE")

	// Students
	s.router.Handle("/hostels", h.GetHostels(s.pool)).Methods("GET")
	s.router.Handle("/students/{id}", m.Authenticator(h.GetStudent(s.pool))).Methods("GET")
}
