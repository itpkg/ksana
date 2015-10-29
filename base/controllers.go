package base

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (p *BaseEngine) Mount(rt *mux.Router) {

	//rt.Static("/assets", "assets")

	//----------------------------------------
	rt.HandleFunc("/users/sign_in", func(http.ResponseWriter, *http.Request) {
	}).Methods("POST")

	rt.HandleFunc("/users/sign_up", func(http.ResponseWriter, *http.Request) {
	}).Methods("POST")
	rt.HandleFunc("/users/sign_out", func(http.ResponseWriter, *http.Request) {
	}).Methods("DELETE")
	rt.HandleFunc("/users/confirm", func(http.ResponseWriter, *http.Request) {
	}).Methods("POST")
	rt.HandleFunc("/users/unlock", func(http.ResponseWriter, *http.Request) {
	}).Methods("POST")
	rt.HandleFunc("/users/forgot_password", func(http.ResponseWriter, *http.Request) {
	}).Methods("POST")
	rt.HandleFunc("/users/reset_password", func(http.ResponseWriter, *http.Request) {
	}).Methods("POST")
	rt.HandleFunc("/users/profile", func(http.ResponseWriter, *http.Request) {
	}).Methods("POST")

}
