import React, { useState, useEffect } from 'react';
import { useDispatch } from 'react-redux';
import { Coin } from '../../types/coin';
import { OrderType, CreateOrderRequest } from '../../types/order';
import { createOrder } from '../../services/trading';
import { formatCurrency } from '../../utils/formatters';
import { validateOrderAmount, validateOrderTotal } from '../../utils/validators';
import ConfirmationModal from '../modals/ConfirmationModal';

interface OrderFormProps {
  coin: Coin;
  type: OrderType;
  balance: number;
  onSuccess?: () => void;
  onError?: (error: Error) => void;
}

const OrderForm: React.FC<OrderFormProps> = ({
  coin,
  type,
  balance,
  onSuccess,
  onError,
}) => {
  const [amount, setAmount] = useState<string>('');
  const [total, setTotal] = useState<number>(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [showConfirmation, setShowConfirmation] = useState(false);
  const [validationErrors, setValidationErrors] = useState<{
    amount: string | null;
    total: string | null;
  }>({ amount: null, total: null });

  useEffect(() => {
    const calculatedTotal = Number(amount) * coin.price;
    setTotal(isNaN(calculatedTotal) ? 0 : calculatedTotal);
  }, [amount, coin.price]);

  const validateForm = (): boolean => {
    const amountError = validateOrderAmount(amount);
    const totalError = validateOrderTotal(total, balance);

    setValidationErrors({
      amount: amountError,
      total: totalError,
    });

    return !amountError && !totalError;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    if (!validateForm()) {
      return;
    }

    setShowConfirmation(true);
  };

  const handleConfirmOrder = async () => {
    setShowConfirmation(false);
    setLoading(true);

    try {
      const orderRequest: CreateOrderRequest = {
        coinId: coin.id,
        type,
        amount: Number(amount),
        price: coin.price,
      };

      await createOrder(orderRequest);
      setAmount('');
      onSuccess?.();
    } catch (err) {
      const error = err as Error;
      setError(error.message);
      onError?.(error);
    } finally {
      setLoading(false);
    }
  };

  const getConfirmationMessage = () => {
    return `Are you sure you want to ${type.toLowerCase()} ${amount} ${coin.symbol} for ${formatCurrency(total)}?`;
  };

  return (
    <>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-700">
            Amount ({coin.symbol})
          </label>
          <div className="mt-1 relative rounded-md shadow-sm">
            <input
              type="number"
              value={amount}
              onChange={(e) => setAmount(e.target.value)}
              min="0"
              step="0.000001"
              className={`focus:ring-indigo-500 focus:border-indigo-500 block w-full pl-4 pr-12 sm:text-sm border-gray-300 rounded-md ${
                validationErrors.amount ? 'border-red-300' : ''
              }`}
              placeholder="0.00"
              disabled={loading}
            />
          </div>
          {validationErrors.amount && (
            <p className="mt-1 text-sm text-red-600">{validationErrors.amount}</p>
          )}
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700">
            Price per {coin.symbol}
          </label>
          <div className="mt-1 relative rounded-md shadow-sm">
            <input
              type="text"
              value={formatCurrency(coin.price)}
              disabled
              className="block w-full pl-4 pr-12 sm:text-sm border-gray-300 rounded-md bg-gray-50"
            />
          </div>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700">
            Total (USD)
          </label>
          <div className="mt-1 relative rounded-md shadow-sm">
            <input
              type="text"
              value={formatCurrency(total)}
              disabled
              className="block w-full pl-4 pr-12 sm:text-sm border-gray-300 rounded-md bg-gray-50"
            />
          </div>
        </div>

        {error && (
          <div className="text-red-600 text-sm mt-2">{error}</div>
        )}

        <button
          type="submit"
          disabled={loading || !amount || Number(amount) <= 0}
          className={`w-full py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white ${
            type === 'BUY'
              ? 'bg-green-600 hover:bg-green-700'
              : 'bg-red-600 hover:bg-red-700'
          } focus:outline-none focus:ring-2 focus:ring-offset-2 ${
            type === 'BUY' ? 'focus:ring-green-500' : 'focus:ring-red-500'
          } disabled:opacity-50 disabled:cursor-not-allowed`}
        >
          {loading ? 'Processing...' : `${type} ${coin.symbol}`}
        </button>
      </form>

      <ConfirmationModal
        isOpen={showConfirmation}
        onClose={() => setShowConfirmation(false)}
        onConfirm={handleConfirmOrder}
        title={`Confirm ${type}`}
        message={getConfirmationMessage()}
        confirmText={type}
        type={type === 'BUY' ? 'info' : 'warning'}
      />
    </>
  );
};

export default OrderForm; 