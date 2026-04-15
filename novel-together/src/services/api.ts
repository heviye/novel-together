import axios from 'axios';
import '../types';

// 模拟器用: http://10.0.2.2:3000/api
// 真机用: 改为电脑局域网IP，如 http://192.168.1.XX:3000/api
const API_URL = 'http://10.0.2.2:3000/api';

const api = axios.create({
  baseURL: API_URL,
  timeout: 10000,
});

api.interceptors.request.use((config) => {
  const token = globalThis.__token;
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export const setToken = (token: string | null | undefined) => {
  globalThis.__token = token ?? undefined;
};

export const getToken = (): string | undefined => globalThis.__token;

// Auth
export const authApi = {
  register: (username: string, email: string, password: string) =>
    api.post('/auth/register', { username, email, password }),
  
  login: (email: string, password: string) =>
    api.post('/auth/login', { email, password }),
};

// Novels
export const novelApi = {
  list: (page = 1, limit = 20) =>
    api.get(`/novels?page=${page}&limit=${limit}`),
  
  create: (title: string, description: string) =>
    api.post('/novels', { title, description }),
  
  get: (id: string) =>
    api.get(`/novels/${id}`),
  
  getChapters: (id: string) =>
    api.get(`/novels/${id}/chapters`),
};

// Chapters
export const chapterApi = {
  create: (novelId: string, content: string) =>
    api.post(`/chapters/novels/${novelId}/chapters`, { content }),
  
  get: (id: string) =>
    api.get(`/chapters/${id}`),
  
  like: (id: string) =>
    api.post(`/chapters/${id}/like`),
  
  unlike: (id: string) =>
    api.delete(`/chapters/${id}/like`),
  
  getLikes: (id: string) =>
    api.get(`/chapters/${id}/likes`),
  
  comment: (id: string, content: string) =>
    api.post(`/chapters/${id}/comments`, { content }),
  
  getComments: (id: string) =>
    api.get(`/chapters/${id}/comments`),
};

// Users
export const userApi = {
  get: (id: string) =>
    api.get(`/users/${id}`),
  
  update: (id: string, bio?: string, avatar_url?: string) =>
    api.put(`/users/${id}`, { bio, avatar_url }),
  
  follow: (id: string) =>
    api.post(`/users/${id}/follow`),
  
  unfollow: (id: string) =>
    api.delete(`/users/${id}/follow`),
  
  followers: (id: string) =>
    api.get(`/users/${id}/followers`),
  
  following: (id: string) =>
    api.get(`/users/${id}/following`),
};

export default api;