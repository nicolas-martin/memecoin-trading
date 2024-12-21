import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { Coin } from '../../types/coin';
import { getTopCoins, getCoinById } from '../../services/api';

interface CoinState {
  coins: Coin[];
  selectedCoin: Coin | null;
  loading: boolean;
  error: string | null;
}

const initialState: CoinState = {
  coins: [],
  selectedCoin: null,
  loading: false,
  error: null,
};

export const fetchTopCoins = createAsyncThunk(
  'coins/fetchTop',
  async (limit: number) => {
    const response = await getTopCoins(limit);
    return response;
  }
);

export const fetchCoinById = createAsyncThunk(
  'coins/fetchById',
  async (id: string) => {
    const response = await getCoinById(id);
    return response;
  }
);

const coinSlice = createSlice({
  name: 'coins',
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchTopCoins.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchTopCoins.fulfilled, (state, action) => {
        state.loading = false;
        state.coins = action.payload;
      })
      .addCase(fetchTopCoins.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch coins';
      })
      .addCase(fetchCoinById.fulfilled, (state, action) => {
        state.selectedCoin = action.payload;
      });
  },
});

export default coinSlice.reducer; 