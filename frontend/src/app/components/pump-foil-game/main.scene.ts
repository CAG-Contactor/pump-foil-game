import {Actor, CollisionType, Color, Engine} from "excalibur";
import {Port} from "./port.actor";

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

export const Ports = {
  Port1: "port1",
  Port2: "port2",
  Port3: "port3",
  Port4: "port4",
  Port5: "port5",
  Port6: "port6",
  Port7: "port7",
}

export function addPorts(game: Engine) {
  [
    {name: Ports.Port1, x: 100, y: 500, rotation: 0},
    {name: Ports.Port2, x: 50, y: 300, rotation: 0},
    {name: Ports.Port3, x: 200, y: 100, rotation: 90},
    {name: Ports.Port4, x: 400, y: 300 , rotation: 180},
    {name: Ports.Port5, x: 550, y: 500 , rotation: 90},
    {name: Ports.Port6, x: 700, y: 400 , rotation: 0},
    {name: Ports.Port7, x: 550, y: 100 , rotation: 0},
  ].forEach(({name, x, y, rotation}) => {
    game.add(new Port(name, x, y, rotation));
  })
}
