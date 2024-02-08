import { useState } from 'react'
import styled from 'styled-components'

import AnimatedButtonWB from '../animated/buttons/AnimatedButtonWB'
import FullContainer from '../applets/FullContainer'

interface SubmitComponentProps {
    useSubmit: any
    setPage: any
}

const Input = styled.input`
    border: 1px solid #000;
    border-radius: 10px;
    padding: 10px;
    margin: 10px;
    outline: none;
`

const InputArea = styled.textarea`
    border: 1px solid #000;
    border-radius: 10px;
    padding: 10px;
    margin: 10px;
    outline: none;
    width: 200px;
    height: 100px;
`

const SubmitForm = styled.form`
    display: flex;
    flex-direction: column;
`

const SelectContainer = styled.select`
    border: 1px solid #000;
    border-radius: 10px;
    padding: 10px;
    margin: 10px;
    outline: none;
`

/* A way to pass props to a component. */
interface MySelectProps {
    values: string[]
    onChange: (e: any) => void
    [key: string]: any
}

/**
 * We're creating a function called MySelect that takes in an object with the properties values,
 * onChange, and props, and returns a SelectContainer with the onChange property set to the onChange
 * property of the object passed in, and the props property set to the props property of the object
 * passed in, and a list of options with the key set to the index of the value in the values array, the
 * value set to the value, and the text set to the value
 * @param {MySelectProps}  - values - an array of values that will be used to populate the select
 * options
 * @returns A SelectContainer component with an onChange prop and props.
 */
const MySelect = ({ values, onChange, ...props }: MySelectProps) => {
    return (
        <SelectContainer onChange={onChange} {...props}>
            {values.map((value, index) => (
                <option key={index} value={value}>
                    {value}
                </option>
            ))}
        </SelectContainer>
    )
}

const Label = styled.label`
    color: #fff;
`

const SubmitField = styled.div`
    display: flex;
    flex-direction: column;
`

/* It's a function that takes in an object with the properties useSubmit and setPage, and returns a
FullContainer component with the title set to "Create (2/2)", the goBack prop set to a function that
calls setPage with the argument 0, and a SubmitForm component with the onSubmit prop set to a
function that calls e.preventDefault(), and then calls useSubmit with an object with the properties
name, description, and public, where public is set to true if visibility is equal to "public", and
false otherwise. */
const SubmitComponent = ({ useSubmit, setPage }: SubmitComponentProps) => {
    const [name, setName] = useState('')
    const [description, setDescription] = useState('')
    const [visibility, setVisibility] = useState('private')

    const submitApplet = (e: any) => {
        e.preventDefault()
        useSubmit({
            name,
            description,
            public: visibility === 'public' ? true : false,
        })
    }

    return (
        <FullContainer title="Create (2/2)" goBack={() => setPage(0)}>
            <SubmitForm onSubmit={submitApplet}>
                <SubmitField>
                    <Label htmlFor="name">Name</Label>
                    <Input
                        onChange={(e) => setName(e.currentTarget.value)}
                        value={name}
                        type="text"
                        placeholder="Name"
                    />
                </SubmitField>
                <SubmitField>
                    <Label color="#fff" htmlFor="description">
                        Description
                    </Label>
                    <InputArea
                        onChange={(e) => setDescription(e.currentTarget.value)}
                        value={description}
                        placeholder="Description"
                    />
                </SubmitField>
                <SubmitField>
                    <Label color="#fff" htmlFor="visibility">
                        Visibility
                    </Label>
                    <MySelect
                        values={['private', 'public']}
                        onChange={(e) => setVisibility(e.currentTarget.value)}
                    />
                </SubmitField>
                {description.length > 0 && name.length > 0 && (
                    <AnimatedButtonWB
                        bgColor="#fff"
                        color="#222222"
                        type="submit"
                    >
                        Submit
                    </AnimatedButtonWB>
                )}
            </SubmitForm>
        </FullContainer>
    )
}

export default SubmitComponent
