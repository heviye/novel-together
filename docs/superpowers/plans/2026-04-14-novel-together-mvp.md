# Novel Together MVP Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 构建一个开放社区的协作写作平台 MVP，支持用户注册登录、创建小说、写新章节、点赞、评论、关注功能

**Architecture:** 
- 前端：React Native (Expo) 移动应用
- 后端：Node.js + Express REST API + PostgreSQL
- 前后端分离，通过 HTTP API 通信

**Tech Stack:**
- Frontend: React Native, Expo SDK 54, React Navigation
- Backend: Node.js, Express, PostgreSQL, JWT, bcrypt

---

## 文件结构

```
novel-together/           # 前端应用
├── App.tsx              # 应用入口
├── src/
│   ├── api/             # API 客户端
│   ├── screens/         # 页面
│   ├── components/      # 组件
│   └── types/          # 类型定义
backend/                 # 后端服务
├── src/
│   ├── index.ts        # 服务入口
│   ├── routes/        # 路由
│   ├── models/        # 数据模型
│   ├── middleware/   # 中间件
│   └── db/           # 数据库连接
└── package.json
```

---

## Task 1: 项目初始化与后端搭建

**Files:**
- Create: `backend/package.json`
- Create: `backend/tsconfig.json`
- Create: `backend/src/index.ts`
- Create: `backend/src/db/index.ts`
- Modify: `novel-together/package.json`

- [ ] **Step 1: 创建后端项目结构和依赖**

```bash
mkdir -p backend/src/{routes,models,middleware,db}
cd backend
npm init -y
npm install express cors dotenv pg bcryptjs jsonwebtoken uuid
npm install -D typescript @types/node @types/express @types/cors @types/bcryptjs @types/jsonwebtoken @types/uuid ts-node nodemon
```

- [ ] **Step 2: 创建 tsconfig.json**

```json
{
  "compilerOptions": {
    "target": "ES2020",
    "module": "commonjs",
    "outDir": "./dist",
    "rootDir": "./src",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true
  },
  "include": ["src/**/*"]
}
```

- [ ] **Step 3: 创建数据库连接**

```typescript
// backend/src/db/index.ts
import { Pool } from 'pg';

export const pool = new Pool({
  host: process.env.DB_HOST || 'localhost',
  port: parseInt(process.env.DB_PORT || '5432'),
  database: process.env.DB_NAME || 'novel_together',
  user: process.env.DB_USER || 'postgres',
  password: process.env.DB_PASSWORD || 'postgres',
});

export const initDatabase = async () => {
  await pool.query(`
    CREATE TABLE IF NOT EXISTS users (
      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
      username VARCHAR(50) UNIQUE NOT NULL,
      email VARCHAR(255) UNIQUE NOT NULL,
      password_hash VARCHAR(255) NOT NULL,
      bio TEXT,
      avatar_url TEXT,
      created_at TIMESTAMP DEFAULT NOW()
    );

    CREATE TABLE IF NOT EXISTS novels (
      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
      title VARCHAR(255) NOT NULL,
      description TEXT,
      author_id UUID REFERENCES users(id),
      status VARCHAR(20) DEFAULT 'active',
      created_at TIMESTAMP DEFAULT NOW(),
      updated_at TIMESTAMP DEFAULT NOW()
    );

    CREATE TABLE IF NOT EXISTS chapters (
      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
      novel_id UUID REFERENCES novels(id) ON DELETE CASCADE,
      chapter_number INTEGER NOT NULL,
      author_id UUID REFERENCES users(id),
      content TEXT NOT NULL,
      created_at TIMESTAMP DEFAULT NOW()
    );

    CREATE TABLE IF NOT EXISTS likes (
      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
      user_id UUID REFERENCES users(id),
      chapter_id UUID REFERENCES chapters(id) ON DELETE CASCADE,
      created_at TIMESTAMP DEFAULT NOW(),
      UNIQUE(user_id, chapter_id)
    );

    CREATE TABLE IF NOT EXISTS comments (
      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
      user_id UUID REFERENCES users(id),
      chapter_id UUID REFERENCES chapters(id) ON DELETE CASCADE,
      content TEXT NOT NULL,
      created_at TIMESTAMP DEFAULT NOW()
    );

    CREATE TABLE IF NOT EXISTS follows (
      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
      follower_id UUID REFERENCES users(id),
      following_id UUID REFERENCES users(id),
      created_at TIMESTAMP DEFAULT NOW(),
      UNIQUE(follower_id, following_id)
    );
  `);
};
```

