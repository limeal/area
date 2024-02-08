import { configureStore } from '@reduxjs/toolkit'
import { serviceApi } from './api'

/* Creating a Redux store with the serviceApi reducer and middleware. */
const store = configureStore({
    reducer: {
        [serviceApi.reducerPath]: serviceApi.reducer,
    },
    middleware: (getDefaultMiddleware) =>
        getDefaultMiddleware().concat(serviceApi.middleware),
})

export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch
export default store
