package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/viper"
	"github.com/swisscdn/OpenfiatServer/app/config"
	db "github.com/swisscdn/OpenfiatServer/app/database"
	"github.com/swisscdn/OpenfiatServer/app/handler"
	"github.com/swisscdn/OpenfiatServer/app/routes"
	"github.com/swisscdn/OpenfiatServer/app/token"
)

func main() {
	con, err := LoadConfig("./")
	if err != nil {
		log.Fatal("can not load config:", err)
	}
	tokenMaker, err := token.NewPasetoMaker(con.TokenSymmetricKey)
	if err != nil {
		log.Fatalf("can not create toekn maker: %v", err)
	}

	//session, err := util.GetAstraDBSession(config)
	client, err := config.GetAstraDBClient(con)
	if err != nil {
		log.Fatal("can not connect to DB:", err)
	}
	// Load the timezone GMT+5:30
	ser := &handler.Server{
		Config:     con,
		TokenMaker: tokenMaker,
		Store: &db.DB{
			DB: client,
		},
	}
	handler.InitServer(ser)
	if con.DBMigrateUp {
		db.MigrateUpDB(&db.DB{DB: client})
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", con.ListenPort),
		Handler: routes.Routes(tokenMaker),
	}

	//start the server
	fmt.Printf("server started on port %s", con.ListenPort)
	err = srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}

}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config config.Config, err error) {
	//viper.AddConfigPath(path)
	//viper.SetConfigFile("app.env")
	viper.SetConfigFile("./application.env")
	//viper.SetConfigName("app")
	//viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Println(err)
		return
	}

	err = viper.Unmarshal(&config)
	return
}
