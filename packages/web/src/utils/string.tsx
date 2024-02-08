/**
 * Capitalize the first letter of a string.
 * @param {string} str - string - The string to capitalize.
 * @returns The first character of the string is being capitalized and then the rest of the string is
 * being added to it.
 */
const capitalize = (str: string) => {
    return str.charAt(0).toUpperCase() + str.slice(1)
}

/**
 * It takes a string, splits it into an array of words, capitalizes each word, and then joins the array
 * back into a string
 * @param {string} str - string - The string to be converted
 * @returns A function that takes a string and returns a string.
 */
const applyPretty = (str: string) => {
    const words = str.split('_')
    return words.map((word) => capitalize(word)).join(' ')
}

/**
 * It takes a string, splits it into an array of words, removes the first word, capitalizes each word,
 * and then joins the words back together
 * @param {string} str - The string to be formatted.
 * @returns A function that takes a string and returns a string.
 */
const applyPrettySettings = (str: string) => {
    const words = str.split(':')
    words.shift()
    return words.map((word) => capitalize(word)).join(' ')
}

export default capitalize
export { applyPretty, applyPrettySettings }
