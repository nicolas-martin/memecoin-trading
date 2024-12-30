import axios from 'axios';

const API_URL = 'http://localhost:8080/api/v1';

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
  const response = await axios.get(`${API_URL}/memecoins`);
  return response.data;
};

export const getMemeCoinDetail = async (id: string): Promise<MemeCoin> => {
  const response = await axios.get(`${API_URL}/memecoins/${id}`);
  return response.data;
};

export const updateMemeCoins = async (): Promise<void> => {
  await axios.post(`${API_URL}/memecoins/update`);
}; 