import { createChannel, createClient } from "nice-grpc";

import type { UserServiceClient } from "../grpc/generated/user";
import { UserServiceDefinition } from "../grpc/generated/user";

let _client: UserServiceClient;

function getClient(): UserServiceClient {
  if (!_client) {
    const address = process.env.API_URL;
    if (!address) {
      throw new Error("API URL is not set");
    }
    const channel = createChannel(address);

    _client = createClient(UserServiceDefinition, channel);
  }
  return _client;
}

export async function getUser(id: number) {
  return await getClient().getUser({id});
}