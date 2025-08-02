package domain

// обработчик "error"
type SocketErr struct {
	Code uint8		`json:"code"`	// уникальный номер ошибки
	Type string		`json:"type"`	// тип ошибки, например: "battle-invite-error", "serv-pomer"
	Data any		`json:"data"`	// либо сообщение либо json обьект, зависит от ошибки
}

type SocketErrors struct {}

func (s *SocketErrors) Unexpected(msg string) SocketErr {
	return SocketErr{Code: 1, Type: "unexpected", Data: msg}
}

func (s *SocketErrors) Battle_Rejected(id string) SocketErr {
	type Data struct {
		ID string	`json:"id"`
	}
	return SocketErr{Code: 1, Type: "battle-rejected", Data: Data{ID: id}}
}

func (s *SocketErrors) Opponent_Offline(id string) SocketErr {
	type Data struct {
		ID string	`json:"id"`
	}
	return SocketErr{Code: 2, Type: "opponent-offline", Data: Data{ID: id}}
}

func (s *SocketErrors) Battle_Already_Requested(id string) SocketErr {
	type Data struct {
		ID string	`json:"id"`	// id юзера которому отправлен запрос
	}
	return SocketErr{Code: 3, Type: "battle-already-requested", Data: Data{ID: id}}
}

func (s *SocketErrors) Connection_Error(errText string) SocketErr {
	type Data struct {
		ErrorMsg string	`json:"msg"`
	}
	return SocketErr{Code: 4, Type: "connection-error", Data: Data{ErrorMsg: errText}}
}

func (s *SocketErrors) JWT_Invalid() SocketErr {
	return SocketErr{Code: 5, Type: "invalid-jwt-token"}
}