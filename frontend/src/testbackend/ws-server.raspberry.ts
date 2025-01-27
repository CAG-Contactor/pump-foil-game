import readline from 'readline';
import {WebSocket, WebSocketServer} from 'ws';

const wss = new WebSocketServer({port: 3001});

class Client {
  constructor(public ws: WebSocket, public id: string) {
  }

  send(msg: any) {
    this.ws.send(JSON.stringify(msg));
  }

  handle(object: any) {
    console.log('received: %s', object);
  }
}

let client: Client | undefined;

wss.on('connection', (ws) => {
  console.log('connected');
  client = new Client(ws, '1');

  ws.on('error', console.error);

  ws.on('message', function message(data) {
    console.log('received: %s', data);
    if (client) {
      client.handle(data);
    }
  });

  ws.on('close', function message(data) {
    console.log('closed: %s', data);
    client = undefined;
  });

  ws.send(JSON.stringify({connected: 'OK'}));
});

class KeyPressHandler {
  private turn: "left" | "right" | null = null;

  constructor() {
    readline.emitKeypressEvents(process.stdin);
    process.stdin.setRawMode(true);

    process.stdin.on('keydown', (chunk, key) => {
      console.log("keydown", key);
      switch (key?.name) {
        case 'right':
          this.turn = "right";
          client!.send({"pump": false, "turn": "right", "timestamp": new Date().getTime()});
          break;
        case 'left':
          this.turn = "left";
          client!.send({"pump": false, "turn": "left", "timestamp": new Date().getTime()});
          break;
        default:
          break;
      }
    });
    process.stdin.on('keyup', (chunk, key) => {
      console.log("keyup", key);
      switch (key?.name) {
        case 'right':
        case 'left':
          this.turn = null;
          client!.send({"pump": false, "turn": null, "timestamp": new Date().getTime()});
          break;
        default:
          break;
      }
    });
    process.stdin.on('keypress', (chunk, key) => {
      console.log("keypress", key);
      switch (key?.name) {
        case 'q':
          process.exit();
          break;
        case 'space':
          client!.send({"pumping": true, "turn": this.turn, "timestamp": new Date().getTime()});
          break;
        case "right":
        case "left":
          this.turn = null;
          client!.send({"pumping": false, "turn": key.name, "timestamp": new Date().getTime()});
          break;
        default:
          break;
      }
    });
  }

  start() {
    console.log('started key press handler');
  }
}

setInterval(() => {
  if (client) {
    client.send({type: "Ping", ping: new Date()});
  }
}, 60000);

const keyPressHandler = new KeyPressHandler();
keyPressHandler.start();

console.log('WebSocket server running on ws://localhost:3001');
