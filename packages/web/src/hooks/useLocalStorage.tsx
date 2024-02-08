import { useEffect, useState } from 'react'

/**
 * If we're in the browser, get the value from localStorage, otherwise return the default value.
 * @param {string} key - The key to store the value under.
 * @param {string} defaultValue - The value to return if the key doesn't exist in localStorage.
 * @returns The value of the key in localStorage, or the defaultValue if the key is not found.
 */
const getStorageValue = (key: string, defaultValue: string) => {
    if (typeof window === 'undefined') return defaultValue
    const v = localStorage.getItem(key)
    if (!v) return defaultValue
    return v
}

/**
 * It returns a stateful value and a function to update it.
 *
 * The only difference between this and useState is that we're using localStorage to persist the value
 * between page refreshes
 * @param {string} key - string - The key to store the value in localStorage
 * @param {string} defaultValue - The default value to use if the localStorage value is not found.
 * @returns A function that takes in a key and a default value and returns an array with two values.
 * The first value is the value of the key in local storage or the default value if the key is not in
 * local storage. The second value is a function that sets the value of the key in local storage.
 */
export const useLocalStorage = (key: string, defaultValue: string) => {
    const [storage, setStorage] = useState(() => {
        return getStorageValue(key, defaultValue)
    })

    useEffect(() => {
        // storing input name
        localStorage.setItem(key, storage)
    }, [key, storage])

    return [storage, setStorage]
}
