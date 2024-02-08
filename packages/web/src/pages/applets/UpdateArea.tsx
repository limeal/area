/* import { useState, useEffect } from 'react'
import { useGetAppletQuery, useModifyAppletMutation } from '../../redux/api'
import { IService, IAction } from '../../interfaces'
import { useNavigate, useParams, useSearchParams } from 'react-router-dom'
import ChooseService from './forms/ChooseService'
import ChooseAreaSettings from './forms/ChooseAreaSettings'
import ChooseArea from './forms/ChooseArea'
import { toast } from 'react-toastify'
import { POSITION } from 'react-toastify/dist/utils'

interface UpdateAreaProps {
    services: IService[]
    type: 'action' | 'reaction'
}

const UpdateArea = ({ services, type }: UpdateAreaProps) => {
    const { id } = useParams()
    const navigate = useNavigate()

    const { error: appletError } = useGetAppletQuery(id, {
        refetchOnMountOrArgChange: true,
    })

    const [
        useUpdateActivity,
        { error: updateError, isSuccess, isLoading: modifyLoading },
    ] = useModifyAppletMutation()

    const [selectedService, setSelectedService] = useState<IService | null>(
        null
    )
    const [selectedArea, setSelectedArea] = useState<IAction | null>(null)

    const ModifyStateWithoutSetting = (a: IAction) => {
        useUpdateActivity({
            id,
            ...{
                service: selectedService?.name,
                area_type: type,
                area_item: a.name,
            },
        })
    }

    const ModifyState = (settings: any) => {
        useUpdateActivity({
            id,
            ...{
                service: selectedService?.name,
                area_type: type,
                area_item: selectedArea?.name,
                area_settings: settings,
            },
        })
    }

    useEffect(() => {
        if (isSuccess) navigate('/applets/' + id)
    }, [isSuccess])

    useEffect(() => {
        if (appletError && 'data' in appletError)
            toast.error(JSON.stringify(appletError.data), {
                position: toast.POSITION.TOP_CENTER,
            })
        if (updateError && 'data' in updateError)
            toast.error(JSON.stringify(updateError.data), {
                position: toast.POSITION.TOP_CENTER,
            })
    }, [appletError, updateError])

    if (!selectedService)
        return (
            <ChooseService
                services={services}
                goBack={() => navigate('/applets')}
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
                    type == 'action'
                        ? selectedService.actions
                        : selectedService.reactions
                }
                callback={(a) =>
                    a && a.store && Object.keys(a.store).length > 0
                        ? setSelectedArea(a)
                        : ModifyStateWithoutSetting(a)
                }
            />
        )
    return (
        <ChooseAreaSettings
            type={type}
            services={services}
            selectedService={selectedService}
            no_working_fields={
                updateError && 'data' in updateError
                    ? JSON.stringify(updateError.data)
                    : null
            }
            goBack={() => setSelectedArea(null)}
            area={selectedArea}
            callback={(settings) => ModifyState(settings)}
            loading={modifyLoading}
        />
    )
}

export default UpdateArea
 */

export {}
