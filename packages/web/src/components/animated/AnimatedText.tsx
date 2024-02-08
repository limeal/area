import { motion } from 'framer-motion'

interface AnimatedTextProps {
    text: string
    style?: any
    delay?: number
}

const AnimatedText = ({ text, style, delay = 0.5 }: AnimatedTextProps) => {
    const sentenceOpt = {
        hidden: { opacity: 1 },
        visible: {
            opacity: 1,
            transition: {
                delay,
                staggerChildren: 0.08,
            },
        },
    }

    const letterOpt = {
        hidden: { opacity: 0, y: 50 },
        visible: {
            opacity: 1,
            y: 0,
        },
    }

    return (
        <motion.h3
            variants={sentenceOpt}
            initial="hidden"
            style={style}
            animate="visible"
        >
            {text.split('').map((char, index) => (
                <motion.span key={index} variants={letterOpt}>
                    {char}
                </motion.span>
            ))}
        </motion.h3>
    )
}

export default AnimatedText
