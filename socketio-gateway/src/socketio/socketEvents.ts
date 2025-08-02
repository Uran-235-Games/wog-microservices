import { Server, Socket } from "socket.io";
import { SocketErr, ClientBattleRequest, BattleObj } from "./types";

export class ClientEvents {
  constructor(private io: Server) {}

  private sendAll(event: string, data: any, sockets: (string | Socket)[]) {
    for (const socket of sockets) {
      if (typeof socket === "string") {
        console.log("отправка ошибки в комнату");
        this.io.to(socket).emit(event, data);
      } else if (socket.emit) {
        console.log("отправка ошибки сокету");
        socket.emit(event, data);
      } else {
        console.error("в sendAll передан инвалидный тип сокета");
      }
    }
  }

  Error(err: SocketErr, ...sockets: (string | Socket)[]) {
    this.sendAll("error", err, sockets);
  }

  BattleRequest(r: ClientBattleRequest, ...sockets: (string | Socket)[]) {
    this.sendAll("battle-request", r, sockets);
  }

  ActiveBattle(gameId: string, ...sockets: (string | Socket)[]) {
    this.sendAll("active-battle", gameId, sockets);
  }

  BattleUpdate(g: BattleObj, roomId: string) {
    this.io.to(roomId).emit("battle-update", g);
  }
}