- [ ] **Step 4: 创建服务入口**

```typescript
// backend/src/index.ts
import express from 'express';
import cors from 'cors';
import dotenv from 'dotenv';
import { initDatabase } from './db';

dotenv.config();

const app = express();
app.use(cors());
app.use(express.json());

const PORT = process.env.PORT || 3000;

app.get('/health', (req, res) => {
  res.json({ status: 'ok' });
});

const start = async () => {
  try {
    await initDatabase();
    app.listen(PORT, () => {
      console.log(`Server running on port ${PORT}`);
    });
  } catch (e) {
    console.error('Failed to start server:', e);
    process.exit(1);
  }
};

start();
```

- [ ] **Step 5: 更新前端 package.json 添加导航依赖**

```bash
cd novel-together
npm install @react-navigation/native @react-navigation/native-stack react-native-screens react-native-safe-area-context axios
```

- [ ] **Step 6: 提交**

```bash
git add backend/ novel-together/package.json
git commit -m "feat: initial backend setup and database schema"
```

---

## Task 2: 认证系统 (注册/登录)

**Files:**
- Create: `backend/src/routes/auth.ts`
- Create: `backend/src/middleware/auth.ts`
- Modify: `backend/src/index.ts`

- [ ] **Step 1: 创建认证路由**

```typescript
// backend/src/routes/auth.ts
import { Router, Request, Response } from 'express';
import bcrypt from 'bcryptjs';
import jwt from 'jwt-simple';
import { pool } from '../db';

const router = Router();

router.post('/register', async (req: Request, res: Response) => {
  try {
    const { username, email, password } = req.body;
    if (!username || !email || !password) {
      return res.status(400).json({ error: 'Missing required fields' });
    }

    const passwordHash = await bcrypt.hash(password, 10);
    const result = await pool.query(
      'INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id, username, email, created_at',
      [username, email, passwordHash]
    );
    res.json(result.rows[0]);
  } catch (e: any) {
    if (e.code === '23505') {
      return res.status(400).json({ error: 'Username or email already exists' });
    }
    res.status(500).json({ error: 'Internal server error' });
  }
});

router.post('/login', async (req: Request, res: Response) => {
  try {
    const { email, password } = req.body;
    if (!email || !password) {
      return res.status(400).json({ error: 'Missing required fields' });
    }

    const result = await pool.query('SELECT * FROM users WHERE email = $1', [email]);
    if (result.rows.length === 0) {
      return res.status(401).json({ error: 'Invalid credentials' });
    }

    const user = result.rows[0];
    const valid = await bcrypt.compare(password, user.password_hash);
    if (!valid) {
      return res.status(401).json({ error: 'Invalid credentials' });
    }

    const token = jwt.encode({ userId: user.id, username: user.username }, process.env.JWT_SECRET || 'secret');
    res.json({ token, user: { id: user.id, username: user.username, email: user.email } });
  } catch (e) {
    res.status(500).json({ error: 'Internal server error' });
  }
});

export default router;
```

- [ ] **Step 2: 创建认证中间件**

```typescript
// backend/src/middleware/auth.ts
import { Request, Response, NextFunction } from 'express';
import jwt from 'jwt-simple';

export interface AuthRequest extends Request {
  userId?: string;
}

export const authMiddleware = (req: AuthRequest, res: Response, next: NextFunction) => {
  const token = req.headers.authorization?.split(' ')[1];
  if (!token) {
    return res.status(401).json({ error: 'No token provided' });
  }

  try {
    const decoded = jwt.decode(token, process.env.JWT_SECRET || 'secret');
    req.userId = decoded.userId;
    next();
  } catch (e) {
    res.status(401).json({ error: 'Invalid token' });
  }
};
```

- [ ] **Step 3: 在 index.ts 中注册路由**

```typescript
import authRoutes from './routes/auth';

app.use('/api/auth', authRoutes);
```

- [ ] **Step 4: 创建前端认证 API 客户端**

```typescript
// novel-together/src/api/auth.ts
import axios from 'axios';

const API_URL = 'http://localhost:3000/api';

export const authApi = {
  register: (username: string, email: string, password: string) =>
    axios.post(`${API_URL}/auth/register`, { username, email, password }),
  
  login: (email: string, password: string) =>
    axios.post(`${API_URL}/auth/login`, { email, password }),
};
```

