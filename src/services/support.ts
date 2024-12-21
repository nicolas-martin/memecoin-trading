import api from './api';

export interface SupportTicket {
  id: string;
  subject: string;
  description: string;
  status: 'open' | 'in_progress' | 'resolved' | 'closed';
  priority: 'low' | 'medium' | 'high';
  createdAt: string;
  updatedAt: string;
  messages: TicketMessage[];
}

export interface TicketMessage {
  id: string;
  ticketId: string;
  content: string;
  sender: 'user' | 'support';
  createdAt: string;
}

export const createTicket = async (data: {
  subject: string;
  description: string;
  priority: SupportTicket['priority'];
}): Promise<SupportTicket> => {
  const response = await api.post<SupportTicket>('/support/tickets', data);
  return response.data;
};

export const getTickets = async (): Promise<SupportTicket[]> => {
  const response = await api.get<SupportTicket[]>('/support/tickets');
  return response.data;
};

export const addTicketMessage = async (
  ticketId: string,
  content: string
): Promise<TicketMessage> => {
  const response = await api.post<TicketMessage>(`/support/tickets/${ticketId}/messages`, {
    content,
  });
  return response.data;
}; 