import api from './api';
import { Order, CreateOrderRequest } from '../types/order';

export const createOrder = async (order: CreateOrderRequest): Promise<Order> => {
  const response = await api.post<Order>('/orders', order);
  return response.data;
};

export const getOrders = async (
  limit: number = 10,
  offset: number = 0
): Promise<{ orders: Order[]; total: number }> => {
  const response = await api.get<{ orders: Order[]; total: number }>(
    `/orders?limit=${limit}&offset=${offset}`
  );
  return response.data;
};

export const cancelOrder = async (orderId: string): Promise<void> => {
  await api.delete(`/orders/${orderId}`);
}; 