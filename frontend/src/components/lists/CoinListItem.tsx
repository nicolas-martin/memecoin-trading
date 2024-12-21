interface CoinListItemProps {
  coin: {
    name: string;
    symbol: string;
    price: string;
    change: string;
    icon: string;
    color: string;
  };
  onTrade: () => void;
}

const CoinListItem = ({ coin, onTrade }: CoinListItemProps) => {
  return (
    <div 
      className="flex items-center py-4 border-b border-gray-800 active:bg-[#1C1C1E] transition-colors"
      onClick={onTrade}
    >
      <div className={`w-10 h-10 ${coin.color} rounded-full flex items-center justify-center mr-3`}>
        <span className="text-white font-medium">{coin.icon}</span>
      </div>
      <div className="flex-1">
        <h3 className="font-medium text-white">{coin.name}</h3>
        <span className="text-sm text-gray-400">{coin.symbol}</span>
      </div>
      <div className="text-right">
        <div className="font-medium text-white">{coin.price}</div>
        <div className={coin.change.startsWith('+') ? 'text-green-400' : 'text-red-400'}>
          {coin.change}
        </div>
      </div>
    </div>
  );
};

export default CoinListItem; 