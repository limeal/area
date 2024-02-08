import { motion } from 'framer-motion'
import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { toast } from 'react-toastify'
import styled from 'styled-components'

import { ArcElement, Chart as ChartJS, Legend, Tooltip } from 'chart.js'
import { Doughnut } from 'react-chartjs-2'
import { MdClose, MdDelete, MdDraw, MdLogout, MdSettings } from 'react-icons/md'

import ConfirmModal from '../components/modals/ConfirmModal'
import { IApplet, IArea } from '../interfaces'
import {
    useDeleteProfileMutation,
    useFetchServiceApiEndpointMutation,
    useGetAppletReactionsMutationMutation,
    useGetAppletsQuery,
    useGetAvatarQuery,
    useGetProfileQuery,
    useModifyProfileMutation,
    useUpdateAvatarMutation,
} from '../redux/api'

/* Registering the ArcElement, Tooltip, and Legend components with the ChartJS library. */
ChartJS.register(ArcElement, Tooltip, Legend)

const Container = styled.div`
    height: 100%;
    font-size: 1.5rem;
    background-color: #222222;
    display: flex;
    flex-direction: column;
    align-items: center;
`

const Header = styled.div`
    display: flex;
    flex-direction: row;
    align-items: center;
    color: white;
    background-color: #b9090b;
    justify-content: space-between;
    flex: 2;
    width: 100%;
`

const HeaderLeft = styled.div`
    display: flex;
    flex-direction: row;
    align-items: center;
`

const Avatar = styled(motion.div)`
    border-radius: 50%;
    background-color: white;
    width: 120px;
    height: 120px;
    overflow: hidden;
    background-size: cover;
    background-position: center;
    cursor: pointer;
    margin: 20px;
`

const UserCard = styled.div`
    display: flex;
    flex-direction: column;
    align-items: flex-start;
`

const Username = styled.p`
    font-size: 2rem;
    font-weight: bold;
    padding-bottom: 10px;
    border-bottom: 1px solid white;
    display: flex;
    align-items: center;
`

const EditUserButton = styled(motion.button)`
    outline: none;
    border: none;
    margin-left: 10px;
    cursor: pointer;
    background: transparent;
    color: white;
`

const Logout = styled(motion.button)`
    cursor: pointer;
    margin-right: 77px;
    background: transparent;
    border: none;
    color: white;
`

const Email = styled.p`
    padding-top: 10px;
    font-size: 1.2rem;
`

const HeaderRight = styled.div``

const BodyWrapper = styled.div`
    max-width: 1600px;
    width: 100%;
    flex: 5;
`

const Body = styled.div`
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    padding: 20px;
    align-items: center;
    height: 100%;
`

const WebhookContainer = styled.div`
    margin: 20px;
    height: 100%;
    display: flex;
    align-items: center;
`

const WebhookForm = styled.form`
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: space-evenly;
    height: 300px;
    background-color: #4f518c;
    padding: 20px;
    border-radius: 2rem;
`

const WebhookLogo = styled.img``

const WebhookSelector = styled.select`
    width: 100%;
    font-size: 1.2rem;
    padding: 10px;
`

const WebhookMessageInput = styled.input`
    width: 100%;
    font-size: 1.2rem;
    padding: 10px;
    outline: none;
`

const WebhookSubmitButton = styled.button`
    font-size: 1.2rem;
    padding: 10px;
    outline: none;
    border: none;
    cursor: pointer;
`

const BodyMain = styled.div`
    margin: 20px;
    height: 100%;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    color: white;
`

const SpecialButton = styled(motion.button)`
    display: flex;
    flex-direction: row;
    align-items: center
    color: white;
    font-size: 1.5rem;
    padding: 1.6rem 4.4rem;
    border-radius: 2rem;
    background-color: ${(props) => props.theme.backgroundColor || '#b9090b'};
    cursor: pointer;
    border: none;
    outline: none;
    margin-bottom: 20px;
    color: white;
    width: 100%;
`

/* Creating an interface called ProfileProps. It is saying that the ProfileProps interface has a
property called updateToken that is of type any. */
interface ProfileProps {
    updateToken: any
}

/**
 * It takes an array of applets and an array of reactions, and returns an object that can be used to
 * create a chart
 * @param {IApplet[]} applets - IApplet[] - an array of applets
 * @param {IArea[]} reactions - IArea[] - this is the array of reactions that you want to get the data
 * from.
 * @returns An object with the labels and datasets for the chart.
 */
