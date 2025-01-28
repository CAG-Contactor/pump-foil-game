import {
  Actor, Circle, CircleCollider,
  Collider,
  CollisionContact,
  CollisionType,
  Color,
  Engine,
  EventEmitter,
  GameEvent,
  Side,
  Text
} from "excalibur";
import {ActorEvents} from "excalibur/build/dist/Actor";
import {PortName} from "./constants";

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


function labelX(rotationDegrees: number) {
  if (rotationDegrees >= 0 && rotationDegrees <= 90) {
    return -45;
  } else if (rotationDegrees <= 180) {
    return 50;
  } else if (rotationDegrees <= 270) {
    return 45;
  } else {
    return -50;
  }
}

function labelY(rotationDegrees: number) {
  if (rotationDegrees >= 0 && rotationDegrees <= 90) {
    return 0;
  } else if (rotationDegrees <= 180) {
    return 0;
  } else if (rotationDegrees <= 270) {
    return 0;
  } else {
    return 0;
  }
}

export class Port extends Actor {
  private _activeTarget = false
  starboardBuoy: Actor
  portBuoy: Actor
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
    this.starboardBuoy = new Actor({
      name: "right stick",
      x: 30,
      y: 0,
      radius: 5,
      collisionType: CollisionType.Fixed,
      color: ACTIVE_STARBOARD_BUOY,
    })
    this.addChild(this.starboardBuoy);
    this.portBuoy = new Actor({
        name: "left stick",
        x: -30,
        y: 0,
        radius: 5,
        collisionType: CollisionType.Fixed,
        color: ACTIVE_PORT_BUOY,
      })
    this.addChild(this.portBuoy);
    const textActor = new Actor({
      x: labelX(rotationDegrees),
      y: labelY(rotationDegrees),
      rotation: -rotationDegrees * (Math.PI / 180),
      collisionType: CollisionType.PreventCollision,
    })
    const nameText = new Text({text: name, color: Color.White});
    textActor.graphics.use(nameText);
    this.addChild(textActor);
  }

  get activeTarget(): boolean {
    return this._activeTarget;
  }

  set activeTarget(value: boolean) {
    this._activeTarget = value;
  }

  override onCollisionStart(self: Collider, other: Collider, side: Side, contact: CollisionContact) {
    super.onCollisionStart(self, other, side, contact);
    console.log("Port collision", other.owner.name);
    if (this.activeTarget) {
      portEvents.emit(PortEvents.PortPassed, new PortPassedEvent(this));
    }
  }

  override update(engine: Engine, delta: number) {
    super.update(engine, delta);
    if (this._activeTarget) {
      (this.starboardBuoy.graphics.current as Circle).color = ACTIVE_STARBOARD_BUOY;
      (this.starboardBuoy.graphics.current as Circle).radius = 7;
      (this.portBuoy.graphics.current as Circle).color = ACTIVE_PORT_BUOY;
      (this.portBuoy.graphics.current as Circle).radius = 7;
    } else {
      (this.starboardBuoy.graphics.current as Circle).color = INACTIVE_STARBOARD_BUOY;
      (this.starboardBuoy.graphics.current as Circle).radius = 5;
      (this.portBuoy.graphics.current as Circle).color = INACTIVE_PORT_BUOY;
      (this.portBuoy.graphics.current as Circle).radius = 5;
    }
  }
}

export class Ports {
  private current: string = PortName.StartPort;

  constructor(readonly portsList: Array<Port>) {
  }

  reset() {
    this.current = PortName.StartPort
    this.setActivePort(this.current)
  }

  setActivePort(portName:string) {
    this.portsList.forEach(port =>
      port.activeTarget = port.name === portName
    )
  }

  activateNextAfter(portName: string) {
    if (portName === PortName.FinishPort) {
      return;
    }
    const index = this.portsList.findIndex(port => port.name === portName);
    if (index >= 0) {
      this.setActivePort(this.portsList[index + 1].name)
    }
  }
}



const ACTIVE_STARBOARD_BUOY = Color.Green
const ACTIVE_PORT_BUOY = Color.Red
const INACTIVE_STARBOARD_BUOY = Color.fromHex("99cd99FF")
const INACTIVE_PORT_BUOY = Color.fromHex("CF6565FF")
