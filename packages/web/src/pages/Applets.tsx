import { motion } from 'framer-motion'
import { useEffect, useState } from 'react'
import { toast } from 'react-toastify'
import styled from 'styled-components'

// Components
import { MdAdd, MdMobileFriendly } from 'react-icons/md'
import { useNavigate } from 'react-router-dom'
import AppletCard from '../components/applets/AppletCard'
import { IApplet, IService } from '../interfaces'
import { useGetAppletsQuery } from '../redux/api'

const Container = styled.div`
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    background-color: #222222;
    color: #fff;
    height: 100%;
`

const Title = styled(motion.h1)`
    font-size: 4rem;
    margin: 0;
    padding: 0;
    margin-left: 10px;
`

const AppletsContainer = styled.ul`
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: center;
    list-style: none;
`

const NewAppletCard = styled(motion.li)`
    height: 342px;
    width: 18em;
    background-color: #fff;
    border-radius: 10px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    margin: 10px;
    transition: all 0.2s ease-in-out;
    color: #222222;
    &:hover {
        transform: scale(1.05);
    }
`

/* Defining the props that the component will receive. */
interface AppletsProps {
    services: IService[]
}

/* It's a React component that displays the applets of the user. */
const Applets = ({ services }: AppletsProps) => {
    const navigate = useNavigate()
    const { data, error } = useGetAppletsQuery(null, {
        refetchOnMountOrArgChange: true,
    })
    const [applets, setApplets] = useState(data)

    useEffect(() => {
        if (error && 'data' in error && error.data)
            toast.error(JSON.stringify(error.data), {
                position: toast.POSITION.TOP_CENTER,
            })
    }, [error])

    useEffect(() => {
        if (data) setApplets(data)
    }, [data])

    return (
        <Container>
            <div
                style={{
                    display: 'flex',
                    flexDirection: 'row',
                    alignItems: 'center',
                    width: '100%',
                    flex: 2,
                    backgroundColor: '#b9090b',
                }}
            >
                <MdMobileFriendly size={50} />
                <Title
                    initial={{ opacity: 0, scale: 0.5 }}
                    animate={{ opacity: 1, scale: 1 }}
                    transition={{ duration: 0.5 }}
                >
                    Applets
                </Title>
            </div>
            <div
                style={{
                    display: 'flex',
                    width: '100%',
                    justifyContent: 'center',
                    flex: 8,
                }}
            >
                <AppletsContainer>
                    {applets &&
                        applets.map((applet: IApplet, index: number) => (
                            <NewAppletCard key={index}>
                                <AppletCard
                                    services={services}
                                    appletInfos={applet}
                                    mode="my"
                                />
                            </NewAppletCard>
                        ))}
                    <NewAppletCard onClick={() => navigate('/create')}>
                        <div
                            style={{
                                padding: '10px',
                            }}
                        >
                            <MdAdd size={100} />
                        </div>
                    </NewAppletCard>
                </AppletsContainer>
            </div>
            {/* <AnimatedButtonWB
                color="#fff"
                bgColor="#000"
                onClick={() => navigate('/create')}
            >
                New Applet
            </AnimatedButtonWB> */}
        </Container>
    )
}

export default Applets
