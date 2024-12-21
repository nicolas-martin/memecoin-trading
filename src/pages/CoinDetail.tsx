import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { useSelector, useDispatch } from 'react-redux';
import { fetchCoinById } from '../store/slices/coinSlice';
import PriceChart from '../components/charts/PriceChart';
import ChartControls from '../components/charts/ChartControls';
import OrderForm from '../components/trading/OrderForm';
import { getCoinPriceHistory } from '../services/api';
import { CoinPrice } from '../types/coin';
import { formatCurrency, formatPercentage } from '../utils/formatters';

const CoinDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const dispatch = useDispatch();
  const coin = useSelector((state: any) => state.coins.selectedCoin);
  const [priceHistory, setPriceHistory] = useState<CoinPrice[]>([]);
  const [timeframe, setTimeframe] = useState<'24h' | '7d' | '30d' | '1y'>('24h');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (id) {
      dispatch(fetchCoinById(id));
    }
  }, [dispatch, id]);

  useEffect(() => {
    const fetchPriceHistory = async () => {
      setLoading(true);
      try {
        const history = await getCoinPriceHistory(id!, timeframe);
        setPriceHistory(history);
      } catch (error) {
        console.error('Failed to fetch price history:', error);
      } finally {
        setLoading(false);
      }
    };

    if (id) {
      fetchPriceHistory();
    }
  }, [id, timeframe]);

  if (!coin) {
    return <div>Loading...</div>;
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <div className="lg:col-span-2">
          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center justify-between mb-6">
              <div className="flex items-center">
                <img
                  src={coin.logoUrl}
                  alt={coin.name}
                  className="w-12 h-12 rounded-full"
                />
                <div className="ml-4">
                  <h1 className="text-2xl font-bold">{coin.name}</h1>
                  <p className="text-gray-500">{coin.symbol}</p>
                </div>
              </div>
              <div className="text-right">
                <div className="text-2xl font-bold">
                  {formatCurrency(coin.price)}
                </div>
                <div
                  className={`text-sm ${
                    coin.priceChange24h >= 0 ? 'text-green-600' : 'text-red-600'
                  }`}
                >
                  {formatPercentage(coin.priceChange24h)}
                </div>
              </div>
            </div>

            <div className="mb-4">
              <ChartControls
                timeframe={timeframe}
                onTimeframeChange={setTimeframe}
              />
            </div>

            <PriceChart
              data={priceHistory}
              symbol={coin.symbol}
              loading={loading}
              timeframe={timeframe}
            />
          </div>
        </div>

        <div className="lg:col-span-1">
          <div className="bg-white rounded-lg shadow p-6">
            <div className="mb-6">
              <h2 className="text-lg font-semibold mb-4">Trade</h2>
              <OrderForm
                coin={coin}
                type="BUY"
                balance={1000} // Replace with actual user balance
                onSuccess={() => {
                  // Handle success
                }}
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default CoinDetail; 