- [ ] **Step 5: 提交**

```bash
git add backend/src/routes/auth.ts backend/src/middleware/auth.ts backend/src/index.ts novel-together/src/api/
git commit -m "feat: add authentication system"
```

---

## Task 3: 用户系统 (个人资料/关注)

**Files:**
- Create: `backend/src/routes/users.ts`
- Modify: `backend/src/index.ts`

- [ ] **Step 1: 创建用户路由**

```typescript
// backend/src/routes/users.ts
import { Router, Response } from 'express';
import { pool } from '../db';
import { authMiddleware, AuthRequest } from '../middleware/auth';

const router = Router();

router.get('/:id', async (req: Request, res: Response) => {
  try {
    const { id } = req.params;
    const result = await pool.query(
      'SELECT id, username, bio, avatar_url, created_at FROM users WHERE id = $1',
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
    await pool.query(
      'UPDATE users SET bio = COALESCE($1, bio), avatar_url = COALESCE($2, avatar_url) WHERE id = $3',
      [bio, avatar_url, id]
    );
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
    await pool.query(
      'INSERT INTO follows (follower_id, following_id) VALUES ($1, $2) ON CONFLICT DO NOTHING',
      [req.userId, id]
    );
    res.json({ success: true });
  } catch (e) {
    res.status(500).json({ error: 'Internal server error' });
  }
});

router.delete('/:id/follow', authMiddleware, async (req: AuthRequest, res: Response) => {
  try {
    const { id } = req.params;
    await pool.query('DELETE FROM follows WHERE follower_id = $1 AND following_id = $2', [req.userId, id]);
    res.json({ success: true });
  } catch (e) {
    res.status(500).json({ error: 'Internal server error' });
  }
});

router.get('/:id/followers', async (req: Request, res: Response) => {
  try {
    const { id } = req.params;
    const result = await pool.query(
      `SELECT u.id, u.username, u.avatar_url FROM users u
       JOIN follows f ON f.follower_id = u.id WHERE f.following_id = $1`,
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
    const result = await pool.query(
      `SELECT u.id, u.username, u.avatar_url FROM users u
       JOIN follows f ON f.following_id = u.id WHERE f.follower_id = $1`,
      [id]
    );
    res.json(result.rows);
  } catch (e) {
    res.status(500).json({ error: 'Internal server error' });
  }
});

export default router;
```

- [ ] **Step 2: 注册用户路由**

```typescript
import userRoutes from './routes/users';
app.use('/api/users', userRoutes);
```

- [ ] **Step 3: 提交**

```bash
git add backend/src/routes/users.ts backend/src/index.ts
git commit -m "feat: add user system with follow functionality"
```

---

## Task 4: 小说系统

**Files:**
- Create: `backend/src/routes/novels.ts`
- Modify: `backend/src/index.ts`

- [ ] **Step 1: 创建小说路由**

```typescript
// backend/src/routes/novels.ts
import { Router, Response } from 'express';
import { pool } from '../db';
import { authMiddleware, AuthRequest } from '../middleware/auth';

const router = Router();

router.get('/', async (req: Request, res: Response) => {
  try {
    const page = parseInt(req.query.page as string) || 1;
    const limit = parseInt(req.query.limit as string) || 20;
    const offset = (page - 1) * limit;

    const result = await pool.query(
      `SELECT n.*, u.username as author_username,
       (SELECT COUNT(*) FROM chapters WHERE novel_id = n.id) as chapter_count
       FROM novels n
       JOIN users u ON n.author_id = u.id
       WHERE n.status = 'active'
       ORDER BY n.updated_at DESC
       LIMIT $1 OFFSET $2`,
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

    const result = await pool.query(
      'INSERT INTO novels (title, description, author_id) VALUES ($1, $2, $3) RETURNING *',
      [title, description, req.userId]
    );
    res.json(result.rows[0]);
  } catch (e) {
    res.status(500).json({ error: 'Internal server error' });
  }
});

router.get('/:id', async (req: Request, res: Response) => {
  try {
    const { id } = req.params;
    const result = await pool.query(
      `SELECT n.*, u.username as author_username FROM novels n JOIN users u ON n.author_id = u.id WHERE n.id = $1`,
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
    const result = await pool.query(
      `SELECT c.id, c.chapter_number, c.created_at, u.username as author_username
       FROM chapters c
       JOIN users u ON c.author_id = u.id
       WHERE c.novel_id = $1
       ORDER BY c.chapter_number`,
      [id]
    );
    res.json(result.rows);
  } catch (e) {
    res.status(500).json({ error: 'Internal server error' });
  }
});

export default router;
```

