package battle

import (
	"errors"
	"log/slog"
	"strconv"
	"time"

	// "wog-server/internal/logger"
	"wog-server/domain"
	"wog-server/internal/db"
	"wog-server/internal/logger"
)

var (
	ErrInvalidRequest = errors.New("Ошибка при создании запроса")
	ErrUserNotFound = errors.New("Юзера с таким id не существует")
)

type BattleService struct {
	gameRepo db.IGameRepo
	userRepo db.IUserRepo
}

func NewBattleService(gameRepo db.IGameRepo) *BattleService {
	return &BattleService{gameRepo: gameRepo}
}

// создает Request и сохраняет его в редис
func (s *BattleService) CreateRequest(gs domain.BattleSettings, senderId string, targetId string) (*domain.BattleRequest, error) {
	id := strconv.FormatInt(time.Now().UnixMilli(), 16)

	senderU, err := s.userRepo.Get(senderId)
	if err != nil {
		return nil, err
	}
	if senderU == nil {
		return nil, ErrUserNotFound
	}

	request := domain.BattleRequest{
		Status: "request",
		Id: id,
		Sender: domain.PublicUser{
			Id: senderU.GetStrId(),
			Name: senderU.Name,
		},
		Target: targetId,
		Game: gs,
	}

	_, err = s.gameRepo.SaveRequestRedis(request)
	if err != nil {
		logger.Log.Error("Ошибка сохранения BattleRequest в редис", slog.String("err", err.Error()))
		return nil, errors.New("Ошибка сохранения запроса")
	}
	return &request, nil
}

// создает обьект игры и сохраняет в редис
func (s *BattleService) CreateGame(br *domain.BattleRequest) (*domain.BattleObj, error) {
	// boardArr := create2DArray(settings.Size.X, settings.Size.Y)
	targetU, err := s.userRepo.Get(br.Target)
	if err != nil {
		return nil, err
	}
	if targetU == nil {
		return nil, ErrUserNotFound
	}

	players := domain.Players{
		White: br.Sender,
		Black: domain.PublicUser{
			Id: targetU.GetStrId(),
			Name: targetU.Name,
		},
	}
	game := domain.Game{
		Board: []string{"test", "katka"},
		Moves: []domain.Move{domain.Move{X: 2, Y: 2, Color: "niga"}},
		BattleSettings: br.Game,
	}
	battleObj := domain.BattleObj{
		Status: "active",
		Id: br.Id,
		Players: players,
		Game: game,
	}

	if _, err = s.gameRepo.SaveGameRedis(battleObj); err != nil {
		return nil, err
	}

	return &battleObj, nil
}

// возвращает nil в случае ошибки при обработке запроса либо при его отсутствии
func (s *BattleService) GetClientRequest(reqId string) (*domain.ClientBattleRequest) {
	r, err := s.gameRepo.GetRequestRedis(reqId)
	if err != nil {
		logger.Log.Error("Ошибка получения BattleRequest из Redis", slog.String("err", err.Error()))
		return nil
	}
	if r == nil {
		return nil
	}

	creq := domain.ClientBattleRequest{
		Sender: r.Sender,
		Game: r.Game,
	}

	return &creq
}

func create2DArray(rows uint8, cols uint8) domain.Board {
	arr := make([]string, rows*cols)
	return arr
}