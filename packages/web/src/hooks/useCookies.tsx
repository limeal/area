import { useLayoutEffect, useState } from 'react'

import constate from 'constate'
import Cookies from 'universal-cookie'

const globalCookies = new Cookies()

const setGlobalCookie = globalCookies.set.bind(globalCookies)
const removeGlobalCookie = globalCookies.remove.bind(globalCookies)

/**
 * It returns an object with a setCookie function, a removeCookie function, and a cookies object
 * @returns {
 *         setCookie: setGlobalCookie,
 *         removeCookie: removeGlobalCookie,
 *         cookies,
 *     }
 */
function useCookies() {
    const [cookies, setCookies] = useState(globalCookies.getAll())

    useLayoutEffect(() => {
        function onChange() {
            const newCookies = globalCookies.getAll()
            setCookies(newCookies)
        }
        globalCookies.addChangeListener(onChange)
        return () => {
            globalCookies.removeChangeListener(onChange)
        }
    }, [])

    return {
        setCookie: setGlobalCookie,
        removeCookie: removeGlobalCookie,
        cookies,
    }
}

export const [CookiesProvider, useCookiesContext] = constate(useCookies)
