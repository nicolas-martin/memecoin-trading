import { useState, useEffect } from 'react';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Filler,
  ChartOptions,
  ScaleType,
} from 'chart.js';
import { Line } from 'react-chartjs-2';
import { fetchHistoricalPrices, PriceData } from '../../services/dexscreener';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Filler
);

interface PriceChartProps {
  pairAddress: string;
  isPositive: boolean;
}

const timeframes = ['1H', '24H', '1W', '1M', '1Y'] as const;
type Timeframe = typeof timeframes[number];

const PriceChart = ({ pairAddress, isPositive }: PriceChartProps) => {
  const [selectedTimeframe, setSelectedTimeframe] = useState<Timeframe>('24H');
  const [priceData, setPriceData] = useState<PriceData[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const color = isPositive ? '#34D399' : '#EF4444';

  useEffect(() => {
    const loadPriceData = async () => {
      try {
        setIsLoading(true);
        setError(null);
        const data = await fetchHistoricalPrices(pairAddress, selectedTimeframe);
        setPriceData(data);
      } catch (err) {
        setError('Failed to load price data');
        console.error(err);
      } finally {
        setIsLoading(false);
      }
    };

    loadPriceData();
  }, [pairAddress, selectedTimeframe]);

  const chartData = {
    labels: priceData.map(d => new Date(d.timestamp).toISOString()),
    datasets: [
      {
        data: priceData.map(d => d.price),
        borderColor: color,
        borderWidth: 2,
        tension: 0.4,
        pointRadius: 0,
        fill: true,
        backgroundColor: (context: any) => {
          const ctx = context.chart.ctx;
          const gradient = ctx.createLinearGradient(0, 0, 0, 300);
          gradient.addColorStop(0, `${color}20`);
          gradient.addColorStop(1, `${color}00`);
          return gradient;
        },
      },
    ],
  };

  const options: ChartOptions<'line'> = {
    responsive: true,
    maintainAspectRatio: false,
    interaction: {
      intersect: false,
      mode: 'index',
    },
    plugins: {
      tooltip: {
        enabled: true,
        backgroundColor: '#1C1C1E',
        titleColor: '#FFFFFF',
        bodyColor: '#FFFFFF',
        padding: 12,
        cornerRadius: 8,
        displayColors: false,
        callbacks: {
          label: function(context) {
            if (context.parsed.y !== null) {
              return `$${context.parsed.y.toLocaleString()}`;
            }
            return '';
          },
          title: function(tooltipItems) {
            if (tooltipItems.length > 0) {
              const item = tooltipItems[0];
              return new Date(item.label).toLocaleString();
            }
            return '';
          },
        },
      },
      legend: {
        display: false,
      },
    },
    scales: {
      x: {
        type: 'category',
        display: false,
      },
      y: {
        type: 'linear' as const,
        display: true,
        position: 'right' as const,
        grid: {
          display: false,
        },
        ticks: {
          color: '#6B7280',
          callback: function(value) {
            if (typeof value === 'number') {
              return `$${value.toLocaleString()}`;
            }
            return '';
          },
        },
      },
    },
  };

  return (
    <div>
      {/* Timeframe Selector */}
      <div className="flex justify-end space-x-2 mb-4">
        {timeframes.map((timeframe) => (
          <button
            key={timeframe}
            onClick={() => setSelectedTimeframe(timeframe)}
            disabled={isLoading}
            className={`px-3 py-1 rounded-lg text-sm font-medium transition-colors ${
              selectedTimeframe === timeframe
                ? 'bg-[#2C2C2E] text-white'
                : 'text-gray-400'
            } ${isLoading ? 'opacity-50' : ''}`}
          >
            {timeframe}
          </button>
        ))}
      </div>

      {/* Chart */}
      <div className="h-64">
        {isLoading ? (
          <div className="flex items-center justify-center h-full">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-white"></div>
          </div>
        ) : error ? (
          <div className="flex items-center justify-center h-full text-red-400">
            {error}
          </div>
        ) : (
          <Line data={chartData} options={options} />
        )}
      </div>
    </div>
  );
};

export default PriceChart; 