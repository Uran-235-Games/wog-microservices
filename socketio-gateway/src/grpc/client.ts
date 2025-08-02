import path from 'path';
import grpc from '@grpc/grpc-js';
import protoLoader from '@grpc/proto-loader';
import { YourServiceClient } from '@grpc/types';

const PROTO_PATH = path.join(__dirname, 'proto', 'auth.proto');

const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
  keepCase: true,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true
});

const protoDescriptor = grpc.loadPackageDefinition(packageDefinition) as any;
const client = new protoDescriptor.yourservice.YourService(
  'localhost:50051',
  grpc.credentials.createInsecure()
) as YourServiceClient;

export default client;
