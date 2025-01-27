import {Actor, Color, Engine, Font, Text} from "excalibur";

export class TimerActor extends Actor {
  startTimeMs: number = 0;
  splitTimeMs?: number;
  finishTimeMs?: number;
  private running = false;
  private timeText;

  constructor(x: number, y: number, private readonly game: Engine) {
    super({
      name: "Timer",
      x,
      y,
      width: 40,
      height: 10,
      z: 1
    });
    this.timeText = new Text({
      text: `${(0).toFixed(2)}`,
      font: new Font({ size: 30, family: "Monospace" }),
      color: Color.White
    });
    this.graphics.use(this.timeText);
  }

  start() {
    this.startTimeMs = this.game.clock.now()
    this.running = true
  }

  stop() {
    this.finishTimeMs = this.game.clock.now()
    this.running = false;
  }

  reset() {
    this.startTimeMs = 0;
    this.finishTimeMs = undefined;
    this.splitTimeMs = undefined;
    this.running = false;
  }

  override update(engine: Engine, delta: number) {
    const now = this.game.clock.now();
    let timeSec: string = "-";
    if (!this.running && this.finishTimeMs != undefined) {
      timeSec = ((this.finishTimeMs - this.startTimeMs) / 1000).toFixed(2) ?? "-";
    } else if (this.running) {
      timeSec = ((now - this.startTimeMs) / 1000).toFixed(2);
    } else {
      timeSec = `${(0).toFixed(2)}`
    }
    const splitTimeSec: string = this.splitTimeMs != null ? ((this.splitTimeMs - this.startTimeMs)/1000).toFixed(2) : "-.--"
    this.timeText.text = `${timeSec} (${splitTimeSec})`;
  }

  splitTime() {
    this.splitTimeMs = this.game.clock.now()
  }
}
