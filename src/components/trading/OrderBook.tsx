import React from 'react';
import { Order } from '../../types/order';
import { formatCurrency, formatDate } from '../../utils/formatters';

interface OrderBookProps {
  orders: Order[];
  onCancelOrder?: (orderId: string) => void;
  loading?: boolean;
}

const OrderBook: React.FC<OrderBookProps> = ({
  orders,
  onCancelOrder,
  loading = false,
}) => {
  if (loading) {
    return <div>Loading orders...</div>;
  }

  return (
    <div className="overflow-x-auto">
      <table className="min-w-full divide-y divide-gray-200">
        <thead className="bg-gray-50">
          <tr>
            <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Type
            </th>
            <th scope="col" className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
              Amount
            </th>
            <th scope="col" className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
              Price
            </th>
            <th scope="col" className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
              Total
            </th>
            <th scope="col" className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
              Status
            </th>
            <th scope="col" className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
              Date
            </th>
            <th scope="col" className="relative px-6 py-3">
              <span className="sr-only">Actions</span>
            </th>
          </tr>
        </thead>
        <tbody className="bg-white divide-y divide-gray-200">
          {orders.map((order) => (
            <tr key={order.id}>
              <td className={`px-6 py-4 whitespace-nowrap text-sm font-medium ${
                order.type === 'BUY' ? 'text-green-600' : 'text-red-600'
              }`}>
                {order.type}
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-right text-sm text-gray-900">
                {order.amount}
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-right text-sm text-gray-900">
                {formatCurrency(order.price)}
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-right text-sm text-gray-900">
                {formatCurrency(order.total)}
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-right text-sm">
                <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                  order.status === 'COMPLETED'
                    ? 'bg-green-100 text-green-800'
                    : order.status === 'FAILED'
                    ? 'bg-red-100 text-red-800'
                    : 'bg-yellow-100 text-yellow-800'
                }`}>
                  {order.status}
                </span>
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-right text-sm text-gray-500">
                {formatDate(order.createdAt)}
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                {order.status === 'PENDING' && onCancelOrder && (
                  <button
                    onClick={() => onCancelOrder(order.id)}
                    className="text-indigo-600 hover:text-indigo-900"
                  >
                    Cancel
                  </button>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default OrderBook; 