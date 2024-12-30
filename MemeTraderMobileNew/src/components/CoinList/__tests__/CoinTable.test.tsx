import React from 'react';
import { screen } from '@testing-library/react-native';
import { render } from '@testing-library/react-native/build/pure';
import { fireEvent } from '@testing-library/react-native/build/fireEvent';
import { Image } from 'react-native';
import { CoinTable } from '../CoinTable';

const mockCoins = [
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

describe('CoinTable', () => {
  it('renders correctly', () => {
    const onCoinPress = jest.fn();
    const { getByText } = render(
      <CoinTable coins={mockCoins} onCoinPress={onCoinPress} />
    );

    // Check if header is rendered
    expect(getByText('#')).toBeTruthy();
    expect(getByText('Name')).toBeTruthy();
    expect(getByText('Price')).toBeTruthy();
    expect(getByText('24h %')).toBeTruthy();
    expect(getByText('Volume')).toBeTruthy();
    expect(getByText('Market Cap')).toBeTruthy();

    // Check if coin data is rendered
    expect(getByText('Dogecoin')).toBeTruthy();
    expect(getByText('DOGE')).toBeTruthy();
    expect(getByText('$0.10')).toBeTruthy();
    expect(getByText('+10.00%')).toBeTruthy();
    expect(getByText('$1.00B')).toBeTruthy();
    expect(getByText('$500.00M')).toBeTruthy();

    expect(getByText('Shiba Inu')).toBeTruthy();
    expect(getByText('SHIB')).toBeTruthy();
    expect(getByText('1.00e-5')).toBeTruthy();
    expect(getByText('-10.00%')).toBeTruthy();
    expect(getByText('$500.00M')).toBeTruthy();
    expect(getByText('$250.00M')).toBeTruthy();
  });

  it('handles coin press correctly', () => {
    const onCoinPress = jest.fn();
    const { getByText } = render(
      <CoinTable coins={mockCoins} onCoinPress={onCoinPress} />
    );

    fireEvent.press(getByText('Dogecoin'));
    expect(onCoinPress).toHaveBeenCalledWith(mockCoins[0]);

    fireEvent.press(getByText('Shiba Inu'));
    expect(onCoinPress).toHaveBeenCalledWith(mockCoins[1]);
  });

  it('handles empty coin list', () => {
    const onCoinPress = jest.fn();
    const { getByText, queryByText } = render(
      <CoinTable coins={[]} onCoinPress={onCoinPress} />
    );

    // Header should still be rendered
    expect(getByText('#')).toBeTruthy();
    expect(getByText('Name')).toBeTruthy();
    expect(getByText('Price')).toBeTruthy();
    expect(getByText('24h %')).toBeTruthy();
    expect(getByText('Volume')).toBeTruthy();
    expect(getByText('Market Cap')).toBeTruthy();

    // No coin data should be rendered
    expect(queryByText('Dogecoin')).toBeNull();
    expect(queryByText('Shiba Inu')).toBeNull();
  });
}); 