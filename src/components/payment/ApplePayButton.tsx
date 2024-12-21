import React, { useState } from 'react';
import { formatCurrency } from '../../utils/formatters';

interface ApplePayButtonProps {
  amount: number;
  onSuccess: (transactionId: string) => void;
  onError: (error: Error) => void;
}

const ApplePayButton: React.FC<ApplePayButtonProps> = ({
  amount,
  onSuccess,
  onError,
}) => {
  const [loading, setLoading] = useState(false);

  const handlePayment = async () => {
    if (!window.ApplePaySession) {
      onError(new Error('Apple Pay is not available'));
      return;
    }

    try {
      setLoading(true);

      const paymentRequest = {
        countryCode: 'US',
        currencyCode: 'USD',
        merchantCapabilities: ['supports3DS'],
        supportedNetworks: ['visa', 'masterCard', 'amex'],
        total: {
          label: 'Add Funds',
          amount: amount.toFixed(2),
          type: 'final',
        },
      };

      const session = new ApplePaySession(3, paymentRequest);

      session.onvalidatemerchant = async (event) => {
        try {
          const merchantSession = await validateMerchant(event.validationURL);
          session.completeMerchantValidation(merchantSession);
        } catch (error) {
          console.error('Merchant validation failed:', error);
          session.abort();
          onError(error as Error);
        }
      };

      session.onpaymentauthorized = async (event) => {
        try {
          const result = await processPayment(event.payment);
          session.completePayment(ApplePaySession.STATUS_SUCCESS);
          onSuccess(result.transactionId);
        } catch (error) {
          console.error('Payment processing failed:', error);
          session.completePayment(ApplePaySession.STATUS_FAILURE);
          onError(error as Error);
        }
      };

      session.oncancel = () => {
        setLoading(false);
      };

      session.begin();
    } catch (error) {
      setLoading(false);
      onError(error as Error);
    }
  };

  return (
    <button
      onClick={handlePayment}
      disabled={loading}
      className={`w-full flex items-center justify-center px-4 py-3 border border-transparent text-base font-medium rounded-md shadow-sm text-white bg-black hover:bg-gray-900 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-black ${
        loading ? 'opacity-50 cursor-not-allowed' : ''
      }`}
    >
      {loading ? (
        <span>Processing...</span>
      ) : (
        <>
          <span className="mr-2">Pay with</span>
          <svg className="h-6 w-auto" viewBox="0 0 24 24" fill="currentColor">
            <path d="M12.5 2.5c-.7 0-1.3.6-1.3 1.3v16.4c0 .7.6 1.3 1.3 1.3s1.3-.6 1.3-1.3V3.8c0-.7-.6-1.3-1.3-1.3zm-5 4c-.7 0-1.3.6-1.3 1.3v8.4c0 .7.6 1.3 1.3 1.3s1.3-.6 1.3-1.3V7.8c0-.7-.6-1.3-1.3-1.3zm10 0c-.7 0-1.3.6-1.3 1.3v8.4c0 .7.6 1.3 1.3 1.3s1.3-.6 1.3-1.3V7.8c0-.7-.6-1.3-1.3-1.3z" />
          </svg>
        </>
      )}
    </button>
  );
};

export default ApplePayButton; 