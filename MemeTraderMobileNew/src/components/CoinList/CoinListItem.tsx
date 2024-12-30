import React from 'react';
import { View, Text, Image, StyleSheet, TouchableOpacity } from 'react-native';
import { MemeCoin } from '../../services/api';

interface CoinListItemProps {
  coin: MemeCoin;
  onPress: (coin: MemeCoin) => void;
}

export const CoinListItem: React.FC<CoinListItemProps> = ({ coin, onPress }) => {
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
    <TouchableOpacity
      style={styles.container}
      onPress={() => onPress(coin)}
      activeOpacity={0.7}
    >
      <View style={styles.leftContainer}>
        {coin.LogoURL ? (
          <Image source={{ uri: coin.LogoURL }} style={styles.logo} />
        ) : (
          <View style={[styles.logo, styles.placeholderLogo]} />
        )}
        <View style={styles.nameContainer}>
          <Text style={styles.name} numberOfLines={1}>
            {coin.name}
          </Text>
          <Text style={styles.symbol}>{coin.symbol}</Text>
        </View>
      </View>
      <View style={styles.rightContainer}>
        <Text style={styles.price}>${formatPrice(coin.price)}</Text>
        <Text style={styles.marketCap}>{formatMarketCap(coin.marketCap)}</Text>
        <Text
          style={[
            styles.priceChange,
            coin.priceChangePercentage24h >= 0
              ? styles.positiveChange
              : styles.negativeChange,
          ]}
        >
          {coin.priceChangePercentage24h >= 0 ? '+' : ''}
          {coin.priceChangePercentage24h.toFixed(2)}%
        </Text>
      </View>
    </TouchableOpacity>
  );
};

const styles = StyleSheet.create({
  container: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: 16,
    backgroundColor: '#FFFFFF',
    borderBottomWidth: 1,
    borderBottomColor: '#E5E5E5',
  },
  leftContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    flex: 1,
  },
  logo: {
    width: 40,
    height: 40,
    borderRadius: 20,
  },
  placeholderLogo: {
    backgroundColor: '#E5E5E5',
  },
  nameContainer: {
    marginLeft: 12,
    flex: 1,
  },
  name: {
    fontSize: 16,
    fontWeight: '600',
    color: '#1F2937',
  },
  symbol: {
    fontSize: 14,
    color: '#6B7280',
    marginTop: 2,
  },
  rightContainer: {
    alignItems: 'flex-end',
  },
  price: {
    fontSize: 16,
    fontWeight: '600',
    color: '#1F2937',
  },
  marketCap: {
    fontSize: 14,
    color: '#6B7280',
    marginTop: 2,
  },
  priceChange: {
    fontSize: 14,
    fontWeight: '500',
    marginTop: 2,
  },
  positiveChange: {
    color: '#10B981',
  },
  negativeChange: {
    color: '#EF4444',
  },
}); 