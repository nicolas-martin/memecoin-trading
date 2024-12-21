import { useState, useEffect } from 'react';
import { fetchHistoricalPrices, HistoricalPrice } from '../services/dexscreener';

export function useCurrentPrice(pairAddress: string) {
  const [price, setPrice] = useState<number | null>(null);
  const [change, setChange] = useState<string>('0.00%');
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchPrice = async () => {
      try {
        setIsLoading(true);
        setError(null);
        const prices = await fetchHistoricalPrices(pairAddress, '24H');
        if (prices && prices.length > 0) {
          const currentPrice = prices[prices.length - 1].price;
          const previousPrice = prices[0].price;
          setPrice(currentPrice);
          const priceChange = previousPrice > 0 
            ? ((currentPrice - previousPrice) / previousPrice * 100).toFixed(2)
            : '0.00';
          setChange(priceChange.startsWith('-') ? priceChange + '%' : '+' + priceChange + '%');
        } else {
          setPrice(null);
          setChange('0.00%');
        }
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to fetch price');
        setPrice(null);
        setChange('0.00%');
      } finally {
        setIsLoading(false);
      }
    };

    if (pairAddress) {
      fetchPrice();
      const interval = setInterval(fetchPrice, 30000); // Refresh every 30 seconds
      return () => clearInterval(interval);
    }
  }, [pairAddress]);

  return { price, change, isLoading, error };
} 