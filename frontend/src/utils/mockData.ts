import { PairData } from '../services/dexscreener';

const MOCK_PAIRS: Record<string, Partial<PairData>> = {
  'BTC': {
    baseToken: {
      name: 'Bitcoin',
      symbol: 'BTC',
      address: '0x7130d2a12b9bcbfae4f2634d864a1ee1ce3ead9c'
    },
    priceUsd: '45000',
    priceChange: {
      h1: 0.5,
      h24: 2.3,
      d7: 5.1
    },
    volume: {
      h24: 1200000000,
      h6: 300000000,
      h1: 50000000,
      m5: 4000000
    },
    liquidity: {
      usd: 500000000,
      base: 11000,
      quote: 500000000
    },
    marketCap: 850000000000,
    fdv: 950000000000,
    info: {
      imageUrl: 'https://assets.coingecko.com/coins/images/1/large/bitcoin.png',
      description: 'Bitcoin is a decentralized cryptocurrency.',
      websites: [{ url: 'https://bitcoin.org' }],
      socials: [
        { platform: 'twitter', handle: 'bitcoin' }
      ]
    }
  },
  'ETH': {
    baseToken: {
      name: 'Ethereum',
      symbol: 'ETH',
      address: '0x2170ed0880ac9a755fd29b2688956bd959f933f8'
    },
    priceUsd: '2800',
    priceChange: {
      h1: -0.2,
      h24: 1.5,
      d7: 3.2
    },
    volume: {
      h24: 800000000,
      h6: 200000000,
      h1: 30000000,
      m5: 2500000
    },
    liquidity: {
      usd: 300000000,
      base: 107000,
      quote: 300000000
    },
    marketCap: 350000000000,
    fdv: 380000000000,
    info: {
      imageUrl: 'https://assets.coingecko.com/coins/images/279/large/ethereum.png',
      description: 'Ethereum is a decentralized computing platform.',
      websites: [{ url: 'https://ethereum.org' }],
      socials: [
        { platform: 'twitter', handle: 'ethereum' }
      ]
    }
  }
};

export function generateMockPair(symbol: string = 'UNKNOWN'): PairData {
  const mockData = MOCK_PAIRS[symbol] || {};
  const randomPrice = Math.floor(Math.random() * 1000) + 1;
  const randomChange = (Math.random() * 10) - 5; // -5 to +5

  return {
    chainId: 'bsc',
    dexId: 'pancakeswap',
    url: `https://dexscreener.com/bsc/${symbol.toLowerCase()}`,
    pairAddress: `0x${Math.random().toString(16).slice(2, 42)}`,
    baseToken: {
      address: `0x${Math.random().toString(16).slice(2, 42)}`,
      name: mockData.baseToken?.name || `${symbol} Token`,
      symbol: mockData.baseToken?.symbol || symbol,
    },
    quoteToken: {
      address: '0xe9e7cea3dedca5984780bafc599bd69add087d56',
      name: 'BUSD',
      symbol: 'BUSD',
    },
    priceNative: mockData.priceUsd || randomPrice.toString(),
    priceUsd: mockData.priceUsd || randomPrice.toString(),
    liquidity: mockData.liquidity || {
      usd: randomPrice * 1000000,
      base: randomPrice * 100,
      quote: randomPrice * 1000000
    },
    volume: mockData.volume || {
      h24: randomPrice * 100000,
      h6: randomPrice * 25000,
      h1: randomPrice * 4000,
      m5: randomPrice * 300
    },
    priceChange: mockData.priceChange || {
      h1: randomChange * 0.2,
      h24: randomChange,
      d7: randomChange * 2
    },
    fdv: mockData.fdv || randomPrice * 1000000,
    marketCap: mockData.marketCap || randomPrice * 800000,
    info: mockData.info || {
      imageUrl: `https://via.placeholder.com/64/1C1C1E/FFFFFF?text=${symbol}`,
      description: `${symbol} is a cryptocurrency token.`,
      websites: [{ url: `https://${symbol.toLowerCase()}.org` }],
      socials: [
        { platform: 'twitter', handle: symbol.toLowerCase() }
      ]
    }
  };
}

export function generateMockPairs(count: number = 10): PairData[] {
  const knownSymbols = Object.keys(MOCK_PAIRS);
  const pairs: PairData[] = [];

  // Add known pairs first
  knownSymbols.forEach(symbol => {
    pairs.push(generateMockPair(symbol));
  });

  // Add random pairs to reach desired count
  for (let i = pairs.length; i < count; i++) {
    const randomSymbol = `TOKEN${i}`;
    pairs.push(generateMockPair(randomSymbol));
  }

  return pairs;
} 