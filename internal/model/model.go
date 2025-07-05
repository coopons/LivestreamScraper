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