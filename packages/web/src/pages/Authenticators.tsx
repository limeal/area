import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { toast } from 'react-toastify'

import SelectorContainer from '../components/applets/SelectorContainer'
import { IAuthenticator, IMeta } from '../interfaces'
import { getAvatarM } from '../utils/more'

// API
import { useEffect } from 'react'
import { MdDoNotDisturbAlt } from 'react-icons/md'
import {
    useDeleteAuthorizationMutation,
    useGetAuthorizationsQuery,
} from '../redux/api'

/* Defining the props that the component will receive. */
interface ServicesProps {
    authenticators: IAuthenticator[]
}

/* A React component that is used to remove an authorization. */
const Services = ({ authenticators }: ServicesProps) => {
    const navigate = useNavigate()

    // Query
    const { data, error, refetch } = useGetAuthorizationsQuery(null, {
        refetchOnMountOrArgChange: true,
        refetchOnReconnect: true,
    })
    const [authorizations, setAuthorizations] = useState<any>([])
    const [meta, setMeta] = useState<any>([])

    const [useDelete, { error: perror, isSuccess: deleteSuccess }] =
        useDeleteAuthorizationMutation()

    const handleServiceClick = (auth: IAuthenticator) => {
        const metar = meta.find(
            (meta: IMeta) => meta.authenticator === auth.name
        )

        if (metar.applets === null) {
            useDelete(auth.name)
        } else {
            toast.error(
                'You cannot remove this authorization: Used by (' +
                    metar.applets +
                    ')',
                {
                    position: toast.POSITION.TOP_CENTER,
                }
            )
        }
    }

    useEffect(() => {
        if (data) {
            setAuthorizations(data.authorizations)
            setMeta(data.meta)
        }
    }, [data])

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
        if (deleteSuccess) {
            refetch()
            toast.success('Successfully remove authorization', {
                position: toast.POSITION.TOP_CENTER,
            })
        }
    }, [deleteSuccess, perror, error])

    return (
        <SelectorContainer
            title="Remove an authorization"
            hasSearchBar={true}
            items={authenticators.filter((auth: IAuthenticator) => {
                const authr = authorizations.find(
                    (authr: any) => authr.name === auth.name
                )
                if (!authr) return false
                if (authr.permanent) return false
                return true
            })}
            cardTheme={(auth: IAuthenticator) => {
                const metar = meta.find(
                    (meta: IMeta) => meta.authenticator === auth.name
                )

                if (metar.applets == null) {
                    return {
                        backgroundColor: 'white',
                        cursor: 'pointer',
                        color: '#000',
                    }
                } else {
                    // Unauthorized
                    return {
                        backgroundColor: '#e0e0e0',
                        cursor: 'not-allowed',
                        color: '#000',
                    }
                }
            }}
            getCardElements={(auth: IAuthenticator) => (
                <>
                    {meta.find(
                        (meta: IMeta) => meta.authenticator === auth.name
                    ).applets == null ? (
                        <img
                            src={getAvatarM(auth.more, auth.name)}
                            alt={auth.name}
                            width="50px"
                            height="50px"
                        />
                    ) : (
                        <MdDoNotDisturbAlt size="50px" color="red" />
                    )}
                    <h2>{auth.name}</h2>
                </>
            )}
            onClickCard={(auth) => handleServiceClick(auth)}
            goBack={() => navigate('/profile')}
        />
    )
}

export default Services
