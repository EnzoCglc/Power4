package controllers

import (
	"log"
	"net/http"
	"power4/models"
	"power4/utils"
)

type ProfilData struct {
	User       *models.User
	WinRate    float64
	TotalGames int
	History    []models.History
}

func calculateWinRate(wins, losses int) float64 {
	totalGames := wins + losses
	if totalGames == 0 {
		return 0.0
	}
	return (float64(wins) / float64(totalGames)) * 100
}

func Profil(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("username")
	if err != nil {
		log.Println("User is not login")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	username := cookie.Value
	user, err := models.GetUserByUsername(username)
	if err != nil {
		log.Println("Error fetching user from database:", err)
		http.Error(w, "Error loading data", http.StatusInternalServerError)
		return
	}

	if user == nil {
		log.Println("User not exits in database:", username)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	winRate := calculateWinRate(user.Win, user.Losses)
	totalGames := user.Win + user.Losses

	history, err := models.GetHistoryByPlayer(username)
	if err != nil {
		log.Printf("Error to get history for %s : %s", username, err)
		return
	}

	data := ProfilData{
		User:       user,
		WinRate:    winRate,
		TotalGames: totalGames,
		History:    history,
	}

	utils.Render(w, "profil.html", data)

	log.Println("info de l'user ", user, "- WinRate:", winRate, "% - Total games:", totalGames)
}
