import { motion } from 'framer-motion'
import { useEffect, useState } from 'react'
import { toast } from 'react-toastify'
import styled from 'styled-components'
import Modal from '../Modal'

import { MdLoop } from 'react-icons/md'
import { IAuthenticator } from '../../interfaces'
import { useAuthAccountMutation } from '../../redux/api'
import { getAvatarM } from '../../utils/more'
import AnimatedButtonHL from '../animated/buttons/AnimatedButtonHL'
import AnimatedButtonWB from '../animated/buttons/AnimatedButtonWB'
import encrypt from '../../utils/encrypt'

const FormContainer = styled(motion.form)`
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: space-around;
    width: 100%;
    height: 100%;
`

const FormTitle = styled(motion.h1)`
    font-size: 1.5rem;
    margin-bottom: 1rem;
`

const FormMiddle = styled(motion.div)`
    display: flex;
    flex-direction: column;
`

const FormBottom = styled(motion.div)`
    display: flex;
    flex-direction: column;
    align-items: center;
`

const FormInput = styled(motion.input)`
    width: 100%;
    padding: 1rem;
    margin: 0.5rem 0;
    border-radius: 2rem;
    border: none;
    background: #222222;
    color: #fff;
    outline: none;
`

const ButtonContainer = styled(motion.div)`
    display: flex;
    flex-direction: row;
    justify-content: space-around;
    width: 100%;
`

/* Defining the props that are common to both the BasicAuthentication and ServiceAuthentication
components. */
interface CommonProps {
    useAuth: any
    mode: 'login' | 'register' | 'external'
    setMode: any
    loading: boolean
}

/* A React component that is used to login or register a user. */
const BasicAuthentication = ({ useAuth, mode, setMode }: CommonProps) => {
    const onSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault()

        const data: any = Object.fromEntries([
            ...new FormData(e.currentTarget).entries(),
        ])

        if (mode === 'register' && data.password !== data.confirm_password)
            return

        const encoded_password = encrypt(
            data.password + data.email,
            data.email.length * 3
        )

        useAuth({
            mode,
            ...{
                email: data.email,
                encoded_password: encoded_password,
            },
        })
    }

    return (
        <>
            <FormContainer onSubmit={onSubmit}>
                <FormTitle>
                    {mode === 'login'
                        ? 'AreaAuth - Login'
                        : 'AreaAuth - Sign Up'}
                </FormTitle>
                <FormMiddle>
                    <FormInput name="email" placeholder="Email" type="email" />
                    <FormInput
                        name="password"
                        placeholder="Password"
                        type="password"
                    />
                    {mode === 'register' ? (
                        <FormInput
                            name="confirm_password"
                            placeholder="Confirm Password"
                            type="password"
                        />
                    ) : null}
                </FormMiddle>
                <FormBottom>
                    <AnimatedButtonWB color="#000" bgColor="#fff" type="submit">
                        {mode === 'login' ? 'Log in' : 'Sign Up'}
                    </AnimatedButtonWB>
                </FormBottom>
            </FormContainer>
            <ButtonContainer>
                <AnimatedButtonWB
                    color="#9ac2c9"
                    bgColor="transparent"
                    onClick={() => setMode('external')}
                >
                    External Authentication
                </AnimatedButtonWB>
                <AnimatedButtonHL
                    color="#fff"
                    hlColor="#139A43"
                    onClick={() =>
                        setMode(mode === 'login' ? 'register' : 'login')
                    }
                >
                    {mode === 'login'
                        ? "Doesn't have an account ?"
                        : 'Already have an account ?'}
                </AnimatedButtonHL>
            </ButtonContainer>
        </>
    )
}

const Container = styled(motion.div)`
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: space-around;
    width: 100%;
    height: 100%;
`

const Title = styled(motion.h1)`
    font-size: 1.5rem;
    margin-bottom: 1rem;
`

const List = styled(motion.ul)`
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: space-around;
    margin: 0;
    padding: 0;
    height: 50%;
`
const ListItem = styled(motion.li)`
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: space-around;
    width: 100%;
    height: 100%;
`

const OAuthButtonContainer = styled(motion.div)`
    display: flex;
    flex-direction: row;
    justify-content: space-around;
    align-items: center;
    width: 100%;
    height: 100%;
`

/* Defining the props that the ServiceAuthentication component will receive. */
interface ServiceAuthenticationProps {
    useAuth: any
    setMode: (mode: any) => void
    authenticators: IAuthenticator[]
    loading: boolean
}

