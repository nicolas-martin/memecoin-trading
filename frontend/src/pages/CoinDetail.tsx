import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { fetchCoinData, CoinData } from '../services/dexscreener';
import { ChevronLeftIcon } from '@heroicons/react/24/outline';
import MiniChart from '../components/charts/MiniChart';
import PriceChart from '../components/charts/PriceChart';
import TradeModal from '../components/modals/TradeModal';
import LoadingSpinner from '../components/LoadingSpinner';
import ErrorMessage from '../components/ErrorMessage';

const MOCK_COIN_DATA: Record<string, CoinData> = {
  '0x7130d2a12b9bcbfae4f2634d864a1ee1ce3ead9c': {
    name: 'Bitcoin',
    symbol: 'BTC',
    pairAddress: '0x7130d2a12b9bcbfae4f2634d864a1ee1ce3ead9c',
    price: '$45,000.00',
    change: '+2.3%',
    icon: '₿',
    color: 'bg-[#F7931A]',
    marketCap: 850000000000,
    volume24h: 28000000000,
    priceHistory: [],
  },
  '0x123': {
    name: 'Dash',
    symbol: 'DASH',
    pairAddress: '0x123',
    price: '$120.00',
    change: '-1.2%',
    icon: 'D',
    color: 'bg-blue-500',
    marketCap: 1200000000,
    volume24h: 500000000,
    priceHistory: [],
  },
  '0x456': {
    name: 'Pundi X',
    symbol: 'NPXS',
    pairAddress: '0x456',
    price: '$0.002',
    change: '-1.2%',
    icon: 'P',
    color: 'bg-yellow-500',
    marketCap: 500000000,
    volume24h: 25000000,
    priceHistory: [],
  },
};

const CoinDetail = () => {
  const { pairAddress } = useParams();
  const navigate = useNavigate();
  const [coinData, setCoinData] = useState<CoinData | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showTradeModal, setShowTradeModal] = useState(false);

  useEffect(() => {
    const loadCoinData = async () => {
      if (!pairAddress) return;
      try {
        setIsLoading(true);
        // Use mock data instead of fetching
        const mockData = MOCK_COIN_DATA[pairAddress];
        if (mockData) {
          setCoinData(mockData);
        } else {
          setError('Coin not found');
        }
      } catch (err) {
        setError('Failed to load coin data');
        console.error(err);
      } finally {
        setIsLoading(false);
      }
    };

    loadCoinData();
  }, [pairAddress]);

  if (isLoading) {
    return <LoadingSpinner />;
  }

  if (error || !coinData) {
    return <ErrorMessage message={error || 'Failed to load coin data'} />;
  }

  return (
    <div className="min-h-screen bg-black text-white safe-area-top">
      {/* Status Bar Space */}
      <div className="h-safe-top bg-black" />

      {/* Fixed Header */}
      <div className="sticky top-0 z-10 bg-black">
        <div className="px-4 py-4">
          <div className="flex items-center">
            <button 
              onClick={() => navigate(-1)}
              className="w-8 h-8 flex items-center justify-center rounded-full bg-[#1C1C1E] active:opacity-70 mr-3"
            >
              <ChevronLeftIcon className="w-5 h-5" />
            </button>
            <div className="flex items-center">
              <div className={`w-8 h-8 ${coinData.color} rounded-full flex items-center justify-center mr-2`}>
                <span className="text-white text-sm font-medium">{coinData.icon}</span>
              </div>
              <h1 className="text-xl font-bold">{coinData.name}</h1>
            </div>
          </div>
        </div>
      </div>

      {/* Scrollable Content */}
      <div className="flex-1 overflow-auto pb-safe-bottom">
        {/* Price Section */}
        <div className="px-4 py-4">
          <div className="flex justify-between items-start mb-4">
            <div>
              <span className="text-3xl font-bold">{coinData.price}</span>
              <div className={`text-sm ${coinData.change.startsWith('+') ? 'text-green-400' : 'text-red-400'}`}>
                {coinData.change}
              </div>
            </div>
            <button 
              onClick={() => setShowTradeModal(true)}
              className="bg-green-500 px-6 py-2 rounded-lg font-medium active:opacity-90"
            >
              Trade
            </button>
          </div>

          {/* Price Chart */}
          <div className="bg-[#1C1C1E] rounded-xl p-4 mb-6">
            <PriceChart
              pairAddress={coinData.pairAddress}
              isPositive={coinData.change.startsWith('+')}
            />
          </div>

          {/* Stats Grid */}
          <div className="grid grid-cols-2 gap-4">
            {[
              { label: 'Market Cap', value: `$${coinData.marketCap.toLocaleString()}` },
              { label: '24h Volume', value: `$${coinData.volume24h.toLocaleString()}` },
              { label: 'Circulating Supply', value: `${(coinData.marketCap / parseFloat(coinData.price.slice(1))).toFixed(2)} ${coinData.symbol}` },
              { label: 'Max Supply', value: `${(coinData.marketCap / parseFloat(coinData.price.slice(1)) * 1.1).toFixed(2)} ${coinData.symbol}` },
            ].map((stat, index) => (
              <div key={index} className="bg-[#1C1C1E] p-4 rounded-xl">
                <div className="text-gray-400 text-sm mb-1">{stat.label}</div>
                <div className="font-medium">{stat.value}</div>
              </div>
            ))}
          </div>

          {/* Price History */}
          <div className="mt-6">
            <h2 className="text-lg font-bold mb-3">Price History</h2>
            <div className="space-y-3">
              {[
                { period: '1H', change: '+0.5%' },
                { period: '24H', change: coinData.change },
                { period: '7D', change: '-3.1%' },
                { period: '1M', change: '+15.7%' },
                { period: '1Y', change: '+124.3%' },
              ].map((item, index) => (
                <div 
                  key={index}
                  className="flex justify-between items-center py-3 border-b border-gray-800"
                >
                  <span className="text-gray-400">{item.period}</span>
                  <span className={item.change.startsWith('+') ? 'text-green-400' : 'text-red-400'}>
                    {item.change}
                  </span>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>

      {/* Trade Modal */}
      {showTradeModal && (
        <TradeModal
          coin={coinData}
          isOpen={showTradeModal}
          onClose={() => setShowTradeModal(false)}
        />
      )}
    </div>
  );
};

export default CoinDetail; 