import { CoinPrice } from '../types/coin';

const calculateMA = (data: CoinPrice[], period: number) => {
  const prices = data.map((d) => d.price);
  const ma = prices.map((_, index) => {
    if (index < period - 1) return null;
    const slice = prices.slice(index - period + 1, index + 1);
    return slice.reduce((sum, price) => sum + price, 0) / period;
  });

  return {
    label: `MA(${period})`,
    data: data.map((price, index) => ({
      x: new Date(price.timestamp),
      y: ma[index],
    })),
    borderColor: `rgba(255, 99, 132, ${period === 20 ? 0.8 : 0.5})`,
    borderWidth: 1,
    pointRadius: 0,
    order: 2,
  };
};

const calculateEMA = (data: CoinPrice[], period: number) => {
  const prices = data.map((d) => d.price);
  const multiplier = 2 / (period + 1);
  const ema = prices.reduce((acc: number[], price, index) => {
    if (index === 0) {
      acc.push(price);
      return acc;
    }
    acc.push(price * multiplier + acc[index - 1] * (1 - multiplier));
    return acc;
  }, []);

  return {
    label: `EMA(${period})`,
    data: data.map((price, index) => ({
      x: new Date(price.timestamp),
      y: ema[index],
    })),
    borderColor: `rgba(54, 162, 235, ${period === 20 ? 0.8 : 0.5})`,
    borderWidth: 1,
    pointRadius: 0,
    order: 3,
  };
};

const calculateRSI = (data: CoinPrice[], period: number = 14) => {
  const prices = data.map((d) => d.price);
  const deltas = prices.map((price, index) => {
    if (index === 0) return 0;
    return price - prices[index - 1];
  });

  const gains = deltas.map((delta) => (delta > 0 ? delta : 0));
  const losses = deltas.map((delta) => (delta < 0 ? -delta : 0));

  const avgGain = gains.reduce((acc, gain, i) => {
    if (i < period) {
      acc.push(i === period - 1 ? gains.slice(0, period).reduce((sum, g) => sum + g, 0) / period : 0);
    } else {
      acc.push((acc[i - 1] * (period - 1) + gain) / period);
    }
    return acc;
  }, [] as number[]);

  const avgLoss = losses.reduce((acc, loss, i) => {
    if (i < period) {
      acc.push(i === period - 1 ? losses.slice(0, period).reduce((sum, l) => sum + l, 0) / period : 0);
    } else {
      acc.push((acc[i - 1] * (period - 1) + loss) / period);
    }
    return acc;
  }, [] as number[]);

  const rsi = avgGain.map((gain, i) => {
    if (i < period) return null;
    const rs = gain / avgLoss[i];
    return 100 - (100 / (1 + rs));
  });

  return {
    label: 'RSI(14)',
    data: data.map((price, index) => ({
      x: new Date(price.timestamp),
      y: rsi[index],
    })),
    borderColor: 'rgba(153, 102, 255, 1)',
    borderWidth: 1,
    pointRadius: 0,
    yAxisID: 'rsi',
    order: 4,
  };
};

export const calculateIndicators = (
  data: CoinPrice[],
  activeIndicators: string[]
) => {
  const indicators = [];

  if (activeIndicators.includes('MA')) {
    indicators.push(calculateMA(data, 20));
    indicators.push(calculateMA(data, 50));
  }

  if (activeIndicators.includes('EMA')) {
    indicators.push(calculateEMA(data, 12));
    indicators.push(calculateEMA(data, 26));
  }

  if (activeIndicators.includes('RSI')) {
    indicators.push(calculateRSI(data));
  }

  return indicators;
}; 