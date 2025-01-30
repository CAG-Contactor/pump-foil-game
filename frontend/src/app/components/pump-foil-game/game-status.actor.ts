import {Actor, CollisionType, Color, Engine, Font, Text} from "excalibur";
import {GameStatus} from "./pump-foil-game";

export class GameStatusActor extends Actor {
  private statusText: Text;

  constructor(x: number, y: number, private readonly game: Engine) {
    super({
      name: "GameStatus",
      x,
      y,
      width: 40,
      height: 10,
      z: 1
    });
    this.statusText = new Text({
      text: "Stopped",
      font: new Font({ size: 30, family: "Monospace" }),
      color: Color.White
    });
    this.graphics.use(this.statusText);
  }

  ready() {
    this.statusText.text = "Ready"
  }

  stop() {
    this.statusText.text = "Stopped"
  }

  pumping() {
    this.statusText.text = "Pumping"
  }

  finish() {
    this.statusText.text = "Finished"
  }
}
