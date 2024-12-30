Pump Foil Game
==============
Spelet styrs via kontrollenhet i backend. 
Events skickas till och från frontend via Websocket.

Livscykel
---------
Följande sekvensdiagram beskriver livscykeln för spelet.

```plantuml
actor Admin
actor Pumper
participant GameFrontend
participant GameBackend

Admin -> GameBackend: InitGame
GameBackendkalle@company.comeFrontend: InitGame(userid, name)
"name": "Kalle"
note over GameFrontend
Nollställ spel:
- nollställ tidmätning
- ställ spelaren vid startlinjen
- visa spelplan
endnote
loop until game_finished
  Pumper -> GameBackend: ControllerUpdate(frequency, tilt)
  GameBackend -> GameFrontend: ControllerUpdate(frequency, tilt)
  note over GameFrontend
    uppdatera speltillstånd:
    - ändra position för spelaren
  endnote
  alt if mållinje passerad (i korrekt ordning)
   GameFrontend -> GameBackend: EndGame(splitTime, finishTime)
   note over GameFrontend
     game_finished = true
   endnote  
  else 
   Admin -> GameBackend: AbortGame
   GameBackend -> GameFrontend: AbortGame  
   note over GameFrontend
     game_finished = true
   endnote  
  end
end
note over GameFrontend
  Inget spel pågår
  Visa resultat i 20 sekunder
endnote
GameFrontend -> GameBackend: FetchResultList
GameFrontend <-- GameBackend: ResultList
note over GameFrontend
  Visa sedan resultatlista
endnote  
```

Protokoll
---------
* InitGame(username)
* ControllerUpdate(frequency, tilt)
* EndGame(splitTime, finishTime)
* AbortGame
* FetchResultList
* ResultList

Events
------
### InitGame
```json
{
  "type": "InitGame",
  "userid": "kalle@company.com",
  "name": "Kalle"
}
```

### ControllerUpdate
Frekvensen anges i Hz och lutningen är ett grader värde mellan -90 och 90.
```json
{
  "type": "ControllerUpdate",
  "frequency": 5,
  "tilt": 10
}
```

### EndGame
SplitTime och FinishTime anges i millisekunder.
```json
{
  "type": "EndGame",
  "splitTime": 10,
  "finishTime": 20
}
```

### AbortGame
```json
{
  "type": "AbortGame"
}
```

### FetchResultList
```json
{
  "type": "FetchResultList"
}
```

### ResultList
```json
{
  "type": "ResultList",
  "results": [
    {
      "username": "player1",
      "splitTime": 10,
      "finishTime": 20
    },
    {
      "username": "player2",
      "splitTime": 15,
      "finishTime": 25
    }
  ]
}
```
