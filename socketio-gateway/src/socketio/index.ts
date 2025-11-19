import { Server } from "socket.io";
import { ClientEvents } from "@socketio/socketEvents";
import { SocketErrors, SocketData, ServerBattleRequest, ClientBattleRequest, BattleObj } from "@socketio/types";
import { match_making_service, user_service } from "@grpc/client"
import { ValidateJWTRequest, ValidateJWTResponse } from "@root/grpc/user-service/main";
import { ServiceError } from "@grpc/grpc-js";
import { GetRequestResponse, GetRequestRequest, GameConfig } from "@grpc/match_making-service/main";

function validateJwtAsync(token: string): Promise<ValidateJWTResponse> {
  return new Promise((resolve, reject) => {
    user_service.validateJwt({ token }, (err: ServiceError | null, res: ValidateJWTResponse) => {
      if (err) {
        reject(err);
      } else {
        resolve(res);
      }
    });
  });
}

function GetBattleRequestsAsync(uid: number): Promise<GetRequestResponse> {
  return new Promise((resolve, reject) => {
    match_making_service.getRequest({ uid }, (err: ServiceError | null, res: GetRequestResponse) => {
      if (err) {
        reject(err);
      } else {
        resolve(res);
      }
    });
  });
}

export function setupSocketIO(io: Server) {
  const emit = new ClientEvents(io);
  const errors = SocketErrors;

  io.on("connection", async (socket) => {
    const rawQuery = socket.handshake.query;
    const token = rawQuery.auth as string;

    let uid: string;
    try {
      uid = (await validateJwtAsync(token)).userId
    } catch (err) {
      socket.disconnect();
      emit.Error(errors.JWT_Invalid(), socket);
      return;
    }

    const data: SocketData = { uid, game: { id: "", color: "", opponentID: "" } };
    socket.data = data;
    socket.join(uid);

    // Проверка запросов на игру через match_making_service
    await (async function CheckBattleRequests() {
      let list: GetRequestResponse | null = null;
      try {
        list = await GetBattleRequestsAsync(Number(uid));
      } catch (err) {
        return console.error("Ошибка получения игр юзера: ", uid);
      }

      if (list && !list.requests) return console.log("Юзер uid: ", uid, " не имеет активных игр");

      for (let battle of list.requests) {
        // TODO
        emit.BattleRequest();
      }
    })();

    // Проверка активных игр через battle_service
    await (async function CheckActiveGames() {
      const games = await GetBattleRequestsAsync(Number(uid));
      if (!games) return console.log("юзер (uid: ", uid, ") не имеет активных игр");
      for (const game of games) {
        console.log(game);
      }
    })();

    // redis check
    (async () => {
      const uBattleData = await deps.userRepo.getRedis(uid);
      if (uBattleData?.GameId) emit.ActiveBattle(uBattleData.GameId, socket);
      if (uBattleData?.Requests?.length) {
        for (const reqId of uBattleData.Requests) {
          const r = deps.battleSrvc.getClientRequest(reqId);
          if (r) emit.BattleRequest(r, socket);
        }
      }
    })();

    socket.on("battle-request", async (r: ServerBattleRequest) => {
      const sData = socket.data as SocketData;
      if (sData.game.opponentID) {
        emit.Error(errors.Battle_Already_Requested(sData.game.opponentID), socket);
        return;
      }

      const requestObj = await deps.battleSrvc.createRequest(r.Game, sData.uid, r.Target);
      if (!requestObj) {
        emit.Error(errors.Unexpected("Ошибка запроса"), socket);
        return;
      }

      sData.game.opponentID = r.Target;
      const res: ClientBattleRequest = {
        Sender: requestObj.Sender,
        Game: r.Game,
      };

      if (isOnline(io, r.Target)) emit.BattleRequest(res, r.Target);
    });

    socket.on("battle-confirm", async (gameId: string) => {
      const reqObj = await deps.gameRepo.getRequestRedis(gameId);
      if (!reqObj) {
        emit.Error(errors.Unexpected("Такого запроса не существует"), socket);
        return;
      }

      const battleObj = await deps.battleSrvc.createGame(reqObj);
      emit.ActiveBattle(battleObj.Id, battleObj.Players.Black.Id, battleObj.Players.White.Id);
    });

    socket.on("battle-connect", async (battleId: string) => {
      socket.join(battleId);
      const clients = await io.in(battleId).allSockets();
      if (clients.size < 2) return;

      const battleObj: BattleObj = await deps.gameRepo.getGameRedis(battleId);
      const sData = socket.data as SocketData;
      sData.game.id = battleId;

      if (battleObj.Players.Black.Id === sData.uid) {
        sData.game.color = "b";
        sData.game.opponentID = battleObj.Players.White.Id;
      } else {
        sData.game.color = "w";
        sData.game.opponentID = battleObj.Players.Black.Id;
      }

      emit.BattleUpdate(battleObj, battleId);
    });

    socket.on("battle-move", (g: BattleObj) => {
      emit.BattleUpdate(g, g.Id);
    });

    socket.on("disconnect", async (reason) => {
      const data = socket.data as SocketData;
      socket.leave(data.uid);
      const uRedisData = await deps.userRepo.getRedis(data.uid);
      await deps.userRepo.saveRedis(data.uid, {
        GameId: uRedisData.GameId,
        Requests: uRedisData.Requests,
      });
    });
  });
}

function isOnline(io: Server, id: string): boolean {
  const adapter = io?.sockets?.adapter;
  if (!adapter) return false;

  const room = adapter.rooms.get(id);
  return room !== undefined && room.size > 0;
}