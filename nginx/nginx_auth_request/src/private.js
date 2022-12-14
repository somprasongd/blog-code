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

  app.get('/api/resource/admin', (req, res) => {
    const user = req.user;
    if (!user.isAdmin) return res.sendStatus(403);
    res.json({ message: `User "${user.username}" can access admin` });
  });

  app.get('/api/resource/private', (req, res) => {
    const user = req.user;
    res.json({ message: `User "${user.username}" can access private` });
  });

  app.listen('3002', console.log('private server start on port 3002'));
}

main();
