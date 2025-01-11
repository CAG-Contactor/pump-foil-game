import {GameController} from "./pump-foil-game";
import {GameMessage, ServerSocketBase} from "./server-socket-base";

const wsUrl = localStorage.getItem("game-control-server-socket-url") ?? "ws://localhost:3001/game-control"

export class GameControlServerSocket extends ServerSocketBase {
  constructor(private readonly gameController: GameController) {
    super(wsUrl)
  }

  protected handleMessage(message: GameMessage) {
    switch (message.type) {
      case "ControllerUpdate":
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
