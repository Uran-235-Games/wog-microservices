export interface SocketErr {
  type: string;
  message: string;
}

export interface ServerBattleRequest {
  Game: any;
  Target: string;
}

export interface ClientBattleRequest {
  Sender: string;
  Game: any;
}

export interface BattleObj {
  Id: string;
  Players: {
    Black: { Id: string };
    White: { Id: string };
  };
  // другие поля
}

export interface SocketData {
  uid: string;
  game: {
    id: string;
    color: string;
    opponentID: string;
  };
}

export const SocketErrors = {
  Connection_Error: (msg: string) => ({ type: "connection", message: msg }),
  JWT_Invalid: () => ({ type: "auth", message: "JWT is invalid" }),
  Battle_Already_Requested: (id: string) => ({ type: "battle", message: `Already requested: ${id}` }),
  Unexpected: (msg: string) => ({ type: "unexpected", message: msg }),
};
