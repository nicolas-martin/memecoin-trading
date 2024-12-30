import React, { useEffect } from 'react';
import {
  View,
  Text,
  ScrollView,
  ActivityIndicator,
  StyleSheet,
  TouchableOpacity,
} from 'react-native';
import { useDispatch, useSelector } from 'react-redux';
import { useRoute, RouteProp } from '@react-navigation/native';
import { RootState, AppDispatch } from '../store';
import { fetchMemeCoinDetail } from '../store/memeCoinsSlice';
import { LineChart } from 'react-native-chart-kit';
import { Dimensions } from 'react-native';

type RootStackParamList = {
  Main: undefined;
  CoinDetail: { coinId: string };
};

type CoinDetailRouteProp = RouteProp<RootStackParamList, 'CoinDetail'>;

export const CoinDetailScreen = () => {
  const dispatch = useDispatch<AppDispatch>();
  const route = useRoute<CoinDetailRouteProp>();
  const { coinId } = route.params;
  const { selectedCoin, loading } = useSelector(
    (state: RootState) => state.memeCoins
  );

  useEffect(() => {
    dispatch(fetchMemeCoinDetail(coinId));
  }, [dispatch, coinId]);

  if (loading || !selectedCoin) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" color="#4F46E5" />
      </View>
    );
  }

  const formatPrice = (price: number) => {
    return price < 0.01 ? price.toExponential(2) : price.toFixed(2);
  };

  const formatMarketCap = (marketCap: number) => {
    if (marketCap >= 1e9) {
      return `$${(marketCap / 1e9).toFixed(2)}B`;
    } else if (marketCap >= 1e6) {
      return `$${(marketCap / 1e6).toFixed(2)}M`;
    } else if (marketCap >= 1e3) {
      return `$${(marketCap / 1e3).toFixed(2)}K`;
    }
    return `$${marketCap.toFixed(2)}`;
  };

  return (
    <ScrollView style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.name}>{selectedCoin.name}</Text>
        <Text style={styles.symbol}>{selectedCoin.symbol}</Text>
      </View>

      <View style={styles.priceContainer}>
        <Text style={styles.price}>${formatPrice(selectedCoin.price)}</Text>
        <Text
          style={[
            styles.priceChange,
            selectedCoin.priceChangePercentage24h >= 0
              ? styles.positiveChange
              : styles.negativeChange,
          ]}
        >
          {selectedCoin.priceChangePercentage24h >= 0 ? '+' : ''}
          {selectedCoin.priceChangePercentage24h.toFixed(2)}%
        </Text>
      </View>

      <View style={styles.statsContainer}>
        <View style={styles.statItem}>
          <Text style={styles.statLabel}>Market Cap</Text>
          <Text style={styles.statValue}>
            {formatMarketCap(selectedCoin.marketCap)}
          </Text>
        </View>
        <View style={styles.statItem}>
          <Text style={styles.statLabel}>24h Volume</Text>
          <Text style={styles.statValue}>
            ${selectedCoin.volume24h.toFixed(2)}
          </Text>
        </View>
      </View>

      <View style={styles.chartContainer}>
        <LineChart
          data={{
            labels: ['1h', '2h', '3h', '4h', '5h', '6h'],
            datasets: [
              {
                data: [
                  selectedCoin.price * 0.95,
                  selectedCoin.price * 0.98,
                  selectedCoin.price * 1.02,
                  selectedCoin.price * 0.97,
                  selectedCoin.price * 1.01,
                  selectedCoin.price,
                ],
              },
            ],
          }}
          width={Dimensions.get('window').width - 32}
          height={220}
          chartConfig={{
            backgroundColor: '#FFFFFF',
            backgroundGradientFrom: '#FFFFFF',
            backgroundGradientTo: '#FFFFFF',
            decimalPlaces: 2,
            color: (opacity = 1) => `rgba(79, 70, 229, ${opacity})`,
            style: {
              borderRadius: 16,
            },
          }}
          bezier
          style={styles.chart}
        />
      </View>

      <View style={styles.buttonsContainer}>
        <TouchableOpacity
          style={[styles.button, styles.buyButton]}
          onPress={() => {
            // Implement buy functionality
          }}
        >
          <Text style={styles.buttonText}>Buy</Text>
        </TouchableOpacity>
        <TouchableOpacity
          style={[styles.button, styles.sellButton]}
          onPress={() => {
            // Implement sell functionality
          }}
        >
          <Text style={styles.buttonText}>Sell</Text>
        </TouchableOpacity>
      </View>
    </ScrollView>
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
  header: {
    padding: 16,
    borderBottomWidth: 1,
    borderBottomColor: '#E5E5E5',
  },
  name: {
    fontSize: 24,
    fontWeight: '700',
    color: '#1F2937',
  },
  symbol: {
    fontSize: 16,
    color: '#6B7280',
    marginTop: 4,
  },
  priceContainer: {
    padding: 16,
    flexDirection: 'row',
    alignItems: 'center',
  },
  price: {
    fontSize: 32,
    fontWeight: '700',
    color: '#1F2937',
  },
  priceChange: {
    fontSize: 18,
    fontWeight: '600',
    marginLeft: 12,
  },
  positiveChange: {
    color: '#10B981',
  },
  negativeChange: {
    color: '#EF4444',
  },
  statsContainer: {
    flexDirection: 'row',
    padding: 16,
    borderBottomWidth: 1,
    borderBottomColor: '#E5E5E5',
  },
  statItem: {
    flex: 1,
  },
  statLabel: {
    fontSize: 14,
    color: '#6B7280',
  },
  statValue: {
    fontSize: 16,
    fontWeight: '600',
    color: '#1F2937',
    marginTop: 4,
  },
  chartContainer: {
    padding: 16,
  },
  chart: {
    marginVertical: 8,
    borderRadius: 16,
  },
  buttonsContainer: {
    flexDirection: 'row',
    padding: 16,
    gap: 12,
  },
  button: {
    flex: 1,
    paddingVertical: 12,
    borderRadius: 8,
    alignItems: 'center',
  },
  buyButton: {
    backgroundColor: '#10B981',
  },
  sellButton: {
    backgroundColor: '#EF4444',
  },
  buttonText: {
    color: '#FFFFFF',
    fontSize: 16,
    fontWeight: '600',
  },
}); 