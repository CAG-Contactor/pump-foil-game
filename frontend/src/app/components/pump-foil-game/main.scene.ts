import {Actor, CollisionType, Color, Engine} from "excalibur";
import {PortName} from "./constants";
import {Port, Ports} from "./port.actor";

export function addWalls(game: Engine) {
  [
    {name: "left", x: 0, y: 0, width: 5, height: game.canvasHeight*2},
    {name: "top", x: 0, y: 0, width: game.canvasWidth*2, height: 5},
    {name: "bottom", x: 0, y: game.canvasHeight, width: game.canvasWidth*2, height: 5},
    {name: "right", x: game.canvasWidth, y: 0, width: 5, height: game.canvasHeight*2},
  ].forEach(({name, x, y, width, height}) => {
    console.log(name, x, y, width, height)
    game.add(new Actor({
      name,
      x,
      y,
      height,
      width,
      collisionType: CollisionType.Fixed,
      color: Color.Yellow
    }));

  })
}

export function addPorts(game: Engine) {
  [
    {name: PortName.StartPort, x: 100, y: 100, rotation: 90},
    {name: PortName.Port2, x: 350, y: 135, rotation: 90},
    {name: PortName.Port3, x: 600, y: 75, rotation: 90},
    {name: PortName.Port4, x: 750, y: 300 , rotation: 180},
    {name: PortName.Port5, x: 550, y: 500 , rotation: 270},
    {name: PortName.Port6, x: 320, y: 500 , rotation: 270},
    {name: PortName.FinishPort, x: 100, y: 500 , rotation: 270},
  ].forEach(({name, x, y, rotation}) => {
    game.add(new Port(name, x, y, rotation));
  })
}

export function createPorts(): Ports {
  return new Ports([
    {name: PortName.StartPort, x: 100, y: 100, rotation: 90},
    {name: PortName.Port2, x: 350, y: 135, rotation: 90},
    {name: PortName.Port3, x: 600, y: 75, rotation: 90},
    {name: PortName.Port4, x: 750, y: 300 , rotation: 180},
    {name: PortName.Port5, x: 550, y: 500 , rotation: 270},
    {name: PortName.Port6, x: 320, y: 500 , rotation: 270},
    {name: PortName.FinishPort, x: 100, y: 500 , rotation: 270},
  ].map(({name, x, y, rotation}) => {
    const port = new Port(name, x, y, rotation);
    if (name === PortName.StartPort) {
      port.activeTarget = true;
    }
    return port;
  }));
}
