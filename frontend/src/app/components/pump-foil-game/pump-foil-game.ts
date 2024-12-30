import {Engine} from "excalibur";
import {Surfer} from "./surfer.actor";
import {addPorts, addWalls, Ports} from "./main.scene";
import {portEvents, PortEvents, PortPassedEvent} from "./port.actor";
import {AbortGameMessage, ControllerUpdateMessage, InitGameMessage, ServerSocket} from "./server-socket";


export class PumpFoilGame {
  private startTimeStamp = 0;
  private endTimeStamp = 0;
  private readonly serverSocket: ServerSocket ;

  constructor() {
    this.serverSocket = new ServerSocket(this);
    this.serverSocket.connectWebSocket();
  }

  init(): void {
    // If no dimensions are specified the game will fit to the screen.
    const game = new Engine({
      canvasElementId: "pump-foil-game",
    });
    game.add(new Surfer(100,550, 0))
    addWalls(game);
    addPorts(game);

    portEvents.on(PortEvents.PortPassed, (event: PortPassedEvent) => {
      console.log("Port Passed", event.port.name, game.clock.now())
      if (event.port.name === Ports.Port1) {
        this.startTimeStamp = game.clock.now();
      }
      if (event.port.name === Ports.Port7) {
        this.endTimeStamp = game.clock.now();
        console.log("Time taken", (this.endTimeStamp - this.startTimeStamp)/1000.0, "seconds")
      }
    });

    // Start the engine to begin the game.
    game.start();
  }

  armGame(gameEvent: InitGameMessage): void {
    console.log('armGame', gameEvent);
  }

  handleControllerEvent(gameEvent: ControllerUpdateMessage) {
    console.log('handleControllerEvent', gameEvent);
  }

  abortGame(gameEvent: AbortGameMessage) {
    console.log('abortGame', gameEvent);
  }
}

