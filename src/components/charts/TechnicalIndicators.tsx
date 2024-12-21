import React from 'react';

type IndicatorType = 'MA' | 'EMA' | 'RSI' | 'MACD' | 'BB';

interface TechnicalIndicatorsProps {
  activeIndicators: IndicatorType[];
  onToggleIndicator: (indicator: IndicatorType) => void;
}

const TechnicalIndicators: React.FC<TechnicalIndicatorsProps> = ({
  activeIndicators,
  onToggleIndicator,
}) => {
  const indicators = [
    { type: 'MA' as const, label: 'Moving Average' },
    { type: 'EMA' as const, label: 'Exponential MA' },
    { type: 'RSI' as const, label: 'RSI' },
    { type: 'MACD' as const, label: 'MACD' },
    { type: 'BB' as const, label: 'Bollinger Bands' },
  ];

  return (
    <div className="flex flex-wrap gap-2">
      {indicators.map(({ type, label }) => (
        <button
          key={type}
          onClick={() => onToggleIndicator(type)}
          className={`px-3 py-1 text-sm font-medium rounded-md ${
            activeIndicators.includes(type)
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

export default TechnicalIndicators; 