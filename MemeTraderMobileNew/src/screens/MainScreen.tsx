import React, { useEffect } from 'react';
import {
  View,
  FlatList,
  ActivityIndicator,
  RefreshControl,
  StyleSheet,
} from 'react-native';
import { useDispatch, useSelector } from 'react-redux';
import { useNavigation } from '@react-navigation/native';
import { RootState, AppDispatch } from '../store';
import { fetchMemeCoins, updateMemeCoinsList } from '../store/memeCoinsSlice';
import { CoinListItem } from '../components/CoinList/CoinListItem';
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
  const { coins, loading } = useSelector((state: RootState) => state.memeCoins);

  useEffect(() => {
    dispatch(fetchMemeCoins());
  }, [dispatch]);

  const handleRefresh = () => {
    dispatch(updateMemeCoinsList());
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

  return (
    <View style={styles.container}>
      <FlatList
        data={coins}
        renderItem={({ item }) => (
          <CoinListItem coin={item} onPress={handleCoinPress} />
        )}
        keyExtractor={(item) => item.id}
        refreshControl={
          <RefreshControl refreshing={loading} onRefresh={handleRefresh} />
        }
      />
    </View>
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