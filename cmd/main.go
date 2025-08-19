package main

import (
	"log"

	pkg "github.com/4Noyis/my-library/pkg/database"
)

func main() {

	err := pkg.ConnectMongoDB()
	if err != nil {
		log.Fatal("failed to connect to mongodb:", err)
	}
	defer pkg.DisconnectMongoDB()

	pkg.AddItemMongoDB()

}
