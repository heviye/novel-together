import { Router, Request, Response } from 'express';
import bcrypt from 'bcryptjs';
import jwt from 'jsonwebtoken';
import { pool, generateId } from '../db';

const router = Router();

router.post('/register', async (req: Request, res: Response) => {
  try {
    const { username, email, password } = req.body;
    if (!username || !email || !password) {
      return res.status(400).json({ error: 'Missing required fields' });
    }

    if (password.length < 6) {
      return res.status(400).json({ error: 'Password must be at least 6 characters' });
    }

    const passwordHash = await bcrypt.hash(password, 10);
    
    try {
      pool.query(
        'INSERT INTO users (id, username, email, password_hash) VALUES (?, ?, ?, ?)',
        [generateId(), username, email, passwordHash]
      );
    } catch (e: any) {
      if (e.message?.includes('UNIQUE')) {
        return res.status(400).json({ error: 'Username or email already exists' });
      }
      throw e;
    }
    
    const user = { id: generateId(), username, email };
    res.json(user);
  } catch (e) {
    res.status(500).json({ error: 'Internal server error' });
  }
});

router.post('/login', async (req: Request, res: Response) => {
  try {
    const { email, password } = req.body;
    if (!email || !password) {
      return res.status(400).json({ error: 'Missing required fields' });
    }

    const result = pool.query('SELECT * FROM users WHERE email = ?', [email]);
    if (result.rows.length === 0) {
      return res.status(401).json({ error: 'Invalid credentials' });
    }

    const user = result.rows[0] as any;
    const valid = await bcrypt.compare(password, user.password_hash);
    if (!valid) {
      return res.status(401).json({ error: 'Invalid credentials' });
    }

    const token = jwt.sign({ userId: user.id, username: user.username }, process.env.JWT_SECRET!, { expiresIn: '7d' });
    res.json({ token, user: { id: user.id, username: user.username, email: user.email } });
  } catch (e) {
    res.status(500).json({ error: 'Internal server error' });
  }
});

export default router;