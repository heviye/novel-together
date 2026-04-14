import Database from 'better-sqlite3';
import { v4 as uuidv4 } from 'uuid';

const dbPath = process.env.DB_PATH || './novel_together.db';

export const db = new Database(dbPath);

db.pragma('journal_mode = WAL');

export const initDatabase = () => {
  db.exec(`
    CREATE TABLE IF NOT EXISTS users (
      id TEXT PRIMARY KEY,
      username TEXT UNIQUE NOT NULL,
      email TEXT UNIQUE NOT NULL,
      password_hash TEXT NOT NULL,
      bio TEXT,
      avatar_url TEXT,
      created_at TEXT DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE IF NOT EXISTS novels (
      id TEXT PRIMARY KEY,
      title TEXT NOT NULL,
      description TEXT,
      author_id TEXT REFERENCES users(id),
      status TEXT DEFAULT 'active',
      created_at TEXT DEFAULT CURRENT_TIMESTAMP,
      updated_at TEXT DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE IF NOT EXISTS chapters (
      id TEXT PRIMARY KEY,
      novel_id TEXT REFERENCES novels(id) ON DELETE CASCADE,
      chapter_number INTEGER NOT NULL,
      author_id TEXT REFERENCES users(id),
      content TEXT NOT NULL,
      created_at TEXT DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE IF NOT EXISTS likes (
      id TEXT PRIMARY KEY,
      user_id TEXT REFERENCES users(id),
      chapter_id TEXT REFERENCES chapters(id) ON DELETE CASCADE,
      created_at TEXT DEFAULT CURRENT_TIMESTAMP,
      UNIQUE(user_id, chapter_id)
    );

    CREATE TABLE IF NOT EXISTS comments (
      id TEXT PRIMARY KEY,
      user_id TEXT REFERENCES users(id),
      chapter_id TEXT REFERENCES chapters(id) ON DELETE CASCADE,
      content TEXT NOT NULL,
      created_at TEXT DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE IF NOT EXISTS follows (
      id TEXT PRIMARY KEY,
      follower_id TEXT REFERENCES users(id),
      following_id TEXT REFERENCES users(id),
      created_at TEXT DEFAULT CURRENT_TIMESTAMP,
      UNIQUE(follower_id, following_id)
    );
  `);
};

// Adapter to make SQLite work like pg Pool
export const pool = {
  query: (sql: string, params: any[] = []) => {
    const stmt = db.prepare(sql);
    const isSelect = sql.trim().toUpperCase().startsWith('SELECT');
    
    if (isSelect) {
      const rows = params.length > 0 ? stmt.all(...params) : stmt.all();
      return { rows };
    } else {
      const result = params.length > 0 ? stmt.run(...params) : stmt.run();
      return { 
        rows: [{ ...result, changes: result.changes, lastInsertRowid: result.lastInsertRowid }],
        rowCount: result.changes 
      };
    }
  }
};

export const generateId = () => uuidv4();