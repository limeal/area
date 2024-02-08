import { motion } from 'framer-motion'
import styled from 'styled-components'
import { applyPretty } from '../../utils/string'

const Container = styled.div`
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    grid-gap: 10px;
    justify-content: center;
    align-items: center;
    padding: 10px;
    max-height: 550px;
    overflow: scroll;
`

const Card = styled(motion.button)`
    display: flex;
    flex-direction: column;
    background-color: ${(props) => props.theme.backgroundColor || '#fff'};
    align-items: center;
    justify-content: space-evenly;
    height: ${(props) => props.theme.height || '250px'};
    width: ${(props) => props.theme.width || '250px'};
    border-radius: 10px;
    margin: 10px;
    padding: 10px;
    color: ${(props) => props.theme.color || '#000'};
    cursor: ${(props) => props.theme.cursor || 'pointer'};
    outline: none;
    border: ${(props) => props.theme.border || '5px solid #000'};
`

/* Defining the props that the component will take in. */
interface ItemGridProps {
    items: any[]
    query?: string
    cardTheme: (item: any) => React.CSSProperties
    onClickCard: (item: any) => void
    getCardElements: (item: any) => React.ReactNode
    [key: string]: any
}

/* A React component that takes in a bunch of props and returns a styled component. */
const ItemGrid = ({
    items,
    query,
    cardTheme,
    getCardElements,
    onClickCard,
    ...props
}: ItemGridProps) => {
    return (
        <Container {...props}>
            {items
                .filter((item: any) => {
                    if (!query || query === '') return true
                    if (!item.name) return true
                    const dname = applyPretty(item.name)
                    return dname.toLowerCase().includes(query.toLowerCase())
                })
                .map((item: any, index: number) => (
                    <Card
                        key={index}
                        theme={cardTheme(item)}
                        whileHover={{
                            scale: 1.1,
                            transition: { duration: 0.2 },
                        }}
                        whileTap={{
                            scale: 0.95,
                            transition: { duration: 0.3 },
                        }}
                        onClick={() => onClickCard(item)}
                    >
                        {getCardElements(item)}
                    </Card>
                ))}
        </Container>
    )
}

export default ItemGrid
