package socketio

import (
	"log"
	"log/slog"

	"net/url"
	"wog-server/domain"
	"wog-server/internal/db"
	"wog-server/internal/lib/jwt"
	"wog-server/internal/logger"
	"wog-server/internal/service/battle"
	
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

type socketData struct {
	uid string	// userID
	game struct {
		id string
		color string
		opponentID string
	}
}

var server *socketio.Server
var errors domain.SocketErrors

func SetupSocketIO(r *gin.Engine, jwtLib *jwt.JWTSrv, userRepo db.IUserRepo, gameRepo db.IGameRepo, battleSrvc *battle.BattleService) {
	server = socketio.NewServer(nil)
	emit := ClientEvents{server: server}

	authSocket := func(s socketio.Conn) {
		log.Println("3")
		rawQuery := s.URL().RawQuery
		values, err := url.ParseQuery(rawQuery)
		if err != nil {
			logger.Log.Debug("Error parsing socket query params", slog.String("err", err.Error()))
			emit.Error(errors.Connection_Error("Error parsing socket query params"), s)
			s.Close()
			return
		}
		
		token := values.Get("auth")
		uid, err := jwtLib.Validate(token)
		if err != nil {
			logger.Log.Debug("Error socket jwt validation", slog.String("err", err.Error()))
			emit.Error(errors.JWT_Invalid(), s)
			s.Close()
			return
		}
	
		data := socketData{ uid: uid }
		s.SetContext(data)
		logger.Log.Debug("New socket connected", slog.String("uid", data.uid))
		s.Join(uid)
		
		go func() {
			if uBattleData, err := userRepo.GetRedis(uid); err != nil {
				logger.Log.Error("Ошибка получения UserRedisData из редис", slog.String("err", err.Error()))
			} else if uBattleData != nil {
				// проверка состоит ли юзер в игре
				if uBattleData.GameId != "" {
					emit.ActiveBattle(uBattleData.GameId)
				}
				// проверка наличия запросов на игру
				if len(uBattleData.Requests) > 0 {
					for _, reqId := range uBattleData.Requests {
						if r := battleSrvc.GetClientRequest(reqId); r != nil {
							emit.BattleRequest(r, s)
						}
					}
				}
			}
		}()
	
		// TODO: хранение в редисе ограничено по времени, после истечения сохранять в mongoDB
	}

	server.OnConnect("/", func(s socketio.Conn) error {
		logger.Log.Debug("New noname socket Connected")
		go authSocket(s)
		return nil
	})

	server.OnEvent("/", "battle-request", func(s socketio.Conn, r domain.ServerBattleRequest) {
		// возврат ошибки в случае наличия активного запроса на игру
		sData := s.Context().(socketData)
		if sData.game.opponentID != "" {
			emit.Error(errors.Battle_Already_Requested(sData.game.opponentID))
			return
		}

		requestObj, err := battleSrvc.CreateRequest(r.Game, sData.uid, r.Target)
		if err != nil {
			emit.Error(errors.Unexpected(err.Error()))
			return
		}

		sData.game.opponentID = r.Target
		res := domain.ClientBattleRequest{
			Sender: requestObj.Sender,
			Game: r.Game,
		}

		// отправка запроса если оппонент в сети
		if isOnline(r.Target) {
			logger.Log.Debug("Оппонент в сети, отправка запроса на игру")
			emit.BattleRequest(&res, r.Target)
		}
	})

	server.OnEvent("/", "battle-confirm", func(s socketio.Conn, gameId string) {
		// все данные заполнены и для Black и для White

		// получения BattleRequest из Redis
		reqObj, err := gameRepo.GetRequestRedis(gameId)
		if err != nil {
			emit.Error(errors.Unexpected(err.Error()))
			return
		}
		if reqObj == nil {
			emit.Error(errors.Unexpected("Такого запроса на игру не существует"))
			return
		}
		
		// формирования обьекта игры
		battleObj, err := battleSrvc.CreateGame(reqObj)

		emit.ActiveBattle(battleObj.Id, battleObj.Players.Black.Id, battleObj.Players.White.Id)
		// после получения "active-battle" клиенты отправляют "battle-connect"
	})

	server.OnEvent("/", "battle-connect", func(s socketio.Conn, battleId string) {
		// подключение к комнате игры и ожидание обоих игроков
		s.Join(battleId)
		if rLen := server.RoomLen("/", battleId); rLen < 2 { return }

		battleObj, err := gameRepo.GetGameRedis(battleId)
		if err != nil {
			log.Fatalf("Ошибка получения battleObj из редис: %s", err.Error())
			emit.Error(errors.Unexpected("Ошибка получения battleObj из редис"))
			return
		}

		// установка значений в обьект сокета: gameId, color, opponentId
		sData := s.Context().(socketData)
		sData.game.id = battleId
		sData.game.opponentID = battleObj.Players.Black.Id
		sData.game.color = "b"
		if sData.game.opponentID != sData.uid {
			sData.game.opponentID = battleObj.Players.White.Id
			sData.game.color = "w"
		}

		emit.BattleUpdate(battleObj, battleId)
	})

	server.OnEvent("/", "battle-move", func(s socketio.Conn, g domain.BattleObj) {
		emit.BattleUpdate(&g, g.Id)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		context := s.Context()
		if context == nil {return}
		data := s.Context().(socketData)
		
		logger.Log.Debug("Socket disconnected", slog.String("reason", reason), slog.String("uid", data.uid))
		s.Leave(data.uid)
		server.ClearRoom("/", data.uid)

		// сохранение в редис domain.UserRedisData
		if uRedisData, err := userRepo.GetRedis(data.uid); err != nil {
			logger.Log.Error("Ошибка получения UserRedisData из редис", slog.String("err", err.Error()))
		} else {
			userRepo.SaveRedis(data.uid, domain.UserRedisData{
				GameId: uRedisData.GameId,
				Requests: uRedisData.Requests,
			})
		}
	})

	go server.Serve()
	r.GET("/socket.io/*any", gin.WrapH(server))
	r.POST("/socket.io/*any", gin.WrapH(server))
}

// проверка онлайн ли сокет
func isOnline(id string) bool {
	if rLen := server.RoomLen("/", id); rLen < 1 {
		return false
	}
	return true
}

/*
	Обработчики на клиенте:
		"error"
		"battle-request"
		"battle-confirm"
		"battle"

	1. Юзер формирует обьект BattleRequest с настройками игры и id соперника
	2. Отправляет на "battle-request"
	3. на "battle-confirm" приходит обьект с игрой
	   или на "error" приходит тип ошибки "battle-rejected"
	4. катка начинается:
	   обеим клиентам приходит "battle" с обьектом игры для отображения катки
	5. ход игрока:
	   клиент формирует обьект хода и отправляет серверу на "battle"
*/