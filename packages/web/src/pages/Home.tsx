import styled from 'styled-components'
import AnimatedButtonWB from '../components/animated/buttons/AnimatedButtonWB'

const HeroContainer = styled.section`
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    font-size: 2rem;
    color: #fff;
    height: 100%;
    box-shadow: 0 5px 16px rgba(0, 0, 0, 0.2);
    background-image: url('/images/bg.svg');
    background-repeat: no-repeat;
    background-position: top;
    background-size: cover;
    background-color: #222222;
`

const HeroTitle = styled.h1`
    padding: 0;
    margin: 0;
    font-weight: bold;
    font-size: 3.5rem;
    font-weight: 800;
    line-height: 4rem;
`

const HeroContent = styled.div`
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: space-evenly;
    padding-bottom: 1rem;
    text-align: center;
    height: 300px;
    width: 500px;
`

const HeroDescription = styled.p`
    padding: 0;
    margin: 0;
    font-size: 1.5rem;
`

/* Defining the props that the component will receive. */
interface HomeProps {
    openAuthModal: () => void
    storage: any
}

/**
 * A functional component that renders the hero section of the home page.
 * @param {HomeProps}  - `color` - the color of the text
 */
const Hero = ({ openAuthModal, storage }: HomeProps) => {
    return (
        <HeroContainer>
            <HeroContent>
                <HeroTitle>Every thing works better together</HeroTitle>
                <HeroDescription>
                    Quickly and easily automate your favorite apps and devices.
                </HeroDescription>
            </HeroContent>
            {storage === '' && (
                <AnimatedButtonWB
                    color="#000"
                    bgColor="white"
                    onClick={openAuthModal}
                    style={{
                        fontSize: '1.8rem',
                        padding: '1.6rem 4.4rem',
                        borderRadius: '5rem',
                        minWidth: '300px',
                    }}
                >
                    Start today
                </AnimatedButtonWB>
            )}
        </HeroContainer>
    )
}

const Container = styled.div`
    height: 100%;
    overflow: hidden;
`

/**
 * `Home` is a function that takes in an object with two properties, `openAuthModal` and `storage`, and
 * returns a `Container` component with a `Hero` component inside of it
 * @param {HomeProps}  - HomeProps
 * @returns A React component
 */
const Home = ({ openAuthModal, storage }: HomeProps) => {
    return (
        <Container>
            <Hero openAuthModal={openAuthModal} storage={storage} />
        </Container>
    )
}

export default Home
