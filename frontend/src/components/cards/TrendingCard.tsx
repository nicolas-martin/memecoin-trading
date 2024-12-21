import MiniChart from '../charts/MiniChart';
import { useNavigate } from 'react-router-dom';

interface TrendingCardProps {
  coin: {
    name: string;
    symbol: string;
    icon: string;
    color: string;
    change: string;
    previousClose: string;
    open: string;
    marketCap: string;
  };
  className?: string;
}

const TrendingCard = ({ coin, className = '' }: TrendingCardProps) => {
  const navigate = useNavigate();
  const mockChartData = [10, 41, 35, 51, 49, 62, 69, 91, 148];
  const isPositive = coin.change.startsWith('+');

  return (
    <div 
      onClick={() => navigate(`/coins/${coin.symbol.toLowerCase()}`)}
      className={`bg-[#1C1C1E] rounded-xl p-4 active:opacity-90 transition-opacity touch-callout-none cursor-pointer ${className}`}
    >
      <div className="flex items-center space-x-3 mb-3">
        <div className={`w-8 h-8 ${coin.color} rounded-full flex items-center justify-center`}>
          <span className="text-white text-lg">{coin.icon}</span>
        </div>
        <div>
          <h3 className="font-semibold">{coin.name}</h3>
          <span className="text-sm text-gray-400">{coin.symbol}</span>
        </div>
        <div className="ml-auto">
          <span className={isPositive ? 'text-green-400' : 'text-red-400'}>
            {coin.change}
          </span>
        </div>
      </div>
      
      <MiniChart 
        data={mockChartData} 
        color={isPositive ? '#34D399' : '#EF4444'} 
        isPositive={isPositive} 
      />

      <div className="space-y-2 mt-3">
        <div className="flex justify-between text-sm text-gray-400">
          <span>Previous Close</span>
          <span>{coin.previousClose}</span>
        </div>
        <div className="flex justify-between text-sm text-gray-400">
          <span>Open</span>
          <span>{coin.open}</span>
        </div>
        <div className="flex justify-between text-sm text-gray-400">
          <span>Market Cap</span>
          <span>{coin.marketCap}</span>
        </div>
      </div>
    </div>
  );
};

export default TrendingCard; 