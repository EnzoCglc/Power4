package controllers

import (
	"net/http"
	"power4/utils"
)

func Profil(w http.ResponseWriter, r *http.Request) {
	utils.Render(w, "profil.html", nil)
}
