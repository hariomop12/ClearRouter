import axios from 'axios';

// In production, API calls go through nginx proxy
// In development, use localhost:8080
const API_BASE_URL = import.meta.env.VITE_API_URL || 
  (import.meta.env.PROD ? '' : 'http://localhost:8080');

// Create axios instance
const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add token to requests if available
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Handle token expiration
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// Auth types
export interface LoginRequest {
  email: string;
  password: string;
}

export interface SignupRequest {
  name: string;
  email: string;
  password: string;
}

export interface User {
  id: string;
  name: string;
  email: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface SignupResponse {
  message: string;
}

// Auth API calls
export const authAPI = {
  login: async (data: LoginRequest): Promise<AuthResponse> => {
    const response = await api.post('/auth/login', data);
    return response.data;
  },

  signup: async (data: SignupRequest): Promise<SignupResponse> => {
    const response = await api.post('/auth/signup', data);
    return response.data;
  },

  verify: async (token: string): Promise<{message: string}> => {
    const response = await api.get(`/auth/verify?token=${token}`);
    return response.data;
  }
};

// User API calls
export const userAPI = {
  updateUsername: async (name: string): Promise<{message: string, user: User}> => {
    const response = await api.put('/user/username', { name });
    return response.data;
  },

  deleteAccount: async (): Promise<{message: string}> => {
    const response = await api.delete('/user/account');
    return response.data;
  }
};

// Helper functions
export const setAuthToken = (token: string) => {
  localStorage.setItem('token', token);
  api.defaults.headers.common['Authorization'] = `Bearer ${token}`;
};

export const removeAuthToken = () => {
  localStorage.removeItem('token');
  localStorage.removeItem('user');
  delete api.defaults.headers.common['Authorization'];
};

export default api;