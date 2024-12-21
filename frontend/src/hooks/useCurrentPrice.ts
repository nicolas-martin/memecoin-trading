import { useState, useEffect } from 'react';
import { fetchHistoricalPrices } from '../services/dexscreener';

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
        if (prices.length > 0) {
          const currentPrice = prices[prices.length - 1].price;
          const previousPrice = prices[0].price;
          setPrice(currentPrice);
          setChange(((currentPrice - previousPrice) / previousPrice * 100).toFixed(2) + '%');
        }
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to fetch price');
      } finally {
        setIsLoading(false);
      }
    };

    if (pairAddress) {
      fetchPrice();
      const interval = setInterval(fetchPrice, 30000); // Update every 30 seconds
      return () => clearInterval(interval);
    }
  }, [pairAddress]);

  return { price, change, isLoading, error };
} 