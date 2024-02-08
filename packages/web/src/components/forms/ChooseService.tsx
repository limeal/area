import { useEffect, useRef, useState } from 'react'
import { toast } from 'react-toastify'
import styled from 'styled-components'

// Local
import SelectorContainer from '../applets/SelectorContainer'
import { IService } from '../../interfaces'
import {
    useGetServiceAuthorizationsQuery,
    usePostAuthorizationMutation,
} from '../../redux/api'
import { getAvatar } from '../../utils/more'

const Description = styled.p`
    overflow: hidden;
    max-height: 125px;
`

/* Defining the props that the component will receive. */
interface ChooseServiceProps {
    services: IService[]
    filter?: 'all' | 'action' | 'reaction'
    callback: (service: IService) => void
    goBack: any
}

/* A React component that is used to choose a service. */
const ChooseService = ({
    services,
    filter,
    callback,
    goBack,
}: ChooseServiceProps) => {
    const popup = useRef<Window | null>()
    const [tmpService, setTmpService] = useState<IService | null>(null)
    const [useAuthorize, { error: perror, isLoading: pisLoading, isSuccess }] =
        usePostAuthorizationMutation()
    const { data, isLoading, error } = useGetServiceAuthorizationsQuery(null, {
        refetchOnMountOrArgChange: true,
        refetchOnReconnect: true,
    })

    useEffect(() => {
        if (error && 'data' in error) {
            toast.error(JSON.stringify(error.data), {
                position: toast.POSITION.TOP_CENTER,
            })
        }
        if (perror && 'data' in perror) {
            toast.error(JSON.stringify(perror.data), {
                position: toast.POSITION.TOP_CENTER,
            })
        }
        if (isSuccess && data) {
            if (tmpService) callback(tmpService)
        }
        if (pisLoading || isSuccess) {
            toast.promise(
                new Promise((res, rej) => {
                    if (isSuccess) return res(0)
                    if (pisLoading) return rej()
                }),
                {
                    pending: 'Adding the service...',
                    success: 'Successfully added Service !',
                },
                {
                    position: toast.POSITION.TOP_CENTER,
                }
            )
        }
    }, [isSuccess, pisLoading, error, perror])

    const handleServiceClick = (service: IService) => {
        const redirectUri = 'http://localhost:8081'

        if (data && data[service.name]) {
            // Already authorized
            callback(service)
            return
        }

        if (service.authenticator === null) {
            toast.error('Service is not available for oauth2', {
                position: toast.POSITION.TOP_CENTER,
            })
            return
        }

        setTmpService(service)
        window.handleAuthorization = (data: any) =>
            useAuthorize({
                ...{
                    authenticator: service.authenticator?.name,
                    code: data.code,
                    redirect_uri: redirectUri,
                },
            })

        popup.current = window.open(
            service.authenticator?.authorization_uri +
                `&redirect_uri=${redirectUri}`,
            'Login with ' + service.name,
            'width=800,height=600'
        )
    }

    if (isLoading) return <div>Loading...</div>
    if (error) return <div>Error</div>

    return (
        <SelectorContainer
            title="Choose a service"
            hasSearchBar={true}
            items={services.filter((service) => {
                if (filter === 'all') return true
                if (filter === 'action')
                    return (
                        service.actions?.filter((a) => {
                            if (a.wip && process.env.NODE_ENV === 'production')
                                return false
                            return true
                        }).length > 0
                    )
                if (filter === 'reaction')
                    return (
                        service.reactions?.filter((a) => {
                            if (a.wip && process.env.NODE_ENV === 'production')
                                return false
                            return true
                        }).length > 0
                    )
                return false
            })}
            cardTheme={(service: IService) => {
                if (data && data[service.name])
                    return {
                        backgroundColor: '#fff',
                        cursor: 'pointer',
                    }
                return {
                    backgroundColor: '#2e2c2f',
                    cursor: 'help',
                    color: '#fff',
                }
            }}
            getCardElements={(service: IService) => (
                <>
                    <img
                        src={getAvatar(service)}
                        alt={service.name}
                        width="50px"
                        height="50px"
                    />
                    <h2>{service.name}</h2>
                    <Description>{service.description}</Description>
                </>
            )}
            onClickCard={(service) => handleServiceClick(service)}
            goBack={goBack}
        />
    )
}

export default ChooseService
