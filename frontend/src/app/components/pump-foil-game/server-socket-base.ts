export abstract class ServerSocketBase {
  private socket?: WebSocket;

  constructor(private readonly webSocketUrl: string) {
    this.connectWebSocket();
  }

  private connectWebSocket(): void {
    console.log('Connecting to WebSocket server:',this.webSocketUrl);
    this.socket = new WebSocket(this.webSocketUrl);

    this.socket.onopen = (event) => {
      const socket = this.socket;
      if (socket) {
        console.log('WebSocket to server: %s; connection opened: %s', this.webSocketUrl, event);
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
        console.log('WebSocket to server: %s; connection closed: %s', this.webSocketUrl, event);
        this.socket = undefined;
        setTimeout(() => {
          this.connectWebSocket();
        }, 5000);
      }
    };

    this.socket.onerror = (error) => {
      const socket = this.socket;
      if (socket) {
        console.error('WebSocket to server: %s; error:', this.webSocketUrl, error);
      }
    };
  }

  protected send(message: string): void {
    this.socket?.send(message);
  }

  protected abstract handleMessage(message: GameMessage): void;
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

export type GameAdminMessage =
  InitGameMessage
  | EndGameMessage
  | AbortGameMessage
  | FetchResultlistMessage
  | ResultlistMessage
  | PingMessage;

export type GameMessage =
  GameAdminMessage
  | ControllerUpdateMessage;
