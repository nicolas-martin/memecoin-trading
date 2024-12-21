import { useState } from 'react';
import { Link } from 'react-router-dom';
import TrendingCard from '../components/cards/TrendingCard';
import TabBar from '../components/navigation/TabBar';
import CoinListItem from '../components/lists/CoinListItem';
import TradeModal from '../components/modals/TradeModal';

const TRENDING_COINS = [
  {
    name: 'Bitcoin',
    symbol: 'BTC',
    icon: '₿',
    color: 'bg-[#F7931A]',
    change: '+1.2%',
    previousClose: '$356.08B',
    open: '$123.88B',
    marketCap: '$123.88B',
  },
  // Add more trending coins
];

const TABS = [
  { id: 'coins', label: 'Coins' },
  { id: 'watchlist', label: 'Watchlist' },
  { id: 'recent', label: 'Recently Added' },
  { id: 'rated', label: 'Top Rated' },
];

const COIN_LIST = [
  { name: 'Bitcoin', symbol: 'BTC', price: '$8,907.02', change: '+1.2%', icon: '₿', color: 'bg-[#F7931A]' },
  { name: 'Dash', symbol: 'DASH', price: '$8,907.02', change: '-1.2%', icon: 'D', color: 'bg-blue-500' },
  { name: 'Pundi X', symbol: 'NPXS', price: '$8,907.02', change: '-1.2%', icon: 'P', color: 'bg-yellow-500' },
  // Add more coins
];

const Portfolio = () => {
  const [activeTab, setActiveTab] = useState('coins');
  const [selectedCoin, setSelectedCoin] = useState<null | typeof COIN_LIST[0]>(null);

  return (
    <div className="min-h-screen bg-black text-white">
      {/* Header */}
      <div className="px-4 pt-12 pb-4">
        <div className="flex justify-between items-center">
          <h1 className="text-3xl font-bold">Markets</h1>
          <button className="p-2">
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </button>
        </div>
      </div>

      {/* Trending Section */}
      <div className="px-4 mb-6">
        <h2 className="text-2xl font-bold mb-4">Trending</h2>
        <div className="flex space-x-4 overflow-x-auto pb-4 hide-scrollbar">
          {TRENDING_COINS.map((coin, index) => (
            <TrendingCard key={index} coin={coin} />
          ))}
        </div>
      </div>

      {/* Tab Navigation */}
      <div className="px-4">
        <TabBar tabs={TABS} activeTab={activeTab} onTabChange={setActiveTab} />
      </div>

      {/* Coin List */}
      <div className="px-4">
        {COIN_LIST.map((coin, index) => (
          <CoinListItem
            key={index}
            coin={coin}
            onTrade={() => setSelectedCoin(coin)}
          />
        ))}
      </div>

      {selectedCoin && (
        <TradeModal
          coin={selectedCoin}
          isOpen={!!selectedCoin}
          onClose={() => setSelectedCoin(null)}
        />
      )}
    </div>
  );
};

export default Portfolio; 