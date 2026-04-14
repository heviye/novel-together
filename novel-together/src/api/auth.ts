import axios from 'axios';

const API_URL = 'http://localhost:3000/api';

export const authApi = {
  register: (username: string, email: string, password: string) =>
    axios.post(`${API_URL}/auth/register`, { username, email, password }),
  
  login: (email: string, password: string) =>
    axios.post(`${API_URL}/auth/login`, { email, password }),
};
