import React from 'react';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  TimeScale,
} from 'chart.js';
import { Line } from 'react-chartjs-2';
import 'chartjs-adapter-date-fns';
import { formatCurrency } from '../../utils/formatters';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  TimeScale
);

interface PortfolioValue {
  timestamp: string;
  value: number;
}

interface PortfolioChartProps {
  data: PortfolioValue[];
  loading?: boolean;
  timeframe?: '24h' | '7d' | '30d' | '1y';
}

const PortfolioChart: React.FC<PortfolioChartProps> = ({
  data,
  loading = false,
  timeframe = '24h',
}) => {
  if (loading) {
    return <div className="h-[400px] flex items-center justify-center">Loading chart...</div>;
  }

  const chartData = {
    datasets: [
      {
        label: 'Portfolio Value',
        data: data.map((point) => ({
          x: new Date(point.timestamp),
          y: point.value,
        })),
        borderColor: 'rgb(99, 102, 241)',
        backgroundColor: 'rgba(99, 102, 241, 0.5)',
        tension: 0.1,
        fill: true,
      },
    ],
  };

  const options = {
    responsive: true,
    maintainAspectRatio: false,
    interaction: {
      mode: 'index' as const,
      intersect: false,
    },
    scales: {
      x: {
        type: 'time' as const,
        time: {
          unit: timeframe === '24h' ? 'hour' : timeframe === '7d' ? 'day' : 'week',
        },
        title: {
          display: true,
          text: 'Time',
        },
      },
      y: {
        title: {
          display: true,
          text: 'Value (USD)',
        },
        ticks: {
          callback: (value: number) => formatCurrency(value),
        },
        min: 0,
      },
    },
    plugins: {
      legend: {
        display: false,
      },
      tooltip: {
        callbacks: {
          label: (context: any) => {
            return `Value: ${formatCurrency(context.parsed.y)}`;
          },
        },
      },
    },
  };

  return (
    <div className="h-[400px]">
      <Line data={chartData} options={options} />
    </div>
  );
};

export default PortfolioChart; 