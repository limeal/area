/**
 * It takes an array of characters and a shift value and returns a map of characters to characters
 * @param {any} alphabets - an array of alphabets
 * @param {number} shift - the amount of shift to be applied to the alphabets
 * @returns A function that takes in a string and a shift number and returns a string.
 */
const createMAp = (alphabets: any, shift: number) => {
    return alphabets.reduce(
        (charsMap: any, currentChar: any, charIndex: any) => {
            const copy = { ...charsMap }
            let ind = (charIndex + shift) % alphabets.length
            if (ind < 0) {
                ind += alphabets.length
            }
            copy[currentChar] = alphabets[ind]
            return copy
        },
        {}
    )
}

/**
 * It takes a string and a shift value, and returns a new string with each character shifted by the
 * shift value
 * @param {any} org - The original string that you want to encrypt.
 * @param [shift=0] - The number of characters to shift the alphabet by.
 * @returns A function that takes in a string and a shift value and returns a string.
 */
const encrypt = (org: any, shift = 0) => {
    const alphabets =
        'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890'.split(
            ''
        )
    const map = createMAp(alphabets, shift)
    return org
        .toLowerCase()
        .split('')
        .map((char: any) => map[char] || char)
        .join('')
}

export default encrypt
