import { useEffect, useRef, useState } from 'react'
import { FaTrash } from 'react-icons/fa'
import { MdUpdate } from 'react-icons/md'
import { useNavigate, useParams } from 'react-router-dom'
import { toast } from 'react-toastify'
import styled from 'styled-components'

import AnimatedSwitch from '../../components/animated/AnimatedSwitch'
import AnimatedButtonWB from '../../components/animated/buttons/AnimatedButtonWB'
import { IApplet, IService } from '../../interfaces'
import {
    useDeleteAppletMutation,
    useGetAppletQuery,
    useGetAppletReactionsQuery,
    useStartAppletMutation,
    useStopAppletMutation,
    useUpdateAppletActivityMutation,
} from '../../redux/api'
import { getAvatar, getColor } from '../../utils/more'

const Container = styled.div`
    display: flex;
    flex-direction: column;
    background: ${(props) => props.theme.background || '#222222'};
    height: 100%;
`

const AppletTerminalContainer = styled.div`
    display: flex;
    flex-direction: column;
    justify-content: flex-end;
    height: 400px;
    width: 800px;
    background-color: #222222;
    box-shadow: 0 0 10px 0 rgba(0, 0, 0, 0.2);
    border-radius: 10px;
    border: 1rem solid #fff;
    padding: 10px;
    color: #fff;
`

const Header = styled.div`
    display: flex;
    flex-direction: column;
    align-items: start;
    width: 100%;
    padding: 15px 0;
`

const Body = styled.div`
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: space-around;
    height: 100%;
    width: 100%;
`

const Title = styled.h2`
    font-size: 3rem;
    color: #000;
`

const Description = styled.p`
    font-size: 0.8rem;
    color: #000;
`

const ElementContainer = styled.div`
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: space-around;
`

const ActionContainer = styled.div`
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: space-around;
    height: 400px;
`

/**
 * It's a React component that displays a set of buttons and a switch to control an applet
 * @param  - `id` is the id of the applet
 * @returns A function that returns a component
 */
const AppletActions = ({
    id,
    applet,
    refetch,
}: {
    id: string
    applet: IApplet
    refetch: any
}) => {
    const [
        useUpdateActivity,
        { isLoading: isUpdating, isSuccess: successSwitch },
    ] = useUpdateAppletActivityMutation()
    const [
        useStartApplet,
        {
            error: ErrorStarting,
            isLoading: appletStarting,
            isSuccess: appletStarted,
        },
    ] = useStartAppletMutation()
    const [
        useStopApplet,
        {
            error: ErrorStopping,
            isLoading: appletStopping,
            isSuccess: appletStopped,
        },
    ] = useStopAppletMutation()
    const [useDelete, { error: deleteError, isSuccess: DeleteSuccess }] =
        useDeleteAppletMutation()
    const navigate = useNavigate()
    const [active, setActive] = useState(applet.active || false)

    useEffect(() => {
        if (appletStarting || appletStarted) {
            toast.promise(
                new Promise((res, rej) => {
                    if (appletStarted) return res(0)
                    if (appletStarting) return rej()
                }),
                {
                    pending: 'Applet starting...',
                    success: 'Applet started successfully !',
                },
                {
                    position: toast.POSITION.TOP_CENTER,
                }
            )
        } else if (appletStopping || appletStopped) {
            toast.promise(
                new Promise((res, rej) => {
                    if (appletStopped) return res(0)
                    if (appletStopping) return rej()
                }),
                {
                    pending: 'Applet stopping...',
                    success: 'Applet stopping successfully !',
                },
                {
                    position: toast.POSITION.TOP_CENTER,
                }
            )
        } else if (isUpdating || successSwitch) {
            toast.promise(
                new Promise((res, rej) => {
                    if (successSwitch) return res(0)
                    if (isUpdating) return rej()
                }),
                {
                    pending: 'Applet updating...',
                    success:
                        'Applet is now ' + (active ? 'inactive' : 'active'),
                },
                {
                    position: toast.POSITION.TOP_CENTER,
                }
            )
        }

        if (appletStarted) refetch()
        else if (appletStopped) refetch()
    }, [
        appletStarted,
        appletStopped,
        isUpdating,
        successSwitch,
        appletStarting,
        appletStopping,
    ])

    useEffect(() => {
        if (DeleteSuccess) {
            toast.success('Applet deleted successfully', {
                position: toast.POSITION.TOP_CENTER,
            })
            navigate('/applets')
        }
    }, [DeleteSuccess])

    useEffect(() => {
        if (deleteError && 'data' in deleteError && deleteError.data)
            toast.error(JSON.stringify(deleteError.data), {
                position: toast.POSITION.TOP_CENTER,
            })
        if (ErrorStarting && 'data' in ErrorStarting && ErrorStarting.data)
            toast.error(JSON.stringify(ErrorStarting.data), {
                position: toast.POSITION.TOP_CENTER,
            })
        if (ErrorStopping && 'data' in ErrorStopping && ErrorStopping.data)
            toast.error(JSON.stringify(ErrorStopping.data), {
                position: toast.POSITION.TOP_CENTER,
            })
    }, [ErrorStarting, ErrorStopping, deleteError])

    useEffect(() => {
        if (applet) setActive(applet.active || false)
    }, [applet])

    useEffect(() => {
        if (successSwitch) setActive(!active)
    }, [successSwitch])

    return (
        <ActionContainer>
            {applet.status === 'stopped' ? (
                <AnimatedButtonWB
                    color="#fff"
                    bgColor="blue"
                    onClick={() => useStartApplet(id)}
                >
                    <MdUpdate />
                    Start
                </AnimatedButtonWB>
            ) : (
                <AnimatedButtonWB
                    color="#fff"
                    bgColor="blue"
                    onClick={() => useStopApplet(id)}
                >
                    <MdUpdate />
                    Stop
                </AnimatedButtonWB>
            )}
            {/* TODO: Patch update
            <AnimatedButtonWB
                color="#fff"
                bgColor="blue"
                onClick={() => navigate(`/applets/${id}/update`)}
            >
                <MdUpdate />
                Update
            </AnimatedButtonWB> */}
            <AnimatedButtonWB
                color="#fff"
                bgColor="#b90b90"
                onClick={() => useDelete(id)}
            >
                <FaTrash />
                Delete
            </AnimatedButtonWB>
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
                onClick={() =>
                    useUpdateActivity({
                        id,
                        active: !active,
                    })
                }
            />
        </ActionContainer>
    )
}

