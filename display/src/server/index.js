const express = require('express');
const os = require('os');
const fs = require('fs');

const app = express();

app.use(express.static('dist'));
app.use(function(req, res, next) {
  res.header("Access-Control-Allow-Origin", "*");
  res.header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept");
  next();
});

app.get('/api/getUsername', (req, res) => res.send({ username: os.userInfo().username }));
app.get('/api/getLog/:path', (req, res) => {
  const logname = req.params.path;

  fs.readFile(`logs/${logname}.txt`, 'utf8', (err, data) => {
    let lines = data.split('\n');
    let parsedLines = lines.map(l => {
      try {
        return JSON.parse(l);
      } catch {
        return null;
      }
    });
    let goodLines = parsedLines.filter(x => x != null);

    res.send(JSON.stringify(goodLines));
  });
});


app.listen(process.env.PORT || 8080, () => console.log(`Listening on port ${process.env.PORT || 8080}!`));
