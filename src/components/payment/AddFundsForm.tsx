import React, { useState } from 'react';
import ApplePayButton from './ApplePayButton';
import { addFunds } from '../../services/payment';
import { formatCurrency } from '../../utils/formatters';

const PRESET_AMOUNTS = [50, 100, 250, 500, 1000];

interface AddFundsFormProps {
  onSuccess: () => void;
  onError: (error: Error) => void;
}

const AddFundsForm: React.FC<AddFundsFormProps> = ({ onSuccess, onError }) => {
  const [amount, setAmount] = useState<string>('');
  const [loading, setLoading] = useState(false);

  const handlePresetAmount = (value: number) => {
    setAmount(value.toString());
  };

  const handlePaymentSuccess = async (transactionId: string) => {
    try {
      setLoading(true);
      await addFunds(Number(amount), 'apple_pay', transactionId);
      setAmount('');
      onSuccess();
    } catch (error) {
      onError(error as Error);
    } finally {
      setLoading(false);
    }
  };

  const isValidAmount = Number(amount) >= 10 && Number(amount) <= 10000;

  return (
    <div className="space-y-6">
      <div>
        <label htmlFor="amount" className="block text-sm font-medium text-gray-700">
          Amount (USD)
        </label>
        <div className="mt-1">
          <input
            type="number"
            name="amount"
            id="amount"
            min="10"
            max="10000"
            value={amount}
            onChange={(e) => setAmount(e.target.value)}
            className="shadow-sm focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md"
            placeholder="Enter amount"
          />
        </div>
        {!isValidAmount && amount && (
          <p className="mt-2 text-sm text-red-600">
            Please enter an amount between $10 and $10,000
          </p>
        )}
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">
          Quick Select
        </label>
        <div className="grid grid-cols-3 gap-2">
          {PRESET_AMOUNTS.map((value) => (
            <button
              key={value}
              type="button"
              onClick={() => handlePresetAmount(value)}
              className={`px-3 py-2 text-sm font-medium rounded-md ${
                Number(amount) === value
                  ? 'bg-indigo-600 text-white'
                  : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
              }`}
            >
              {formatCurrency(value)}
            </button>
          ))}
        </div>
      </div>

      <div className="pt-4">
        <ApplePayButton
          amount={Number(amount)}
          onSuccess={handlePaymentSuccess}
          onError={onError}
        />
      </div>
    </div>
  );
};

export default AddFundsForm; 