import {Actor, Canvas, Collider, CollisionContact, CollisionType, Engine, Input, Side, vec} from "excalibur";
import {PumpControlUpdateMessage} from "./server-socket-base";

export class Surfer extends Actor {
  running = false;
  private speed = 0.0;
  private direction = 0.0; // 0-359 degrees
  private readonly canvas = new Canvas({
    width: 40,
    height: 10,
    cache: false, // If true draw once until flagged dirty again, otherwise draw every time
    draw: (ctx: CanvasRenderingContext2D) => {
      ctx.fillStyle = 'red';
      ctx.fillRect(0, 0, 30, 10);
      ctx.fillStyle = 'white';
      ctx.fillRect(31, 0, 9, 10);
    }
  });
  private pumpCounter: number = 0;
  private turn: "left"|"right"|null = null;
  private colliding: boolean = false;

  constructor(x: number = 400, y: number = 300, initialDirection: number = 180) {
    super({
      name: "Surfer",
      x,
      y,
      width: 40,
      height: 10,
      z: 10
    });
    this.body.collisionType = CollisionType.Passive;
    this.direction = initialDirection;
  }

  override onInitialize(engine: ex.Engine) {
    // set as the "default" drawing
    this.graphics.use(this.canvas);
  }

  override update(engine: Engine, delta: number) {
    super.update(engine, delta);
    // Inner helpers
    const updateFromController = () => {
      if (this.pumpCounter > 0) {
        this.speed = this.speed + increase(this.speed, this.pumpCounter);
        console.log("speed", this.speed, "pump", this.pumpCounter)
        this.pumpCounter = 0;
      }

      switch (this.turn) {
        case "left":
          this.direction = (this.direction - 5) % 360;
          this.turn = null;
          break;
        case "right":
          this.direction = (this.direction + 5) % 360;
          this.turn = null;
          break;
        default:
          break;
      }
    }

    // Method body
    if (this.running) {
      if (this.speed > 0.0) {
        this.speed = this.speed - 0.1;
      } else if (this.speed < 0.0) {
        this.speed = this.speed + 0.1;
      }

      updateFromController();
    }

    const currentDirection = vec(
      Math.sin(this.direction * (Math.PI / 180)),
      -(Math.cos(this.direction * (Math.PI / 180)))
    )
    this.vel = currentDirection.normalize().scale(this.speed);
    this.rotation = currentDirection.toAngle();
  }

  override onCollisionStart(self: Collider, other: Collider, side: Side, contact: CollisionContact) {
    super.onCollisionStart(self, other, side, contact);
    if (other.owner instanceof Actor && other.owner.body.collisionType === CollisionType.Passive) {
      return;
    }

    // Bounce back
    this.speed = -15;
    this.colliding = true;
  }


  override onCollisionEnd(self: Collider, other: Collider, side: Side, lastContact: CollisionContact) {
    super.onCollisionEnd(self, other, side, lastContact);
    this.colliding = false;
  }

  start() {
    console.log("start");
    this.running = true;
  }

  handleControlChange(pumpControlMessage: PumpControlUpdateMessage) {
    if (pumpControlMessage.pump) {
      this.pumpCounter = this.pumpCounter + 1;
    }
    if (pumpControlMessage.turn != null) {
      this.turn = pumpControlMessage.turn;
    }
  }

  reset(xPos: number, yPos: number) {
    console.log("reset");
    this.running = false;
    this.pos.setTo(xPos, yPos);
    this.speed = 0.0;
    this.direction = 0.0;
  }
}
function increase(currentSpeed: number, pumpCounter: number) {
    if (currentSpeed < 5) {
      return pumpCounter*5
    } else if (currentSpeed < 10) {
      return pumpCounter*10
    } else if (currentSpeed < 20) {
      return pumpCounter*10
    } else if (currentSpeed < 30) {
      return pumpCounter*10
    } else if (currentSpeed < 40) {
      return pumpCounter*5;
    } else if (currentSpeed < 50) {
      return pumpCounter*2;
    } else if (currentSpeed < 60) {
      return 0;
    } else {
      return -3;
    }
}

