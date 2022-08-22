const express = require('express');

function main() {
  const app = express();
  app.use((req, res, next) => {
    const idToken = req.headers['x-id-token'];
    if (!idToken) return res.sendStatus(401);
    const decoded = Buffer.from(idToken, 'base64').toString('ascii');
    const user = JSON.parse(decoded);
    req.user = user;
    next();
  });

  app.get('/api/admin', (req, res) => {
    const user = req.user;
    if (!user.isAdmin) return res.sendStatus(403);
    res.json({ message: 'ok admin' });
  });

  app.get('/api/private', (req, res) => {
    const user = req.user;
    res.json({ message: 'ok private' });
  });

  app.listen('3001', console.log('private server start on port 3001'));
}

main();
