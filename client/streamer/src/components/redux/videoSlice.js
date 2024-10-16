import { createSlice } from '@reduxjs/toolkit';

const videoSlice = createSlice({
  name: 'videos',
  initialState: { videoList: [] },
  reducers: {
    setVideos: (state, action) => {
      state.videoList = action.payload;
    }
  }
});

export const { setVideos } = videoSlice.actions;
export default videoSlice.reducer;
