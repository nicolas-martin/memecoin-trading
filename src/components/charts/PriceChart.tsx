import React, { useState, useRef } from 'react';
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
  ChartOptions,
} from 'chart.js';
import { Line, getElementAtEvent } from 'react-chartjs-2';
import zoomPlugin from 'chartjs-plugin-zoom';
import annotationPlugin from 'chartjs-plugin-annotation';
import 'chartjs-adapter-date-fns';
import { CoinPrice } from '../../types/coin';
import { formatCurrency, formatDate } from '../../utils/formatters';
import { calculateIndicators } from '../../utils/indicators';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  TimeScale,
  zoomPlugin,
  annotationPlugin
);

interface PriceChartProps {
  data: CoinPrice[];
  symbol: string;
  loading?: boolean;
  timeframe?: '24h' | '7d' | '30d' | '1y';
  activeIndicators?: ('MA' | 'EMA' | 'RSI' | 'MACD' | 'BB')[];
}

const PriceChart: React.FC<PriceChartProps> = ({
  data,
  symbol,
  loading = false,
  timeframe = '24h',
  activeIndicators = [],
}) => {
  const chartRef = useRef<ChartJS>(null);
  const [selectedPoint, setSelectedPoint] = useState<{
    price: number;
    timestamp: string;
  } | null>(null);

  if (loading) {
    return <div className="h-[400px] flex items-center justify-center">Loading chart...</div>;
  }

  const indicators = calculateIndicators(data, activeIndicators);

  const datasets = [
    {
      label: `${symbol} Price`,
      data: data.map((price) => ({
        x: new Date(price.timestamp),
        y: price.price,
      })),
      borderColor: 'rgb(75, 192, 192)',
      backgroundColor: 'rgba(75, 192, 192, 0.5)',
      tension: 0.1,
      order: 1,
    },
    ...indicators,
  ];

  const options: ChartOptions<'line'> = {
    responsive: true,
    maintainAspectRatio: false,
    interaction: {
      mode: 'index',
      intersect: false,
    },
    scales: {
      x: {
        type: 'time',
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
          text: 'Price (USD)',
        },
        ticks: {
          callback: (value) => formatCurrency(value as number),
        },
      },
    },
    plugins: {
      legend: {
        display: true,
        position: 'top',
      },
      tooltip: {
        callbacks: {
          label: (context) => {
            return `${context.dataset.label}: ${formatCurrency(context.parsed.y)}`;
          },
        },
      },
      zoom: {
        pan: {
          enabled: true,
          mode: 'x',
        },
        zoom: {
          wheel: {
            enabled: true,
          },
          pinch: {
            enabled: true,
          },
          mode: 'x',
        },
      },
      annotation: selectedPoint
        ? {
            annotations: {
              selectedPoint: {
                type: 'point',
                xValue: new Date(selectedPoint.timestamp),
                yValue: selectedPoint.price,
                backgroundColor: 'rgba(255, 99, 132, 0.25)',
                borderColor: 'rgb(255, 99, 132)',
                borderWidth: 1,
                radius: 4,
              },
            },
          }
        : undefined,
    },
  };

  const handleClick = (event: React.MouseEvent<HTMLCanvasElement>) => {
    if (!chartRef.current) return;

    const elements = getElementAtEvent(chartRef.current, event);
    if (elements.length > 0) {
      const { index } = elements[0];
      const point = data[index];
      setSelectedPoint({
        price: point.price,
        timestamp: point.timestamp,
      });
    }
  };

  const handleResetZoom = () => {
    if (chartRef.current) {
      chartRef.current.resetZoom();
    }
  };

  return (
    <div className="relative h-[400px]">
      <div className="absolute top-2 right-2 z-10">
        <button
          onClick={handleResetZoom}
          className="px-3 py-1 text-sm font-medium rounded-md bg-gray-100 text-gray-700 hover:bg-gray-200"
        >
          Reset Zoom
        </button>
      </div>
      <Line
        ref={chartRef}
        data={{ datasets }}
        options={options}
        onClick={handleClick}
      />
      {selectedPoint && (
        <div className="absolute bottom-2 left-2 bg-white p-2 rounded-md shadow-md">
          <div className="text-sm font-medium">
            {formatDate(selectedPoint.timestamp)}
          </div>
          <div className="text-sm">
            Price: {formatCurrency(selectedPoint.price)}
          </div>
        </div>
      )}
    </div>
  );
};

export default PriceChart; 
export default PriceChart; 