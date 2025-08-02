import express from 'express';
import { createServer } from 'http';
import { Server } from 'socket.io';
import { setupSocketIO } from '@socketio/index';
import config from '@root/config';

const app = express();
const httpServer = createServer(app);
const io = new Server(httpServer, {
  cors: { origin: "*" }
});



setupSocketIO(io);

app.get('/', (_req, res) => {
  res.send('Server is running');
});

httpServer.listen(config.server.port, () => {
  console.log(`Server running on http://localhost:${config.server.port}`);
});
