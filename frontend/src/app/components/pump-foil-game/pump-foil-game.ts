import {Engine} from "excalibur";
import {GameAdminServerSocket} from "./game-admin-server-socket";
import {GameControlServerSocket} from "./game-control-server-socket";
import {addPorts, addWalls, Ports} from "./main.scene";
import {portEvents, PortEvents, PortPassedEvent} from "./port.actor";
import {AbortGameMessage, ControllerUpdateMessage, InitGameMessage} from "./server-socket-base";
import {Surfer} from "./surfer.actor";


export interface GameAdminstration {
  armGame(gameEvent: InitGameMessage): void;
  abortGame(gameEvent: AbortGameMessage): void;
}

export interface GameController {
  handleControllerEvent(gameEvent: ControllerUpdateMessage): void;
}

export class PumpFoilGame implements GameAdminstration, GameController {
  private static readonly SURFER_X_INIT = 100;
  private static readonly SURFER_Y_INIT = 550;

  private startTimeStamp = 0;
  private endTimeStamp = 0;
  private readonly gameControlServerSocket: GameControlServerSocket ;
  private readonly gameAdminServerSocket: GameAdminServerSocket ;
  private game!: Engine<any>;
  private surfer!: Surfer;

  constructor() {
    this.gameControlServerSocket = new GameControlServerSocket(this);
    this.gameAdminServerSocket = new GameAdminServerSocket(this);
  }

  init(): void {
    // If no dimensions are specified, the game will fit to the screen.
    this.game = new Engine({
      canvasElementId: "pump-foil-game",
    });
    this.surfer = new Surfer(PumpFoilGame.SURFER_X_INIT,PumpFoilGame.SURFER_Y_INIT, 0);
    this.game.add(this.surfer)
    addWalls(this.game);
    addPorts(this.game);

    portEvents.on(PortEvents.PortPassed, (event: PortPassedEvent) => {
      console.log("Port Passed", event.port.name, this.game.clock.now())
      if (event.port.name === Ports.Port1) {
        this.startTimeStamp = this.game.clock.now();
      }
      if (event.port.name === Ports.Port7) {
        this.endTimeStamp = this.game.clock.now();
        console.log("Time taken", (this.endTimeStamp - this.startTimeStamp)/1000.0, "seconds")
      }
    });

    // Start the engine to begin the game.
    this.game.start();
  }

  private resetGame() {
    this.surfer.reset(PumpFoilGame.SURFER_X_INIT, PumpFoilGame.SURFER_Y_INIT)
  }

  armGame(gameEvent: InitGameMessage): void {
    console.log('armGame', gameEvent);
    this.resetGame();
  }

  abortGame(gameEvent: AbortGameMessage) {
    console.log('abortGame', gameEvent);
    this.resetGame();
  }

  handleControllerEvent(gameEvent: ControllerUpdateMessage) {
    console.log('handleControllerEvent', gameEvent);
  }
}

