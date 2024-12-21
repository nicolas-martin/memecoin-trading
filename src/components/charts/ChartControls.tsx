import React from 'react';

interface ChartControlsProps {
  timeframe: '24h' | '7d' | '30d' | '1y';
  onTimeframeChange: (timeframe: '24h' | '7d' | '30d' | '1y') => void;
}

const ChartControls: React.FC<ChartControlsProps> = ({
  timeframe,
  onTimeframeChange,
}) => {
  const timeframes = [
    { value: '24h', label: '24H' },
    { value: '7d', label: '7D' },
    { value: '30d', label: '30D' },
    { value: '1y', label: '1Y' },
  ] as const;

  return (
    <div className="flex space-x-2">
      {timeframes.map(({ value, label }) => (
        <button
          key={value}
          onClick={() => onTimeframeChange(value)}
          className={`px-3 py-1 text-sm font-medium rounded-md ${
            timeframe === value
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

export default ChartControls; 