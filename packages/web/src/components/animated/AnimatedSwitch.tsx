import { motion } from 'framer-motion'
import styled from 'styled-components'

const Container = styled(motion.div)`
    width: 100px;
    height: 40px;
    border-radius: 100px;
    padding: 10px;
    display: flex;
    cursor: pointer;
    background-color: ${(props) => (props.theme.isOn ? '#22cc88' : '#dddddd')};
    justify-content: ${(props) =>
        props.theme.isOn ? 'flex-end' : 'flex-start'};
`

const Element = styled(motion.div)`
    width: 40px;
    height: 40px;
    background-color: #ffffff;
    border-radius: 200px;
    box-shadow: 1px 2px 3px rgba(0, 0, 0, 0.02);
`

interface SwitchProps {
    isOn: boolean
    [key: string]: any
}

const AnimatedSwitch = ({ isOn, ...props }: SwitchProps) => {
    return (
        <Container theme={{ isOn }} animate {...props}>
            <Element animate />
        </Container>
    )
}

export default AnimatedSwitch
