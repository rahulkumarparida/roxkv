import axios from 'axios';

export const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
});

export const chatApi = axios.create({
  baseURL: import.meta.env.VITE_CHAT_API_BASE_URL || import.meta.env.VITE_API_BASE_URL,
});

export default api;
