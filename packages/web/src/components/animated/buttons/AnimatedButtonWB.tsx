import { motion } from 'framer-motion'
import React from 'react'
import styled from 'styled-components'

interface AnimatedButtonWBProps {
    color: string
    bgColor: string
    type?: 'submit' | 'reset' | 'button'
    children: React.ReactNode
    [key: string]: any
}

const Button = styled(motion.button)`
    padding: 1rem 2rem;
    border-radius: 2rem;
    border: none;
    cursor: pointer;
    background-color: ${(props) => props.theme.bgColor};
    color: ${(props) => props.color};
`

const AnimatedButtonWB = ({
    color,
    bgColor,
    type,
    children,
    ...props
}: AnimatedButtonWBProps) => {
    return (
        <Button
            color={color}
            theme={{ bgColor }}
            type={type}
            whileHover={{
                scale: 1.1,
                transition: { duration: 0.2 },
            }}
            whileTap={{
                scale: 0.9,
                transition: { duration: 0.3 },
            }}
            {...props}
        >
            {children}
        </Button>
    )
}

export default AnimatedButtonWB
