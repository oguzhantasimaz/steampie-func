package domain

type Stats struct {
	SteamId   string `json:"steam_id"`
	GameCount int    `json:"games_count"`
	//Categories []Category `json:"categories"`
	Genres []Genre `json:"genres"`
	Games  []Game  `json:"games"`
}

type Game struct {
	Name     string `json:"name"`
	PlayTime int    `json:"play_time"`
}

type Genre struct {
	Name     string   `json:"name"`
	PlayTime int      `json:"play_time"`
	Games    []string `json:"games"`
}

//type Category struct {
//	Name     string `json:"name"`
//	PlayTime int    `json:"play_time"`
//	Games    []string `json:"games"`
//}
