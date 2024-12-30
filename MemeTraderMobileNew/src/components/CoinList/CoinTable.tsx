import React from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  Image,
  Dimensions,
  TextStyle,
  ViewStyle,
  ImageStyle,
  Platform,
} from 'react-native';
import { MemeCoin } from '../../services/api';

interface CoinTableProps {
  coins: MemeCoin[];
  onCoinPress: (coin: MemeCoin) => void;
}

export const CoinTable: React.FC<CoinTableProps> = ({ coins = [], onCoinPress }) => {
  const formatPrice = (price: number) => {
    return price < 0.01 ? price.toExponential(2) : price.toFixed(2);
  };

  const formatNumber = (num: number) => {
    if (num >= 1e9) {
      return `$${(num / 1e9).toFixed(2)}B`;
    } else if (num >= 1e6) {
      return `$${(num / 1e6).toFixed(2)}M`;
    } else if (num >= 1e3) {
      return `$${(num / 1e3).toFixed(2)}K`;
    }
    return `$${num.toFixed(2)}`;
  };

  const renderHeader = () => (
    <View style={styles.headerRow}>
      <Text style={[styles.headerCell, styles.rankCell]}>#</Text>
      <Text style={[styles.headerCell, styles.nameCell]}>Name</Text>
      <Text style={[styles.headerCell, styles.priceCell]}>Price</Text>
      <Text style={[styles.headerCell, styles.changeCell]}>24h %</Text>
      <Text style={[styles.headerCell, styles.volumeCell]}>Volume</Text>
      <Text style={[styles.headerCell, styles.marketCapCell]}>Market Cap</Text>
    </View>
  );

  const renderRow = (coin: MemeCoin, index: number) => (
    <TouchableOpacity
      key={coin.id}
      style={styles.row}
      onPress={() => onCoinPress(coin)}
      activeOpacity={0.7}
    >
      <Text style={[styles.cell, styles.rankCell]}>{index + 1}</Text>
      <View style={[styles.nameContainer, { width: styles.nameCell.width }]}>
        <Image
          source={coin.logoUrl ? { uri: coin.logoUrl } : require('../../../assets/placeholder.png')}
          style={styles.logo}
          onError={(e) => {
            console.log('Image loading error:', e.nativeEvent.error);
          }}
        />
        <View style={styles.nameTextContainer}>
          <Text style={styles.coinName} numberOfLines={1}>
            {coin.name}
          </Text>
          <Text style={styles.symbol}>{coin.symbol.toUpperCase()}</Text>
        </View>
      </View>
      <Text style={[styles.cell, styles.priceCell]}>
        ${formatPrice(coin.price)}
      </Text>
      <Text
        style={[
          styles.cell,
          styles.changeCell,
          coin.priceChangePercentage24h >= 0
            ? styles.positiveChange
            : styles.negativeChange,
        ]}
      >
        {coin.priceChangePercentage24h >= 0 ? '+' : ''}
        {coin.priceChangePercentage24h.toFixed(2)}%
      </Text>
      <Text style={[styles.cell, styles.volumeCell]}>
        {formatNumber(coin.volume24h)}
      </Text>
      <Text style={[styles.cell, styles.marketCapCell]}>
        {formatNumber(coin.marketCap)}
      </Text>
    </TouchableOpacity>
  );

  // Safety check for coins array
  const safeCoins = Array.isArray(coins) ? coins : [];

  return (
    <View style={styles.container}>
      {renderHeader()}
      {safeCoins.length === 0 ? (
        <View style={styles.errorContainer}>
          <Text style={styles.errorText}>No data available</Text>
        </View>
      ) : (
        <ScrollView horizontal showsHorizontalScrollIndicator={false}>
          <View>
            {safeCoins.map((coin, index) => renderRow(coin, index))}
          </View>
        </ScrollView>
      )}
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#FFFFFF',
  },
  headerRow: {
    flexDirection: 'row',
    paddingVertical: 12,
    paddingHorizontal: 16,
    backgroundColor: '#F3F4F6',
    borderBottomWidth: 1,
    borderBottomColor: '#E5E5E5',
  },
  row: {
    flexDirection: 'row',
    paddingVertical: 12,
    paddingHorizontal: 16,
    borderBottomWidth: 1,
    borderBottomColor: '#E5E5E5',
    alignItems: 'center',
  },
  headerCell: {
    fontWeight: '600',
    fontSize: 12,
    color: '#6B7280',
  },
  cell: {
    fontSize: 14,
    color: '#1F2937',
  },
  rankCell: {
    width: 40,
  },
  nameCell: {
    width: 150,
  },
  priceCell: {
    width: 100,
    textAlign: 'right',
  },
  changeCell: {
    width: 80,
    textAlign: 'right',
  },
  volumeCell: {
    width: 100,
    textAlign: 'right',
  },
  marketCapCell: {
    width: 100,
    textAlign: 'right',
  },
  nameContainer: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  logo: {
    width: 24,
    height: 24,
    borderRadius: 12,
    marginRight: 8,
  },
  nameTextContainer: {
    flex: 1,
  },
  coinName: {
    fontSize: 14,
    fontWeight: '500',
    color: '#1F2937',
  },
  symbol: {
    fontSize: 12,
    color: '#6B7280',
  },
  positiveChange: {
    color: '#10B981',
  },
  negativeChange: {
    color: '#EF4444',
  },
  errorContainer: {
    padding: 20,
    alignItems: 'center',
  },
  errorText: {
    color: '#6B7280',
    fontSize: 14,
  },
}); 