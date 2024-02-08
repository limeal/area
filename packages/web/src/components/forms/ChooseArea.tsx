import SelectorContainer from '../applets/SelectorContainer'
import { IAction, IService } from '../../interfaces'
import capitalize, { applyPretty } from '../../utils/string'

import { getAvatar, getColor } from '../../utils/more'
import { MdWarning } from 'react-icons/md'
import ConfirmModal from '../modals/ConfirmModal'
import { useState } from 'react'

/* Defining the props that the component will receive. */
interface ChooseAreaProps {
    type: 'action' | 'reaction'
    selectedService: IService
    elements: IAction[]
    callback: (element: IAction) => void
    goBack: any
}

/* A React component that is being exported. */
const ChooseArea = ({
    selectedService,
    type,
    elements,
    callback,
    goBack,
}: ChooseAreaProps) => {
    const [open, setOpen] = useState(false)
    const [selectedArea, setSelectedArea] = useState<IAction>()

    return (
        <>
            <ConfirmModal
                message="You are about to select an action/reaction that is still in development. This means that it can crash your app or that it doesn't work as expected. Are you sure you want to continue?"
                callback={() => callback(selectedArea as IAction)}
                open={open}
                close={() => setOpen(false)}
                height="270px"
            />
            <SelectorContainer
                title={capitalize(selectedService.name)}
                titleIcon={
                    <img
                        src={getAvatar(selectedService)}
                        alt={selectedService.name}
                        style={{ width: '100px', height: '100px' }}
                    />
                }
                hasSearchBar={true}
                items={elements.filter((element) => {
                    console.log(process.env.NODE_ENV)
                    if (element.wip && process.env.NODE_ENV !== 'development')
                        return false
                    return true
                })}
                cardTheme={() => {
                    return {
                        backgroundColor: getColor(selectedService),
                        height: '180px',
                        color: '#fff',
                        border: 'none',
                    }
                }}
                getCardElements={(area) => (
                    <>
                        {area.wip && (
                            <MdWarning
                                style={{ color: '#fff', fontSize: '2rem' }}
                            />
                        )}
                        <p>{capitalize(type)}</p>
                        <h2>{applyPretty(area.name)}</h2>
                        <p>{area.description}</p>
                    </>
                )}
                onClickCard={(area) => {
                    if (area.wip) {
                        setSelectedArea(area)
                        setOpen(true)
                        return
                    }
                    return callback(area)
                }}
                goBack={goBack}
            />
        </>
    )
}

export default ChooseArea
