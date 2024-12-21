import React from 'react';
import { render, screen } from '@testing-library/react';
import LeaderboardTable from './LeaderboardTable';
import { LeaderboardEntry } from '../../types/leaderboard';

describe('LeaderboardTable', () => {
  const mockEntries: LeaderboardEntry[] = [
    {
      id: '1',
      userId: '1',
      username: 'trader1',
      profit: 1000,
      profitPercentage: 10,
      rank: 1,
      timeFrame: '24h',
      createdAt: new Date().toISOString(),
    },
  ];

  it('renders loading state', () => {
    render(<LeaderboardTable entries={[]} loading={true} />);
    expect(screen.getByText('Loading leaderboard...')).toBeInTheDocument();
  });

  it('renders entries correctly', () => {
    render(<LeaderboardTable entries={mockEntries} loading={false} />);
    expect(screen.getByText('trader1')).toBeInTheDocument();
    expect(screen.getByText('1')).toBeInTheDocument();
    expect(screen.getByText('$1,000.00')).toBeInTheDocument();
    expect(screen.getByText('10.00%')).toBeInTheDocument();
  });

  it('renders empty state when no entries', () => {
    render(<LeaderboardTable entries={[]} loading={false} />);
    expect(screen.queryByRole('row')).toBeNull();
  });
}); 