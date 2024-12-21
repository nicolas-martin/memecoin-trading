export type OrderType = 'BUY' | 'SELL';
export type OrderStatus = 'PENDING' | 'COMPLETED' | 'FAILED';

export interface Order {
  id: string;
  userId: string;
  coinId: string;
  type: OrderType;
  amount: number;
  price: number;
  total: number;
  status: OrderStatus;
  createdAt: string;
  updatedAt: string;
}

export interface CreateOrderRequest {
  coinId: string;
  type: OrderType;
  amount: number;
  price: number;
} 