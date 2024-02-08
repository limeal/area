/* import styled from 'styled-components'
import { motion } from 'framer-motion'
import { useParams } from 'react-router-dom'
import { useNavigate } from 'react-router-dom'
import { MdClose } from 'react-icons/md'

// Interface
import { IService } from '../../interfaces'

// Redux
import { useGetAppletQuery } from '../../redux/api'

// Components
import AnimatedButtonWB from '../../components/animated/buttons/AnimatedButtonWB'
import FullContainer from '../../components/applets/FullContainer'
import { getColor, getAvatar } from '../../utils/more'

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
    background-color: ${(props) => props.theme.bgColor || '#222222'};
    border-radius: 2rem;
    align-items: center;
    justify-content: space-around;
    margin: 10px;
    width: 700px;
    height: 150px;
`

const ButtonContainer = styled.div`
    display: flex;
    flex-direction: column;
    justify-content: space-evenly;
    height: 100%;
`

interface AreaComponentProps {
    services: IService[]
    area: any
    type: string
    id: string
    refetch: () => void
}

const AreaComponent = ({ services, area, type, id }: AreaComponentProps) => {
    const navigate = useNavigate()
    const service = services.find((s) => s.name == area?.service)

    if (!area) return <p>Error</p>

    return (
        <AreaBox
            theme={{
                bgColor: getColor(service),
            }}
        >
            <img
                src={getAvatar(service)}
                alt={service?.name}
                width="50px"
                height="50px"
            />
            {area.service}: {area.name}
            <ButtonContainer>
                <AnimatedButtonWB
                    color="#fff"
                    bgColor="#222222"
                    onClick={(e: any) =>
                        navigate(`/applets/${id}/update/set-${type}`)
                    }
                >
                    Update
                </AnimatedButtonWB>
            </ButtonContainer>
        </AreaBox>
    )
}

interface UpdateAppletProps {
    services: IService[]
}

const UpdateApplet = ({ services }: UpdateAppletProps) => {
    const { id } = useParams()
    const {
        data: applet,
        error: appletError,
        isLoading: appletLoading,
        refetch,
    } = useGetAppletQuery(id)
    const navigate = useNavigate()

    if (appletError || !applet) return <div>Error...</div>

    return (
        <FullContainer
            title="Update Applet"
            backIcon={<MdClose fontSize={'1.5rem'} />}
            goBack={() => navigate('/applets/' + id)}
        >
            <AreaContainer>
                <AreaComponent
                    services={services}
                    id={id || ''}
                    area={{
                        service: applet.action.split(';')[0] || '',
                        name: applet.action.split(';')[1] || '',
                    }}
                    type="action"
                    refetch={refetch}
                />
                <AreaComponent
                    services={services}
                    id={id || ''}
                    area={{
                        service: applet.reaction_linked.split(';')[0] || '',
                        name: applet.reaction_linked.split(';')[1] || '',
                    }}
                    type="reaction"
                    refetch={refetch}
                />
            </AreaContainer>
        </FullContainer>
    )
}

export default UpdateApplet
 */
export {}
