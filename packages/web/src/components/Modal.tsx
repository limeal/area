import { motion } from 'framer-motion'
import React from 'react'
import styled from 'styled-components'
//styled.div`
const Container = styled(motion.div)`
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100vh;
    display: flex;
    z-index: 2;
    background: rgba(0, 0, 0, 0.2);
    justify-content: center;
    align-items: center;
`
//    background: #1c1b37;

const Backdrop = styled(Container)``

const ContentContainer = styled(motion.div)`
    width: ${(props) => props.theme.width || '350px'};
    height: ${(props) => props.theme.height || '500px'};
    background: #0e1116;
    border: 1px solid white;
    border-radius: 2rem;
    z-index: 3;
    display: flex;
    flex-direction: column;
    align-items: center;
    position: relative;
    padding: 0 2rem;
    box-shadow: 0 5px 16px rgba(0, 0, 0, 0.2);
    color: #fff;
`

/* Defining the props that the Modal component will take. */
interface ModalProps {
    open: boolean
    close: () => void
    children: React.ReactNode
    [key: string]: any
}

// BackDrop

/* A function that takes in a ModalProps object and returns a modal. */
function Modal({ open, close, children, ...props }: ModalProps) {
    if (!open) return <></>

    return (
        <Container>
            <Backdrop
                onClick={() => close()}
                initial={{ opacity: 0, scale: 1 }}
                animate={{ opacity: 1, scale: 1 }}
                transition={{ duration: 0.5 }}
            />
            <ContentContainer
                {...props}
                initial={{ scale: 0 }}
                animate={{ scale: 1 }}
                transition={{
                    type: 'spring',
                    stiffness: 260,
                    damping: 20,
                }}
            >
                {children}
            </ContentContainer>
        </Container>
    )
}

export default Modal
