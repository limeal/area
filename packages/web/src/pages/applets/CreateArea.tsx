import { useEffect, useState } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'

import { toast } from 'react-toastify'
import { IAction, IService } from '../../interfaces'
import {
    useAddStateToNewAppletMutation,
    useGetNewAppletQuery,
} from '../../redux/api'
import ChooseArea from '../../components/forms/ChooseArea'
import ChooseAreaSettings from '../../components/forms/ChooseAreaSettings'
import ChooseService from '../../components/forms/ChooseService'

/* Defining the props that the component will receive. */
interface CreateAreaProps {
    services: IService[]
    type: 'action' | 'reaction'
}

/**
 * It's a React component that allows you to create an applet
 * @param {CreateAreaProps}  - `services` - an array of services that the user has connected to their
 * account.
 */
const CreateArea = ({ services, type }: CreateAreaProps) => {
    const [searchParams] = useSearchParams()
    const { error } = useGetNewAppletQuery(null, {
        refetchOnMountOrArgChange: true,
    })
    const [
        useAddState,
        { error: stError, isSuccess, isLoading: updateLoading },
    ] = useAddStateToNewAppletMutation()
    const navigate = useNavigate()
    const [selectedService, setSelectedService] = useState<IService | null>(
        services.find(
            (s) =>
                searchParams.get('service') &&
                s.name === searchParams.get('service')
        ) || null
    )
    const [selectedArea, setSelectedArea] = useState<IAction | null>(() => {
        const item = searchParams.get('item')
        if (!selectedService) return null
        if (!item) return null
        if (type === 'action')
            return selectedService.actions.find((a) => a.name === item) || null
        return selectedService.reactions.find((a) => a.name === item) || null
    })

    const AddStateWithoutSetting = (a: IAction) => {
        useAddState({
            ...{
                service: selectedService?.name,
                area_type: type,
                area_item: a.name,
            },
        })
    }

    const AddState = (settings: any) => {
        useAddState({
            ...{
                service: selectedService?.name,
                area_type: type,
                area_item: selectedArea?.name,
                area_settings: settings,
            },
        })
    }

    useEffect(() => {
        if (isSuccess) navigate('/create')
    }, [isSuccess])

    useEffect(() => {
        if (error && 'data' in error)
            toast.error(JSON.stringify(error.data), {
                position: toast.POSITION.TOP_CENTER,
            })
        if (stError && 'data' in stError && stError.status != 406)
            toast.error(JSON.stringify(stError.data), {
                position: toast.POSITION.TOP_CENTER,
            })
    }, [error, stError])

    if (!selectedService)
        return (
            <ChooseService
                services={services}
                goBack={() => navigate('/create')}
                filter={type}
                callback={(s) => setSelectedService(s)}
            />
        )

    if (!selectedArea)
        return (
            <ChooseArea
                type={type}
                goBack={() => setSelectedService(null)}
                selectedService={selectedService}
                elements={
                    type === 'action'
                        ? selectedService.actions
                        : selectedService.reactions
                }
                callback={(a) =>
                    a && a.store && Object.keys(a.store).length > 0
                        ? setSelectedArea(a)
                        : AddStateWithoutSetting(a)
                }
            />
        )
    return (
        <ChooseAreaSettings
            type={type}
            services={services}
            selectedService={selectedService}
            no_working_fields={
                stError && 'data' in stError
                    ? JSON.stringify(stError.data)
                    : null
            }
            goBack={() => setSelectedArea(null)}
            area={selectedArea}
            callback={(settings) => AddState(settings)}
            loading={updateLoading}
        />
    )
}

export default CreateArea
