import { useState, useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Link } from 'react-router-dom';
import TrendingCard from '../components/cards/TrendingCard';
import TabBar from '../components/navigation/TabBar';
import CoinListItem from '../components/lists/CoinListItem';
import TradeModal from '../components/modals/TradeModal';
import { fetchTrendingCoins, setSelectedCoin } from '../store/slices/coinSlice';
import type { RootState, AppDispatch } from '../store';
import { PairData } from '../services/dexscreener';

const TABS = [
  { id: 'coins', label: 'Coins' },
  { id: 'watchlist', label: 'Watchlist' },
  { id: 'recent', label: 'Recently Added' },
  { id: 'rated', label: 'Top Rated' },
];

const COIN_LIST = [
  {
    name: 'Bitcoin',
    symbol: 'BTC',
    price: '$8,907.02',
    change: '+1.2%',
    icon: '₿',
    color: 'bg-[#F7931A]',
    pairAddress: '0x7130d2a12b9bcbfae4f2634d864a1ee1ce3ead9c',
  },
  { name: 'Dash', symbol: 'DASH', price: '$8,907.02', change: '-1.2%', icon: 'D', color: 'bg-blue-500', pairAddress: '0x123' },
  { name: 'Pundi X', symbol: 'NPXS', price: '$8,907.02', change: '-1.2%', icon: 'P', color: 'bg-yellow-500', pairAddress: '0x456' },
];

const Portfolio = () => {
  const dispatch = useDispatch<AppDispatch>();
  const [activeTab, setActiveTab] = useState('coins');
  const { trendingCoins, loading, error, selectedCoin } = useSelector((state: RootState) => state.coins);

  useEffect(() => {
    dispatch(fetchTrendingCoins());
  }, [dispatch]);

  return (
    <div className="min-h-screen bg-black text-white safe-area-top">
      {/* Status Bar Space */}
      <div className="h-safe-top bg-black" />

      {/* Fixed Header */}
      <div className="sticky top-0 z-10 bg-black">
        <div className="px-4 py-4">
          <div className="flex justify-between items-center">
            <h1 className="text-2xl font-bold">Markets</h1>
            <button className="w-8 h-8 flex items-center justify-center rounded-full bg-[#1C1C1E] active:opacity-70">
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </button>
          </div>
        </div>

        {/* Tab Navigation - Sticky */}
        <div className="px-4">
          <TabBar tabs={TABS} activeTab={activeTab} onTabChange={setActiveTab} />
        </div>
      </div>

      {/* Scrollable Content */}
      <div className="flex-1 overflow-auto pb-safe-bottom">
        {/* Trending Section */}
        <div className="px-4 mb-6">
          <h2 className="text-2xl font-bold mb-4">Trending</h2>
          {loading ? (
            <div className="flex justify-center py-8">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-white"></div>
            </div>
          ) : error ? (
            <div className="text-red-400 text-center py-8">{error}</div>
          ) : (
            <div className="flex space-x-4 overflow-x-auto pb-4 hide-scrollbar -mx-4 px-4">
              {trendingCoins.map((pair: PairData) => (
                <TrendingCard 
                  key={pair.pairAddress}
                  coin={{
                    name: pair.baseToken.name,
                    symbol: pair.baseToken.symbol,
                    icon: pair.info?.imageUrl || pair.baseToken.symbol.charAt(0),
                    color: 'bg-[#1C1C1E]',
                    change: `${pair.priceChange.h24 >= 0 ? '+' : ''}${pair.priceChange.h24.toFixed(1)}%`,
                    previousClose: `$${parseFloat(pair.priceUsd).toLocaleString()}`,
                    open: `$${(parseFloat(pair.priceUsd) * (1 - pair.priceChange.h24/100)).toLocaleString()}`,
                    marketCap: `$${pair.marketCap.toLocaleString()}`,
                  }}
                  className="w-[280px] flex-shrink-0"
                />
              ))}
            </div>
          )}
        </div>

        {/* Coin List */}
        <div className="px-4 pb-32"> {/* Extra padding for bottom nav */}
          {COIN_LIST.map((coin, index) => (
            <Link to={`/coins/${coin.pairAddress}`} key={index}>
              <CoinListItem
                coin={coin}
                onTrade={() => dispatch(setSelectedCoin(coin))}
              />
            </Link>
          ))}
        </div>
      </div>

      {/* Trade Modal */}
      {selectedCoin && (
        <TradeModal
          coin={selectedCoin}
          isOpen={!!selectedCoin}
          onClose={() => dispatch(setSelectedCoin(null))}
        />
      )}
    </div>
  );
};

export default Portfolio; 