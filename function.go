package steampie

import (
	"SteamPie/domain"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

var ApiKey = os.Getenv("STEAM_API_KEY")
var SecretKey = os.Getenv("SECRET_KEY")

func init() {
	functions.HTTP("SteamPieHTTP", steamPieHTTP)
}

// steamPieHTTP is an HTTP Cloud Function.
func steamPieHTTP(w http.ResponseWriter, r *http.Request) {
	var Id struct {
		SteamId   string `json:"steamId"`
		SecretKey string `json:"secretKey"`
	}
	if err := json.NewDecoder(r.Body).Decode(&Id); err != nil {
		fmt.Fprint(w, "Error: ", err)
		return
	}

	if Id.SteamId == "" {
		fmt.Fprint(w, "Error: ", "SteamId is empty")
		return
	}

	if Id.SecretKey != SecretKey {
		fmt.Fprintf(w, "Error: ", "SecretKey is wrong")
		return
	}

	var games *domain.Games
	var gameInfos []*domain.GameInfo
	var gameInfoResp *json.RawMessage
	var genreStats map[string]int
	var categoryStats map[string]int

	genreStats = make(map[string]int)
	categoryStats = make(map[string]int)

	games, _ = GetGamesRequest(ApiKey, Id.SteamId)

	games = FilterGames(games)

	var stats *domain.Stats
	stats = &domain.Stats{
		SteamId: Id.SteamId,
	}

	for i, game := range games.Response.Games {
		fmt.Println(i, " Loading game info for", game.Name)
		gameInfoResp, _ = GetGameInfoRequest(fmt.Sprintf("%d", game.Appid))

		var gameInfo *domain.GameInfo
		if err := json.Unmarshal(*gameInfoResp, &gameInfo); err != nil {
			panic(err)
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

	fmt.Fprint(w, stats)
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

	fmt.Println(resp.StatusCode, resp.Status)

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
