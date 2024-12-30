import { configureStore } from '@reduxjs/toolkit';
import memeCoinsReducer from './memeCoinsSlice';

export const store = configureStore({
  reducer: {
    memeCoins: memeCoinsReducer,
  },
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch; 