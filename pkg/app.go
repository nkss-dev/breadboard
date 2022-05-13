package pkg

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"NKSS-backend/config"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

func (a *App) Initialise(config *config.Config) {
	dbURI := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		config.DB.User,
		config.DB.Password,
		config.DB.Host,
		config.DB.Port,
		config.DB.Database,
	)

	db, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{})
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
	panic("unimplemented!")
}
