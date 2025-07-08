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
	Platform 		string `json:"platform"` // e.g., "twitch", "youtube", "kick"
}

// Data from yt-dlp --dump-json of live channels
type YtStream struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	UserName     string   `json:"uploader"`
	UserID     	 string	  `json:"channel_id"`
	ViewCount    int      `json:"concurrent_view_count"`
	LiveStatus   bool     `json:"is_live"`
	Language     string   `json:"language"`
	Thumbnail    string   `json:"thumbnail"`
	StartTime    int64    `json:"start_time"`
	Categories	 []string `json:"categories"`
}

// Data from the Kick API
type KickStream struct {
	BroadcasterUserID int `json:"broadcaster_user_id"`
	ChannelID         int `json:"channel_id"`
	Slug              string `json:"slug"`
	StreamTitle       string `json:"stream_title"`
	Language          string `json:"language"`
	HasMatureContent  bool   `json:"has_mature_content"`
	ViewerCount       int    `json:"viewer_count"`
	Thumbnail         string `json:"thumbnail"`
	StartedAt         string `json:"started_at"` // ISO8601 string; you can also use time.Time with parsing

	Category struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Thumbnail string `json:"thumbnail"`
	} `json:"category"`
}