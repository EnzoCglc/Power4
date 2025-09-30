package controllers

import (
	"net/http"
	"power4/utils"
)

func LoginPage(w http.ResponseWriter, r *http.Request) {
	// action := r.FormValue("play")
	utils.Render(w, "loginPage.html", nil)
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	// action := r.FormValue("play")
	utils.Render(w, "registerPage.html", nil)
}
