import {Routes} from '@angular/router';
import {PumpFoilGameComponent} from "./components/pump-foil-game/pump-foil-game.component";

export const routes: Routes = [
  {path: "", redirectTo: "pump-foil-game", pathMatch: "full"},
  {path: "pump-foil-game", component: PumpFoilGameComponent},
];
