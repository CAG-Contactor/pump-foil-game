import express from "express";
import bodyParser from 'body-parser';

const PORT = 3000;
const app = express();
app.use(bodyParser.json());

app.get('/api/ping', (req, res) => {
  res.send({ping: 'pong'});
});

app.listen(PORT, () => {
  console.log(`Server is running on http://localhost:${PORT}`);
});
