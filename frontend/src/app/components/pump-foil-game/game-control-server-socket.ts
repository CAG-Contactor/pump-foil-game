import {GameController} from "./pump-foil-game";
import {GameMessage, PumpControlUpdateMessage, ServerSocketBase} from "./server-socket-base";

const RASPBERRY_SERVER= "localhost:3001";
const wsUrl = localStorage.getItem("game-control-server-socket-url") ?? `ws://${RASPBERRY_SERVER}/game-control`

export class GameControlServerSocket extends ServerSocketBase {
  constructor(private readonly gameController: GameController) {
    super(wsUrl)
  }


  protected override messageJsonToObject(data: string): GameMessage {
    const jsObject = JSON.parse(data);
    if (jsObject.pump != null || jsObject.turn != null) {
      return {type: "PumpControlUpdate", ...jsObject} as PumpControlUpdateMessage;
    } else {
      return super.messageJsonToObject(data);
    }
  }

  protected handleMessage(message: GameMessage) {
    switch (message.type) {
      case "PumpControlUpdate":
        this.gameController.handleControllerEvent(message);
        break;
      case "Ping":
        console.log('Ping', message);
        break;
      default:
        console.log('Unexpected', message);
    }
  }
}
