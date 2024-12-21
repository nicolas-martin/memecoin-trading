export const validateOrderAmount = (amount: string, minAmount: number = 0): string | null => {
  const numAmount = Number(amount);
  
  if (!amount) {
    return 'Amount is required';
  }
  
  if (isNaN(numAmount)) {
    return 'Amount must be a valid number';
  }
  
  if (numAmount <= minAmount) {
    return `Amount must be greater than ${minAmount}`;
  }
  
  return null;
};

export const validateOrderTotal = (total: number, balance: number): string | null => {
  if (total <= 0) {
    return 'Total must be greater than 0';
  }
  
  if (total > balance) {
    return 'Insufficient balance';
  }
  
  return null;
}; 