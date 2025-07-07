package model

// Data parsed from the twitch api curl
type Stream struct {
	ID 				string `json:"id"`
	UserID 			string `json:"user_id"`
	UserName 		string `json:"user_name"`
	Title 			string `json:"title"`
	GameID 			string `json:"game_id"`
	GameName 		string `json:"game_name"`
	Language 		string `json:"language"`
	ViewerCount 	int    `json:"viewer_count"`
	StartedAt 		string `json:"started_at"`
	ThumbnailURL 	string `json:"thumbnail_url"`
	IsMature 		bool   `json:"is_mature"`
}

// Data from yt-dlp --dump-json of live channels
type YtStream struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	UserName     string   `json:"uploader"`
	UserID     	 string	  `json:"channel_id"`
	ViewCount    int      `json:"view_count"`
	LiveStatus   bool     `json:"is_live"`
	Language     string   `json:"language"`
	Thumbnail    string   `json:"thumbnail"`
	StartTime    int64    `json:"start_time"`
	Categories	 []string `json:"categories"`
}