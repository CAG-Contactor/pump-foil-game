import express from 'express';
import cors from "cors"
import { createServer } from 'http'; // HTTP-server
import { WebSocketServer } from 'ws'; // WebSocket-server
import readline from 'readline';

const PORT = 3003; // Gemensam port för både REST och WebSocket

// Skapa en Express-applikation
const app = express();
app.use(express.json()); // Middleware för JSON
app.use(cors());

// Skapa en klass för att hantera WebSocket-klienter
class Client {
  constructor(public ws: WebSocket, public id: string) {}

  send(msg: any) {
    this.ws.send(JSON.stringify(msg)); // Skickar JSON-meddelande till WebSocket-klienten
  }

  handle(message: any) {
    console.log('Received WebSocket message:', message);
  }
}

let client: Client | undefined;

// HTTP-endpoints
app.post('/api/v1/game-finish', (req, res) => {
  console.log('Received game-finish request', req.body);
  res.send();
});


// Skapa en HTTP-server med Express
const httpServer = createServer(app);

// Skapa en WebSocketServer och anslut den till HTTP-servern
const wss = new WebSocketServer({ server: httpServer });

wss.on('connection', (ws: any) => {
  console.log('New WebSocket client connected');
  client = new Client(ws, '1');

  ws.on('message', (data: any) => {
    console.log('Received WebSocket message:', data.toString());
    if (client) {
      client.handle(data.toString());
    }
  });

  ws.on('close', () => {
    console.log('WebSocket client disconnected');
    client = undefined;
  });

  ws.send(JSON.stringify({ type: "WSConnectionStatus", connected: 'OK' })); // Skicka anslutningsbekräftelse
});

// Tangenttryckningshanterare
class KeyPressHandler {
  constructor() {
    readline.emitKeypressEvents(process.stdin);
    process.stdin.setRawMode(true);

    process.stdin.on('keypress', (chunk, key) => {
      switch (key?.name) {
        case 'i':
          if (client) {
            console.log('InitGame via keyboard');
            client.send({ type: "InitGame", userid: "user@example.com", name: "Player1" });
          }
          break;
        case 'a':
          if (client) {
            console.log('AbortGame via keyboard');
            client.send({ type: "AbortGame" });
          }
          break;
      }
    });
  }

  start() {
    console.log('KeyPress handler started');
  }
}

const keyPressHandler = new KeyPressHandler();
keyPressHandler.start();

// Skicka regelbundna ping-meddelanden till WebSocket-klienten
setInterval(() => {
  if (client) {
    client.send({ type: "Ping", timestamp: new Date() });
  }
}, 10000);

// Starta HTTP-servern som hanterar både REST och WebSocket
httpServer.listen(PORT, () => {
  console.log(`Server is running on http://localhost:${PORT} (REST & WebSocket)`);
});
