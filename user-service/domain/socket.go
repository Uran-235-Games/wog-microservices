package domain

// socket.io to server event: battle-request
type ServerBattleRequest struct {
	// player id
	Target string			`json:"target"`
	Game BattleSettings		`json:"game"`
}

// socket.io to client event: battle-request
type ClientBattleRequest struct {
	// player id
	Sender 		PublicUser			`json:"sender"`
	Game 		BattleSettings		`json:"game"`
}

// socket.io to client event: error
// эти ошибки чекать в файле domain/socket_errors.go