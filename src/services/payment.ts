import api from './api';

interface PaymentResult {
  transactionId: string;
  amount: number;
  status: 'completed' | 'failed';
}

export const validateMerchant = async (validationURL: string): Promise<any> => {
  const response = await api.post('/payments/apple-pay/validate', {
    validationURL,
  });
  return response.data;
};

export const processPayment = async (paymentData: any): Promise<PaymentResult> => {
  const response = await api.post<PaymentResult>('/payments/apple-pay/process', {
    paymentData,
  });
  return response.data;
};

export const addFunds = async (amount: number, paymentMethod: string, transactionId: string): Promise<void> => {
  await api.post('/user/funds/add', {
    amount,
    paymentMethod,
    transactionId,
  });
}; 