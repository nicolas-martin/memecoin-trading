import { useState } from 'react';
import { XMarkIcon } from '@heroicons/react/24/outline';

interface TradeModalProps {
  coin: {
    name: string;
    symbol: string;
    price: string;
    icon: string;
    color: string;
  };
  isOpen: boolean;
  onClose: () => void;
}

const TradeModal = ({ coin, isOpen, onClose }: TradeModalProps) => {
  const [amount, setAmount] = useState('');
  const [type, setType] = useState<'buy' | 'sell'>('buy');

  if (!isOpen) return null;

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    // Handle trade submission
    console.log({ type, amount, coin });
    onClose();
  };

  return (
    <div className="fixed inset-0 z-50 bg-black bg-opacity-75">
      <div className="min-h-screen px-4 text-center">
        <div className="fixed inset-0" onClick={onClose} />

        <div className="inline-block w-full max-w-md p-6 my-8 text-left align-middle transition-all transform bg-[#1C1C1E] rounded-2xl shadow-xl">
          <div className="flex justify-between items-center mb-4">
            <div className="flex items-center space-x-3">
              <div className={`w-10 h-10 ${coin.color} rounded-full flex items-center justify-center`}>
                <span className="text-white text-lg">{coin.icon}</span>
              </div>
              <div>
                <h3 className="text-xl font-semibold text-white">{coin.name}</h3>
                <span className="text-gray-400">{coin.symbol}</span>
              </div>
            </div>
            <button onClick={onClose} className="text-gray-400 hover:text-white">
              <XMarkIcon className="w-6 h-6" />
            </button>
          </div>

          <div className="mb-6">
            <div className="flex justify-between mb-2">
              <span className="text-gray-400">Current Price</span>
              <span className="text-white">{coin.price}</span>
            </div>
          </div>

          <form onSubmit={handleSubmit}>
            <div className="flex gap-2 mb-4">
              <button
                type="button"
                className={`flex-1 py-2 px-4 rounded-lg font-medium ${
                  type === 'buy'
                    ? 'bg-green-500 text-white'
                    : 'bg-[#2C2C2E] text-gray-400'
                }`}
                onClick={() => setType('buy')}
              >
                Buy
              </button>
              <button
                type="button"
                className={`flex-1 py-2 px-4 rounded-lg font-medium ${
                  type === 'sell'
                    ? 'bg-red-500 text-white'
                    : 'bg-[#2C2C2E] text-gray-400'
                }`}
                onClick={() => setType('sell')}
              >
                Sell
              </button>
            </div>

            <div className="mb-4">
              <label className="block text-gray-400 mb-2">Amount</label>
              <div className="relative">
                <input
                  type="number"
                  value={amount}
                  onChange={(e) => setAmount(e.target.value)}
                  className="w-full bg-[#2C2C2E] text-white px-4 py-3 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="0.00"
                />
                <span className="absolute right-4 top-1/2 transform -translate-y-1/2 text-gray-400">
                  USD
                </span>
              </div>
            </div>

            <button
              type="submit"
              className={`w-full py-3 px-4 rounded-lg font-medium ${
                type === 'buy' ? 'bg-green-500' : 'bg-red-500'
              } text-white`}
            >
              {type === 'buy' ? 'Buy' : 'Sell'} {coin.symbol}
            </button>
          </form>
        </div>
      </div>
    </div>
  );
};

export default TradeModal; 