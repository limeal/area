import { IService, IServiceMore } from '../interfaces'
import { baseUrl } from '../redux/api'

/**
 * "If the service is null or undefined, return null. Otherwise, if the service has a more property,
 * return that. Otherwise, if the service has an authenticator property that has a more property,
 * return that. Otherwise, return null."
 *
 * The above function is a good example of how TypeScript can help you write code that is more correct
 * @param {IService | null | undefined} service - The service that is being rendered.
 * @returns The more property of the service object, or the more property of the authenticator property
 * of the service object, or null.
 */
const getMore = (service: IService | null | undefined) => {
    if (service === null || service === undefined) {
        return null
    }

    if (service.more !== null) {
        return service.more
    }

    if (service.authenticator !== null) {
        return service.authenticator.more
    }

    return null
}

/**
 * It returns the color property of the more object.
 * @param {IServiceMore} more - IServiceMore
 * @returns The color property of the more object.
 */
const getColorM = (more: IServiceMore) => {
    return more.color
}

/**
 * It returns the avatar image for a service.
 * @param {IServiceMore} more - IServiceMore
 * @param {string} name - The name of the avatar, which is the same as the name of the image file.
 * @returns A string
 */
const getAvatarM = (more: IServiceMore, name: string) => {
    return baseUrl + '/assets/' + name + '.png'
}

/**
 * If the service is null, return #222222, otherwise return the color of the service.
 * @param {IService | null | undefined} service - The service object that you want to get the color of.
 * @returns The color of the service
 */
const getColor = (service: IService | null | undefined) => {
    const more = getMore(service)

    if (more === null) {
        return '#222222'
    }

    return more.color
}

/**
 * If the service is null, return a placeholder image. Otherwise, return the image from the assets
 * folder
 * @param {IService | null | undefined} service - The service object
 * @returns A string
 */
const getAvatar = (service: IService | null | undefined) => {
    const more = getMore(service)

    if (more === null) {
        return 'https://via.placeholder.com/50'
    }

    return baseUrl + '/assets/' + service?.name + '.png'
}

export { getColor, getAvatar, getColorM, getAvatarM }
