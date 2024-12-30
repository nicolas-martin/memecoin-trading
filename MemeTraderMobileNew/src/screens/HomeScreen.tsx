import React, { useEffect, useState } from 'react';
import { View, StyleSheet, ActivityIndicator } from 'react-native';
import { CoinTable } from '../components/CoinList/CoinTable';
import { MemeCoin, getTopMemeCoins } from '../services/api';
import { useNavigation } from '@react-navigation/native';

export const HomeScreen: React.FC = () => {
  const [coins, setCoins] = useState<MemeCoin[]>([]);
  const [loading, setLoading] = useState(true);
  const navigation = useNavigation();

  useEffect(() => {
    fetchMemeCoins();
  }, []);

  const fetchMemeCoins = async () => {
    try {
      const data = await getTopMemeCoins();
      setCoins(data);
    } catch (error) {
      console.error('Error fetching meme coins:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleCoinPress = (coin: MemeCoin) => {
    navigation.navigate('CoinDetail', { coin });
  };

  if (loading) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" color="#0891b2" />
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <CoinTable coins={coins} onCoinPress={handleCoinPress} />
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#FFFFFF',
  },
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: '#FFFFFF',
  },
}); 