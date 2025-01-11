import {GameAdminstration} from "./pump-foil-game";
import {EndGameMessage, FetchResultlistMessage, GameMessage, ServerSocketBase} from "./server-socket-base";

const wsUrl = localStorage.getItem("game-admin-server-socket-url") ?? "ws://localhost:3001/game-admin"

export class GameAdminServerSocket extends ServerSocketBase {
  constructor(private readonly gameAdministration: GameAdminstration) {
    super(wsUrl)
  }

  endGame(splitTime: number, finishTime: number): void {
    const message: EndGameMessage = {type: "EndGame", splitTime, finishTime};
    this.send(JSON.stringify(message));
  }

  fetchResutlist(): void {
    const message: FetchResultlistMessage = {type: "FetchResultlist"};
    this.send(JSON.stringify(message));
  }

  protected handleMessage(message: GameMessage) {
    switch (message.type) {
      case "InitGame":
        this.gameAdministration.armGame(message);
        break;
      case "AbortGame":
        this.gameAdministration.abortGame(message);
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
