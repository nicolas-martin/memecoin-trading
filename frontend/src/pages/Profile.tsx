import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { CogIcon } from '@heroicons/react/24/outline';
import PortfolioChart from '../components/charts/PortfolioChart';
import AddFundsModal from '../components/modals/AddFundsModal';
import { getIconForCoin } from '../utils/coinIcons';

const Profile = () => {
  const [showAddFunds, setShowAddFunds] = useState(false);
  const navigate = useNavigate();

  // Mock user data - in a real app, this would come from your auth context/redux
  const user = {
    name: 'John Doe',
    balance: 10000.00,
    profitLoss: 2345.67,
    profitLossPercentage: 23.45,
  };

  // Mock holdings with icon data
  const holdings = [
    { 
      coin: 'Bitcoin', 
      symbol: 'BTC', 
      amount: 0.5, 
      value: 15000, 
      change: '+2.3%',
      ...getIconForCoin('BTC')
    },
    { 
      coin: 'Ethereum', 
      symbol: 'ETH', 
      amount: 2.5, 
      value: 8000, 
      change: '-1.2%',
      ...getIconForCoin('ETH')
    },
    { 
      coin: 'Solana', 
      symbol: 'SOL', 
      amount: 15.5, 
      value: 3000, 
      change: '+5.2%',
      ...getIconForCoin('SOL')
    }
  ];

  return (
    <div className="min-h-screen bg-black text-white safe-area-top">
      {/* Status Bar Space */}
      <div className="h-safe-top bg-black" />

      {/* Header */}
      <div className="sticky top-0 z-10 bg-black">
        <div className="px-4 py-4">
          <div className="flex justify-between items-center">
            <h1 className="text-2xl font-bold">Profile</h1>
            <button 
              onClick={() => navigate('/settings')}
              className="w-8 h-8 flex items-center justify-center rounded-full bg-[#1C1C1E] active:opacity-70"
            >
              <CogIcon className="w-5 h-5" />
            </button>
          </div>
        </div>
      </div>

      {/* Content */}
      <div className="px-4 py-4">
        {/* Balance Card */}
        <div className="bg-[#1C1C1E] rounded-xl p-4 mb-6">
          <div className="flex justify-between items-start mb-4">
            <div>
              <h2 className="text-gray-400 text-sm mb-1">Total Balance</h2>
              <span className="text-3xl font-bold">${user.balance.toLocaleString()}</span>
            </div>
            <button 
              onClick={() => setShowAddFunds(true)}
              className="bg-green-500 px-4 py-2 rounded-lg text-sm font-medium active:opacity-90"
            >
              Add Funds
            </button>
          </div>
          <div className="flex items-center">
            <span className={`text-sm ${user.profitLoss >= 0 ? 'text-green-400' : 'text-red-400'}`}>
              {user.profitLoss >= 0 ? '+' : '-'}${Math.abs(user.profitLoss).toLocaleString()}
            </span>
            <span className={`text-xs ml-2 ${user.profitLoss >= 0 ? 'text-green-400' : 'text-red-400'}`}>
              ({user.profitLossPercentage}%)
            </span>
          </div>
        </div>

        {/* Portfolio Chart */}
        <div className="bg-[#1C1C1E] rounded-xl p-4 mb-6">
          <h2 className="text-lg font-bold mb-4">Portfolio Value</h2>
          <PortfolioChart />
        </div>

        {/* Holdings */}
        <div className="bg-[#1C1C1E] rounded-xl p-4">
          <h2 className="text-lg font-bold mb-4">Your Holdings</h2>
          <div className="space-y-4">
            {holdings.map((holding, index) => (
              <div 
                key={index} 
                className="flex items-center justify-between py-3 border-b border-gray-800 last:border-b-0"
              >
                <div className="flex items-center">
                  <div className={`w-10 h-10 ${holding.color} rounded-full flex items-center justify-center mr-3`}>
                    <span className="text-white font-medium">{holding.icon}</span>
                  </div>
                  <div>
                    <h3 className="font-medium">{holding.coin}</h3>
                    <span className="text-sm text-gray-400">{holding.amount} {holding.symbol}</span>
                  </div>
                </div>
                <div className="text-right">
                  <div className="font-medium">${holding.value.toLocaleString()}</div>
                  <div className={holding.change.startsWith('+') ? 'text-green-400' : 'text-red-400'}>
                    {holding.change}
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Add Funds Modal */}
      <AddFundsModal
        isOpen={showAddFunds}
        onClose={() => setShowAddFunds(false)}
      />
    </div>
  );
};

export default Profile; 