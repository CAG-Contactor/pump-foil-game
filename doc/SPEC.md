Spec, The Pump Foiler
=====================
Se även [Gamla specen](https://docs.google.com/document/d/1iw420msEF6ePYX12Zz6Wj7ceX9JUx-VIVOwcCjhF-Ug/edit?usp=sharing)

Översikt
--------
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

Anders -> adminfe: lägg till ny tävlande
note over adminfe
  Visa en sida för ny tävlande
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
Brödrost -> mainbe: namn&epost
mainbe -> adminfe: / namn&epost
adminfe --> Anders: visa namn&epost
Anders -> adminfe: bekräfta
adminfe -> mainbe: registrera tävlande  
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

note over adminfe
  Visa en sida för att lägga till i kö
  - lista på tävlande
  - knapp "Lägg till i kö"
endnote
pumpfoil_db -> mainbe: hämta tävlande
mainbe -> adminfe: tävlande
Anders -> adminfe: lägg till i kö
adminfe -> mainbe: lägg till i kö
mainbe -> pumpfoil_db: lägg till i kö
```

Starta spel för tävlande
------------------------
```plantuml
actor Anders
actor Tävlande
participant "Admin Frontend" as adminfe
participant "Main Backend" as mainbe
participant "Game Frontend" as gamefe
database pumpfoil_db

adminfe -> mainbe: tävlande i kö
mainbe -> pumpfoil_db: hämta tävlande i kö
pumpfoil_db -> mainbe: tävlande i kö
mainbe --> adminfe: tävlande i kö
note over adminfe
  Visa en sida för att starta spel
  - lista på tävlande i kö
  - knapp "Starta spel"
endnote
Anders -> adminfe: starta spel
adminfe -> mainbe: starta spel
mainbe -> gamefe: starta spel
mainbe -> pumpfoil_db: Sätt status pumping för tävlande
note over gamefe
  Spel visar klart för start
end note

== avsluta spel ==
```
