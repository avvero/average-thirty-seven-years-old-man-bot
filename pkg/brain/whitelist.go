package brain

import (
	"github.com/avvero/the_gamers_guild_bot/internal/utils"
	"strconv"
)

type Whitelist struct {
}

func (this Whitelist) Check(chatId int64, text string) (forbidden bool, response string) {
	if !utils.Contains([]string{"0", "-1001733786877", "245851441", "-578279468", "-831690505"}, strconv.FormatInt(chatId, 10)) {
		return true, "Mr Moony presents his compliments to Professor Snape, and begs him to keep his " +
			"abnormally large nose out of other people’s business."
	} else {
		return false, ""
	}
}
