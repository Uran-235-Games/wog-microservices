import { Socket } from 'socket.io';

export function handleSocketConnection(socket: Socket) {
  console.log('Client connected:', socket.id);

  socket.on('message', (data) => {
    console.log('Message received:', data);
    socket.broadcast.emit('message', data);
  });

  socket.on('disconnect', () => {
    console.log('Client disconnected:', socket.id);
  });
}
