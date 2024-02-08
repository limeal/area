import Modal from '../Modal'
import AnimatedButtonWB from '../animated/buttons/AnimatedButtonWB'

interface ConfirmModalProps {
    message: string
    open: boolean
    close: any
    callback: any
    height?: string
}

/* A function that returns a modal. */
const ConfirmModal = ({
    message,
    open,
    close,
    callback,
    height = '200px',
}: ConfirmModalProps) => {
    return (
        <Modal
            open={open}
            close={close}
            style={{
                width: '350px',
                height,
            }}
        >
            <h2
                style={{
                    color: 'white',
                    fontSize: '1.2rem',
                    textAlign: 'center',
                    padding: '1rem',
                }}
            >
                {message}
            </h2>
            <div
                style={{
                    display: 'flex',
                    flexDirection: 'row',
                    justifyContent: 'space-around',
                    width: '100%',
                    padding: '1rem',
                }}
            >
                <AnimatedButtonWB
                    color="red"
                    bgColor="white"
                    onClick={() => {
                        callback()
                        close()
                    }}
                >
                    Confirm
                </AnimatedButtonWB>
                <AnimatedButtonWB
                    color="white"
                    bgColor="red"
                    onClick={() => {
                        close()
                    }}
                >
                    Cancel
                </AnimatedButtonWB>
            </div>
        </Modal>
    )
}

export default ConfirmModal
