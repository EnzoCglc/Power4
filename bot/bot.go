package bot

import (
	"power4/models"
)

func BotMove(game *models.GridPage, level int , player int) int {
	switch level {
	case 1:
		return Lvl1Bot(game, player)
	case 2:
		return Lvl2Bot(game, player)
	case 3:
		return Lvl3Bot(game, player) 
	case 4:
		return Lvl1Bot(game, player) // TODO: Mettre le niv 4
	case 5:
		return Lvl1Bot(game, player) // TODO: Mettre le niv 5
	default:
		return Lvl1Bot(game, player)
	}
}