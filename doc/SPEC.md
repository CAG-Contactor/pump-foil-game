Spec, The Pump Foiler
=====================
Se även [Gamla specen](https://docs.google.com/document/d/1iw420msEF6ePYX12Zz6Wj7ceX9JUx-VIVOwcCjhF-Ug/edit?usp=sharing)

Översikt
--------
### Appar
```plantuml
rectangle gamefe [
  Game Frontend
]

rectangle adminfe [
  Admin Frontend
]

rectangle pumpinterface [
  Pumpberry
  ===
  Läsa av accelerometer
  Räkna om till sidolutning och hastighet
  Skicka via websocket 
  - sidolutning
  - hastighet
  - tidsserie vertikalposition
]

rectangle mainbe [
  Main Backend
  ===
  Tävling
  ----
  - Registrera ny tävlande
  - Ställa tävlande i kö
  - Starta spel för tävlande
    - Välja i kö och starta  
    - Skicka event till Game FE 
      - Tävlande
  - Ta emot registrering av resultat (från Game FE)
    - Mellantid
    - Sluttid
    - Tävlande 
  
  Admin
  ---
  - Ta bort från kö
  - Ta bort resultat
  - Avbryt pågående spel och ta bort
  - Avbryt pågående spel och starta om
]

rectangle toaster [
  Daniels brödrost
  ==
  Toastberry som läser QR-kod
  Skickar läst data till Main Backend
]


gamefe -- pumpinterface
adminfe -- mainbe
gamefe -- mainbe
toaster -- mainbe
```
### Datamodell

```plantuml
class PumpfoilApplication
class Contestant {
  email
  namn
}
class QueueItem {
  timeAdded
}
class LeaderboardItem {
  splitTime
  endTime
}
class GameState

PumpfoilApplication *-down-> "0,n\nContestants" Contestant
PumpfoilApplication *-down-> "0,n\nQueue" QueueItem
PumpfoilApplication *-down-> "1" GameState
PumpfoilApplication *-down-> "0,n\nLeaderboard" LeaderboardItem
QueueItem --> "1\nthePumper" Contestant
GameState --> "0,n\ncurrentPumper" Contestant
LeaderboardItem "1"  -- "1" Contestant
```

API, Main Backend
-----------------
Registrera tävlande
-------------------
```plantuml
actor Anders
actor Tävlande
participant "Admin Frontend" as adminfe
participant Brödrost
participant "Main Backend" as mainbe
database pumpfoil_db

Anders -> adminfe: visa sida för att lägga till ny tävlande
Anders <-- adminfe: visar sida
Anders -> adminfe: tryck knapp "Registrera ny tävlande"
Anders <-- adminfe: visar formulär
note over adminfe
  inmatningsfält för
  - namn
  - epost
  meddelande: väntar på uppgifter...
  samt knapp "Spara"
endnote
Tävlande -> Brödrost: visa badge
note over Brödrost
  Läser av badge
  Anropar Main Backend
endnote
Brödrost -> mainbe: POST /badge-data\nnamn & epost
mainbe -> adminfe: badge-data (WS-event)\nnamn & epost
note over adminfe
  - fyller i fält för namn & epost
  - knapp "Bekräfta"
endnote  
adminfe --> Anders: visar namn & epost
Anders -> adminfe: trycker knapp "Bekräfta"
adminfe -> mainbe: POST /register-contestan\nnamn & epost  
mainbe -> pumpfoil_db: spara ny tävlande

```
Lägg tävlande i kö
-------------------
```plantuml 
actor Anders
actor Tävlande
participant "Admin Frontend" as adminfe
participant "Main Backend" as mainbe
database pumpfoil_db

Anders -> adminfe: visa sida för att lägga till i kö
adminfe -> mainbe: hämta tävlande
mainbe -> pumpfoil_db: läs tävlande
mainbe <-- pumpfoil_db: tävlande
mainbe -> pumpfoil_db: läs kö
mainbe <-- pumpfoil_db: kö
adminfe <-- mainbe: tävlande förutom de i kö
Anders <-- adminfe: visar sida
note over adminfe
  Visar en sida för att lägga till i kö
  - lista på alla tävlande (som inte är i kö)
  - knapp på varje tävlande "Lägg till i kö"
endnote
Anders -> adminfe: klicka på "Lägg till i kö"\nför en tävlande
adminfe -> mainbe: lägg till i kö\ntävlande-ID
mainbe -> pumpfoil_db: läs upp kö
note over mainbe
  Lägg till tävlande i kö
  
endnote
mainbe -> pumpfoil_db: spara uppdaterad kö
adminfe <-- mainbe: lista med tävlande\nköad tävlande exkluderad
Anders <-- adminfe: visar lista med tävlande\nköad tävlande exkluderad
```

Starta spel för tävlande
------------------------
```plantuml
actor Anders
actor Tävlande
participant "Admin Frontend" as adminfe
participant "Game Frontend" as gamefe
participant "Main Backend" as mainbe
database pumpfoil_db

Anders -> adminfe: visa sida med kö
adminfe -> mainbe: hämta kö
mainbe -> pumpfoil_db: läs upp kö
pumpfoil_db --> mainbe: kö
mainbe --> adminfe: kö
note over adminfe
  Visa en sida för att starta spel
  - lista på tävlande i kö
  - knapp "Starta spel" på varje post i kön
endnote
Anders -> adminfe: klick på "starta spel"\nför en tävlande
adminfe -> mainbe: starta spel\ntävlande-ID
mainbe -> gamefe: starta spel\ntävlande
mainbe -> pumpfoil_db: Sätt status pumping för tävlande
note over gamefe
  Spel visar klart för start
end note

== avsluta spel ==
```

Genomför och avsluta spel
-------------------------
```plantuml
actor Anders
actor Tävlande
participant "Admin Frontend" as adminfe
participant "Game Frontend" as gamefe
participant "Main Backend" as mainbe
participant Pumpberry
database pumpfoil_db

Tävlande -> Pumpberry: pumpar
Pumpberry -> gamefe: pump_event\n(hastighet,lutning,graf)
note over gamefe
  Startlinje passeras: FE startar klockan
  Klocka medeltid
  Klocka sluttid
endnote
gamefe -> mainbe: game-finish\n(splitTime,endTime)
mainbe -> pumpfoil_db: hämta GameState.currentPumper
note over mainbe
  Skapa LeaderBoardItem för GameState.currentPumper
  Sätt splitTime och endTime i leaderBoardItem
  Nulla GameState.currentPumper
endnote  
mainbe -> pumpfoil_db: spara GameState, LeaderBoardItem
```

Avbryt
------
Om man inte blir klar inom 2 min, eller om Anders vill avbryta.
```plantuml
actor Anders
actor Tävlande
participant "Admin Frontend" as adminfe
participant "Game Frontend" as gamefe
participant "Main Backend" as mainbe
participant Pumpberry
database pumpfoil_db

alt
Anders -> adminfe: avbryt
adminfe -> mainbe: abort
else
note over gamefe
timeout
endnote
gamefe -> mainbe: timeout
end
mainbe -> pumpfoil_db: hämta GameState
note over mainbe
  Nulla GameState.currentPumper
end note
mainbe -> pumpfoil_db: spara GameState, LeaderBoardItem
gamefe <- mainbe: abort 
note over gamefe
  Avbryter pågående spel och 
  sätter tillbaka spelet till 
  grundtillstånd 
endnote   

```

Ta bort från kö
---------------
```plantuml
actor Anders
participant "Admin Frontend" as adminfe
participant "Main Backend" as mainbe
database pumpfoil_db

Anders -> adminfe: visa kö
== Hämta och visa kö ==
Anders -> adminfe: klicka "ta bort" på en tävlande
adminfe -> mainbe: remove\ncontestantID
note over mainbe
  Ta bort QueueItem för Contestant
  från DB
endnote
mainbe -> pumpfoil_db: uppdatera kö      


```


