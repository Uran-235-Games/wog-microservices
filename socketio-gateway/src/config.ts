import dotenv from 'dotenv';
dotenv.config();

const config = {
  server: {
    port: Number(process.env.PORT) || 1486,
  },
  grpc: {
    url: process.env.GRPC_URL || 'localhost:50051',
  }
};

export default config;
