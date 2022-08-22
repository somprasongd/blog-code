const express = require('express');
const cors = require('cors');
const jwt = require('jsonwebtoken');

const secret = 'secretjwt';

function main() {
  const app = express();
  app.use(cors());
  app.post('/api/auth/signin', (req, res) => {
    const isAdmin = req.query.isAdmin === 'true';

    const token = jwt.sign(
      {
        id: 1,
        username: isAdmin ? 'admin' : 'user',
        isAdmin,
      },
      secret
    );

    res.json({ token });
  });

  app.get('/api/auth/verify', (req, res) => {
    const token = req.headers.authorization
      ? req.headers.authorization.split(' ')
      : undefined;
    if (!token || token.length < 2) return res.sendStatus(401);
    try {
      const user = jwt.verify(token[1], secret);
      res.set(
        'X-Id-Token',
        Buffer.from(JSON.stringify(user)).toString('base64')
      );
      res.status(200).json(user);
    } catch (error) {
      res.sendStatus(401);
    }
  });

  app.listen('3000', console.log('auth server start on port 3000'));
}

main();