/**
 * It renders a list of buttons that will open a new window to the third party authentication service
 * @param e - React.MouseEvent<HTMLButtonElement, MouseEvent>
 * @param {IAuthenticator} authenticator - IAuthenticator
 */
const ServiceAuthentication = ({
    useAuth,
    setMode,
    authenticators,
    loading,
}: ServiceAuthenticationProps) => {
    const onClick = async (
        e: React.MouseEvent<HTMLButtonElement, MouseEvent>,
        authenticator: IAuthenticator
    ) => {
        e.preventDefault()
        e.stopPropagation()
        const redirectUri = 'http://localhost:8081'

        window.handleAuthorization = (data: any) => {
            useAuth({
                mode: 'external',
                ...{
                    authenticator: authenticator.name,
                    code: data.code,
                    redirect_uri: redirectUri,
                },
            })
        }

        window.open(
            authenticator.authorization_uri +
                '&redirect_uri=http://localhost:8081',
            'Login with ' + authenticator.name,
            'width=800,height=600'
        )
    }

    return (
        <Container>
            <Title>AreaAuth - Third Party</Title>
            <List>
                {authenticators
                    .filter((a: IAuthenticator) => a.enabled)
                    .map((a: IAuthenticator, index: number) => {
                        return (
                            <ListItem key={index}>
                                <AnimatedButtonWB
                                    color="#fff"
                                    bgColor={a.more.color}
                                    disabled={loading}
                                    onClick={(e: any) => onClick(e, a)}
                                >
                                    <OAuthButtonContainer
                                        style={{
                                            filter: loading
                                                ? 'blur(0.4px)'
                                                : '',
                                        }}
                                    >
                                        {loading ? (
                                            <>
                                                <MdLoop className="spin" />
                                                Loading..
                                            </>
                                        ) : (
                                            <>
                                                <img
                                                    src={getAvatarM(
                                                        a.more,
                                                        a.name
                                                    )}
                                                    alt={a.name}
                                                    width="30"
                                                    height="30"
                                                />
                                                {a.name}
                                            </>
                                        )}
                                    </OAuthButtonContainer>
                                </AnimatedButtonWB>
                            </ListItem>
                        )
                    })}
            </List>
            <ButtonContainer>
                <AnimatedButtonHL
                    color="#fff"
                    hlColor="#139A43"
                    onClick={() => setMode('register')}
                >
                    Switch to register
                </AnimatedButtonHL>
                <AnimatedButtonHL
                    color="#fff"
                    hlColor="#139A43"
                    onClick={() => setMode('login')}
                >
                    Switch to login
                </AnimatedButtonHL>
            </ButtonContainer>
        </Container>
    )
}

/* Defining the props that the AuthModal component will receive. */
interface AuthModalProps {
    open: boolean
    close: () => void
    authenticators: IAuthenticator[]
    updateToken: any
}

/**
 * It renders a modal with a form to login or register, and it also renders a list of buttons to login
 * with external services
 * @param {AuthModalProps}  - `open`: a boolean that determines whether the modal is open or not
 */
const AuthModal = ({
    open,
    close,
    authenticators,
    updateToken,
}: AuthModalProps) => {
    const [mode, setMode] = useState<'login' | 'register' | 'external'>(
        'external'
    )
    const [useAuth, { data, isLoading, error, isSuccess }] =
        useAuthAccountMutation()

    useEffect(() => {
        if (isLoading || isSuccess) {
            toast.promise(
                new Promise((res, rej) => {
                    if (isSuccess) return res(0)
                    if (isLoading) return rej()
                }),
                {
                    pending: 'Logging in...',
                    success: 'Logged in successfully !',
                },
                {
                    position: toast.POSITION.TOP_CENTER,
                }
            )
        }
        if (error && 'data' in error) {
            toast.error(JSON.stringify(error.data), {
                position: toast.POSITION.TOP_CENTER,
            })
        }
        if (isSuccess) {
            updateToken(data)
            close()
        }
    }, [isSuccess, error, isLoading])

    return (
        <Modal open={open} close={close}>
            {mode == 'external' ? (
                <ServiceAuthentication
                    useAuth={useAuth}
                    setMode={setMode}
                    authenticators={authenticators}
                    loading={isLoading}
                />
            ) : (
                <BasicAuthentication
                    useAuth={useAuth}
                    mode={mode}
                    setMode={setMode}
                    loading={isLoading}
                />
            )}
        </Modal>
    )
}

export default AuthModal
