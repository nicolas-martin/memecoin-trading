import axios from 'axios';
import { Platform } from 'react-native';

// Helper function to get the base API URL based on platform
const getBaseUrl = () => {
  if (Platform.OS === 'web') {
    // For web development with Expo
    const hostname = window.location.hostname;
    const port = window.location.port;
    // Use the same port as the frontend is running on
    return `http://${hostname}:8080`;
  }
  
  if (Platform.OS === 'android') {
    return 'http://10.0.2.2:8080';
  }
  
  // iOS or default
  return 'http://localhost:8080';
};

const API_URL = `${getBaseUrl()}/api/v1`;

// Configure axios defaults
const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Accept': 'application/json',
    'Content-Type': 'application/json',
  },
  timeout: 10000, // 10 seconds timeout
});

// Add request/response interceptors for debugging
api.interceptors.request.use(
  request => {
    console.log('Starting API Request:', {
      method: request.method?.toUpperCase(),
      url: request.url,
      baseURL: request.baseURL,
      headers: request.headers
    });
    return request;
  },
  error => {
    console.error('Request Error:', error);
    return Promise.reject(error);
  }
);

api.interceptors.response.use(
  response => {
    console.log('API Response:', {
      status: response.status,
      url: response.config.url,
      data: response.data
    });
    return response;
  },
  error => {
    if (error.response) {
      // The request was made and the server responded with a status code
      // that falls out of the range of 2xx
      console.error('API Error Response:', {
        status: error.response.status,
        data: error.response.data,
        headers: error.response.headers,
        url: error.config?.url
      });
    } else if (error.request) {
      // The request was made but no response was received
      console.error('API No Response:', {
        request: error.request,
        url: error.config?.url
      });
    } else {
      // Something happened in setting up the request that triggered an Error
      console.error('API Request Setup Error:', {
        message: error.message,
        url: error.config?.url
      });
    }
    return Promise.reject(error);
  }
);

export interface MemeCoin {
  id: string;
  symbol: string;
  name: string;
  logoUrl: string;
  price: number;
  marketCap: number;
  volume24h: number;
  priceChange24h: number;
  priceChangePercentage24h: number;
  contractAddress: string;
}

export interface PriceHistory {
  timestamp: number;
  price: number;
}

export const getTopMemeCoins = async (): Promise<MemeCoin[]> => {
  try {
    const response = await api.get('/memecoins');
    return response.data;
  } catch (error) {
    console.error('Error fetching top meme coins:', error);
    throw error;
  }
};

export const getMemeCoinDetail = async (id: string): Promise<MemeCoin> => {
  try {
    const response = await api.get(`/memecoins/${id}`);
    return response.data;
  } catch (error) {
    console.error('Error fetching meme coin detail:', error);
    throw error;
  }
};

export const updateMemeCoins = async (): Promise<void> => {
  try {
    await api.post('/memecoins/update');
  } catch (error) {
    console.error('Error updating meme coins:', error);
    throw error;
  }
}; 