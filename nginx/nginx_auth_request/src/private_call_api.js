const express = require('express');
const axios = require('axios');

function main() {
  const app = express();
  app.use(async (req, res, next) => {
    const response = await axios.post(
      'http://localhost:3000/api/auth/verify',
      {},
      {
        headers: {
          Authorization: req.headers.authorization,
        },
      }
    );
    const user = response.data;
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

  app.listen('3001', console.log('private server start on port 3001'));
}

main();
