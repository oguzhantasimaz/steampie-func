package steampie

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/oguzhantasimaz/steampie-func/domain"
)

var ApiKey = os.Getenv("STEAM_API_KEY")

func init() {
	functions.HTTP("SteamPieHTTP", steamPieHTTP)
}

// steamPieHTTP is an HTTP Cloud Function.
func steamPieHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	var Id struct {
		SteamId string `json:"steamId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&Id); err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	if Id.SteamId == "" {
		json.NewEncoder(w).Encode("SteamId is required")
		return
	}

	var games *domain.Games
	var gameInfos []*domain.GameInfo
	var gameInfoResp *json.RawMessage
	var genreStats map[string]int
	var categoryStats map[string]int

	genreStats = make(map[string]int)
	categoryStats = make(map[string]int)

	var err error
	games, err = GetGamesRequest(ApiKey, Id.SteamId)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	games = FilterGames(games)

	var stats *domain.Stats
	stats = &domain.Stats{
		SteamId:   Id.SteamId,
		GameCount: games.Response.GameCount,
	}

	for _, game := range games.Response.Games {
		gameInfoResp, err = GetGameInfoRequest(fmt.Sprintf("%d", game.Appid))
		if err != nil {
			json.NewEncoder(w).Encode(err)
			return
		}

		var gameInfo *domain.GameInfo
		if err = json.Unmarshal(*gameInfoResp, &gameInfo); err != nil {
			json.NewEncoder(w).Encode(err)
			return
		}

		for _, genre := range gameInfo.Data.Genres {
			genreStats[genre.Description] += game.PlaytimeForever / 60
		}

		for _, category := range gameInfo.Data.Categories {
			categoryStats[category.Description] += game.PlaytimeForever / 60
		}

		stats.Games = append(stats.Games, domain.Game{
			Name:     game.Name,
			PlayTime: game.PlaytimeForever / 60,
		})

		gameInfos = append(gameInfos, gameInfo)
		time.Sleep(101 * time.Millisecond)
	}

	for genre, playTime := range genreStats {
		stats.Genres = append(stats.Genres, domain.Genre{
			Name:     genre,
			PlayTime: playTime,
		})
	}

	for category, playTime := range categoryStats {
		stats.Categories = append(stats.Categories, domain.Category{
			Name:     category,
			PlayTime: playTime,
		})
	}

	//return as json
	json.NewEncoder(w).Encode(stats)
}

func GetGamesRequest(key, steamId string) (*domain.Games, error) {
	client := &http.Client{}

	url := fmt.Sprintf("https://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key=%s&steamid=%s&format=json&include_appinfo=true&include_played_free_games=true", key, steamId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var body []byte
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var games *domain.Games
	if err = json.Unmarshal(body, &games); err != nil {
		return nil, err
	}

	return games, nil
}

func GetGameInfoRequest(appid string) (*json.RawMessage, error) {
	client := &http.Client{}

	url := fmt.Sprintf("https://store.steampowered.com/api/appdetails?appids=%s&l=english", appid)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var body []byte
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	//fmt.Println(string(body))

	var raw map[string]*json.RawMessage
	if err = json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	return raw[appid], nil
}

func FilterGames(games *domain.Games) *domain.Games {
	var filteredGames *domain.Games

	filteredGames = &domain.Games{}

	//sort games by playtime
	for i := 0; i < len(games.Response.Games); i++ {
		for j := 0; j < len(games.Response.Games); j++ {
			if games.Response.Games[i].PlaytimeForever > games.Response.Games[j].PlaytimeForever {
				games.Response.Games[i], games.Response.Games[j] = games.Response.Games[j], games.Response.Games[i]
			}
		}
	}

	//filter most played games limit 50
	for _, game := range games.Response.Games {
		if game.PlaytimeForever > 0 {
			filteredGames.Response.Games = append(filteredGames.Response.Games, game)
		}

		if len(filteredGames.Response.Games) == 25 {
			break
		}
	}

	return filteredGames
}
