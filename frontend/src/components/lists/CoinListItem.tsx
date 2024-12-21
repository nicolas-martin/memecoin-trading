import { useNavigate } from 'react-router-dom';
import { useCurrentPrice } from '../../hooks/useCurrentPrice';
import { formatPrice } from '../../utils/format';

interface CoinListItemProps {
  coin: {
    name: string;
    symbol: string;
    price: string;
    change: string;
    icon: string;
    color: string;
    pairAddress: string;
  };
  onTrade: () => void;
}

const CoinListItem = ({ coin, onTrade }: CoinListItemProps) => {
  const navigate = useNavigate();
  const { price, change, isLoading } = useCurrentPrice(coin.pairAddress);

  return (
    <div 
      className="flex items-center py-4 border-b border-gray-800 active:bg-[#1C1C1E] transition-colors"
      onClick={(e) => {
        e.stopPropagation();
        navigate(`/coins/${coin.symbol.toLowerCase()}`);
      }}
    >
      <div className={`w-10 h-10 ${coin.color} rounded-full flex items-center justify-center mr-3`}>
        <span className="text-white font-medium">{coin.icon}</span>
      </div>
      <div className="flex-1">
        <h3 className="font-medium text-white">{coin.name}</h3>
        <span className="text-sm text-gray-400">{coin.symbol}</span>
      </div>
      <div className="text-right">
        <div className="font-medium text-white">
          {isLoading ? '...' : formatPrice(price || 0)}
        </div>
        <div className={change.startsWith('+') ? 'text-green-400' : 'text-red-400'}>
          {isLoading ? '...' : change}
        </div>
      </div>
      <button
        className="ml-4 px-4 py-2 bg-[#2C2C2E] rounded-lg text-sm font-medium active:opacity-90"
        onClick={(e) => {
          e.stopPropagation();
          onTrade();
        }}
      >
        Trade
      </button>
    </div>
  );
};

export default CoinListItem; 