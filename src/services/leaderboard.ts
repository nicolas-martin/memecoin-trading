import api from './api';

export interface LeaderboardEntry {
  userId: string;
  username: string;
  profit: number;
  rank: number;
  profitPercentage: number;
  avatar?: string;
}

export type TimeFrame = '24h' | '7d' | '30d';

export const getLeaderboard = async (timeframe: TimeFrame): Promise<LeaderboardEntry[]> => {
  const response = await api.get<LeaderboardEntry[]>(`/leaderboard?timeframe=${timeframe}`);
  return response.data;
}; 