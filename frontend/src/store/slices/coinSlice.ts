import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import { getTrendingPairs, PairData } from '../../services/dexscreener';

interface CoinState {
  trendingCoins: PairData[];
  selectedCoin: PairData | null;
  loading: boolean;
  error: string | null;
}

const initialState: CoinState = {
  trendingCoins: [],
  selectedCoin: null,
  loading: false,
  error: null,
};

export const fetchTrendingCoins = createAsyncThunk<PairData[]>(
  'coins/fetchTrending',
  async () => {
    const coins = await getTrendingPairs();
    return coins;
  }
);

const coinSlice = createSlice({
  name: 'coins',
  initialState,
  reducers: {
    setSelectedCoin: (state, action: PayloadAction<PairData | null>) => {
      state.selectedCoin = action.payload;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchTrendingCoins.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchTrendingCoins.fulfilled, (state, action) => {
        state.loading = false;
        state.trendingCoins = action.payload;
      })
      .addCase(fetchTrendingCoins.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch trending coins';
      });
  },
});

export const { setSelectedCoin } = coinSlice.actions;
export default coinSlice.reducer; 