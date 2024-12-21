import api from './api';
import { Coin } from '../types/coin';

export interface PortfolioHolding {
  coin: Coin;
  amount: number;
  value: number;
  averagePrice: number;
  profitLoss: number;
  profitLossPercentage: number;
}

export interface PortfolioValue {
  timestamp: string;
  value: number;
}

export const getPortfolioHoldings = async (): Promise<PortfolioHolding[]> => {
  const response = await api.get<PortfolioHolding[]>('/portfolio/holdings');
  return response.data;
};

export const getPortfolioHistory = async (timeframe: string): Promise<PortfolioValue[]> => {
  const response = await api.get<PortfolioValue[]>(`/portfolio/history?timeframe=${timeframe}`);
  return response.data;
}; 