- [ ] **Step 2: 注册小说路由**

```typescript
import novelRoutes from './routes/novels';
app.use('/api/novels', novelRoutes);
```

- [ ] **Step 3: 提交**

```bash
git add backend/src/routes/novels.ts backend/src/index.ts
git commit -m "feat: add novel system"
```

---

## Task 5: 章节系统

**Files:**
- Create: `backend/src/routes/chapters.ts`
- Modify: `backend/src/index.ts`

- [ ] **Step 1: 创建章节路由**

```typescript
// backend/src/routes/chapters.ts
import { Router, Response } from 'express';
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
```

- [ ] **Step 2: 注册章节路由**

```typescript
import chapterRoutes from './routes/chapters';
app.use('/api/chapters', chapterRoutes);
```

- [ ] **Step 3: 提交**

```bash
git add backend/src/routes/chapters.ts backend/src/index.ts
git commit -m "feat: add chapter system"
```

---

## Task 6: 互动系统 (点赞/评论)

**Files:**
- Modify: `backend/src/routes/chapters.ts`

- [ ] **Step 1: 添加点赞评论功能**

```typescript
// Add to chapters.ts
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
```

- [ ] **Step 2: 提交**

```bash
git add backend/src/routes/chapters.ts
git commit -m "feat: add like and comment functionality"
```

---

## Task 7: 前端页面

**Files:**
- Create: `novel-together/src/screens/HomeScreen.tsx`
- Create: `novel-together/src/screens/LoginScreen.tsx`
- Create: `novel-together/src/screens/RegisterScreen.tsx`
- Create: `novel-together/src/screens/NovelListScreen.tsx`
- Create: `novel-together/src/screens/NovelDetailScreen.tsx`
- Create: `novel-together/src/screens/WriteChapterScreen.tsx`
- Create: `novel-together/src/screens/ChapterScreen.tsx`
- Create: `novel-together/src/screens/ProfileScreen.tsx`
- Modify: `novel-together/App.tsx`

- [ ] **Step 1: 创建导航和基础结构**

```typescript
// novel-together/App.tsx
import { NavigationContainer } from '@react-navigation/native';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import HomeScreen from './src/screens/HomeScreen';
import LoginScreen from './src/screens/LoginScreen';
import RegisterScreen from './src/screens/RegisterScreen';
import NovelListScreen from './src/screens/NovelListScreen';
import NovelDetailScreen from './src/screens/NovelDetailScreen';
import WriteChapterScreen from './src/screens/WriteChapterScreen';
import ChapterScreen from './src/screens/ChapterScreen';
import ProfileScreen from './src/screens/ProfileScreen';

const Stack = createNativeStackNavigator();

export default function App() {
  return (
    <NavigationContainer>
      <Stack.Navigator>
        <Stack.Screen name="Home" component={HomeScreen} />
        <Stack.Screen name="Login" component={LoginScreen} />
        <Stack.Screen name="Register" component={RegisterScreen} />
        <Stack.Screen name="NovelList" component={NovelListScreen} />
        <Stack.Screen name="NovelDetail" component={NovelDetailScreen} />
        <Stack.Screen name="WriteChapter" component={WriteChapterScreen} />
        <Stack.Screen name="Chapter" component={ChapterScreen} />
        <Stack.Screen name="Profile" component={ProfileScreen} />
      </Stack.Navigator>
    </NavigationContainer>
  );
}
```

- [ ] **Step 2: 创建各个页面**

根据设计规格创建各页面，简单实现核心功能。

- [ ] **Step 3: 提交**

```bash
git add novel-together/App.tsx novel-together/src/screens/
git commit -m "feat: add frontend screens"
```

---

## 实现顺序

按 Task 顺序执行：
1. 项目初始化与后端搭建
2. 认证系统
3. 用户系统
4. 小说系统
5. 章节系统
6. 互动系统
7. 前端页面