import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react'
import {
    IApplet,
    IArea,
    IAuthenticator,
    IMeta,
    IResponse,
    IService,
    IUser,
} from '../interfaces'

export const baseUrl = window.location.origin + '/api'
/* The above code is creating a serviceApi object that is being used to make requests to the server. */
export const serviceApi = createApi({
    reducerPath: 'AreaAPI',
    baseQuery: fetchBaseQuery({
        baseUrl,
        mode: 'cors',
    }),
    endpoints: (builder) => ({
        getAbout: builder.query({
            query: () => ({
                url: '/about.json',
                headers: {
                    Authorization: `Bearer ${
                        localStorage.getItem('token') || ''
                    }`,
                },
            }),
            transformResponse: (response: {
                server: {
                    services: IService[]
                    authenticators: IAuthenticator[]
                }
            }) => response.server,
            transformErrorResponse: (error) => {
                if (error.status === 401) {
                    localStorage.removeItem('token')
                    window.location.assign('/')
                }
                return error
            },
        }),
        // Authentication + Logout
        authAccount: builder.mutation({
            query: ({ mode, ...body }) => ({
                url: `/auth/${mode}`,
                method: 'POST',
                body,
                validateStatus: (status) =>
                    status.status === 201 || status.status === 200,
            }),
            transformResponse: (response: IResponse<{ token: string }>) =>
                response.data.token,
            transformErrorResponse: (error) => {
                if (error.status === 401) {
                    localStorage.removeItem('token')
                    window.location.assign('/')
                }
                return error
            },
        }),
        // Authorization
        getAuthorizations: builder.query({
            query: () => ({
                url: '/authorization',
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('token')}`,
                },
            }),
            transformResponse: (
                response: IResponse<{ authorizations: any; meta: IMeta[] }>
            ) => response.data,
            transformErrorResponse: (error) => {
                if (error.status === 401) {
                    localStorage.removeItem('token')
                    window.location.assign('/')
                }
                return error
            },
        }),
        /* A function that is being called to get the service authorizations. */
        getServiceAuthorizations: builder.query({
            query: () => ({
                url: '/authorization/services',
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('token')}`,
                },
            }),
            transformResponse: (response: IResponse<{ services: any }>) =>
                response.data.services,
            transformErrorResponse: (error) => {
                if (error.status === 401) {
                    localStorage.removeItem('token')
                    window.location.assign('/')
                }
                return error
            },
        }),
        /* Creating a mutation that will be used to post authorization to the server. */
        postAuthorization: builder.mutation({
            query: (body) => ({
                url: '/authorization',
                method: 'POST',
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('token')}`,
                },
                body,
            }),
            transformResponse: (response: IResponse<{ message: string }>) =>
                response.data.message,
            transformErrorResponse: (error) => {
                if (error.status === 401) {
                    localStorage.removeItem('token')
                    window.location.assign('/')
                }
                return error
            },
        }),
        /* Creating a mutation that will delete an authorization. */
        deleteAuthorization: builder.mutation({
            query: (id) => ({
                url: `/authorization/${id}`,
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('token')}`,
                },
                method: 'DELETE',
            }),
            transformResponse: (response: IResponse<{ message: string }>) =>
                response.data.message,
            transformErrorResponse: (error) => {
                if (error.status === 401) {
                    localStorage.removeItem('token')
                    window.location.assign('/')
                }
                return error
            },
        }),
        // Profile
        /* A function that is being called to get the profile of the user. */
        getProfile: builder.query({
            query: () => ({
                url: '/me',
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('token')}`,
                },
            }),
            transformResponse: (response: IResponse<{ account: IUser }>) =>
                response.data.account,
            transformErrorResponse: (error) => {
                if (error.status === 401) {
                    localStorage.removeItem('token')
                    window.location.assign('/')
                }
                return error
            },
        }),
        /* A mutation that is being made to the server. */
        logoutAccount: builder.mutation({
            query: () => ({
                url: '/me/logout',
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('token')}`,
                },
                method: 'POST',
            }),
            transformResponse: (response: IResponse<{ message: string }>) =>
                response.data.message,
            transformErrorResponse: (error) => {
                if (error.status === 401) {
                    localStorage.removeItem('token')
                    window.location.assign('/')
                }
                return error
            },
        }),
        /* A mutation that is modifying the profile of the user. */
        modifyProfile: builder.mutation({
            query: (body) => ({
                url: '/me',
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('token')}`,
                },
                method: 'PUT',
                body,
            }),
            transformResponse: (response: IResponse<{ message: string }>) =>
                response.data.message,
            transformErrorResponse: (error) => {
                if (error.status === 401) {
                    localStorage.removeItem('token')
                    window.location.assign('/')
                }
                return error
            },
        }),
        /* Deleting the profile of the user. */
        deleteProfile: builder.mutation({
            query: () => ({
                url: '/me',
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('token')}`,
                },
                method: 'DELETE',
            }),
            transformResponse: (response: IResponse<{ message: string }>) =>
                response.data.message,
            transformErrorResponse: (error) => {
                if (error.status === 401) {
                    localStorage.removeItem('token')
                    window.location.assign('/')
                }
                return error
            },
        }),
        // Avatar
        /* A function that is being called to get the avatar of the user. */
        getAvatar: builder.query({
            query: () => ({
                url: '/me/avatar',
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('token')}`,
                },
                method: 'GET',
            }),
            transformResponse: (response: IResponse<{ uri: string }>) =>
                response.data.uri,
            transformErrorResponse: (error) => {
                if (error.status === 401) {
                    localStorage.removeItem('token')
                    window.location.assign('/')
                }
                return error
            },
        }),
        /* A mutation that is updating the avatar of the user. */
        updateAvatar: builder.mutation({
            query: (body) => ({
                url: '/me/avatar',
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('token')}`,
                    ContentType: 'multipart/form-data',
                },
                method: 'PUT',
                body,
            }),
            transformResponse: (response: IResponse<{ message: string }>) =>
                response.data.message,
            transformErrorResponse: (error) => {
                if (error.status === 401) {
                    localStorage.removeItem('token')
                    window.location.assign('/')
                }
                return error
            },
        }),
        // Applets
        /* A function that is returning a builder.query object. */
        getApplets: builder.query({
            query: () => ({
                url: '/applet',
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('token')}`,
                },
                validateStatus: (status) =>
                    status.status == 200 || status.status == 204,
            }),
            transformResponse: (
                response: IResponse<{ applets: IApplet[] }>
            ) => {
                if (!response) return []
                return response.data.applets
            },
            transformErrorResponse: (error) => {
                if (error.status === 401) {
                    localStorage.removeItem('token')
                    window.location.assign('/')
                }
                return error
            },
        }),
        /* A function that is being called in the AppletBuilder component. */
        getNewApplet: builder.query({
            query: ({ field }) => ({
                url: '/applet/new' + (field ? `?field=${field}` : ''),
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('token')}`,
                },
            }),
            transformResponse: (
                response: IResponse<{ action?: IArea; reactions?: IArea[] }>
            ) => response.data,
            transformErrorResponse: (error) => {
                if (error.status === 401) {
                    localStorage.removeItem('token')
                    window.location.assign('/')
                }
                return error
            },
        }),
        /* A mutation that is being made to the database. */
        addStateToNewApplet: builder.mutation({
            query: (body) => ({
                url: '/applet/new',
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('token')}`,
                },
                method: 'PUT',
                body,
            }),
            transformResponse: (response: IResponse<{ message: string }>) =>
                response.data.message,
            transformErrorResponse: (error) => {
                if (error.status === 401) {
                    localStorage.removeItem('token')
                    window.location.assign('/')
                }
                return error
            },
        }),
        /* A mutation that is being called in the AppletForm component. */
        submitNewApplet: builder.mutation({
            query: (body) => ({
                url: '/applet/new',
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('token')}`,
                },
                method: 'POST',
                body,
            }),
            transformResponse: (response: IResponse<{ message: string }>) =>
                response.data.message,
            transformErrorResponse: (error) => {
                if (error.status === 401) {
                    localStorage.removeItem('token')
                    window.location.assign('/')
                }
                return error
            },
        }),
        /* Deleting the new applet. */
        deleteNewApplet: builder.mutation({
            query: ({ type, number }) => ({
                url: '/applet/new',
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('token')}`,
                },
                method: 'DELETE',
                params: { type: type, number: number },
            }),
            transformResponse: (response: IResponse<{ message: string }>) =>
                response.data.message,
            transformErrorResponse: (error) => {
                if (error.status === 401) {
                    localStorage.removeItem('token')
                    window.location.assign('/')
                }
                return error
            },
        }),
        /* A function that is being called to get the applet. */
        getApplet: builder.query({
            query: (id) => ({
                url: `/applet/${id}`,
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('token')}`,
                },
                method: 'GET',
            }),
            transformResponse: (response: IResponse<{ applet: IApplet }>) =>
                response.data.applet,
            transformErrorResponse: (error) => {
                if (error.status === 401) {
                    localStorage.removeItem('token')
                    window.location.assign('/')
                }
                return error
            },
        }),
        /* A function that is being called in the AppletReactions component. It is a function that is
        being called in the AppletReactions component. */
        getAppletReactions: builder.query({
            query: (id) => ({
                url: `/applet/${id}/reactions`,
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('token')}`,
                },
                method: 'GET',
            }),
            transformResponse: (response: IResponse<{ reactions: IArea[] }>) =>
                response.data.reactions,
            transformErrorResponse: (error) => {
                if (error.status === 401) {
                    localStorage.removeItem('token')
                    window.location.assign('/')
                }
                return error
            },
        }),
        /* A mutation that is used to get the reactions of an applet. */
        getAppletReactionsMutation: builder.mutation({
            query: (id) => ({
                url: `/applet/${id}/reactions`,
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('token')}`,
                },
                method: 'GET',
            }),
            transformResponse: (response: IResponse<{ reactions: IArea[] }>) =>
                response.data.reactions,
            transformErrorResponse: (error) => {
                if (error.status === 401) {
                    localStorage.removeItem('token')
                    window.location.assign('/')
                }
                return error
            },
        }),
        /* A mutation that is modifying the applet. */
        modifyApplet: builder.mutation({
            query: ({ id, ...body }) => ({
                url: `/applet/${id}`,
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('token')}`,
                },
                method: 'PUT',
                body,
            }),
            transformResponse: (response: IResponse<{ message: string }>) =>
                response.data.message,
            transformErrorResponse: (error) => {
                if (error.status === 401) {
                    localStorage.removeItem('token')
                    window.location.assign('/')
                }
                return error
            },
        }),
        /* A mutation that updates the applet activity. */
        updateAppletActivity: builder.mutation({
            query: ({ id, active }) => ({
                url: `/applet/${id}`,
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('token')}`,
                },
                params: { active },
                method: 'PATCH',
            }),
            transformResponse: (response: IResponse<{ message: string }>) =>
                response.data.message,
            transformErrorResponse: (error) => {
                if (error.status === 401) {
                    localStorage.removeItem('token')
                    window.location.assign('/')
                }
                return error
            },
        }),
        /* Deleting the applet with the given id. */
        deleteApplet: builder.mutation({
            query: (id) => ({
                url: `/applet/${id}`,
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('token')}`,
                },
                method: 'DELETE',
            }),
            transformResponse: (response: IResponse<{ message: string }>) =>
                response.data.message,
            transformErrorResponse: (error) => {
                if (error.status === 401) {
                    localStorage.removeItem('token')
                    window.location.assign('/')
                }
                return error
            },
        }),
        /* A mutation that is being called in the startApplet function. */
        startApplet: builder.mutation({
            query: (id) => ({
                url: `/applet/${id}/start`,
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('token')}`,
                },
                method: 'PUT',
            }),
            transformResponse: (response: IResponse<{ message: string }>) =>
                response.data.message,
            transformErrorResponse: (error) => {
                if (error.status === 401) {
                    localStorage.removeItem('token')
                    window.location.assign('/')
                }
                return error
            },
        }),
        /* A mutation that is being made to the server. */
        stopApplet: builder.mutation({
            query: (id) => ({
                url: `/applet/${id}/stop`,
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('token')}`,
                },
                method: 'PUT',
            }),
            transformResponse: (response: IResponse<{ message: string }>) =>
                response.data.message,
            transformErrorResponse: (error) => {
                if (error.status === 401) {
                    localStorage.removeItem('token')
                    window.location.assign('/')
                }
                return error
            },
        }),
        /* Creating a query object that will be used to make a GET request to the /store endpoint. */
        getStoreApplets: builder.query({
            query: () => ({
                url: '/store',
                method: 'GET',
            }),
            transformResponse: (response: IResponse<{ applets: IApplet[] }>) =>
                response.data.applets,
            transformErrorResponse: (error) => {
                if (error.status === 401) {
                    localStorage.removeItem('token')
                    window.location.assign('/')
                }
                return error
            },
        }),
        // Services API
        /* A mutation that is being made to the database. */
        fetchServiceApiEndpoint: builder.mutation({
            query: ({ service, endpoint, method, body }) => ({
                url: `/services/${service}/api${endpoint}`,
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('token')}`,
                },
                method,
                body: method !== 'GET' ? body : undefined,
            }),
            transformResponse: (
                response: IResponse<{ data: { data: any; fields: string[] } }>
            ) => response.data,
            transformErrorResponse: (error) => {
                if (error.status === 401) {
                    localStorage.removeItem('token')
                    window.location.assign('/')
                }
                return error
            },
        }),
    }),
})

export const {
    useGetAboutQuery,
    useAuthAccountMutation,
    useLogoutAccountMutation,
    useGetAuthorizationsQuery,
    useGetServiceAuthorizationsQuery,
    usePostAuthorizationMutation,
    useDeleteAuthorizationMutation,
    useGetProfileQuery,
    useModifyProfileMutation,
    useDeleteProfileMutation,
    useGetAvatarQuery,
    useUpdateAvatarMutation,
    useGetAppletsQuery,
    useGetAppletReactionsQuery,
    useGetAppletReactionsMutationMutation,
    useGetNewAppletQuery,
    useAddStateToNewAppletMutation,
    useSubmitNewAppletMutation,
    useDeleteNewAppletMutation,
    useGetAppletQuery,
    useModifyAppletMutation,
    useUpdateAppletActivityMutation,
    useDeleteAppletMutation,
    useStartAppletMutation,
    useStopAppletMutation,
    useGetStoreAppletsQuery,
    useFetchServiceApiEndpointMutation,
} = serviceApi
