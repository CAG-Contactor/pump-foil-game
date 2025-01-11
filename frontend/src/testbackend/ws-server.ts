import readline from 'readline';
import {WebSocket, WebSocketServer} from 'ws';
import {
  AbortGameMessage,
  ControllerUpdateMessage,
  InitGameMessage,
  PingMessage
} from "../app/components/pump-foil-game/game-control-server-socket";

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
  private readonly bufferSize: number = 5;
  private readonly fifoBuffer: Array<number> = [];
  private direction = 0.0

  constructor() {
    readline.emitKeypressEvents(process.stdin);
    process.stdin.setRawMode(true);

    process.stdin.on('keypress', (chunk, key) => {
      switch (key?.name) {
        case 'q':
          process.exit();
          break;
        case 'space':
          this.fifoBuffer.push(Date.now());
          if (this.fifoBuffer.length > this.bufferSize) {
            this.fifoBuffer.shift();
          }
          break;
        case 'right':
          this.direction = this.direction + 5;
          break;
        case 'left':
          this.direction = this.direction - 5;
          break;
        case "i":
          if (client) {
            client.send({type: "InitGame", userid: "x@y.z", name: "Kalle Banan"} as InitGameMessage);
          }
          break;
        case "a":
          if (client) {
            client.send({type: "AbortGame"} as AbortGameMessage);
          }
          break;
      }
    });

    setInterval(() => {
      this.calculateKeyPressFrequencyAndSend();
    }, 1000);
  }

  private calculateKeyPressFrequencyAndSend() {
    const prevBufLength = this.fifoBuffer.length;

    const delta_t = this.fifoBuffer.length < 3  ? 0 : (this.fifoBuffer[2] - this.fifoBuffer[0]); // in milliseconds
    if (delta_t > 0) {
      const frequency = 2 / (delta_t / 1000)
      console.log('key press frequency:', frequency, " hz", "delta_t:", delta_t, "delta_n:", this.fifoBuffer.length);
      if (client) {
        client.send({type: "ControllerUpdate", frequency: frequency, tilt: this.direction} as ControllerUpdateMessage);
      }
    } else if (this.fifoBuffer.length > 0) {
      console.log('key press frequency: 0 hz');
      if (client) {
        client.send({type: "ControllerUpdate", frequency: 0.0, tilt: 0.0} as ControllerUpdateMessage);
      }
    }

    if (this.fifoBuffer.length && this.fifoBuffer.length === prevBufLength) {
      this.fifoBuffer.shift();
    }
  }

  start() {
    console.log('started key press handler');
  }
}

setInterval(() => {
  if (client) {
    client.send({type: "Ping", ping: new Date()} as PingMessage);
  }
}, 10000);

const keyPressHandler = new KeyPressHandler();
keyPressHandler.start();

console.log('WebSocket server running on ws://localhost:3001');
