import { Client, ServiceError, credentials } from "@grpc/grpc-js";
import { UserServiceClient } from "@grpc/user-service/main";
import { GoMatchMakingClient } from "./match_making-service/main";

export const user_service = new UserServiceClient(
  "localhost:1488",
  credentials.createInsecure()
)

export const match_making_service = new GoMatchMakingClient(
  "localhost:1489",
  credentials.createInsecure()
)