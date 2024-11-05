import {Actor, Canvas, Collider, CollisionContact, CollisionType, Engine, Input, Side, vec} from "excalibur";
import {Port} from "./port.actor";

export class Surfer extends Actor {
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

    if (this.speed > 0.0) {
      this.speed = this.speed - 0.1;
    } else if (this.speed < 0.0) {
      this.speed = this.speed + 0.1;
    }

    if (engine.input.keyboard.isHeld(Input.Keys.A)) {
      this.direction = (this.direction - 1) % 360;
    }

    if (engine.input.keyboard.isHeld(Input.Keys.D)) {
      this.direction = (this.direction + 1) % 360;
    }

    if (engine.input.keyboard.wasPressed(Input.Keys.L)) {
      this.speed = this.speed + 5;
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
  }
}
