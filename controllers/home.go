package controllers

import (
	"net/http"
	"power4/utils"
)

func Home(w http.ResponseWriter, r *http.Request) {
	utils.Render(w, "index.html", nil)
}