const GetData = (applets: IApplet[], reactions: IArea[]) => {
    const services = new Set()
    const serviceUsage = new Map()

    applets.forEach(async (applet) => {
        const actionService = applet.action.split(';')[0]
        services.add(actionService)
        serviceUsage.set(
            actionService,
            (serviceUsage.get(actionService) || 0) + 1
        )
    })

    reactions.forEach(async (reaction) => {
        const reactionService = reaction.service
        services.add(reactionService)
        serviceUsage.set(
            reactionService,
            (serviceUsage.get(reactionService) || 0) + 1
        )
    })

    return {
        labels: Array.from(services),
        datasets: [
            {
                label: 'Service Usage',
                data: Array.from(services).map((service) =>
                    serviceUsage.get(service)
                ),
                backgroundColor: [
                    'rgba(255, 99, 132, 0.2)',
                    'rgba(54, 162, 235, 0.2)',
                    'rgba(255, 206, 86, 0.2)',
                    'rgba(75, 192, 192, 0.2)',
                    'rgba(153, 102, 255, 0.2)',
                    'rgba(255, 159, 64, 0.2)',
                ],
                borderColor: [
                    'rgba(255, 99, 132, 1)',
                    'rgba(54, 162, 235, 1)',
                    'rgba(255, 206, 86, 1)',
                    'rgba(75, 192, 192, 1)',
                    'rgba(153, 102, 255, 1)',
                    'rgba(255, 159, 64, 1)',
                ],
                borderWidth: 1,
            },
        ],
    }
}

