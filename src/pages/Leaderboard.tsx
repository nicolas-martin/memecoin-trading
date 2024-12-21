import React, { useState, useEffect } from 'react';
import { TimeFrame, getLeaderboard, LeaderboardEntry } from '../services/leaderboard';
import LeaderboardTable from '../components/leaderboard/LeaderboardTable';
import TimeframeSelector from '../components/leaderboard/TimeframeSelector';

const Leaderboard: React.FC = () => {
  const [timeframe, setTimeframe] = useState<TimeFrame>('24h');
  const [entries, setEntries] = useState<LeaderboardEntry[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchLeaderboard = async () => {
      try {
        setLoading(true);
        const data = await getLeaderboard(timeframe);
        setEntries(data);
        setError(null);
      } catch (err) {
        setError('Failed to load leaderboard data');
        console.error('Leaderboard error:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchLeaderboard();
  }, [timeframe]);

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="bg-white rounded-lg shadow-lg p-6">
        <div className="flex items-center justify-between mb-6">
          <h1 className="text-2xl font-bold">Top Traders</h1>
          <TimeframeSelector
            selected={timeframe}
            onChange={setTimeframe}
          />
        </div>

        {error ? (
          <div className="text-red-600 text-center py-4">{error}</div>
        ) : (
          <LeaderboardTable entries={entries} loading={loading} />
        )}
      </div>
    </div>
  );
};

export default Leaderboard; 