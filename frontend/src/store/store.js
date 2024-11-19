import { configureStore, combineReducers } from '@reduxjs/toolkit';
import counterReducer from './reducers/counterReducer.js';
import userReducer from './reducers/userReducer.js';

// Combine reducers
const rootReducer = combineReducers({
  counter: counterReducer,
  user: userReducer
});

// Create Redux store with the combined reducer
const store = configureStore({
  reducer: rootReducer
});

export default store;
