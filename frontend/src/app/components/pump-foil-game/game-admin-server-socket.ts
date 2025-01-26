import {HttpClient} from "@angular/common/http";
import {firstValueFrom} from "rxjs";
import {GameAdminstration} from "./pump-foil-game";
import {EndGameMessage, GameMessage, ServerSocketBase} from "./server-socket-base";

const ADMIN_SERVER = "localhost:3003";
const wsUrl = localStorage.getItem("game-admin-server-socket-url") ?? `ws://${ADMIN_SERVER}/api/v1/ws`

export interface LeaderBoardEntry {
  result: {
    splitTime: number,
    finishTime: number
  }
  contestant: {
    email: string,
    name: string
  }
}

export class GameAdminServerSocket extends ServerSocketBase {
  constructor(private readonly gameAdministration: GameAdminstration, private readonly http: HttpClient) {
    super(wsUrl)
  }

  async finishGame(splitTime: number, endTime: number): Promise<Array<LeaderBoardEntry>> {
    const message: EndGameMessage = {type: "EndGame", splitTime, endTime};
    return firstValueFrom(this.http.post<Array<LeaderBoardEntry>>(`http://${ADMIN_SERVER}/api/v1/game-finish`, message));
  }

  async getResutlist(): Promise<Array<LeaderBoardEntry>> {
    return firstValueFrom(this.http.get<Array<LeaderBoardEntry>>(`http://${ADMIN_SERVER}/api/v1/leaderboard`));
  }

  protected handleMessage(message: GameMessage) {
    switch (message.type) {
      case "InitGame":
        this.gameAdministration.initGame(message);
        break;
      case "AbortGame":
        this.gameAdministration.abortGame(message);
        break;
      case "EndGame":
        this.gameAdministration.endGame(message);
        break;
      case "Ping":
        //console.log('Ping', message);
        break;
      default:
        console.log('Unexpected', message);
    }
  }
}
