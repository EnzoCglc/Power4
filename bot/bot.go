package bot

import (
	"power4/models"
)

// BotMove determines the AI's next move based on the difficulty level.
func BotMove(game *models.GridPage, level int, player int) int {
	switch level {
	case 1:
		return Lvl1Bot(game, player)
	case 2:
		return Lvl2Bot(game, player)
	case 3:
		return Lvl3Bot(game, player)
	case 4:
		return Lvl4Bot(game, player)
	case 5:
		return Lvl5Bot(game, player)
	default:
		return Lvl1Bot(game, player) // Default to easiest level if invalid
	}
}