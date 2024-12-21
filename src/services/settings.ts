import api from './api';

export interface SecuritySettings {
  twoFactorEnabled: boolean;
  emailNotificationsEnabled: boolean;
  loginNotificationsEnabled: boolean;
  lastPasswordChange: string;
  activeSessions: {
    id: string;
    device: string;
    location: string;
    lastActive: string;
  }[];
}

export interface PrivacySettings {
  profileVisibility: 'public' | 'private' | 'friends';
  showPortfolio: boolean;
  showTradeHistory: boolean;
  showLeaderboardPosition: boolean;
}

export const getSecuritySettings = async (): Promise<SecuritySettings> => {
  const response = await api.get<SecuritySettings>('/settings/security');
  return response.data;
};

export const updateSecuritySettings = async (
  settings: Partial<SecuritySettings>
): Promise<SecuritySettings> => {
  const response = await api.put<SecuritySettings>('/settings/security', settings);
  return response.data;
};

export const getPrivacySettings = async (): Promise<PrivacySettings> => {
  const response = await api.get<PrivacySettings>('/settings/privacy');
  return response.data;
};

export const updatePrivacySettings = async (
  settings: Partial<PrivacySettings>
): Promise<PrivacySettings> => {
  const response = await api.put<PrivacySettings>('/settings/privacy', settings);
  return response.data;
}; 