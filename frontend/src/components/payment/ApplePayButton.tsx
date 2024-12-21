import { useEffect, useState } from 'react';

interface ApplePayButtonProps {
  amount: number;
  onSuccess: (paymentResult: any) => void;
  onError: (error: Error) => void;
}

const ApplePayButton = ({ amount, onSuccess, onError }: ApplePayButtonProps) => {
  const [canMakePayments, setCanMakePayments] = useState(false);

  useEffect(() => {
    if (window.ApplePaySession && ApplePaySession.canMakePayments()) {
      setCanMakePayments(true);
    }
  }, []);

  const handlePayment = async () => {
    try {
      const paymentRequest = {
        countryCode: 'US',
        currencyCode: 'USD',
        merchantCapabilities: ['supports3DS'],
        supportedNetworks: ['visa', 'masterCard', 'amex'],
        total: {
          label: 'MemeCoin Trading',
          amount: amount.toFixed(2),
          type: 'final'
        }
      };

      const session = new ApplePaySession(3, paymentRequest);

      session.onvalidatemerchant = async (event) => {
        try {
          const response = await fetch('/api/v1/payments/apple-pay/validate', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ validationURL: event.validationURL }),
          });

          if (!response.ok) throw new Error('Merchant validation failed');
          const merchantSession = await response.json();
          session.completeMerchantValidation(merchantSession);
        } catch (error) {
          console.error('Merchant validation failed:', error);
          session.abort();
          onError(error as Error);
        }
      };

      session.onpaymentauthorized = async (event) => {
        try {
          const response = await fetch('/api/v1/payments/apple-pay/process', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
              payment: event.payment,
              amount: amount,
            }),
          });

          if (!response.ok) throw new Error('Payment processing failed');
          const result = await response.json();
          
          session.completePayment(ApplePaySession.STATUS_SUCCESS);
          onSuccess(result);
        } catch (error) {
          console.error('Payment processing failed:', error);
          session.completePayment(ApplePaySession.STATUS_FAILURE);
          onError(error as Error);
        }
      };

      session.oncancel = () => {
        onError(new Error('Payment cancelled'));
      };

      session.begin();
    } catch (error) {
      console.error('Apple Pay session failed:', error);
      onError(error as Error);
    }
  };

  if (!canMakePayments) return null;

  return (
    <button
      onClick={handlePayment}
      className="w-full h-12 bg-black rounded-lg flex items-center justify-center active:opacity-90"
      style={{
        WebkitAppearance: '-apple-pay-button',
        appearance: '-apple-pay-button',
      }}
    >
      <span style={{ display: 'none' }}>Apple Pay</span>
    </button>
  );
};

export default ApplePayButton; 