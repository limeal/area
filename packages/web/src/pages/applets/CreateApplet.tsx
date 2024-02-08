import { motion } from 'framer-motion'
import { useEffect, useState } from 'react'
import { toast } from 'react-toastify'
import styled from 'styled-components'

// Interface
import { MdArrowRight, MdClose, MdLock } from 'react-icons/md'
import { IArea, IService } from '../../interfaces'

// Redux
import { useNavigate } from 'react-router-dom'
import {
    useDeleteNewAppletMutation,
    useGetNewAppletQuery,
    useSubmitNewAppletMutation,
} from '../../redux/api'

import AnimatedButtonHL from '../../components/animated/buttons/AnimatedButtonHL'
import AnimatedButtonWB from '../../components/animated/buttons/AnimatedButtonWB'
import FullContainer from '../../components/applets/FullContainer'
import ConfirmModal from '../../components/modals/ConfirmModal'
import { getAvatar, getColor } from '../../utils/more'
import SubmitComponent from '../../components/forms/SubmitApplet'
import { applyPretty } from '../../utils/string'

const AreaContainer = styled.div`
    display: flex;
    flex-direction: column;
    align-items: center;
`

const AreaBox = styled(motion.div)`
    display: flex;
    flex-direction: row;
    width: 100%;
    color: #fff;
    background: ${(props) => props.theme.background || '#222222'};
    border-radius: 2rem;
    align-items: center;
    justify-content: space-around;
    margin: 10px;
    width: 700px;
    height: 150px;
`

const AreaItem = styled(AreaBox)`
    cursor: ${(props) => (props.theme.active ? 'pointer' : 'not-allowed')};
    background-color: ${(props) => (props.theme.active ? '#fff' : '#999999')};
    color: ${(props) => (props.theme.active ? '#222222' : '#fff')};
    border: 5px solid ${(props) => (props.theme.active ? '#222222' : '#999999')};
`

/* Defining the props that will be passed to the AreaComponent. */
interface AreaComponentProps {
    services: IService[]
    areas: IArea[] | undefined
    type: string
    default: string | null
    active?: boolean
    refetchApplets: () => void
}

/* A React component that is used to display the areas of the applet. */
const AreaComponent = ({
    services,
    areas,
    type,
    active,
    default: defaultArea,
    refetchApplets,
}: AreaComponentProps) => {
    const navigate = useNavigate()
    const [openConfirm, setOpenConfirm] = useState(false)
    const [useDelete, { error: deleteError, isSuccess: successDelete }] =
        useDeleteNewAppletMutation()
    const [selected, setSelected] = useState<IArea | undefined>()
    const [selectedService, setSelectedService] = useState<IService>()

    const deleteComponent = (index: number) => {
        useDelete({
            type,
            number: index >= 0 ? index.toString() : '',
        })
    }

    useEffect(() => {
        if (active && defaultArea)
            navigate(
                `/create/set-${type}?service=${
                    defaultArea.split(';')[0]
                }&area=${defaultArea.split(';')[1]}`
            )
    }, [active])

    useEffect(() => {
        if (type == 'action' && areas && areas.length > 0)
            setSelectedService(services.find((s) => s.name == areas[0].service))
        setSelected(areas && areas.length > 0 ? areas[0] : undefined)
    }, [areas])

    useEffect(() => {
        if (selected)
            setSelectedService(services.find((s) => s.name == selected.service))
    }, [selected])

    useEffect(() => {
        if (deleteError) toast.error('Error deleting component')
        if (successDelete) {
            toast.success('Component deleted', {
                position: toast.POSITION.TOP_CENTER,
            })
            refetchApplets()
        }
    }, [successDelete, deleteError])

    return areas && areas.length > 0 ? (
        <AreaBox
            theme={{
                background: getColor(selectedService),
            }}
        >
            {type == 'action' ? (
                <>
                    <ConfirmModal
                        message="By deleting the action, you will also delete all reactions. Are you sure you want to delete it?"
                        callback={() => deleteComponent(-1)}
                        open={openConfirm}
                        close={() => setOpenConfirm(false)}
                    />
                    <img
                        src={getAvatar(
                            services.find((s) => s.name == areas[0].service)
                        )}
                        alt={areas[0].service}
                        width="50px"
                        height="50px"
                    />
                    {areas[0].service.charAt(0).toUpperCase() +
                        areas[0].service.slice(1)}
                    : {applyPretty(areas[0].name)}
                    <AnimatedButtonHL
                        color="#b9090b"
                        hlColor="#b9090b"
                        onClick={() => setOpenConfirm(true)}
                    >
                        Remove
                    </AnimatedButtonHL>
                </>
            ) : (
                <>
                    <div
                        style={{
                            display: 'flex',
                            flexDirection: 'column',
                            alignItems: 'center',
                            justifyContent: 'space-evenly',
                            height: '100%',
                        }}
                    >
                        <p
                            style={{
                                fontSize: '1.2rem',
                                fontWeight: 'bold',
                                color: '#fff',
                            }}
                        >
                            N°: {areas.findIndex((a) => a.id == selected?.id)}
                        </p>
                        <AnimatedButtonWB
                            bgColor="#b90b0b"
                            color="#fff"
                            style={{
                                padding: '0.6rem 0.8rem',
                                display: 'flex',
                                alignItems: 'center',
                            }}
                            onClick={() =>
                                deleteComponent(
                                    areas.findIndex((a) => a.id == selected?.id)
                                )
                            }
                        >
                            Remove <MdClose size="1rem" />
                        </AnimatedButtonWB>
                    </div>
                    <select
                        style={{
                            background: 'none',
                            border: 'none',
                            color: '#fff',
                            marginLeft: '10px',
                        }}
                        onChange={(e) =>
                            setSelected(
                                areas.find((a) => a.id == e.target.value)
                            )
                        }
                    >
                        {areas.map((area: IArea, index: number) => (
                            <option key={index} value={area.id}>
                                {area.service.charAt(0).toUpperCase() +
                                    area.service.slice(1)}
                                : {applyPretty(area.name)}
                            </option>
                        ))}
                    </select>
                    <div
                        style={{
                            display: 'flex',
                            flexDirection: 'column',
                            height: '100%',
                            justifyContent: 'space-evenly',
                        }}
                    >
                        <AnimatedButtonWB
                            bgColor="#222222"
                            color="#fff"
                            onClick={() =>
                                active && navigate(`/create/add-${type}`)
                            }
                        >
                            Add
                        </AnimatedButtonWB>
                        <AnimatedButtonHL
                            color="#b9090b"
                            hlColor="#b9090b"
                            onClick={() => deleteComponent(-1)}
                        >
                            Remove All
                        </AnimatedButtonHL>
                    </div>
                </>
            )}
        </AreaBox>
    ) : (
        <AreaItem
            whileHover={{
                scale: 1.1,
                transition: { duration: 0.2 },
            }}
            whileTap={{
                scale: 0.95,
                transition: { duration: 0.3 },
            }}
            theme={{ active }}
            onClick={() => active && navigate(`/create/add-${type}`)}
        >
            <h2>{type == 'action' ? 'If this' : 'Then That'}</h2>
            {active ? (
                <AnimatedButtonWB bgColor="#222222" color="#fff">
                    Add
                </AnimatedButtonWB>
            ) : (
                <MdLock size="2rem" />
            )}
        </AreaItem>
    )
}

