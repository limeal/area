import React, { useState } from 'react'
import styled from 'styled-components'

import FullContainer from './FullContainer'
import ItemGrid from './ItemGrid'

const SearchBar = styled.input`
    width: 30%;
    height: 50px;
    background-color: #fff;
    border-radius: 2rem;
    border: none;
    border: 5px solid #000;
    outline: none;
    padding-left: 10px;
`

/* Defining the props that the component will receive. */
interface SelectorContainerProps {
    title: string
    titleIcon?: React.ReactNode
    items: any[]
    cardTheme: (item: any) => React.CSSProperties
    onClickCard: (item: any) => void
    getCardElements: (item: any) => React.ReactNode
    hasSearchBar: boolean
    goBack: () => void
    closeW?: () => void
}

/* A function that returns a react component. */
const SelectorContainer = (props: SelectorContainerProps) => {
    const [querySearch, setQuerySearch] = useState('')

    return (
        <FullContainer
            goBack={props.goBack}
            title={props.title}
            titleicon={props.titleIcon}
            help={props.closeW}
            helpText={
                'Easteregg: A specific combo in the search bar will give you a surprise!'
            }
        >
            {props.hasSearchBar && props.items.length > 0 && (
                <SearchBar
                    placeholder="Search..."
                    value={querySearch}
                    onChange={(e) => setQuerySearch(e.currentTarget.value)}
                />
            )}
            {querySearch !== 'bebe yoda' ? (
                props.items.length > 0 ? (
                    <ItemGrid
                        items={props.items}
                        query={querySearch}
                        cardTheme={props.cardTheme}
                        onClickCard={props.onClickCard}
                        getCardElements={props.getCardElements}
                    />
                ) : (
                    <p
                        style={{
                            color: 'white',
                            fontSize: '2rem',
                        }}
                    >
                        No items found
                    </p>
                )
            ) : (
                <img
                    src="https://i.pinimg.com/originals/51/42/f9/5142f92e0cf4a3fc8e9062cf5c0911a2.gif"
                    alt="chicken"
                />
            )}
        </FullContainer>
    )
}

export default SelectorContainer
