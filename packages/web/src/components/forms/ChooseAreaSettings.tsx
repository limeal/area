import { useEffect, useState } from 'react'
import { MdAdd, MdLock, MdLoop } from 'react-icons/md'
import styled from 'styled-components'

import AnimatedMenu from '../animated/AnimatedMenu'
import AnimatedButtonWB from '../animated/buttons/AnimatedButtonWB'
import FullContainer from '../applets/FullContainer'
import { IAction, IService, IStore } from '../../interfaces'
import { useGetNewAppletQuery } from '../../redux/api'
import { applyPretty, applyPrettySettings } from '../../utils/string'

// Fields
import InputField from '../fields/InputField'
import SelectField from '../fields/SelectField'
import SelectFieldURI from '../fields/SelectFieldURI'

const ParamList = styled.ul`
    display: flex;
    flex-direction: column;
    justify-content: center;
    list-style: none;
    max-width: 300px;
`

const Hero = styled.div`
    display: flex;
    flex-direction: column;
    justify-content: center;
`

/* Defining the props that the component will receive. */
interface ChooseAreaSettingsProps {
    type: 'action' | 'reaction'
    services: IService[]
    selectedService: IService
    area: IAction
    no_working_fields: any
    callback: (settings: any) => void
    goBack: any
    loading: boolean
}

/* A React component that is used to display the settings for a specific area. */
const ChooseAreaSettings = ({
    type,
    services,
    selectedService,
    area,
    no_working_fields,
    callback,
    goBack,
    loading,
}: ChooseAreaSettingsProps) => {
    const { data } = useGetNewAppletQuery(
        {
            field: 'action',
        },
        {
            refetchOnMountOrArgChange: true,
            refetchOnReconnect: true,
        }
    )

    const [settings, setSettings] = useState<any>({})
    const [staticSettings, setStaticSettings] = useState<any>([])
    const [optionalFields, setOptionalFields] = useState<any>([])
    const [action, setAction] = useState<any>()

    const updateSettings = (key: string, value: any) => {
        if (value === null) {
            delete settings[key]
            setSettings({ ...settings })
            return
        }

        setSettings({ ...settings, [key]: value })
    }

    const updateStaticSettings = (key: string, value: any) => {
        console.log(key, value)
        if (value === null) {
            delete staticSettings[key]
            setStaticSettings({ ...staticSettings })
            return
        }

        setStaticSettings({ ...staticSettings, [key]: value })
    }

    const addOptionalField = (key: string) => {
        setOptionalFields([...optionalFields, key])
    }

    useEffect(() => {
        if (data && data.action) {
            const actionService = services.find(
                (s) => s.name === data.action?.service
            )
            const action = actionService?.actions.find(
                (a) => a.name === data.action?.name
            )
            setAction(action)
        }
    }, [data])

    const getField = (
        name: string,
        value: IStore,
        action: any
    ): React.ReactElement => {
        if (value.type === 'select')
            return (
                <SelectField
                    name={name}
                    value={value}
                    values={value.values}
                    settings={settings}
                    updateSettings={updateSettings}
                    no_working_fields={no_working_fields}
                />
            )
        if (value.type === 'select_uri') {
            return (
                <SelectFieldURI
                    name={name}
                    value={value}
                    service={selectedService}
                    uri={value.values[0]}
                    action={action}
                    settings={settings}
                    updateSettings={updateSettings}
                    no_working_fields={no_working_fields}
                />
            )
        }
        return (
            <InputField
                name={name}
                value={value}
                settings={staticSettings}
                updateSettings={updateStaticSettings}
                no_working_fields={no_working_fields}
                action={action}
            />
        )
    }

    return (
        <FullContainer
            title={applyPretty(area.name)}
            goBack={goBack}
            style={{
                justifyContent: 'space-evenly',
            }}
            helpText={
                'You can add more parameters by clicking on the "Add parameter" button.\n' +
                'You can remove a parameter by clicking on the "Remove" button (only optional).\n' +
                'You can use components by clicking on the "Add component" button. (reactions only)\n'
            }
        >
            <Hero>
                <ParamList>
                    {Object.entries(area.store)
                        .sort((a, b) => {
                            if (a[1].priority < b[1].priority) return -1
                            if (a[1].priority > b[1].priority) return 1
                            if (
                                a[1].need_fields &&
                                a[1].need_fields.includes(b[0])
                            )
                                return 1
                            if (
                                b[1].need_fields &&
                                b[1].need_fields.includes(a[0])
                            )
                                return -1
                            return 0
                        })
                        .filter(([key, value]) => {
                            const store: IStore = value
                            if (optionalFields.includes(key)) return true

                            // Check if field in need fields is filled
                            if (store.need_fields) {
                                for (const field of store.need_fields) {
                                    if (!settings[field]) return false
                                }
                            }

                            if (store.required) return true
                            return false
                        })
                        .map(([key, value], index) => (
                            <li
                                key={index}
                                style={{
                                    color: '#fff',
                                }}
                            >
                                <label
                                    style={{
                                        color: '#fff',
                                        fontSize: '1.2rem',
                                        borderBottom: '1px solid #fff',
                                    }}
                                >
                                    {applyPrettySettings(key)}
                                </label>
                                <p>{value.description || ''}</p>
                                {getField(
                                    key,
                                    value,
                                    type === 'reaction' ? action : null
                                )}
                                {!value.required && (
                                    <button
                                        onClick={() => {
                                            updateSettings(key, null)
                                            setOptionalFields(
                                                optionalFields.filter(
                                                    (field: string) =>
                                                        field !== key
                                                )
                                            )
                                        }}
                                    >
                                        Remove
                                    </button>
                                )}
                            </li>
                        ))}
                </ParamList>
                {optionalFields.length <
                    Object.values(area.store).filter((v) => !v.required)
                        .length && (
                    <AnimatedMenu
                        ButtonContent={
                            <>
                                <MdAdd />
                                <p>Add parameter</p>
                            </>
                        }
                        options={Object.entries(area.store)
                            .filter(([key, value]) => {
                                if (value.required) return false
                                if (optionalFields.includes(key)) return false
                                return true
                            })
                            .map(([key, _]) => ({
                                label: key,
                                onClick: () => addOptionalField(key),
                            }))}
                    />
                )}
            </Hero>
            {Object.entries({ ...settings, ...staticSettings }).length >=
                Object.entries(area.store).filter(
                    ([_, value]) => value.required
                ).length ||
            Object.entries(area.store).filter(([_, value]) => value.required)
                .length === 0 ? (
                <AnimatedButtonWB
                    onClick={() => callback({ ...settings, ...staticSettings })}
                    disabled={loading}
                    bgColor={loading ? '#222222' : '#fff'}
                    color={loading ? '#fff' : '#222222'}
                >
                    {loading ? (
                        <>
                            <MdLoop className="spin" />
                            Submitting...
                        </>
                    ) : (
                        <p>Submit</p>
                    )}
                </AnimatedButtonWB>
            ) : (
                <AnimatedButtonWB bgColor="#f00" color="#fff">
                    <MdLock />
                    Fill all required fields
                </AnimatedButtonWB>
            )}
        </FullContainer>
    )
}

export default ChooseAreaSettings
