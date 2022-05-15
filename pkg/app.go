package pkg

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"NKSS-backend/config"
	"NKSS-backend/pkg/handlers"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialise(config *config.Config) {
	db, err := sql.Open("postgres", config.DatabaseUrl)
	if err != nil {
		log.Fatalln(err)
	}

	a.DB = db
	a.Router = mux.NewRouter().StrictSlash(true)
	a.setRouters()
}

func (a *App) Run() {
	port, port_exists := os.LookupEnv("PORT")
	if !port_exists {
		port = "8081"
	}
	log.Fatal(http.ListenAndServe(":"+port, a.Router))
}

func (a *App) setRouters() {
	a.Router.HandleFunc("/announcements", handlers.GetAnnouncements).Methods("GET")

	a.Router.HandleFunc("/courses", a.passDB(handlers.GetCourses)).Methods("GET")
	a.Router.HandleFunc("/course/{code}", a.passDB(handlers.GetCourse)).Methods("GET")
}

func (a *App) passDB(handler func(db *sql.DB, w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(a.DB, w, r)
	}
}
