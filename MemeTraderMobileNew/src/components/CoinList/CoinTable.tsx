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
} from 'react-native';
import { MemeCoin } from '../../services/api';

interface CoinTableProps {
  coins: MemeCoin[];
  onCoinPress: (coin: MemeCoin) => void;
}

export const CoinTable: React.FC<CoinTableProps> = ({ coins, onCoinPress }) => {
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
          source={{ uri: coin.logoUrl }}
          style={styles.logo}
          defaultSource={require('../../assets/placeholder.png')}
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

  return (
    <View style={styles.container}>
      {renderHeader()}
      <ScrollView horizontal showsHorizontalScrollIndicator={false}>
        <View>
          {coins.map((coin, index) => renderRow(coin, index))}
        </View>
      </ScrollView>
    </View>
  );
};

const windowWidth = Dimensions.get('window').width;

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#FFFFFF',
  } as ViewStyle,
  headerRow: {
    flexDirection: 'row',
    paddingVertical: 12,
    paddingHorizontal: 16,
    backgroundColor: '#F3F4F6',
    borderBottomWidth: 1,
    borderBottomColor: '#E5E5E5',
  } as ViewStyle,
  row: {
    flexDirection: 'row',
    paddingVertical: 12,
    paddingHorizontal: 16,
    borderBottomWidth: 1,
    borderBottomColor: '#E5E5E5',
    alignItems: 'center',
  } as ViewStyle,
  headerCell: {
    fontWeight: '600',
    fontSize: 12,
    color: '#6B7280',
  } as TextStyle,
  cell: {
    fontSize: 14,
    color: '#1F2937',
  } as TextStyle,
  rankCell: {
    width: 40,
  } as TextStyle,
  nameCell: {
    width: 150,
  } as TextStyle,
  priceCell: {
    width: 100,
    textAlign: 'right',
  } as TextStyle,
  changeCell: {
    width: 80,
    textAlign: 'right',
  } as TextStyle,
  volumeCell: {
    width: 100,
    textAlign: 'right',
  } as TextStyle,
  marketCapCell: {
    width: 100,
    textAlign: 'right',
  } as TextStyle,
  nameContainer: {
    flexDirection: 'row',
    alignItems: 'center',
  } as ViewStyle,
  logo: {
    width: 24,
    height: 24,
    borderRadius: 12,
    marginRight: 8,
  } as ImageStyle,
  nameTextContainer: {
    flex: 1,
  } as ViewStyle,
  coinName: {
    fontSize: 14,
    fontWeight: '500',
    color: '#1F2937',
  } as TextStyle,
  symbol: {
    fontSize: 12,
    color: '#6B7280',
  } as TextStyle,
  positiveChange: {
    color: '#10B981',
  } as TextStyle,
  negativeChange: {
    color: '#EF4444',
  } as TextStyle,
}); 