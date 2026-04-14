import { Router, Request, Response } from 'express';
import { pool } from '../db';
import { authMiddleware, AuthRequest } from '../middleware/auth';

const router = Router();

router.post('/novels/:novelId/chapters', authMiddleware, async (req: AuthRequest, res: Response) => {
  try {
    const { novelId } = req.params;
    const { content } = req.body;
    if (!content) {
      return res.status(400).json({ error: 'Content required' });
    }

    const novelResult = await pool.query('SELECT id FROM novels WHERE id = $1', [novelId]);
    if (novelResult.rows.length === 0) {
      return res.status(404).json({ error: 'Novel not found' });
    }

    // Use transaction to prevent race condition on chapter numbers
    await pool.query('BEGIN');
    try {
      const lastChapter = await pool.query(
        'SELECT MAX(chapter_number) as max FROM chapters WHERE novel_id = $1 FOR UPDATE',
        [novelId]
      );
      const nextNumber = (lastChapter.rows[0]?.max || 0) + 1;

      const result = await pool.query(
        `INSERT INTO chapters (novel_id, chapter_number, author_id, content)
         VALUES ($1, $2, $3, $4) RETURNING *`,
        [novelId, nextNumber, req.userId, content]
      );

      await pool.query('UPDATE novels SET updated_at = NOW() WHERE id = $1', [novelId]);
      await pool.query('COMMIT');
      res.json(result.rows[0]);
    } catch (e) {
      await pool.query('ROLLBACK');
      throw e;
    }
  } catch (e) {
    res.status(500).json({ error: 'Internal server error' });
  }
});

router.get('/:id', async (req: Request, res: Response) => {
  try {
    const { id } = req.params;
    const result = await pool.query(
      `SELECT c.*, u.username as author_username, n.title as novel_title
       FROM chapters c
       JOIN users u ON c.author_id = u.id
       JOIN novels n ON c.novel_id = n.id
       WHERE c.id = $1`,
      [id]
    );
    if (result.rows.length === 0) {
      return res.status(404).json({ error: 'Chapter not found' });
    }
    res.json(result.rows[0]);
  } catch (e) {
    res.status(500).json({ error: 'Internal server error' });
  }
});

router.post('/:id/like', authMiddleware, async (req: AuthRequest, res: Response) => {
  try {
    const { id } = req.params;
    await pool.query(
      'INSERT INTO likes (user_id, chapter_id) VALUES ($1, $2) ON CONFLICT DO NOTHING',
      [req.userId, id]
    );
    res.json({ success: true });
  } catch (e) {
    res.status(500).json({ error: 'Internal server error' });
  }
});

router.delete('/:id/like', authMiddleware, async (req: AuthRequest, res: Response) => {
  try {
    const { id } = req.params;
    await pool.query('DELETE FROM likes WHERE user_id = $1 AND chapter_id = $2', [req.userId, id]);
    res.json({ success: true });
  } catch (e) {
    res.status(500).json({ error: 'Internal server error' });
  }
});

router.get('/:id/likes', async (req: Request, res: Response) => {
  try {
    const { id } = req.params;
    const result = await pool.query(
      'SELECT COUNT(*) as count FROM likes WHERE chapter_id = $1',
      [id]
    );
    res.json({ count: parseInt(result.rows[0].count) });
  } catch (e) {
    res.status(500).json({ error: 'Internal server error' });
  }
});

router.post('/:id/comments', authMiddleware, async (req: AuthRequest, res: Response) => {
  try {
    const { id } = req.params;
    const { content } = req.body;
    if (!content) {
      return res.status(400).json({ error: 'Content required' });
    }
    const result = await pool.query(
      'INSERT INTO comments (user_id, chapter_id, content) VALUES ($1, $2, $3) RETURNING *',
      [req.userId, id, content]
    );
    res.json(result.rows[0]);
  } catch (e) {
    res.status(500).json({ error: 'Internal server error' });
  }
});

router.get('/:id/comments', async (req: Request, res: Response) => {
  try {
    const { id } = req.params;
    const result = await pool.query(
      `SELECT c.*, u.username as author_username FROM comments c
       JOIN users u ON c.author_id = u.id WHERE c.chapter_id = $1 ORDER BY c.created_at`,
      [id]
    );
    res.json(result.rows);
  } catch (e) {
    res.status(500).json({ error: 'Internal server error' });
  }
});

export default router;
