package domain

import "go.mongodb.org/mongo-driver/v2/bson"

type PublicUser struct {
	Id		string		`json:"id"`
	Name	string		`json:"name"`
}

type User struct {
	Name      string  `json:"name" bson:"name"`
	Email     string  `json:"email" bson:"email"`
	Password  string  `json:"password" bson:"password"`
	Id	string		`json:"id"`
	Token string `json:"token"`
}

type UserRedisData struct {
	// id активной игры или ""
	GameId string	`json:"game" bson:"game"`

	// хранит список id юзеров отправивших запрос
	Requests []string	`json:"requests" bson:"requests"`
}

type DBUser struct {
	ID        bson.ObjectID `json:"id" bson:"_id"`
	Name      string  `json:"name" bson:"name"`
	Email     string  `json:"email" bson:"email"`
	Password  string  `json:"password" bson:"password"`
	// User
	// BattleData UserBattleData `json:"battle_data" bson:"battle_data" binding:"omitempty"`
}

func (u *DBUser) GetStrId() string {
	return u.ID.Hex()
}