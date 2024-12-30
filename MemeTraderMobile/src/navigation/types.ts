import { NativeStackNavigationProp } from '@react-navigation/native-stack';

export type RootStackParamList = {
  Main: undefined;
  CoinDetail: { coinId: string };
};

export type MainScreenNavigationProp = NativeStackNavigationProp<RootStackParamList, 'Main'>;
export type CoinDetailScreenNavigationProp = NativeStackNavigationProp<RootStackParamList, 'CoinDetail'>; 