/* Defining an interface called CreateAppletProps. It is saying that the CreateAppletProps interface
has a property called services that is an array of IService objects. */
interface CreateAppletProps {
    services: IService[]
}

/**
 * It renders a page with two components that allow the user to select an action and a reaction, and
 * then a button that takes the user to the next page
 * @param {CreateAppletProps}  - `services` - an array of services that are available to the user
 */
const CreateApplet = ({ services }: CreateAppletProps) => {
    const {
        data,
        error,
        refetch: refetchApplets,
    } = useGetNewAppletQuery(
        {
            field: null,
        },
        {
            refetchOnMountOrArgChange: true,
        }
    )
    const [useSubmit, { error: submitError, isSuccess: submitSucess }] =
        useSubmitNewAppletMutation()
    const [page, setPage] = useState(0)
    const navigate = useNavigate()

    useEffect(() => {
        if (submitSucess) {
            toast.success('Applet created!', {
                position: toast.POSITION.TOP_CENTER,
            })
            navigate('/applets')
        }
        if (error && 'data' in error && error.data)
            toast.error(JSON.stringify(error.data), {
                position: toast.POSITION.TOP_CENTER,
            })
        if (submitError && 'data' in submitError && submitError.data)
            toast.error(JSON.stringify(submitError.data), {
                position: toast.POSITION.TOP_CENTER,
            })
    }, [submitSucess, data, error, submitError])

    if (page == 1)
        return <SubmitComponent setPage={setPage} useSubmit={useSubmit} />

    return (
        <FullContainer
            title="Create (1/2)"
            backIcon={<MdClose fontSize={'1.5rem'} />}
            goBack={() => navigate('/applets')}
        >
            <AreaContainer>
                <AreaComponent
                    services={services}
                    areas={
                        data?.action !== undefined ? [data?.action] : undefined
                    }
                    active={true}
                    default={null}
                    type="action"
                    refetchApplets={refetchApplets}
                />
                <AreaComponent
                    services={services}
                    areas={data?.reactions}
                    active={data?.action != null}
                    default={null}
                    type="reaction"
                    refetchApplets={refetchApplets}
                />
            </AreaContainer>
            {data?.action && data?.reactions && data?.reactions?.length > 0 && (
                <AnimatedButtonWB
                    bgColor="#222222"
                    color="#fff"
                    onClick={() => setPage(1)}
                    style={{
                        marginTop: '20px',
                        display: 'flex',
                        flexDirection: 'row',
                        alignItems: 'center',
                        fontSize: '1.5rem',
                    }}
                >
                    Next Step
                    <MdArrowRight size="3rem" />
                </AnimatedButtonWB>
            )}
        </FullContainer>
    )
}

export default CreateApplet
