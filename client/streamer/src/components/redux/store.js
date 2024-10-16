import { configureStore } from '@reduxjs/toolkit';
import adminReducer from './adminSlice';
import videoReducer from './videoSlice';

export const store = configureStore({
  reducer: {
    admin: adminReducer,
    videos: videoReducer,
  },
});
