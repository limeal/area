import { motion } from 'framer-motion'
import { useEffect, useState } from 'react'
import styled from 'styled-components'

import { toast } from 'react-toastify'
import { IService, IStore } from '../../interfaces'
import { useFetchServiceApiEndpointMutation } from '../../redux/api'
import { applyPretty } from '../../utils/string'

const Select = styled(motion.select)`
    border: 1px solid #000;
    border-radius: 10px;
    padding: 10px;
    width: 100%;
    max-width: 300px;
    margin: 10px 0;
`

const Loading = styled.div`
    border: 1px solid #000;
    border-radius: 10px;
    padding: 10px;
    width: 100%;
    max-width: 300px;
    margin: 10px 0;
    text-align: center;
`

const SelectFieldURI = ({
    name,
    value,
    service,
    uri,
    action,
    settings,
    updateSettings,
    no_working_fields,
}: {
    name: string
    value: IStore
    service: IService
    uri: string
    action: any
    settings: any
    updateSettings: (key: string, value: any) => void
    no_working_fields: any
}) => {
    const [values, setValues] = useState([])
    const [useApi, { error }] = useFetchServiceApiEndpointMutation()
    const [components, setComponents] = useState<string[]>(action?.components)

    useEffect(() => {
        setComponents(action?.components)
    }, [action])

    useEffect(() => {
        useApi({
            service: service.name,
            // Replace ${key} with value in settings for key
            endpoint: uri.replace(/\${(.*?)}/g, (_match, key) => {
                return settings[key] || 'default'
            }),
            method: 'GET',
        })
            .unwrap()
            .then((payload: any) => {
                const namePath = payload?.fields[0].split(':')
                const valuePath = payload?.fields[1].split(':')

                const options = payload?.data.map((item: any) => ({
                    name: namePath.reduce(
                        (acc: any, cur: any) => acc[cur],
                        item
                    ),
                    value: valuePath.reduce(
                        (acc: any, cur: any) => acc[cur],
                        item
                    ),
                }))

                if (components) {
                    for (const component of components) {
                        // Component: discord:message:id
                        // Name: req:message:id
                        // Check if component starts with service name and ends with end of name
                        if (
                            component.startsWith(service.name) &&
                            component.endsWith(name.slice(4))
                        ) {
                            options.unshift({
                                name: '[Variable]: ' + applyPretty(component),
                                value: '{{' + component + '}}',
                            })
                        }
                    }
                }

                setValues(options)
            })
            .catch((error) => console.error(error))
    }, [uri, settings, components])

    useEffect(() => {
        if (error && 'data' in error) {
            toast.error(JSON.stringify(error.data))
        }
    }, [error])

    if (values.length === 0) return <Loading>Loading...</Loading>

    return (
        <Select
            required={value.required}
            value={settings[name] || ''}
            animate={
                (no_working_fields &&
                    no_working_fields.includes(name) &&
                    no_working_fields[name] && {
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
                if (e.currentTarget.value === '') {
                    updateSettings(name, null)
                    return
                }
                updateSettings(name, e.currentTarget.value)
            }}
        >
            <option value=""></option>
            {values.map((val: any, index: number) => (
                <option key={index} value={val.value}>
                    {val.name}
                </option>
            ))}
        </Select>
    )
}

export default SelectFieldURI
