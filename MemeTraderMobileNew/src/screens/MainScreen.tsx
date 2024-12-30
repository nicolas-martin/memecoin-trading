import React, { useEffect } from 'react';
import {
  View,
  ActivityIndicator,
  RefreshControl,
  StyleSheet,
  ScrollView,
} from 'react-native';
import { useDispatch, useSelector } from 'react-redux';
import { useNavigation } from '@react-navigation/native';
import { RootState, AppDispatch } from '../store';
import { fetchMemeCoins, updateMemeCoinsList } from '../store/memeCoinsSlice';
import { CoinTable } from '../components/CoinList/CoinTable';
import { MemeCoin } from '../services/api';
import { NativeStackNavigationProp } from '@react-navigation/native-stack';

type RootStackParamList = {
  Main: undefined;
  CoinDetail: { coinId: string };
};

type NavigationProp = NativeStackNavigationProp<RootStackParamList, 'Main'>;

export const MainScreen = () => {
  const dispatch = useDispatch<AppDispatch>();
  const navigation = useNavigation<NavigationProp>();
  const { coins, loading, error } = useSelector((state: RootState) => state.memeCoins);

  useEffect(() => {
    const initializeData = async () => {
      try {
        console.log('Starting data initialization...');
        const updateResult = await dispatch(updateMemeCoinsList()).unwrap();
        console.log('Update result:', updateResult);
        const fetchResult = await dispatch(fetchMemeCoins()).unwrap();
        console.log('Fetch result:', fetchResult);
      } catch (err) {
        console.error('Error initializing data:', err);
      }
    };
    
    initializeData();
  }, [dispatch]);

  const handleRefresh = async () => {
    try {
      console.log('Starting refresh...');
      await dispatch(updateMemeCoinsList()).unwrap();
    } catch (err) {
      console.error('Error refreshing:', err);
    }
  };

  const handleCoinPress = (coin: MemeCoin) => {
    navigation.navigate('CoinDetail', { coinId: coin.id });
  };

  if (loading && coins.length === 0) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" color="#4F46E5" />
      </View>
    );
  }

  if (error) {
    console.error('Error state:', error);
  }

  return (
    <ScrollView
      style={styles.container}
      refreshControl={
        <RefreshControl refreshing={loading} onRefresh={handleRefresh} />
      }
    >
      <CoinTable coins={coins} onCoinPress={handleCoinPress} />
    </ScrollView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#F3F4F6',
  },
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: '#F3F4F6',
  },
}); 