import React from 'react';
import { useNavigate } from 'react-router-dom';
import { Coin } from '../../types/coin';
import { formatCurrency, formatPercentage } from '../../utils/formatters';

interface CoinTableProps {
  coins: Coin[];
  loading: boolean;
}

const CoinTable: React.FC<CoinTableProps> = ({ coins, loading }) => {
  const navigate = useNavigate();

  if (loading) {
    return <div>Loading...</div>;
  }

  return (
    <div className="overflow-x-auto">
      <table className="min-w-full divide-y divide-gray-200">
        <thead className="bg-gray-50">
          <tr>
            <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Name
            </th>
            <th scope="col" className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
              Price
            </th>
            <th scope="col" className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
              24h Change
            </th>
            <th scope="col" className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
              Market Cap
            </th>
            <th scope="col" className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
              Volume (24h)
            </th>
          </tr>
        </thead>
        <tbody className="bg-white divide-y divide-gray-200">
          {coins.map((coin) => (
            <tr
              key={coin.id}
              onClick={() => navigate(`/coins/${coin.id}`)}
              className="hover:bg-gray-50 cursor-pointer"
            >
              <td className="px-6 py-4 whitespace-nowrap">
                <div className="flex items-center">
                  <img className="h-8 w-8 rounded-full" src={coin.logoUrl} alt={coin.name} />
                  <div className="ml-4">
                    <div className="text-sm font-medium text-gray-900">{coin.name}</div>
                    <div className="text-sm text-gray-500">{coin.symbol}</div>
                  </div>
                </div>
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-right text-sm text-gray-900">
                {formatCurrency(coin.price)}
              </td>
              <td className={`px-6 py-4 whitespace-nowrap text-right text-sm ${
                coin.priceChange24h >= 0 ? 'text-green-600' : 'text-red-600'
              }`}>
                {formatPercentage(coin.priceChange24h)}
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-right text-sm text-gray-900">
                {formatCurrency(coin.marketCap)}
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-right text-sm text-gray-900">
                {formatCurrency(coin.volume24h)}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default CoinTable; 