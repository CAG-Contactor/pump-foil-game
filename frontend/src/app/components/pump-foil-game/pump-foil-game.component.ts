import {HttpClient,} from "@angular/common/http";
import {Component, OnInit} from '@angular/core';
import {PumpFoilGame} from "./pump-foil-game";

@Component({
  selector: 'app-pump-foil-game',
  standalone: true,
  imports: [],
  template: `
    <div class="game">
<!--      <h1>The Pump Foil Game</h1>-->
      <div class="canvas-wrapper">
        <canvas id="pump-foil-game"></canvas>
      </div>
    </div>
  `,
  styles: `
    .game {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 2rem;
    }
    .canvas-wrapper {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      height: 100%;
      width: 100%;
    }
  `
})
export class PumpFoilGameComponent implements OnInit {
  private readonly game;

  constructor(private readonly http: HttpClient) {
    this.game = new PumpFoilGame(http);
  }

  ngOnInit(): void {
    this.game.init();
  }
}
