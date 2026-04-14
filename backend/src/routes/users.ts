import { Router, Request, Response } from 'express';
import { pool, generateId } from '../db';
import { authMiddleware, AuthRequest } from '../middleware/auth';

const router = Router();

router.get('/:id', async (req: Request, res: Response) => {
  try {
    const { id } = req.params;
    const result = pool.query(
      'SELECT id, username, bio, avatar_url, created_at FROM users WHERE id = ?',
      [id]
    );
    if (result.rows.length === 0) {
      return res.status(404).json({ error: 'User not found' });
    }
    res.json(result.rows[0]);
  } catch (e) {
    res.status(500).json({ error: 'Internal server error' });
  }
});

router.put('/:id', authMiddleware, async (req: AuthRequest, res: Response) => {
  try {
    const { id } = req.params;
    if (req.userId !== id) {
      return res.status(403).json({ error: 'Forbidden' });
    }
    const { bio, avatar_url } = req.body;
    if (bio !== undefined) {
      pool.query('UPDATE users SET bio = ? WHERE id = ?', [bio, id]);
    }
    if (avatar_url !== undefined) {
      pool.query('UPDATE users SET avatar_url = ? WHERE id = ?', [avatar_url, id]);
    }
    res.json({ success: true });
  } catch (e) {
    res.status(500).json({ error: 'Internal server error' });
  }
});

router.post('/:id/follow', authMiddleware, async (req: AuthRequest, res: Response) => {
  try {
    const { id } = req.params;
    if (req.userId === id) {
      return res.status(400).json({ error: 'Cannot follow yourself' });
    }
    try {
      pool.query(
        'INSERT INTO follows (id, follower_id, following_id) VALUES (?, ?, ?)',
        [generateId(), req.userId, id]
      );
    } catch (e) {
      // Ignore duplicate
    }
    res.json({ success: true });
  } catch (e) {
    res.status(500).json({ error: 'Internal server error' });
  }
});

router.delete('/:id/follow', authMiddleware, async (req: AuthRequest, res: Response) => {
  try {
    const { id } = req.params;
    pool.query('DELETE FROM follows WHERE follower_id = ? AND following_id = ?', [req.userId, id]);
    res.json({ success: true });
  } catch (e) {
    res.status(500).json({ error: 'Internal server error' });
  }
});

router.get('/:id/followers', async (req: Request, res: Response) => {
  try {
    const { id } = req.params;
    const result = pool.query(
      `SELECT u.id, u.username, u.avatar_url FROM users u
       JOIN follows f ON f.follower_id = u.id WHERE f.following_id = ?`,
      [id]
    );
    res.json(result.rows);
  } catch (e) {
    res.status(500).json({ error: 'Internal server error' });
  }
});

router.get('/:id/following', async (req: Request, res: Response) => {
  try {
    const { id } = req.params;
    const result = pool.query(
      `SELECT u.id, u.username, u.avatar_url FROM users u
       JOIN follows f ON f.following_id = u.id WHERE f.follower_id = ?`,
      [id]
    );
    res.json(result.rows);
  } catch (e) {
    res.status(500).json({ error: 'Internal server error' });
  }
});

export default router;