import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { getTopMemeCoins, getMemeCoinDetail, updateMemeCoins, MemeCoin } from '../services/api';

interface MemeCoinsState {
  coins: MemeCoin[];
  selectedCoin: MemeCoin | null;
  loading: boolean;
  error: string | null;
}

const initialState: MemeCoinsState = {
  coins: [],
  selectedCoin: null,
  loading: false,
  error: null,
};

export const fetchMemeCoins = createAsyncThunk(
  'memeCoins/fetchMemeCoins',
  async () => {
    const coins = await getTopMemeCoins();
    return coins;
  }
);

export const fetchMemeCoinDetail = createAsyncThunk(
  'memeCoins/fetchMemeCoinDetail',
  async (id: string) => {
    const coin = await getMemeCoinDetail(id);
    return coin;
  }
);

export const updateMemeCoinsList = createAsyncThunk(
  'memeCoins/updateMemeCoinsList',
  async () => {
    await updateMemeCoins();
    const coins = await getTopMemeCoins();
    return coins;
  }
);

const memeCoinsSlice = createSlice({
  name: 'memeCoins',
  initialState,
  reducers: {
    clearSelectedCoin: (state) => {
      state.selectedCoin = null;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchMemeCoins.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchMemeCoins.fulfilled, (state, action) => {
        state.loading = false;
        state.coins = action.payload;
      })
      .addCase(fetchMemeCoins.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch meme coins';
      })
      .addCase(fetchMemeCoinDetail.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchMemeCoinDetail.fulfilled, (state, action) => {
        state.loading = false;
        state.selectedCoin = action.payload;
      })
      .addCase(fetchMemeCoinDetail.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch coin details';
      })
      .addCase(updateMemeCoinsList.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(updateMemeCoinsList.fulfilled, (state, action) => {
        state.loading = false;
        state.coins = action.payload;
      })
      .addCase(updateMemeCoinsList.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to update meme coins';
      });
  },
});

export const { clearSelectedCoin } = memeCoinsSlice.actions;
export default memeCoinsSlice.reducer; 