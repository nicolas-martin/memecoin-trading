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