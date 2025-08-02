import { Client } from '@grpc/grpc-js';
import { ServiceClientConstructor } from '@grpc/grpc-js/build/src/make-client';

export interface YourServiceClient extends Client {
  GetData(
    call: { id: string },
    callback: (error: Error | null, response: { result: string }) => void
  ): void;
}
