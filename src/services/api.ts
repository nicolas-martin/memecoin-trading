import axios from 'axios';
import { Coin } from '../types/coin';

const API_URL = import.meta.env.VITE_API_BASE_URL;

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

export const getTopCoins = async (limit: number = 50): Promise<Coin[]> => {
  const response = await api.get<Coin[]>(`/coins/top?limit=${limit}`);
  return response.data;
};

export const getCoinById = async (id: string): Promise<Coin> => {
  try {
    const response = await axios.get<Coin>(`${API_URL}/coins/${id}`);
    return response.data;
  } catch (error) {
    throw new Error('Failed to fetch coin details');
  }
};

// Add authentication interceptor
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export default api; 