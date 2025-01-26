import readline from 'readline';
import {WebSocket, WebSocketServer} from 'ws';
import express from 'express';


const wss = new WebSocketServer({port: 3002});

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
  constructor() {
    readline.emitKeypressEvents(process.stdin);
    process.stdin.setRawMode(true);

    process.stdin.on('keypress', (chunk, key) => {
      switch (key?.name) {
        case 'q':
          process.exit();
          break;
        case "i":
          if (client) {
            client.send({type: "InitGame", userid: "x@y.z", name: "Kalle Banan"});
          }
          break;
        case "a":
          if (client) {
            client.send({type: "AbortGame"});
          }
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
}, 10000);

const keyPressHandler = new KeyPressHandler();
keyPressHandler.start();

console.log('WebSocket server running on ws://localhost:3002');
