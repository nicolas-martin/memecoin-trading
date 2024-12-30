import React from 'react';
import { render, waitFor, fireEvent } from '@testing-library/react-native/build';
import { CoinTable } from '../CoinTable';
import { getTopMemeCoins } from '../../../services/api';

// Mock the API service
jest.mock('../../../services/api', () => ({
  getTopMemeCoins: jest.fn(),
}));

// Mock the entire react-native module
jest.mock('react-native', () => ({
  Image: 'Image',
  Text: 'Text',
  View: 'View',
  TouchableOpacity: 'TouchableOpacity',
  ScrollView: 'ScrollView',
  StyleSheet: {
    create: (styles: any) => styles,
  },
  Dimensions: {
    get: () => ({
      width: 375,
      height: 812,
    }),
  },
}));

// Mock the placeholder image
jest.mock('../../../assets/placeholder.png', () => ({
  uri: 'placeholder.png',
}));

const mockApiResponse = [
  {
    id: '1',
    symbol: 'DOGE',
    name: 'Dogecoin',
    logoUrl: 'https://example.com/doge.png',
    price: 0.1,
    marketCap: 1000000000,
    volume24h: 500000000,
    priceChange24h: 0.01,
    priceChangePercentage24h: 10,
    contractAddress: '0x123',
  },
  {
    id: '2',
    symbol: 'SHIB',
    name: 'Shiba Inu',
    logoUrl: 'https://example.com/shib.png',
    price: 0.00001,
    marketCap: 500000000,
    volume24h: 250000000,
    priceChange24h: -0.000001,
    priceChangePercentage24h: -10,
    contractAddress: '0x456',
  },
];

describe('CoinTable Integration', () => {
  beforeEach(() => {
    // Clear mock calls between tests
    jest.clearAllMocks();
    // Setup mock API response
    (getTopMemeCoins as jest.Mock).mockResolvedValue(mockApiResponse);
  });

  it('displays data fetched from the Go service correctly', async () => {
    const onCoinPress = jest.fn();
    const { getByText, getAllByText } = render(
      <CoinTable coins={mockApiResponse} onCoinPress={onCoinPress} />
    );

    // Verify that the table headers are rendered
    expect(getByText('#')).toBeTruthy();
    expect(getByText('Name')).toBeTruthy();
    expect(getByText('Price')).toBeTruthy();
    expect(getByText('24h %')).toBeTruthy();
    expect(getByText('Volume')).toBeTruthy();
    expect(getByText('Market Cap')).toBeTruthy();

    // Verify that the first coin's data from the API is displayed correctly
    await waitFor(() => {
      expect(getByText('Dogecoin')).toBeTruthy();
      expect(getByText('DOGE')).toBeTruthy();
      expect(getByText('$0.10')).toBeTruthy();
      expect(getByText('+10.00%')).toBeTruthy();
      expect(getByText('$1.00B')).toBeTruthy();
    });

    // Verify that the second coin's data from the API is displayed correctly
    await waitFor(() => {
      expect(getByText('Shiba Inu')).toBeTruthy();
      expect(getByText('SHIB')).toBeTruthy();
      expect(getByText('$1.00e-5')).toBeTruthy();
      expect(getByText('-10.00%')).toBeTruthy();
      expect(getByText('$250.00M')).toBeTruthy();
    });

    // Verify that there are two elements with $500.00M (one for volume, one for market cap)
    const elements = getAllByText('$500.00M');
    expect(elements).toHaveLength(2);
  });

  it('handles coin press correctly with API data', async () => {
    const onCoinPress = jest.fn();
    const { getByText } = render(
      <CoinTable coins={mockApiResponse} onCoinPress={onCoinPress} />
    );

    // Verify interaction with the first coin
    fireEvent.press(getByText('Dogecoin'));
    expect(onCoinPress).toHaveBeenCalledWith(mockApiResponse[0]);

    // Verify interaction with the second coin
    fireEvent.press(getByText('Shiba Inu'));
    expect(onCoinPress).toHaveBeenCalledWith(mockApiResponse[1]);
  });

  it('handles empty API response', async () => {
    // Mock empty API response
    (getTopMemeCoins as jest.Mock).mockResolvedValue([]);
    
    const onCoinPress = jest.fn();
    const { getByText, queryByText } = render(
      <CoinTable coins={[]} onCoinPress={onCoinPress} />
    );

    // Verify that headers are still rendered with empty data
    expect(getByText('#')).toBeTruthy();
    expect(getByText('Name')).toBeTruthy();
    expect(getByText('Price')).toBeTruthy();
    expect(getByText('24h %')).toBeTruthy();
    expect(getByText('Volume')).toBeTruthy();
    expect(getByText('Market Cap')).toBeTruthy();

    // Verify that no coin data is rendered
    expect(queryByText('Dogecoin')).toBeNull();
    expect(queryByText('Shiba Inu')).toBeNull();
  });

  it('formats API data correctly', async () => {
    const onCoinPress = jest.fn();
    const { getByText, getAllByText } = render(
      <CoinTable coins={mockApiResponse} onCoinPress={onCoinPress} />
    );

    // Verify price formatting
    await waitFor(() => {
      // Regular price
      expect(getByText('$0.10')).toBeTruthy();
      // Small price in scientific notation
      expect(getByText('$1.00e-5')).toBeTruthy();
    });

    // Verify percentage formatting
    await waitFor(() => {
      // Positive percentage
      expect(getByText('+10.00%')).toBeTruthy();
      // Negative percentage
      expect(getByText('-10.00%')).toBeTruthy();
    });

    // Verify market cap and volume formatting
    await waitFor(() => {
      // Billions
      expect(getByText('$1.00B')).toBeTruthy();
      // Multiple elements with the same value
      const millionElements = getAllByText('$500.00M');
      expect(millionElements).toHaveLength(2);
      expect(getByText('$250.00M')).toBeTruthy();
    });
  });
}); 