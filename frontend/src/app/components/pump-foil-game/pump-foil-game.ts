import {HttpClient} from "@angular/common/http";
import {Engine} from "excalibur";
import {GameAdminServerSocket, LeaderBoardEntry} from "./game-admin-server-socket";
import {GameControlServerSocket} from "./game-control-server-socket";
import {addPorts, addWalls, Ports} from "./main.scene";
import {portEvents, PortEvents, PortPassedEvent} from "./port.actor";
import {AbortGameMessage, EndGameMessage, InitGameMessage, PumpControlUpdateMessage} from "./server-socket-base";
import {Surfer} from "./surfer.actor";


export interface GameAdminstration {
  initGame(gameEvent: InitGameMessage): void;
  abortGame(gameEvent: AbortGameMessage): void;
  endGame(message: EndGameMessage): void;
}

export interface GameController {
  handleControllerEvent(gameEvent: PumpControlUpdateMessage): void;
}

export class PumpFoilGame implements GameAdminstration, GameController {
  private static readonly SURFER_X_INIT = 100;
  private static readonly SURFER_Y_INIT = 550;

  private gameRunning = false;
  private startTimeStamp?: number;
  private endTimeStamp?: number;
  private splitTimeStamp?: number;
  private readonly gameControlServerSocket: GameControlServerSocket ;
  private readonly gameAdminServerSocket: GameAdminServerSocket ;
  private game!: Engine<any>;
  private surfer!: Surfer;
  private leaderBoard: Array<LeaderBoardEntry> = [];

  constructor(private readonly http: HttpClient) {
    this.gameControlServerSocket = new GameControlServerSocket(this);
    this.gameAdminServerSocket = new GameAdminServerSocket(this, this.http);
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

    portEvents.on(PortEvents.PortPassed, async (event: PortPassedEvent) => {
      console.log("Port Passed", event.port.name, this.game.clock.now())
      if (event.port.name === Ports.Port1) {
        this.startTimeStamp = this.game.clock.now();
      }
      if (event.port.name === Ports.Port4) {
        this.splitTimeStamp = this.game.clock.now();
      }
      if (event.port.name === Ports.Port7) {
        this.endTimeStamp = this.game.clock.now();
        this.leaderBoard = await this.finishGame();
      }
    });

    // Start the engine to begin the game.
    this.game.start();
  }

  private stopAndResetGame() {
    this.surfer.reset(PumpFoilGame.SURFER_X_INIT, PumpFoilGame.SURFER_Y_INIT)
    this.gameRunning = false;
    this.startTimeStamp = undefined;
    this.endTimeStamp = undefined;
    this.splitTimeStamp = undefined;
  }

  initGame(gameEvent: InitGameMessage): void {
    console.log('initGame', gameEvent);

    this.stopAndResetGame();
    this.surfer.start();
    this.gameRunning = true;
  }

  abortGame(gameEvent: AbortGameMessage) {
    console.log('abortGame', gameEvent);

    this.stopAndResetGame();
  }

  handleControllerEvent(pumpControlMessage: PumpControlUpdateMessage) {
    console.log('handleControllerEvent', pumpControlMessage);
    this.surfer.handleControlChange(pumpControlMessage);
  }

  endGame(message: EndGameMessage) {
    console.log('endGame', message);
    this.stopAndResetGame();
  }

  private async finishGame(): Promise<Array<LeaderBoardEntry>> {
    console.log("Split time", (this.splitTimeStamp! - this.startTimeStamp!)/1000.0, "seconds")
    console.log("Time taken", (this.endTimeStamp! - this.startTimeStamp!)/1000.0, "seconds")

    return this.gameAdminServerSocket.finishGame(this.splitTimeStamp!-this.startTimeStamp!, this.endTimeStamp!-this.startTimeStamp!);
  }
}
