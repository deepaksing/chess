package types

type MoveResp struct {
	Move_from       string `json:"move_from"`
	Move_to         string `json:"move_to"`
	Move_type       int    `json:"move_type"`
	Player_username string `json:"player_username"`
	Match_id        int    `json:"match_id"`
}
