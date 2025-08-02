package domain

type Color string // "w" or "b"

type Players struct {
	White 		PublicUser 		`json:"white"`
	Black 		PublicUser 		`json:"black"`
}

type Move struct {
	X uint8 		`json:"x"`
	Y uint8 		`json:"y"`
	Color Color		`json:"color"`
}

type BattleSettings struct {
	Size		[2]int8		`json:"size"`
	Handicap	int			`json:"handicap"`
	Komi		float32		`json:"komi"`
}

type BattleRequest struct {
	// request
	Status 		   	string			`json:"status"`
	Id             	string			`json:"id"`
	Sender 			PublicUser		`json:"sender"`
	// player id
	Target 			string			`json:"target"`
	Game 		   	BattleSettings	`json:"game"`
}

type Board []string
type MovesList []Move

type Game struct {
	Board     Board       `json:"board"`
	Moves     MovesList	  `json:"moves"`
	BattleSettings
}

type BattleObj struct {
	// active
	Status 		   string		  `json:"status"`
	Id             string         `json:"id"`
	Players 	   Players        `json:"players"`
	Game 		   Game			  `json:"game"`
}