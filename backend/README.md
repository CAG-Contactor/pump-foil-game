Backend för Pump-Foil-Game
=======

I denna fil beskrivs vad som behövs för att bygga och köra denna backend-tjänst

Beroenden
---------

Denna tjänst är skriven i **[Go](https://go.dev/)** du behöver därför installera go version 1.22 eller senare.

Vi har använt webapplikationsramverket **[gin](https://github.com/gin-gonic/gin)** pga av dess enkelhet och enligt tester dess höga prestanda. Dessutom har detta ramverk en utmärkt integration med swagger via paketet **[gin-swagger](https://github.com/swaggo/gin-swagger)**. Följ länken för att läsa hur detta paket installeras och används. 

Det verkat som den senaste versionen har en bugg i programmet swag. För att fixa detta kan man installera en äldre version av swag:

``go install github.com/swaggo/swag/cmd/swag@v1.8.12``

Bygga ock köra
---

Kör `swag init` om du har gjort ändringar i swagger dokumentationen. Bygg applikationen med `go build` eller `go run backend` om du vill kör igång servern direkt.

Tjänsten publiceras fn på `localhost:8080` det finns ett swagger-ui via url:en <http://localhost:8080/swagger/index.html>.



