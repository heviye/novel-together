import { Router, Request, Response } from 'express';
import { pool, generateId } from '../db';
import { authMiddleware, AuthRequest } from '../middleware/auth';

const router = Router();

router.post('/novels/:novelId/chapters', authMiddleware, async (req: AuthRequest, res: Response) => {
  try {
    const { novelId } = req.params;
    const { content } = req.body;
    if (!content) {
      return res.status(400).json({ error: 'Content required' });
    }

    const novelResult = pool.query('SELECT id FROM novels WHERE id = ?', [novelId]);
    if (novelResult.rows.length === 0) {
      return res.status(404).json({ error: 'Novel not found' });
    }

    const lastChapter = pool.query(
      'SELECT MAX(chapter_number) as max FROM chapters WHERE novel_id = ?',
      [novelId]
    );
    const nextNumber = ((lastChapter.rows[0] as any)?.max || 0) + 1;

    const result = pool.query(
      'INSERT INTO chapters (id, novel_id, chapter_number, author_id, content) VALUES (?, ?, ?, ?, ?)',
      [generateId(), novelId, nextNumber, req.userId, content]
    );

    pool.query('UPDATE novels SET updated_at = CURRENT_TIMESTAMP WHERE id = ?', [novelId]);

    res.json({ id: generateId(), novel_id: novelId, chapter_number: nextNumber, author_id: req.userId, content });
  } catch (e) {
    console.error(e);
    res.status(500).json({ error: 'Internal server error' });
  }
});

router.get('/:id', async (req: Request, res: Response) => {
  try {
    const { id } = req.params;
    const result = pool.query(
      `SELECT c.*, u.username as author_username, n.title as novel_title
       FROM chapters c
       JOIN users u ON c.author_id = u.id
       JOIN novels n ON c.novel_id = n.id
       WHERE c.id = ?`,
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
    try {
      pool.query(
        'INSERT INTO likes (id, user_id, chapter_id) VALUES (?, ?, ?)',
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

router.delete('/:id/like', authMiddleware, async (req: AuthRequest, res: Response) => {
  try {
    const { id } = req.params;
    pool.query('DELETE FROM likes WHERE user_id = ? AND chapter_id = ?', [req.userId, id]);
    res.json({ success: true });
  } catch (e) {
    res.status(500).json({ error: 'Internal server error' });
  }
});

router.get('/:id/likes', async (req: Request, res: Response) => {
  try {
    const { id } = req.params;
    const result = pool.query(
      'SELECT COUNT(*) as count FROM likes WHERE chapter_id = ?',
      [id]
    );
    res.json({ count: (result.rows[0] as any)?.count || 0 });
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
    pool.query(
      'INSERT INTO comments (id, user_id, chapter_id, content) VALUES (?, ?, ?, ?)',
      [generateId(), req.userId, id, content]
    );
    res.json({ id: generateId(), user_id: req.userId, chapter_id: id, content });
  } catch (e) {
    res.status(500).json({ error: 'Internal server error' });
  }
});

router.get('/:id/comments', async (req: Request, res: Response) => {
  try {
    const { id } = req.params;
    const result = pool.query(
      `SELECT c.id, c.content, c.created_at, u.username as author_username FROM comments c
       JOIN users u ON c.author_id = u.id WHERE c.chapter_id = ? ORDER BY c.created_at`,
      [id]
    );
    res.json(result.rows);
  } catch (e) {
    res.status(500).json({ error: 'Internal server error' });
  }
});

export default router;