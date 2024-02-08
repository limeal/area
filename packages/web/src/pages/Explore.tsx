import { motion } from 'framer-motion'
import { useEffect, useState } from 'react'
import styled from 'styled-components'
import { useGetStoreAppletsQuery } from '../redux/api'

import AppletCard from '../components/applets/AppletCard'
import { IApplet, IService } from '../interfaces'

const AppletsContainer = styled.ul`
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: center;
    list-style: none;
`

const Container = styled.div`
    color: #fff;
    background-color: #222222;
    height: 100%;
`

const Title = styled(motion.h1)`
    font-size: 5rem;
    margin: 0;
    padding: 0;
`

/* Defining the props that the component will receive. */
interface ExploreProps {
    services: IService[]
}

/* A React component that is using the `useGetStoreAppletsQuery` hook to fetch data from the GraphQL
API. */
const Explore = ({ services }: ExploreProps) => {
    const { data } = useGetStoreAppletsQuery(null, {
        refetchOnMountOrArgChange: true,
    })
    const [applets, setApplets] = useState(data)

    useEffect(() => {
        if (data) setApplets(data)
    }, [data])

    // List all applets that have public visibility in a grid view 4x4
    return (
        <Container>
            <Title
                initial={{ opacity: 0, scale: 0.5 }}
                animate={{ opacity: 1, scale: 1 }}
                transition={{ duration: 0.5 }}
            >
                Explore
            </Title>
            {applets && applets.length > 0 && (
                <AppletsContainer>
                    {applets.map((applet: IApplet, index: number) => (
                        <li key={index}>
                            <AppletCard
                                services={services}
                                appletInfos={applet}
                                mode="copy"
                            />
                        </li>
                    ))}
                </AppletsContainer>
            )}
        </Container>
    )
}

export default Explore
