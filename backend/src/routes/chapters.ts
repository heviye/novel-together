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

    const lastChapter = await pool.query(
      'SELECT MAX(chapter_number) as max FROM chapters WHERE novel_id = $1',
      [novelId]
    );
    const nextNumber = (lastChapter.rows[0]?.max || 0) + 1;

    const result = await pool.query(
      `INSERT INTO chapters (novel_id, chapter_number, author_id, content)
       VALUES ($1, $2, $3, $4) RETURNING *`,
      [novelId, nextNumber, req.userId, content]
    );

    await pool.query('UPDATE novels SET updated_at = NOW() WHERE id = $1', [novelId]);

    res.json(result.rows[0]);
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

export default router;
