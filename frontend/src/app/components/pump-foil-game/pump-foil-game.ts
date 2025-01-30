import {HttpClient} from "@angular/common/http";
import {Engine} from "excalibur";
import {ConnectionStatusActor} from "./connection-status.actor";
import {PortName} from "./constants";
import {GameAdminServerSocket, LeaderBoardEntry} from "./game-admin-server-socket";
import {GameControlServerSocket} from "./game-control-server-socket";
import {addPorts, addWalls, createPorts} from "./main.scene";
import {portEvents, PortEvents, PortPassedEvent, Ports} from "./port.actor";
import {AbortGameMessage, EndGameMessage, InitGameMessage, PumpControlUpdateMessage} from "./server-socket-base";
import {GameStatusActor} from "./game-status.actor";
import {Surfer} from "./surfer.actor";
import {TimerActor} from "./timer.actor";

export type GameStatus = "Stopped"|"Finished" | "Ready" | "Pumping";

export interface GameAdminstration {
  initGame(gameEvent: InitGameMessage): void;
  abortGame(gameEvent: AbortGameMessage): void;
  endGame(message: EndGameMessage): void;
}

export interface GameController {
  handleControllerEvent(gameEvent: PumpControlUpdateMessage): void;
}

export class PumpFoilGame implements GameAdminstration, GameController {
  private static readonly SURFER_X_INIT = 50;
  private static readonly SURFER_Y_INIT = 100;

  private gameRunning = false;
  private startTimeStamp?: number;
  private endTimeStamp?: number;
  private splitTimeStamp?: number;
  private readonly gameControlServerSocket: GameControlServerSocket ;
  private readonly gameAdminServerSocket: GameAdminServerSocket ;
  private game!: Engine<any>;
  private surfer!: Surfer;
  private timer!: TimerActor;
  private ports!: Ports;
  private leaderBoard: Array<LeaderBoardEntry> = [];
  private status!: GameStatusActor;
  private connectionStatus!: ConnectionStatusActor;

  constructor(private readonly http: HttpClient) {
    this.gameControlServerSocket = new GameControlServerSocket(this);
    this.gameControlServerSocket.attachConnectionStatusListener(status => this.connectionStatus.setPumperConnectionStatus(status === "Open"));
    this.gameAdminServerSocket = new GameAdminServerSocket(this, this.http);
    this.gameAdminServerSocket.attachConnectionStatusListener(status => this.connectionStatus.setAdminConnectionStatus(status === "Open"));
  }

  init(): void {
    if (this.game) {
      return;
    }
    // If no dimensions are specified, the game will fit to the screen.
    this.game = new Engine({
      canvasElementId: "pump-foil-game",
    });

    this.ports = createPorts();
    this.ports.reset();
    this.surfer = new Surfer(PumpFoilGame.SURFER_X_INIT,PumpFoilGame.SURFER_Y_INIT, 90);
    this.timer = new TimerActor(400,25, this.game);
    this.status = new GameStatusActor(100,25, this.game);
    this.connectionStatus = new ConnectionStatusActor(776,10, this.game);

    this.game.add(this.surfer)
    this.game.add(this.timer)
    this.game.add(this.status)
    this.game.add(this.connectionStatus)
    addWalls(this.game);
    this.ports.portsList.forEach(port => this.game.add(port));

    portEvents.on(PortEvents.PortPassed, async (event: PortPassedEvent) => {
      console.log("Port Passed", event.port.name, this.game.clock.now(), this.gameRunning)
      if (!this.gameRunning) {
        return;
      }
      this.ports.activateNextAfter(event.port.name)
      if (!this.gameRunning) {
        return;
      }
      if (event.port.name === PortName.StartPort) {
        this.timer.start();
        this.startTimeStamp = this.timer.startTimeMs;
        this.status.pumping();
      }
      if (event.port.name === PortName.Port4) {
        this.timer.splitTime();
        this.splitTimeStamp = this.timer.splitTimeMs;
      }
      if (event.port.name === PortName.FinishPort) {
        this.timer.stop();
        this.endTimeStamp = this.timer.finishTimeMs;
        this.gameRunning = false;
        this.leaderBoard = await this.finishGame();
        this.status.finish();
      }
    });

    // Start the engine to begin the game.
    this.game.start();
  }

  private stopAndResetGame() {
    this.surfer.reset(PumpFoilGame.SURFER_X_INIT, PumpFoilGame.SURFER_Y_INIT, 90)
    this.gameRunning = false;
    this.startTimeStamp = undefined;
    this.endTimeStamp = undefined;
    this.splitTimeStamp = undefined;
    this.ports.reset()
    this.timer.reset();
    this.status.stop();
  }

  initGame(gameEvent: InitGameMessage): void {
    console.log('initGame', gameEvent);

    this.stopAndResetGame();
    this.surfer.start();
    this.gameRunning = true;
    this.status.ready();
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
    // TODO öppna leaderboard efter 30 s
    // this.stopAndResetGame();
  }

  private async finishGame(): Promise<Array<LeaderBoardEntry>> {
    console.log("Split time", (this.splitTimeStamp! - this.startTimeStamp!)/1000.0, "seconds")
    console.log("Time taken", (this.endTimeStamp! - this.startTimeStamp!)/1000.0, "seconds")

    return this.gameAdminServerSocket.finishGame(this.splitTimeStamp!-this.startTimeStamp!, this.endTimeStamp!-this.startTimeStamp!);
  }
}