/* A React component that displays the user's profile. */
const Profile = ({ updateToken }: ProfileProps) => {
    const [nusername, setNUsername] = useState<string>('')
    const [email, setEmail] = useState<string>('')
    const [mode, setMode] = useState<'normal' | 'editing'>('normal')
    const [username, setUsername] = useState<string>('')
    const [reactions, setReactions] = useState<IArea[]>([])
    const [useGetReactions] = useGetAppletReactionsMutationMutation()
    const [openConfirm, setOpenConfirm] = useState<boolean>(false)

    // Query
    const { data, error } = useGetProfileQuery(null, {
        refetchOnMountOrArgChange: true,
        refetchOnReconnect: true,
    })
    const [
        useModifyProfile,
        {
            error: modifyError,
            isLoading: modifyLoading,
            isSuccess: modifySuccess,
        },
    ] = useModifyProfileMutation()
    const [
        useDeleteProfile,
        {
            error: deleteError,
            isLoading: deleteLoading,
            isSuccess: deleteSuccess,
        },
    ] = useDeleteProfileMutation()
    const { data: applets } = useGetAppletsQuery(null, {
        refetchOnMountOrArgChange: true,
        refetchOnReconnect: true,
    })
    const [useApi] = useFetchServiceApiEndpointMutation()

    const [avatarURI, setAvatarURI] = useState<string>('')
    const { data: avatar, refetch: refetchAvatar } = useGetAvatarQuery(null, {
        refetchOnMountOrArgChange: true,
        refetchOnReconnect: true,
    })
    const [useUpdateAvatar] = useUpdateAvatarMutation()

    const navigate = useNavigate()

    const handleLogout = () => {
        //useLogout(null)
        navigate('/')
        updateToken('')
        toast.success('Logged out successfully', {
            position: toast.POSITION.TOP_CENTER,
        })
    }

    const handleAvatar = (e: any) => {
        const payload = new FormData()
        payload.append('avatar', e.target.files[0])

        useUpdateAvatar(payload)
        refetchAvatar()
    }

    const handleSubmitWebhook = (e: any) => {
        e.preventDefault()

        // Get the form data out of state
        const data = Object.fromEntries([...new FormData(e.target).entries()])

        useApi({
            service: 'webhook',
            endpoint: '/' + data.applet,
            method: 'POST',
            body: data.message,
        })
    }

    useEffect(() => {
        if (applets) {
            for (const applet of applets) {
                useGetReactions(applet.id)
                    .unwrap()
                    .then((reactions: IArea[]) => {
                        setReactions((prev) => [...prev, ...reactions])
                    })
            }
        }
    }, [applets])

    useEffect(() => {
        if (data && 'email' in data && 'username' in data) {
            setEmail(data.email)
            setUsername(data.username)
        }
    }, [data])

    useEffect(() => {
        setAvatarURI(avatar || '')
        console.log(avatar)
    }, [avatar])

    useEffect(() => {
        if (modifyLoading || modifySuccess) {
            toast.promise(
                new Promise((res, rej) => {
                    if (modifySuccess) return res(0)
                    if (modifyLoading) return rej()
                }),
                {
                    pending: 'Modifying username...',
                    success: 'Username changed !',
                },
                {
                    position: toast.POSITION.TOP_CENTER,
                }
            )
        }
        if (deleteLoading || deleteSuccess) {
            toast.promise(
                new Promise((res, rej) => {
                    if (deleteSuccess) return res(0)
                    if (deleteLoading) return rej()
                }),
                {
                    pending: 'Deleting account...',
                    success: 'Account deleted successfully !',
                },
                {
                    position: toast.POSITION.TOP_CENTER,
                }
            )
            if (deleteSuccess) {
                navigate('/')
                updateToken('')
            }
        }
    }, [modifyLoading, deleteLoading, modifySuccess, deleteSuccess])

    useEffect(() => {
        if (error && 'data' in error)
            toast.error(JSON.stringify(error.data), {
                position: toast.POSITION.TOP_CENTER,
            })
        if (modifyError && 'data' in modifyError)
            toast.error(JSON.stringify(modifyError.data), {
                position: toast.POSITION.TOP_CENTER,
            })
        if (deleteError && 'data' in deleteError)
            toast.error(JSON.stringify(deleteError.data), {
                position: toast.POSITION.TOP_CENTER,
            })
    }, [error, modifyError, deleteError])

    return (
        <Container>
            <ConfirmModal
                message="Deleting your account will remove all your data from the server. This action cannot be undone."
                open={openConfirm}
                callback={() => useDeleteProfile(null)}
                close={() => setOpenConfirm(false)}
            />
            <Header>
                <HeaderLeft>
                    <Avatar
                        style={{
                            backgroundImage: `url(${avatarURI})`,
                        }}
                        whileHover={{
                            scale: 1.1,
                            transition: { duration: 0.2 },
                        }}
                    >
                        <input
                            type="file"
                            accept="image/*"
                            style={{
                                opacity: '0',
                                height: '100%',
                                cursor: 'pointer',
                            }}
                            onChange={handleAvatar}
                        />
                    </Avatar>
                    <UserCard>
                        <Username>
                            {mode === 'normal' ? (
                                <>
                                    {username}
                                    <EditUserButton
                                        whileHover={{
                                            scale: 1.1,
                                            transition: { duration: 0.2 },
                                        }}
                                        onClick={() => {
                                            setNUsername(username)
                                            setMode('editing')
                                        }}
                                    >
                                        <MdDraw />
                                    </EditUserButton>
                                </>
                            ) : (
                                <form
                                    onSubmit={(e) => {
                                        e.preventDefault()
                                        useModifyProfile({
                                            username: nusername,
                                        })
                                        setUsername(nusername)
                                        setMode('normal')
                                    }}
                                >
                                    <input
                                        value={nusername}
                                        onChange={(e) =>
                                            setNUsername(e.target.value)
                                        }
                                        style={{
                                            maxWidth: '200px',
                                            overflow: 'hidden',
                                            border: '1px solid rgba(0, 0, 0, 0.5)',
                                            outline: 'none',
                                            color: '#222222',
                                            caretColor: '#222222',
                                        }}
                                    />
                                    <button
                                        type="button"
                                        style={{
                                            background: 'red',
                                            textAlign: 'center',
                                            border: 'none',
                                            color: 'white',
                                        }}
                                        onClick={() => setMode('normal')}
                                    >
                                        <MdClose />
                                    </button>
                                </form>
                            )}
                        </Username>
                        <Email>{email}</Email>
                    </UserCard>
                </HeaderLeft>
                <HeaderRight>
                    <Logout
                        whileHover={{
                            scale: 1.1,
                            transition: { duration: 0.2 },
                        }}
                        onClick={() => handleLogout()}
                    >
                        <MdLogout
                            style={{
                                fontSize: '5rem',
                            }}
                        />
                    </Logout>
                </HeaderRight>
            </Header>
            <BodyWrapper>
                <Body>
                    <WebhookContainer>
                        <WebhookForm onSubmit={handleSubmitWebhook}>
                            <WebhookLogo
                                src="https://cdn.worldvectorlogo.com/logos/webhook-1.svg"
                                alt="webhook"
                                width="100px"
                                height="100px"
                            />
                            <WebhookSelector name="applet">
                                {applets &&
                                    applets
                                        .filter((applet) =>
                                            applet.action.startsWith('webhook')
                                        )
                                        .map((applet) => {
                                            return (
                                                <option
                                                    key={applet.id}
                                                    value={applet.name}
                                                >
                                                    {applet.name}
                                                </option>
                                            )
                                        })}
                            </WebhookSelector>
                            <WebhookMessageInput
                                name="message"
                                placeholder="Webhook message"
                            />
                            <WebhookSubmitButton type="submit">
                                Submit
                            </WebhookSubmitButton>
                        </WebhookForm>
                    </WebhookContainer>
                    <BodyMain>
                        <SpecialButton
                            whileHover={{
                                scale: 1.1,
                                transition: { duration: 0.2 },
                            }}
                            onClick={() => navigate('/authorizations')}
                            theme={{
                                backgroundColor: '#499f68',
                            }}
                        >
                            <MdSettings />
                            Authorizations
                        </SpecialButton>
                        <SpecialButton
                            whileHover={{
                                scale: 1.1,
                                transition: { duration: 0.2 },
                            }}
                            onClick={() => setOpenConfirm(true)}
                        >
                            <MdDelete />
                            Delete Account
                        </SpecialButton>
                        <p>Warning: This action is irreversible</p>
                    </BodyMain>
                    <div>
                        {applets && applets.length > 0 ? (
                            <Doughnut
                                style={{
                                    maxWidth: '500px',
                                    maxHeight: '500px',
                                }}
                                data={GetData(applets, reactions)}
                            />
                        ) : (
                            <p
                                style={{
                                    textAlign: 'center',
                                    color: 'white',
                                }}
                            >
                                No applets
                            </p>
                        )}
                    </div>
                </Body>
            </BodyWrapper>
        </Container>
    )
}

export default Profile
