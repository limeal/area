import { motion } from 'framer-motion'
import { useNavigate } from 'react-router-dom'
import styled from 'styled-components'
import AnimatedButtonHL from './animated/buttons/AnimatedButtonHL'
import AnimatedButtonWB from './animated/buttons/AnimatedButtonWB'

const Container = styled.div`
    width: 100%;
    display: flex;
    justify-content: center;
    background: transparent;
    font-size: 1.2rem;
    background-color: #222222;
`

const NavigationBar = styled.nav`
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    width: 100%;
    max-width: 95%;
`

const Title = styled(motion.h1)`
    cursor: pointer;
    color: #fff;

    &:hover {
        color: #b2ddf7;
    }
`

const TabList = styled.ul`
    display: flex;
    flex-direction: row;
    align-items: center;
    list-style: none;
`

const Tab = styled.li`
    margin: 0 10px;
`

/* Defining the props that the component will receive. */
interface NavigationbarProps {
    isConnected: boolean
    LoginOnClick: () => void
}

// React router doc: https://reactrouter.com/en/main
/* A React component that is a navigation bar. */
const Navigationbar = ({ isConnected, LoginOnClick }: NavigationbarProps) => {
    const navigate = useNavigate()

    const tabs = [
        /* {
            name: 'Explore',
            link: '/explore',
            color: '#b4e1ff',
            backgroundColor: '#fff',
        }, */
        {
            name: 'My Applets',
            link: '/applets',
            color: '#b4e1ff',
            backgroundColor: '#fff',
        },
    ]

    return (
        <Container>
            <NavigationBar>
                <Title
                    onClick={() => navigate('/')}
                    whileHover={{
                        rotate: 360,
                        transition: { duration: 0.5 },
                    }}
                >
                    AREA
                </Title>
                <TabList>
                    {tabs.map(
                        (tab, index) =>
                            (isConnected || tab.name != 'My Applets') && (
                                <Tab key={index}>
                                    <AnimatedButtonHL
                                        color={tab.color}
                                        onClick={() =>
                                            navigate(tab.link, {
                                                replace: true,
                                            })
                                        }
                                        hlColor={tab.backgroundColor}
                                    >
                                        {tab.name}
                                    </AnimatedButtonHL>
                                </Tab>
                            )
                    )}
                    {!isConnected ? (
                        <Tab>
                            <AnimatedButtonWB
                                color="#fff"
                                bgColor="#b9090b"
                                onClick={LoginOnClick}
                            >
                                Login
                            </AnimatedButtonWB>
                        </Tab>
                    ) : (
                        <>
                            <Tab>
                                <AnimatedButtonWB
                                    color="#fff"
                                    bgColor="#b9090b"
                                    onClick={() =>
                                        navigate('/create', { replace: true })
                                    }
                                >
                                    Create
                                </AnimatedButtonWB>
                            </Tab>
                            <Tab>
                                <AnimatedButtonWB
                                    color="#fff"
                                    bgColor="#4f518c"
                                    onClick={() =>
                                        navigate('/profile', { replace: true })
                                    }
                                >
                                    Profile
                                </AnimatedButtonWB>
                            </Tab>
                        </>
                    )}
                </TabList>
            </NavigationBar>
        </Container>
    )
}

export default Navigationbar
