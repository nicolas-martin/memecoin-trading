import React from 'react';
import { LeaderboardEntry } from '../../services/leaderboard';
import { formatCurrency, formatPercentage } from '../../utils/formatters';

interface LeaderboardTableProps {
  entries: LeaderboardEntry[];
  loading?: boolean;
}

const LeaderboardTable: React.FC<LeaderboardTableProps> = ({ entries, loading }) => {
  if (loading) {
    return <div>Loading leaderboard...</div>;
  }

  return (
    <div className="overflow-x-auto">
      <table className="min-w-full divide-y divide-gray-200">
        <thead className="bg-gray-50">
          <tr>
            <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Rank
            </th>
            <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Trader
            </th>
            <th scope="col" className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
              Profit
            </th>
            <th scope="col" className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
              Profit %
            </th>
          </tr>
        </thead>
        <tbody className="bg-white divide-y divide-gray-200">
          {entries.map((entry) => (
            <tr key={entry.userId}>
              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                {entry.rank}
              </td>
              <td className="px-6 py-4 whitespace-nowrap">
                <div className="flex items-center">
                  {entry.avatar && (
                    <img
                      className="h-8 w-8 rounded-full mr-3"
                      src={entry.avatar}
                      alt=""
                    />
                  )}
                  <div className="text-sm font-medium text-gray-900">
                    {entry.username}
                  </div>
                </div>
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-right text-sm text-gray-900">
                {formatCurrency(entry.profit)}
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-right text-sm">
                <span className={entry.profitPercentage >= 0 ? 'text-green-600' : 'text-red-600'}>
                  {formatPercentage(entry.profitPercentage)}
                </span>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default LeaderboardTable; 