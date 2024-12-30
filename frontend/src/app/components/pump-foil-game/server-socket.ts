import {PumpFoilGame} from "./pump-foil-game";

export class ServerSocket {
  private socket?: WebSocket;

  constructor(private readonly game: PumpFoilGame) {
    this.connectWebSocket();
  }

  connectWebSocket(): void {
    console.log('Connecting to WebSocket server');
    this.socket = new WebSocket('ws://localhost:3001/game-control');

    this.socket.onopen = (event) => {
      const socket = this.socket;
      if (socket) {
        console.log('WebSocket connection opened:', event);
        socket.send(JSON.stringify({init: "OK"}));
      }
    };

    this.socket.onmessage = (event) => {
      const socket = this.socket;
      if (socket && event?.data) {
        this.handleMessage(JSON.parse(event.data));
      }
    };

    this.socket.onclose = (event) => {
      const socket = this.socket;
      if (socket) {
        console.log('WebSocket connection closed:', event);
        this.socket = undefined;
        setTimeout(() => {
          this.connectWebSocket();
        }, 5000);
      }
    };

    this.socket.onerror = (error) => {
      const socket = this.socket;
      if (socket) {
        console.error('WebSocket error:', error);
      }
    };
  }

  endGame(splitTime: number, finishTime: number): void {
    const message: EndGameMessage = {type: "EndGame", splitTime, finishTime};
    this.socket?.send(JSON.stringify(message));
  }

  fetchResutlist(): void {
    const message: FetchResultlistMessage = {type: "FetchResultlist"};
    this.socket?.send(JSON.stringify(message));
  }

  private handleMessage(message: GameMessage) {

    switch (message.type) {
      case "InitGame":
        this.game.armGame(message);
        break;
      case "ControllerUpdate":
        this.game.handleControllerEvent(message);
        break;
      case "AbortGame":
        console.log('AbortGame');
        break;
      case "Resultlist":
        console.log('Resultlist', message.resultlist);
        break;
      case "Ping":
        console.log('Ping', message);
        break;
      default:
        console.log('Unexpected', message);
    }
  }
}

export interface PingMessage {
  type: "Ping";
  ping: Date;
}

export interface InitGameMessage {
  type: "InitGame";
  userid: string;
  name: string;
}

export interface ControllerUpdateMessage {
  type: "ControllerUpdate";
  frequency: number; // Hz
  tilt: number; // degrees (-90..90)
}

export interface EndGameMessage {
  type: "EndGame";
  splitTime: number; // millis
  finishTime: number; // millis
}

export interface AbortGameMessage {
  type: "AbortGame";
}

export interface FetchResultlistMessage {
  type: "FetchResultlist";
}

export interface ResultlistMessage {
  type: "Resultlist";
  resultlist: Array<{ userName: string, splitTime: number, finishTime: number }>;
}

type GameMessage =
  InitGameMessage
  | ControllerUpdateMessage
  | EndGameMessage
  | AbortGameMessage
  | FetchResultlistMessage
  | ResultlistMessage
  | PingMessage;
