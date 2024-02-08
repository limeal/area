import { motion } from 'framer-motion'
import React from 'react'
import styled from 'styled-components'

import { MdArrowBack, MdHelp } from 'react-icons/md'
import AnimatedButtonWB from '../animated/buttons/AnimatedButtonWB'

const Container = styled(motion.div)`
    position: absolute;
    background-color: #222222;
    display: flex;
    flex-direction: column;
    top: 0;
    left: 0;
    width: 100%;
    min-height: 100%;
`

const ContainerHeader = styled(motion.div)`
    display: flex;
    flex: 1;
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
    border-bottom: 1px solid #fff;
    margin-bottom: 10px;
    padding: 10px;
`

const ContainerBody = styled(motion.div)`
    display: flex;
    flex: 9;
    flex-direction: column;
    align-items: center;
    justify-content: center;
`

const TitleContainer = styled(motion.div)`
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
    min-width: 300px;
`

const Title = styled(motion.h1)`
    font-size: 3rem;
    font-weight: 800;
    letter-spacing: 0.1rem;
    color: #fff;
`

interface FullContainerProps {
    backIcon?: React.ReactNode
    titleicon?: React.ReactNode
    title: string
    goBack: () => void
    help?: () => void
    children: React.ReactNode
    style?: React.CSSProperties
    helpText?: string
}

/* A React component that takes in props and returns a styled component. */
const FullContainer = (props: FullContainerProps) => {
    return (
        <Container>
            <ContainerHeader>
                <AnimatedButtonWB
                    bgColor="#222222"
                    color="#fff"
                    onClick={() => props.goBack()}
                    style={{
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                    }}
                >
                    {props.backIcon || (
                        <MdArrowBack fontSize={'1.5rem'} color="#fff" />
                    )}
                </AnimatedButtonWB>
                <TitleContainer>
                    {props.titleicon}
                    <Title>{props.title}</Title>
                </TitleContainer>
                <AnimatedButtonWB
                    bgColor="#222222"
                    color="#fff"
                    onClick={() => props.help && props.help()}
                    style={{
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        color: '#000',
                    }}
                >
                    <MdHelp
                        fontSize={'1.5rem'}
                        color="#fff"
                        title={props.helpText || ''}
                    />
                </AnimatedButtonWB>
            </ContainerHeader>
            <ContainerBody style={props.style}>{props.children}</ContainerBody>
        </Container>
    )
}

export default FullContainer
