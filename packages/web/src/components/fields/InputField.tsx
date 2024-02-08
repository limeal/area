import { motion } from 'framer-motion'
import { useEffect, useState } from 'react'
import { MdAdd } from 'react-icons/md'
import styled from 'styled-components'

import { toast } from 'react-toastify'
import { IStore } from '../../interfaces'
import AnimatedMenu from '../animated/AnimatedMenu'

const Input = styled(motion.input)`
    border: 1px solid #000;
    border-radius: 10px;
    padding: 10px;
    margin: 10px 0;
    width: 100%;
`

const TextArea = styled(motion.textarea)`
    border: 1px solid #000;
    border-radius: 10px;
    padding: 10px;
    margin: 10px 0;
    width: 100%;
`

const InputField = ({
    name,
    value,
    settings,
    updateSettings,
    no_working_fields,
    action,
}: {
    name: string
    value: IStore
    settings: any
    updateSettings: (key: string, value: any) => void
    no_working_fields: any
    action: any
}) => {
    const [content, setContent] = useState<string>(settings[name] || '')
    const [components, setComponents] = useState<string[]>()
    const [isOk, setIsOk] = useState<boolean>(true)

    useEffect(() => {
        setComponents(action?.components)
    }, [action])

    const replContent = (ncontent: string) => {
        if (ncontent === '') updateSettings(name, null)
        else if (value.type === 'number')
            updateSettings(name, parseInt(ncontent))
        else updateSettings(name, ncontent)
        setContent(ncontent)
    }

    useEffect(() => {
        const obj = JSON.parse(no_working_fields)
        if (obj && obj['data'][name] == false) {
            toast.error('Please enter a valid value for ' + name + '.')
            setIsOk(false)
        }
    }, [no_working_fields])

    return (
        <div
            style={{
                display: 'flex',
                flexDirection: 'column',
                margin: '10px 0',
                width: '100%',
            }}
        >
            <div
                style={{
                    display: 'flex',
                    width: '100%',
                    alignItems: 'center',
                }}
            >
                {value.type !== 'long_string' ? (
                    <Input
                        autoFocus={true}
                        type={value.type === 'email' ? 'email' : 'text'}
                        title="Enter text..."
                        placeholder={'Enter text...'}
                        required={value.required}
                        value={content}
                        animate={
                            (!isOk && {
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
                            replContent(e.currentTarget.value)
                        }}
                    />
                ) : (
                    <TextArea
                        autoFocus={true}
                        title="Enter text..."
                        placeholder={'Enter text...'}
                        required={value.required}
                        value={content}
                        animate={
                            (!isOk && {
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
                            replContent(e.currentTarget.value)
                        }}
                    />
                )}
            </div>
            {components && components.length > 0 && (
                <AnimatedMenu
                    ButtonContent={
                        <>
                            <MdAdd />
                            Add component
                        </>
                    }
                    options={components.map((component) => ({
                        label: component,
                        onClick: () => {
                            replContent(content + '{{' + component + '}}')
                        },
                    }))}
                />
            )}
        </div>
    )
}

export default InputField
