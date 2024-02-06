package main

import (
	"mymodul/service"
	"net/http"
)

func main() {
	mux := http.NewServeMux() //новый мультиплексор, что бы слушать больше чем один ендпоинт
	srv := service.New()      //новый сервис

	mux.HandleFunc("/vote", srv.Vote)   //все, что приходит на этот ендпоинт, обрабатывается в методе Vote
	mux.HandleFunc("/stats", srv.Stats) //все, что приходит на этот ендпоинт, обрабатывается в методе Stats

	http.ListenAndServe(":8000", mux)

}
