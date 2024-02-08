import { motion } from 'framer-motion'
import React from 'react'
import styled from 'styled-components'

const Button = styled(motion.button)`
    cursor: pointer;
    border: none;
    background: none;
    color: #fff;

    &:hover {
        padding-bottom: 10px;
        border-bottom: 2px solid ${(props) => props.color || '#b2ddf7'};
        color: ${(props) => props.theme.backgroundColor || '#b2ddf7'};
    }
`

interface AnimatedButtonHLProps {
    color: string
    hlColor: string
    type?: 'submit' | 'reset' | 'button'
    children: React.ReactNode
    [key: string]: any
}

const AnimatedButtonHL = ({
    color,
    hlColor,
    type,
    children,
    ...props
}: AnimatedButtonHLProps) => {
    return (
        <Button
            color={color}
            theme={{ backgroundColor: hlColor }}
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

export default AnimatedButtonHL
