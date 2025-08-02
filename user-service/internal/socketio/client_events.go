package socketio

import (
	"log"
	"wog-server/domain"
	"wog-server/internal/logger"

	socketio "github.com/googollee/go-socket.io"
)

type ClientEvents struct {
	server *socketio.Server
}

func (s *ClientEvents) sendAll(event string, data any, sockets []any) {
	for _, socket := range sockets {
		switch v := socket.(type) {
		case string:
			log.Println("отправка ошибки в комнату")
			s.server.BroadcastToRoom("/", v, event, data)
		default:
			if conn, ok := v.(socketio.Conn); ok {
				log.Println("отправка ошибки сокету")
				conn.Emit(event, data)
			} else {
				logger.Log.Error("в sendAll передан инвалидный тип сокета")
			}
		}
	}
}

// event: "error"
//
// принимает uid либо socketio.Conn
func (s *ClientEvents) Error(err domain.SocketErr, sockets ...any) {
	s.sendAll("error", err, sockets)
}

// event: "battle-request"
//
// Если уже есть активный запрос:
// на "error" прийдет SocketErrors.Battle_Already_Requested
func (s *ClientEvents) BattleRequest(r *domain.ClientBattleRequest, sockets ...any) {
	s.sendAll("battle-request", r, sockets)
}

// event: "active-battle"
func (s *ClientEvents) ActiveBattle(gameId string, sockets ...any) {
	// после этого события клиент отправляет на сервер "battle-connect"
	s.sendAll("active-battle", gameId, sockets)
}

// event: "battle-update"
//
// Отправка полностью измененного обьекта игры обеим
// игрокам после хода одного из игроков
func (s *ClientEvents) BattleUpdate(g *domain.BattleObj, roomId string) {
	s.server.BroadcastToRoom("/", roomId, "battle-update", g)
}