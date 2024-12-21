import api from './api';

export interface Notification {
  id: string;
  type: 'info' | 'success' | 'warning' | 'error';
  title: string;
  message: string;
  read: boolean;
  createdAt: string;
}

export const getNotifications = async (): Promise<Notification[]> => {
  const response = await api.get<Notification[]>('/notifications');
  return response.data;
};

export const markNotificationAsRead = async (id: string): Promise<void> => {
  await api.put(`/notifications/${id}/read`);
};

export const deleteNotification = async (id: string): Promise<void> => {
  await api.delete(`/notifications/${id}`);
}; 