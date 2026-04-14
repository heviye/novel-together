# Novel Together - 协作写作平台设计

## 1. 项目概述

**项目名称**：Novel Together  
**类型**：移动应用 (React Native / Expo) + 后端 API  
**核心功能**：开放社区的协作写作平台，支持多人"故事接龙"模式  
**目标用户**：所有人

## 2. 核心规则

- 每个用户只能写小说的下一章，不能修改已有章节
- 所有人可以阅读任意小说
- 用户可以对章节点赞、评论
- 用户可以关注其他作者

## 3. 技术架构

### 前端
- React Native (Expo SDK 54)
- React 19.1.0
- React Native 0.81.5

### 后端
- Node.js + Express
- PostgreSQL
- JWT 认证

### 数据模型

```
User (用户)
- id: UUID
- username: string (唯一)
- email: string (唯一)
- password_hash: string
- bio: string (可选)
- avatar_url: string (可选)
- created_at: timestamp

Novel (小说)
- id: UUID
- title: string
- description: string
- author_id: UUID (创建者)
- status: enum (active/completed)
- created_at: timestamp
- updated_at: timestamp

Chapter (章节)
- id: UUID
- novel_id: UUID
- chapter_number: integer
- author_id: UUID (作者)
- content: text
- created_at: timestamp

Like (点赞)
- id: UUID
- user_id: UUID
- chapter_id: UUID
- created_at: timestamp

Comment (评论)
- id: UUID
- user_id: UUID
- chapter_id: UUID
- content: text
- created_at: timestamp

Follow (关注)
- id: UUID
- follower_id: UUID
- following_id: UUID
- created_at: timestamp
```

## 4. API 设计

### 认证
- `POST /api/auth/register` - 注册
- `POST /api/auth/login` - 登录
- `GET /api/auth/me` - 获取当前用户

### 用户
- `GET /api/users/:id` - 获取用户资料
- `PUT /api/users/:id` - 更新用户资料
- `POST /api/users/:id/follow` - 关注用户
- `DELETE /api/users/:id/follow` - 取消关注
- `GET /api/users/:id/followers` - 获取粉丝列表
- `GET /api/users/:id/following` - 获取关注列表

### 小说
- `GET /api/novels` - 小说列表 (分页)
- `POST /api/novels` - 创建小说
- `GET /api/novels/:id` - 获取小说详情
- `GET /api/novels/:id/chapters` - 获取章节列表

### 章节
- `POST /api/novels/:id/chapters` - 创建新章节 (必须是下一章)
- `GET /api/chapters/:id` - 获取章节内容

### 互动
- `POST /api/chapters/:id/like` - 点赞
- `DELETE /api/chapters/:id/like` - 取消点赞
- `GET /api/chapters/:id/likes` - 获取点赞数
- `POST /api/chapters/:id/comments` - 评论
- `GET /api/chapters/:id/comments` - 获取评论列表

## 5. 核心流程

### 创建小说
1. 用户登录
2. 点击"创建小说"
3. 填写标题、简介
4. 系统创建小说并自动创建第一章 (用户填写内容)

### 写新章节
1. 进入小说详情页
2. 点击"写新章节"
3. 检查是否是下一章 (chapter_number = 当前最大 + 1)
4. 填写章节内容
5. 提交

### 互动流程
- 点赞：点击心形图标
- 评论：在章节底部输入评论
- 关注：在用户资料页点击关注

## 6. MVP 实现顺序

1. **用户系统** - 注册/登录/个人资料/关注
2. **小说系统** - 创建小说/查看列表
3. **章节系统** - 写新章节/章节列表
4. **互动系统** - 点赞/评论
5. **首页** - 热门/最新小说

## 7. 验证标准

- [ ] 用户可以注册、登录
- [ ] 用户可以创建小说
- [ ] 用户可以写新章节 (按顺序)
- [ ] 用户可以点赞、评论
- [ ] 用户可以关注其他作者
- [ ] API 返回正确的错误信息