package domain

type GameInfo struct {
	Data struct {
		Type       string `json:"type"`
		Name       string `json:"name"`
		SteamAppid int    `json:"steam_appid"`
		Categories []struct {
			Id          int    `json:"id"`
			Description string `json:"description"`
		} `json:"categories"`
		Genres []struct {
			Id          string `json:"id"`
			Description string `json:"description"`
		} `json:"genres"`
	} `json:"data"`
}
