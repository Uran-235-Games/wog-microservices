package domain

/*
Основные параметры
Полные параметры
Список игроковююю
*/

type BoardSettings struct {
	X_size int32
	Y_size int32
	// standart - 9x9, 19x19 ...
	// custom - 14x9, 1x8 ...
	// full custom - изначально измененный массив доски клиентом
	Type string
}

type BaseSettings struct {
	PlayersCount int32
	BoardSettings BoardSettings
}

// массив доски
type FullSettings struct {

}

type Player struct {
	Name string
	Id string

	Komi float32
	// число от 0 показывающее порядок игрока в очереди ходов
	Color int8
}

type RequestObj struct {
	Id string
	Players []Player
	BaseSettings
	FullSettings
	BoardArray int8
}