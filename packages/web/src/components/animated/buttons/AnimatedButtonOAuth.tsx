import { motion } from 'framer-motion'
import styled from 'styled-components'

interface AnimatedButtonOAuthProps {
    icon?: string
    text: string
    color: string
    bgColor: string
    onClick?: () => void
}

const Button = styled(motion.button)`
    color: ${(props) => props.color};
    background-color: ${(props) => props.theme.bgColor};
    border: none;
    border-radius: 2rem;
    padding: 10px 20px;
    font-size: 1.2rem;
    cursor: pointer;
    margin: 10px;
    box-shadow: 0 5px 16px rgba(0, 0, 0, 0.2);
`

const Icon = styled.img`
    width: 20px;
    height: 20px;
    margin-right: 10px;
    color: white;
`

const AnimatedButtonOAuth = ({
    icon,
    text,
    color,
    bgColor,
    onClick,
}: AnimatedButtonOAuthProps) => {
    return (
        <Button
            color={color}
            theme={{ bgColor }}
            onClick={onClick}
            whileHover={{ scale: 1.1, transition: { duration: 0.2 } }}
            whileTap={{ scale: 0.9, transition: { duration: 0.2 } }}
        >
            {icon && <Icon src={icon} alt={'Image'} />}
            {text}
        </Button>
    )
}

export default AnimatedButtonOAuth
