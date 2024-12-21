import { generateMockPair, generateMockPairs } from '../utils/mockData';

export interface TokenInfo {
  address: string;
  name: string;
  symbol: string;
}

export interface PriceData {
  priceUsd: string;
  priceNative: string;
  priceChange: {
    h1: number;
    h24: number;
    d7: number;
  };
}

export interface VolumeData {
  h24: number;
  h6: number;
  h1: number;
  m5: number;
}

export interface LiquidityData {
  usd: number;
  base: number;
  quote: number;
}

export interface SocialInfo {
  platform: string;
  handle: string;
}

export interface PairData {
  chainId: string;
  dexId: string;
  url: string;
  pairAddress: string;
  baseToken: TokenInfo;
  quoteToken: TokenInfo;
  priceNative: string;
  priceUsd: string;
  liquidity: LiquidityData;
  volume: VolumeData;
  priceChange: PriceData['priceChange'];
  fdv: number;
  marketCap: number;
  info?: {
    imageUrl?: string;
    description?: string;
    websites?: { url: string }[];
    socials?: SocialInfo[];
  };
}

export async function searchPairs(query: string): Promise<PairData[]> {
  return generateMockPairs(10);
}

export async function getPairsByAddress(chainId: string, pairAddress: string): Promise<PairData[]> {
  return [generateMockPair(pairAddress.slice(0, 6))];
}

export async function getPairsByToken(tokenAddresses: string[]): Promise<PairData[]> {
  return tokenAddresses.map(addr => generateMockPair(addr.slice(0, 6)));
}

export async function getTrendingPairs(): Promise<PairData[]> {
  return generateMockPairs(10);
}

export interface HistoricalPrice {
  timestamp: number;
  price: number;
}

export const fetchHistoricalPrices = async (pairAddress: string, timeframe: string): Promise<HistoricalPrice[]> => {
  // Generate mock historical prices
  const now = Date.now();
  const points = 100;
  const basePrice = 1000 + Math.random() * 9000;
  const volatility = 0.02;

  return Array.from({ length: points }, (_, i) => {
    const timeDiff = i * (24 * 60 * 60 * 1000 / points); // Spread over 24 hours
    const randomWalk = basePrice * (1 + (Math.random() - 0.5) * volatility);
    return {
      timestamp: now - timeDiff,
      price: randomWalk
    };
  }).reverse();
};

export interface CoinData {
  name: string;
  symbol: string;
  pairAddress: string;
  price: string;
  change: string;
  icon: string;
  color: string;
  marketCap: number;
  volume24h: number;
  priceHistory: {
    timestamp: number;
    price: number;
  }[];
}

export const fetchCoinData = async (pairAddress: string): Promise<CoinData> => {
  const pairs = await getPairsByAddress('', pairAddress);
  const pair = pairs[0];
  
  return {
    name: pair.baseToken.name,
    symbol: pair.baseToken.symbol,
    pairAddress: pair.pairAddress,
    price: `$${parseFloat(pair.priceUsd).toLocaleString()}`,
    change: `${pair.priceChange.h24 >= 0 ? '+' : ''}${pair.priceChange.h24.toFixed(1)}%`,
    icon: pair.info?.imageUrl || pair.baseToken.symbol.charAt(0),
    color: 'bg-[#1C1C1E]',
    marketCap: pair.marketCap,
    volume24h: pair.volume.h24,
    priceHistory: [],
  };
}; 