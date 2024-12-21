import React from 'react';
import { LeaderboardEntry } from '../../services/leaderboard';
import { formatCurrency, formatPercentage } from '../../utils/formatters';

interface LeaderboardStatsProps {
  entries: LeaderboardEntry[];
}

const LeaderboardStats: React.FC<LeaderboardStatsProps> = ({ entries }) => {
  const totalProfit = entries.reduce((sum, entry) => sum + entry.profit, 0);
  const averageProfit = totalProfit / entries.length;
  const topTrader = entries[0];

  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
      <div className="bg-white rounded-lg shadow p-4">
        <h3 className="text-sm font-medium text-gray-500">Total Profit</h3>
        <p className="mt-1 text-2xl font-semibold text-gray-900">
          {formatCurrency(totalProfit)}
        </p>
      </div>

      <div className="bg-white rounded-lg shadow p-4">
        <h3 className="text-sm font-medium text-gray-500">Average Profit</h3>
        <p className="mt-1 text-2xl font-semibold text-gray-900">
          {formatCurrency(averageProfit)}
        </p>
      </div>

      <div className="bg-white rounded-lg shadow p-4">
        <h3 className="text-sm font-medium text-gray-500">Top Trader</h3>
        <div className="mt-1">
          <p className="text-lg font-semibold text-gray-900">
            {topTrader?.username}
          </p>
          <p className="text-sm text-green-600">
            {formatCurrency(topTrader?.profit || 0)} ({formatPercentage(topTrader?.profitPercentage || 0)})
          </p>
        </div>
      </div>
    </div>
  );
};

export default LeaderboardStats; 