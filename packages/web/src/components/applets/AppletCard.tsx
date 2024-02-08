import { motion } from 'framer-motion'
import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import styled from 'styled-components'

import { IApplet, IService } from '../../interfaces'
import {
    useGetAppletReactionsQuery,
    useUpdateAppletActivityMutation,
} from '../../redux/api'
import { getAvatar, getColor } from '../../utils/more'
import { applyPretty } from '../../utils/string'

import AnimatedSwitch from '../animated/AnimatedSwitch'
import AnimatedButtonWB from '../animated/buttons/AnimatedButtonWB'

const Container = styled(motion.div)`
    display: flex;
    flex-direction: column;
    color: #fff;
    cursor: pointer;
    width: 100%;
    height: 100%;
    padding: 10px;
    border-radius: 10px;
    background: ${(props) => props.theme.background || '#222222'};
    box-shadow: 0 0 10px 0 rgba(0, 0, 0, 0.2);
    justify-content: space-between;
`

const Title = styled.h2`
    font-size: 2.2rem;
    color: #fff;
`

const Description = styled.p`
    font-size: 0.8rem;
    color: #fff;
`

const Header = styled.div`
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: space-around;
    width: 100%;
    padding: 15px 0;
    border-bottom: 1px solid #fff;
`

const ElementContainer = styled.div`
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: space-around;
    width: 100%;
`

const Main = styled.div`
    padding: 15px 0;
    display: flex;
    flex-direction: column;
    justify-content: space-around;
    width: 100%;
`

interface AppletProps {
    services: IService[]
    appletInfos: IApplet
    mode: 'copy' | 'my'
}

/**
 * It takes a service and an element, and returns a styled div with an image and a paragraph
 * @param {any}  - service - the service object
 * @returns A React component
 */
const AppletService = ({ service, elem }: any) => {
    return (
        <ElementContainer>
            <img
                src={getAvatar(service)}
                alt={service?.name}
                style={{ width: '50px', height: '50px' }}
            />
            <p>{applyPretty(elem)}</p>
        </ElementContainer>
    )
}

/* A React component that takes in a service and an element, and returns a styled div with an image and
a paragraph. */
const AppletCard = ({ services, appletInfos, mode }: AppletProps) => {
    const [useUpdateActivity, { status }] = useUpdateAppletActivityMutation()
    const [active, setActive] = useState(appletInfos.active)
    const [actionService] = useState(
        services.find((s) => s.name === appletInfos.action.split(';')[0]) ||
            null
    )
    const { data: reactions } = useGetAppletReactionsQuery(appletInfos.id, {
        refetchOnMountOrArgChange: true,
    })

    const navigate = useNavigate()

    useEffect(() => {
        if (status === 'fulfilled') setActive(!active)
    }, [status])

    return (
        <>
            <Container
                onClick={() => navigate(`/applets/${appletInfos.id}`)}
                theme={{
                    background:
                        'linear-gradient(45deg,' +
                        getColor(actionService) +
                        ',' +
                        reactions
                            ?.map((reaction) => {
                                const reactionService = services.find(
                                    (s) => s.name === reaction.service
                                )
                                return getColor(reactionService)
                            })
                            .join(',') +
                        ')',
                }}
            >
                <div>
                    <Header>
                        <AppletService
                            service={actionService}
                            elem={appletInfos.action.split(';')[1]}
                        />
                        {reactions?.map((reaction, index) => {
                            const reactionService = services.find(
                                (s) => s.name === reaction.service
                            )

                            return (
                                <AppletService
                                    key={index}
                                    service={reactionService}
                                    elem={reaction.name}
                                />
                            )
                        })}
                    </Header>
                    <Main>
                        <Title>Name: {appletInfos.name}</Title>
                        <Description>
                            Description: {appletInfos.description}
                        </Description>
                    </Main>
                </div>
                {mode === 'my' ? (
                    <AnimatedSwitch
                        isOn={active}
                        whileHover={{
                            scale: 1.1,
                            transition: { duration: 0.2 },
                        }}
                        whileTap={{
                            scale: 0.95,
                            transition: { duration: 0.3 },
                        }}
                        onClick={(e: any) => {
                            e.stopPropagation()
                            useUpdateActivity({
                                id: appletInfos.id,
                                active: !active,
                            })
                        }}
                        style={{
                            alignSelf: 'center',
                        }}
                    />
                ) : (
                    <AnimatedButtonWB
                        bgColor="#222222"
                        color="#fff"
                        style={{
                            alignSelf: 'center',
                        }}
                    >
                        Copy
                    </AnimatedButtonWB>
                )}
            </Container>
        </>
    )
}

export default AppletCard
