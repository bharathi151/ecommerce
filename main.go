package main

import (
	"log"
	"net/http"
	"quickcommerce/controller"
	"quickcommerce/core"
	"quickcommerce/repo"
	"quickcommerce/server"
	"quickcommerce/service"
)

func main() {
	r := repo.NewRepo()
	c := core.NewCore(r)
	s := service.NewService(r, c)
	ctrl := controller.NewController(s)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", server.New(ctrl)))
}
