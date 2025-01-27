import {Actor, Collider, CollisionContact, CollisionType, Color, EventEmitter, GameEvent, Side, Text} from "excalibur";
import {ActorEvents} from "excalibur/build/dist/Actor";

export type PortEventsType = {
  PortPassed: PortPassedEvent;
}

export class PortPassedEvent extends GameEvent<Port> {
  constructor(public readonly port: Port) {
    super();
  }
}

export const PortEvents = {
  PortPassed: 'PortPassed'
} as const;

export const portEvents = new EventEmitter<ActorEvents & PortEventsType>();


export class Port extends Actor {
  constructor(name: string, x: number, y: number, rotationDegrees: number) {
    super({
      name,
      x,
      y,
      rotation: rotationDegrees * (Math.PI / 180),
      width: 50,
      height: 2,
      collisionType: CollisionType.Passive,
      color: Color.Gray,
    });
    this.addChild(new Actor({
      name: "right stick",
      x: 30,
      y: 0,
      radius: 5,
      collisionType: CollisionType.Fixed,
      color: Color.Green,
    }));
    this.addChild(new Actor({
      name: "left stick",
      x: -30,
      y: 0,
      radius: 5,
      collisionType: CollisionType.Fixed,
      color: Color.Red,
    }));
    const textActor = new Actor({
      x: 30,
      y: -10,
      collisionType: CollisionType.PreventCollision,
    })
    const nameText = new Text({text: name, color: Color.White});
    textActor.graphics.use(nameText);
    this.addChild(textActor);
  }

  override onCollisionStart(self: Collider, other: Collider, side: Side, contact: CollisionContact) {
    super.onCollisionStart(self, other, side, contact);
    console.log("Port collision", other.owner.name);
    portEvents.emit(PortEvents.PortPassed, new PortPassedEvent(this));
  }

}
