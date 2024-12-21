import { useState } from 'react';
import { XMarkIcon } from '@heroicons/react/24/outline';
import { toast } from 'react-hot-toast';
import ApplePayButton from '../payment/ApplePayButton';

interface AddFundsModalProps {
  isOpen: boolean;
  onClose: () => void;
}

const AMOUNTS = [100, 500, 1000, 5000];

const AddFundsModal = ({ isOpen, onClose }: AddFundsModalProps) => {
  const [amount, setAmount] = useState('');
  const [selectedAmount, setSelectedAmount] = useState<number | null>(null);
  const [isProcessing, setIsProcessing] = useState(false);

  if (!isOpen) return null;

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    // Handle payment processing here
    console.log('Processing payment for:', amount || selectedAmount);
    onClose();
  };

  const handlePaymentSuccess = async (paymentResult: any) => {
    try {
      // Update user's balance in your backend
      const response = await fetch('/api/v1/payments/funds/add', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          amount: amount || selectedAmount,
          paymentResult,
        }),
      });
      
      if (!response.ok) throw new Error('Failed to add funds');
      
      toast.success('Funds added successfully!');
      onClose();
    } catch (error) {
      console.error('Failed to add funds:', error);
      toast.error('Failed to add funds. Please try again.');
    } finally {
      setIsProcessing(false);
    }
  };

  const handlePaymentError = (error: Error) => {
    console.error('Payment failed:', error);
    toast.error('Payment failed. Please try again.');
    setIsProcessing(false);
  };

  return (
    <div className="fixed inset-0 z-50 bg-black bg-opacity-75">
      <div className="min-h-screen px-4 text-center">
        <div className="fixed inset-0" onClick={onClose} />

        <div className="inline-block w-full max-w-md p-6 my-8 text-left align-middle transition-all transform bg-[#1C1C1E] rounded-2xl shadow-xl">
          <div className="flex justify-between items-center mb-6">
            <h3 className="text-xl font-semibold text-white">Add Funds</h3>
            <button onClick={onClose} className="text-gray-400 hover:text-white">
              <XMarkIcon className="w-6 h-6" />
            </button>
          </div>

          <form onSubmit={handleSubmit}>
            {/* Quick amount selection */}
            <div className="grid grid-cols-2 gap-3 mb-6">
              {AMOUNTS.map((amt) => (
                <button
                  key={amt}
                  type="button"
                  onClick={() => {
                    setSelectedAmount(amt);
                    setAmount('');
                  }}
                  className={`py-3 rounded-lg font-medium transition-colors ${
                    selectedAmount === amt
                      ? 'bg-blue-500 text-white'
                      : 'bg-[#2C2C2E] text-gray-300'
                  }`}
                >
                  ${amt}
                </button>
              ))}
            </div>

            {/* Custom amount input */}
            <div className="mb-6">
              <label className="block text-gray-400 mb-2">Custom Amount</label>
              <div className="relative">
                <input
                  type="number"
                  value={amount}
                  onChange={(e) => {
                    setAmount(e.target.value);
                    setSelectedAmount(null);
                  }}
                  className="w-full bg-[#2C2C2E] text-white px-4 py-3 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="0.00"
                  min="0"
                />
                <span className="absolute right-4 top-1/2 transform -translate-y-1/2 text-gray-400">
                  USD
                </span>
              </div>
            </div>

            {/* Payment methods */}
            <div className="space-y-3">
              <ApplePayButton
                amount={Number(amount) || selectedAmount || 0}
                onSuccess={handlePaymentSuccess}
                onError={handlePaymentError}
              />

              <button
                type="submit"
                disabled={isProcessing}
                className="w-full py-3 bg-[#2C2C2E] rounded-lg flex items-center justify-center space-x-2 active:opacity-90"
              >
                <span className="font-medium">Pay with Card</span>
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default AddFundsModal; 