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
	store := repo.NewRepo()
	assigner := core.NewCore(store)
	orderService := service.NewService(store, assigner)
	orderController := controller.NewController(orderService)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", server.New(orderController)))
}