/**
 * It creates a websocket connection to the server, and then displays the last 10 lines of the log file
 * @param  - `id` - The ID of the applet to get logs for.
 * @returns A React component that renders a terminal with the last 10 logs of the applet.
 */
const AppletTerminal = ({ id }: { id: string }) => {
    const [logs, setLogs] = useState<Array<string>>([])
    const socket = useRef<WebSocket>(
        new WebSocket(`ws://localhost:8080/logs/${id}`)
    )

    socket.current.addEventListener('message', (e) => {
        const newLogs = e.data.split('\n')
        setLogs([...logs, ...newLogs])
    })

    return (
        <AppletTerminalContainer>
            {logs.slice(Math.max(logs.length - 10, 0)).map((log, index) => (
                <p key={index}>{log}</p>
            ))}
        </AppletTerminalContainer>
    )
}

/* A React component that is using the useGetAppletQuery hook to fetch the applet data from the server. */
const Applet = ({ services }: { services: IService[] }) => {
    const { id } = useParams()
    const {
        data: applet,
        error: appletError,
        isLoading: appletLoading,
        refetch,
    } = useGetAppletQuery(id, {
        refetchOnMountOrArgChange: true,
    })
    const [actionService, setActionService] = useState(
        services.find((s) => s.name === applet?.action.split(';')[0]) || null
    )
    const { data: reactions, error: reactionsError } =
        useGetAppletReactionsQuery(id, {
            refetchOnMountOrArgChange: true,
        })

    useEffect(() => {
        if (applet) {
            setActionService(
                services.find((s) => s.name === applet?.action.split(';')[0]) ||
                    null
            )
        }
        if (appletError && 'data' in appletError)
            toast.error(JSON.stringify(appletError.data), {
                position: toast.POSITION.TOP_CENTER,
            })

        if (reactionsError && 'data' in reactionsError)
            toast.error(JSON.stringify(reactionsError.data), {
                position: toast.POSITION.TOP_CENTER,
            })
    }, [applet, appletError, reactionsError])

    if (appletLoading || !applet) return <p>Loading...</p>

    return (
        <Container
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
            <Header>
                <ElementContainer>
                    <img
                        src={getAvatar(actionService)}
                        alt={actionService?.name || 'test'}
                        style={{ width: '50px', height: '50px' }}
                    />
                    {reactions?.map((reaction, index: number) => {
                        const reactionService = services.find(
                            (s) => s.name === reaction.service
                        )
                        return (
                            <img
                                key={index}
                                src={getAvatar(reactionService)}
                                alt={reactionService?.name || 'test'}
                                style={{ width: '50px', height: '50px' }}
                            />
                        )
                    })}
                </ElementContainer>
                <Title>Name: {applet.name}</Title>
                <Description>Description: {applet.description}</Description>
            </Header>
            <Body>
                <AppletTerminal id={id || ''} />
                <AppletActions
                    id={id || ''}
                    applet={applet}
                    refetch={refetch}
                />
            </Body>
        </Container>
    )
}

export default Applet
