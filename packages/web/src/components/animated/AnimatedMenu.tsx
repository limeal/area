import { motion } from 'framer-motion'
import React, { useState } from 'react'
import styled from 'styled-components'

const ButtonActivateMenu = styled(motion.button)`
    cursor: pointer;
    display: flex;
    flex-direction: row;
    padding: 10px;
    background-color: #fff;
    outline: none;
    border: none;
`

const Menu = styled(motion.ul)`
    list-style: none;
    display: flex;
    flex-direction: column;
`

const MenuItem = styled(motion.li)``

const MenuItemButton = styled(motion.button)`
    cursor: pointer;
    width: 100%;
    padding: 10px;
`

interface Option {
    label: string
    onClick: () => void
}

interface AnimatedMenuProps {
    ButtonContent: React.ReactNode
    options: Option[]
}

const AnimatedMenu = ({ ButtonContent, options }: AnimatedMenuProps) => {
    const [open, setOpen] = useState(false)

    const handleClick = () => {
        setOpen(!open)
    }

    const handleClose = (option: Option) => {
        setOpen(false)
        option.onClick()
    }

    return (
        <>
            <ButtonActivateMenu
                id="basic-button"
                aria-controls={open ? 'basic-menu' : undefined}
                aria-haspopup="true"
                aria-expanded={open ? 'true' : undefined}
                onClick={handleClick}
            >
                {ButtonContent}
            </ButtonActivateMenu>
            {open && (
                <Menu id="basic-menu">
                    {options.map((option: Option, index: number) => (
                        <MenuItem key={index}>
                            <MenuItemButton onClick={() => handleClose(option)}>
                                {option.label}
                            </MenuItemButton>
                        </MenuItem>
                    ))}
                </Menu>
            )}
        </>
    )
}

export default AnimatedMenu
