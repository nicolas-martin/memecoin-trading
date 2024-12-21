export interface Coin {
  id: string;
  symbol: string;
  name: string;
  price: number;
  marketCap: number;
  volume24h: number;
  priceChange24h: number;
  logoUrl: string;
  contractAddress: string;
}

export interface CoinPrice {
  id: string;
  coinId: string;
  price: number;
  marketCap: number;
  volume24h: number;
  priceChange24h: number;
  timestamp: string;
} 