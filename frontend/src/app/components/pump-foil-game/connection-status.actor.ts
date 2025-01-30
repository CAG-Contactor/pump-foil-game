import {Actor, CollisionType, Color, Engine, Font, Text} from "excalibur";

export class ConnectionStatusActor extends Actor {
  private adminConnectionIndicator: Actor;
  private pumperConnectionIndicator: Actor;

  constructor(x: number, y: number, private readonly game: Engine) {
    super({
      name: "ConnectionStatus",
      x,
      y,
      width: 40,
      height: 10,
      z: 1
    });
    this.adminConnectionIndicator = new Actor({
      name: "adminConnection",
      x: 0,
      y: 0,
      z: 1000,
      radius: 5,
      collisionType: CollisionType.PreventCollision,
      color: DISCONNECTED,
    })
    this.addChild(this.adminConnectionIndicator);
    const adminText = new Text({
      text: "A",
      font: new Font({ size: 10, family: "Monospace" }),
      color: DISCONNECTED
    });
    this.adminConnectionIndicator.graphics.use(adminText);

    this.pumperConnectionIndicator = new Actor({
      name: "pumperConnection",
      x: 13,
      y: 0,
      z: 1000,
      radius: 5,
      collisionType: CollisionType.PreventCollision,
      color: DISCONNECTED,
    })
    this.addChild(this.pumperConnectionIndicator);
    const pumperText = new Text({
      text: "P",
      font: new Font({ size: 10, family: "Monospace" }),
      color: DISCONNECTED
    });
    this.pumperConnectionIndicator.graphics.use(pumperText);

  }

  setAdminConnectionStatus(connected: boolean) {
    this.adminConnectionIndicator.color = connected ? CONNECTED : DISCONNECTED;
  }

  setPumperConnectionStatus(connected: boolean) {
    this.pumperConnectionIndicator.color = connected ? CONNECTED : DISCONNECTED;
  }

}

const DISCONNECTED = Color.Red;
const CONNECTED = Color.Green;
