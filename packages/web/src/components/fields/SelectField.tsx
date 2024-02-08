import { motion } from 'framer-motion'
import styled from 'styled-components'

import { IStore } from '../../interfaces'

const Select = styled(motion.select)`
    border: 1px solid #000;
    border-radius: 10px;
    padding: 10px;
    width: 100%;
    margin: 10px 0;
`

const SelectField = ({
    name,
    value,
    values,
    settings,
    updateSettings,
    no_working_fields,
}: {
    name: string
    value: IStore
    values: any
    settings: any
    updateSettings: (key: string, value: any) => void
    no_working_fields: any
}) => {
    return (
        <Select
            required={value.required}
            defaultValue={settings[name] || false}
            animate={
                (no_working_fields &&
                    no_working_fields.includes(name) && {
                        x: [0, 10, -10, 10, -10, 10, -10, 10, -10, 0],
                        backgroundColor: [
                            '#fff',
                            '#f00',
                            '#f00',
                            '#f00',
                            '#f00',
                            '#f00',
                            '#f00',
                            '#f00',
                            '#f00',
                            '#fff',
                        ],
                        transition: {
                            delay: 0.1,
                        },
                    }) || {
                    x: 0,
                    transition: {
                        delay: 0.1,
                    },
                }
            }
            onChange={(e) => {
                e.preventDefault()
                updateSettings(name, e.currentTarget.value)
            }}
        >
            <option value=""></option>
            {values.map((val: any, index: number) => (
                <option key={index} value={val}>
                    {val}
                </option>
            ))}
        </Select>
    )
}

export default SelectField
