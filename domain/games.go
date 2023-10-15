package domain

type Games struct {
	Response struct {
		GameCount int `json:"game_count"`
		Games     []struct {
			Appid                    int    `json:"appid"`
			Name                     string `json:"name"`
			PlaytimeForever          int    `json:"playtime_forever"`
			ImgIconUrl               string `json:"img_icon_url"`
			PlaytimeWindowsForever   int    `json:"playtime_windows_forever"`
			PlaytimeMacForever       int    `json:"playtime_mac_forever"`
			PlaytimeLinuxForever     int    `json:"playtime_linux_forever"`
			RtimeLastPlayed          int    `json:"rtime_last_played"`
			ContentDescriptorids     []int  `json:"content_descriptorids,omitempty"`
			PlaytimeDisconnected     int    `json:"playtime_disconnected"`
			HasCommunityVisibleStats bool   `json:"has_community_visible_stats,omitempty"`
			HasLeaderboards          bool   `json:"has_leaderboards,omitempty"`
			Playtime2Weeks           int    `json:"playtime_2weeks,omitempty"`
		} `json:"games"`
	} `json:"response"`
}
