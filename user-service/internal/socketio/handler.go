package socketio

import (
	"wog-server/internal/app"

	"fmt"
	"net/http"

	socketio "github.com/doquangtan/socket.io/v4"
)

func Setup(a *app.App) http.Handler {
	io := socketio.New()

	io.OnConnection(func(s *socketio.Socket) {
		fmt.Println("new socket")
		s.Emit("msg", "Ivan Rak")
	})

	return io.HttpHandler()
}