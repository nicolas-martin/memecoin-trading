import { Line } from 'react-chartjs-2';
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
} from 'chart.js';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Filler
);

const PortfolioChart = () => {
  // Mock data - in a real app, this would come from your API
  const mockData = {
    labels: Array.from({ length: 30 }, (_, i) => {
      const date = new Date();
      date.setDate(date.getDate() - (30 - i));
      return date.toLocaleDateString();
    }),
    values: Array.from({ length: 30 }, (_, i) => {
      const baseValue = 10000;
      const trend = Math.sin(i / 5) * 1000; // Create a wavy pattern
      const noise = (Math.random() - 0.5) * 500; // Add some randomness
      return baseValue + trend + noise;
    }),
  };

  const data = {
    labels: mockData.labels,
    datasets: [
      {
        data: mockData.values,
        borderColor: '#34D399', // Green color
        borderWidth: 2,
        tension: 0.4,
        pointRadius: 0,
        fill: true,
        backgroundColor: (context: any) => {
          const ctx = context.chart.ctx;
          const gradient = ctx.createLinearGradient(0, 0, 0, 300);
          gradient.addColorStop(0, 'rgba(52, 211, 153, 0.2)');
          gradient.addColorStop(1, 'rgba(52, 211, 153, 0)');
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
        },
      },
      legend: {
        display: false,
      },
    },
    scales: {
      x: {
        display: false,
      },
      y: {
        display: true,
        position: 'right',
        grid: {
          display: false,
        },
        border: {
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
    <div className="h-64">
      <Line data={data} options={options} />
    </div>
  );
};

export default PortfolioChart; 