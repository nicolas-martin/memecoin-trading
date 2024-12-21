interface CoinIcon {
  icon: string;
  color: string;
}

const COIN_ICONS: Record<string, CoinIcon> = {
  BTC: { icon: '₿', color: 'bg-[#F7931A]' },
  ETH: { icon: 'Ξ', color: 'bg-[#627EEA]' },
  SOL: { icon: '◎', color: 'bg-[#DC1FFF]' },
  DOGE: { icon: 'Ð', color: 'bg-[#C2A633]' },
  USDT: { icon: '₮', color: 'bg-[#26A17B]' },
  // Add more coins as needed
};

export const getIconForCoin = (symbol: string): CoinIcon => {
  return COIN_ICONS[symbol.toUpperCase()] || {
    icon: symbol.charAt(0),
    color: 'bg-gray-500'
  };
};

export const getCoinColor = (symbol: string): string => {
  return COIN_ICONS[symbol.toUpperCase()]?.color || 'bg-gray-500';
}; 