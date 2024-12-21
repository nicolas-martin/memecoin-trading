import React from 'react';
import { TimeFrame } from '../../services/leaderboard';

interface TimeframeSelectorProps {
  selected: TimeFrame;
  onChange: (timeframe: TimeFrame) => void;
}

const TimeframeSelector: React.FC<TimeframeSelectorProps> = ({
  selected,
  onChange,
}) => {
  const timeframes: { value: TimeFrame; label: string }[] = [
    { value: '24h', label: 'Last 24 Hours' },
    { value: '7d', label: 'Last 7 Days' },
    { value: '30d', label: 'Last 30 Days' },
  ];

  return (
    <div className="flex space-x-2">
      {timeframes.map(({ value, label }) => (
        <button
          key={value}
          onClick={() => onChange(value)}
          className={`px-4 py-2 rounded-md text-sm font-medium transition-colors ${
            selected === value
              ? 'bg-indigo-600 text-white'
              : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
          }`}
        >
          {label}
        </button>
      ))}
    </div>
  );
};

export default TimeframeSelector; 