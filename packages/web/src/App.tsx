import { motion } from 'framer-motion'
import queryString from 'query-string'
import { useEffect, useRef, useState } from 'react'
import { Provider } from 'react-redux'
import { BrowserRouter, Route, Routes } from 'react-router-dom'
import { ToastContainer } from 'react-toastify'
import 'react-toastify/dist/ReactToastify.css'
import styled from 'styled-components'
import { CookiesProvider } from './hooks/useCookies'
import store from './redux/store'

// Pages
import Applets from './pages/Applets'
import Explore from './pages/Explore'
import Home from './pages/Home'
import Profile from './pages/Profile'

// Components
import Navigationbar from './components/NavigationBar'
import AnimatedText from './components/animated/AnimatedText'
import LoginModal from './components/modals/AuthModal'

// API
import { useLocalStorage } from './hooks/useLocalStorage'
import Authenticators from './pages/Authenticators'
import Applet from './pages/applets/Applet' /*
import UpdateApplet from './pages/applets/UpdateApplet'
import UpdateArea from './pages/applets/UpdateArea' */
import CreateApplet from './pages/applets/CreateApplet'
import CreateArea from './pages/applets/CreateArea'
import { useGetAboutQuery } from './redux/api'

const Container = styled.div`
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: #0e1116;
    display: flex;
    justify-content: center;
    align-items: center;
    flex-direction: column;
    color: #fff;
    font-size: 1.5rem;
`

const LoadingElement = styled(motion.div)`
    width: 150px;
    height: 150px;
    border-radius: 5rem;
    background: none;
    border: 5px dashed #b2ddf7;
    opacity: 1;
`

const Header = styled.header`
    width: 100%;
    background: #222222;
    flex: 1;
`

const Main = styled.main`
    width: 100%;
    flex: 12;
`

const Footer = styled.footer`
    flex: 2;
    bottom: 0;
    left: 0;
    width: 100%;
    height: 100px;
    background: #0e1116;
    display: flex;
    justify-content: center;
    align-items: center;
    flex-direction: column;
    color: #fff;
    font-size: 1.5rem;
`

/* A loading screen. */
const Loading = () => {
    return (
        <Container>
            <LoadingElement
                animate={{ rotate: 180 }}
                transition={{
                    type: 'spring',
                    duration: 2.1,
                    repeat: Infinity,
                }}
            />
            <AnimatedText text="Loading..." />
        </Container>
    )
}

/**
 * It renders a header, a main area, and a footer
 * @returns The Area component is being returned.
 */
const Area = () => {
    // Fetch service in server (/about.json)
    const [open, setOpen] = useState(false)
    const { isLoading, error, data, refetch, isSuccess } = useGetAboutQuery(
        null,
        {
            refetchOnReconnect: true,
        }
    )
    const [storage, setStorage] = useLocalStorage('token', '')
    const interval = useRef<any>(null)

    useEffect(() => {
        const data = queryString.parse(window.location.search)
        if (
            data &&
            Object.keys(data).length > 0 &&
            Object.keys(data).includes('code')
        ) {
            // Call the handler with the data and close the popup
            window.opener.handleAuthorization(data)
            window.close()
        }
    }, [])

    useEffect(() => {
        if (isLoading && !interval.current) {
            interval.current = setInterval(() => {
                refetch()
            }, 2000)
        }
        if (isSuccess && interval.current) {
            clearInterval(interval.current)
            interval.current = null
        }
    }, [isSuccess, isLoading])

    if (isLoading || error) return <Loading />
    if (!data) return <Loading />

    return (
        <BrowserRouter>
            <ToastContainer
                autoClose={300}
                pauseOnFocusLoss={false}
                closeOnClick={true}
                newestOnTop={true}
            />
            <LoginModal
                updateToken={setStorage}
                open={open}
                close={() => setOpen(false)}
                authenticators={data.authenticators}
            />
            <Header>
                <Navigationbar
                    isConnected={storage !== '' ? true : false}
                    LoginOnClick={() => setOpen(true)}
                />
            </Header>
            <Main>
                <Routes>
                    <Route
                        path="/"
                        element={
                            <Home
                                storage={storage}
                                openAuthModal={() => setOpen(true)}
                            />
                        }
                    />
                    <Route
                        path="/profile"
                        element={<Profile updateToken={setStorage} />}
                    />
                    <Route
                        path="/authorizations"
                        element={
                            <Authenticators
                                authenticators={data.authenticators}
                            />
                        }
                    />
                    <Route
                        path="/applets"
                        element={<Applets services={data.services} />}
                    />
                    <Route
                        path="/applets/:id"
                        element={<Applet services={data.services} />}
                    />
                    {/* Update Applet */}
                    {/* <Route
                        path="/applets/:id/update"
                        element={<UpdateApplet services={data.services} />}
                    />
                    <Route
                        path="/applets/:id/update/set-action"
                        element={
                            <UpdateArea
                                services={data.services}
                                type="action"
                            />
                        }
                    />
                    <Route
                        path="/applets/:id/update/set-reaction"
                        element={
                            <UpdateArea
                                services={data.services}
                                type="reaction"
                            />
                        }
                    /> */}
                    <Route
                        path="/explore"
                        element={<Explore services={data.services} />}
                    />
                    {/* Create Applet */}
                    <Route
                        path="/create"
                        element={<CreateApplet services={data.services} />}
                    />
                    <Route
                        path="/create/add-action"
                        element={
                            <CreateArea
                                services={data.services}
                                type="action"
                            />
                        }
                    />
                    <Route
                        path="/create/add-reaction"
                        element={
                            <CreateArea
                                services={data.services}
                                type="reaction"
                            />
                        }
                    />
                    <Route path="*" element={<h1>404</h1>} />
                </Routes>
            </Main>
            <Footer>
                <p>© 2023 - All rights reserved</p>
                <p>Web App credits @limeal and @p0lar1s</p>
            </Footer>
        </BrowserRouter>
    )
}

/**
 * The App function returns a Provider component that wraps the Area component
 * @returns The App component is being returned.
 */
const App = () => {
    return (
        <CookiesProvider>
            <Provider store={store}>
                <Area />
            </Provider>
        </CookiesProvider>
    )
}

export default App
