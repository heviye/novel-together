import { Router, Request, Response } from 'express';
import { pool, generateId } from '../db';
import { authMiddleware, AuthRequest } from '../middleware/auth';

const router = Router();

router.get('/', async (req: Request, res: Response) => {
  try {
    const page = parseInt(req.query.page as string) || 1;
    const limit = parseInt(req.query.limit as string) || 20;
    const offset = (page - 1) * limit;

    const result = pool.query(
      `SELECT n.*, u.username as author_username,
       (SELECT COUNT(*) FROM chapters WHERE novel_id = n.id) as chapter_count
       FROM novels n
       JOIN users u ON n.author_id = u.id
       WHERE n.status = 'active'
       ORDER BY n.updated_at DESC
       LIMIT ? OFFSET ?`,
      [limit, offset]
    );
    res.json(result.rows);
  } catch (e) {
    res.status(500).json({ error: 'Internal server error' });
  }
});

router.post('/', authMiddleware, async (req: AuthRequest, res: Response) => {
  try {
    const { title, description } = req.body;
    if (!title) {
      return res.status(400).json({ error: 'Title required' });
    }

    const id = generateId();
    pool.query(
      'INSERT INTO novels (id, title, description, author_id) VALUES (?, ?, ?, ?)',
      [id, title, description || null, req.userId]
    );
    res.json({ id, title, description, author_id: req.userId, status: 'active' });
  } catch (e) {
    res.status(500).json({ error: 'Internal server error' });
  }
});

router.get('/:id', async (req: Request, res: Response) => {
  try {
    const { id } = req.params;
    const result = pool.query(
      `SELECT n.*, u.username as author_username FROM novels n JOIN users u ON n.author_id = u.id WHERE n.id = ?`,
      [id]
    );
    if (result.rows.length === 0) {
      return res.status(404).json({ error: 'Novel not found' });
    }
    res.json(result.rows[0]);
  } catch (e) {
    res.status(500).json({ error: 'Internal server error' });
  }
});

router.get('/:id/chapters', async (req: Request, res: Response) => {
  try {
    const { id } = req.params;
    const result = pool.query(
      `SELECT c.id, c.chapter_number, c.created_at, u.username as author_username
       FROM chapters c
       JOIN users u ON c.author_id = u.id
       WHERE c.novel_id = ?
       ORDER BY c.chapter_number`,
      [id]
    );
    res.json(result.rows);
  } catch (e) {
    res.status(500).json({ error: 'Internal server error' });
  }
});

export default router;