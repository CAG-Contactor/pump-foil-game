import express from "express";
import bodyParser from 'body-parser';

const PORT = 3000;
const app = express();
app.use(bodyParser.json());

app.get('/api/ping', (req, res) => {
  res.send({ping: 'pong'});
});

app.get('/api/v1/game-finish', (req, res) => {
  res.send();
});

app.listen(PORT, () => {
  console.log(`Server is running on http://localhost:${PORT}`);
});
