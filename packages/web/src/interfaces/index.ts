export interface IStore {
    priority: number // (0=highest)
    description: string
    required: boolean
    type: string
    values: Array<string>
    need_fields: Array<string>
    value: any
}

export interface IAction {
    name: string
    description: string
    components: Array<string>
    store: object
    wip: boolean
}

export interface IArea {
    id: string
    type: string
    service: string
    name: string
    store: object
}

export interface IServiceMore {
    avatar: boolean
    color: string
}

export interface IService {
    name: string
    description: string
    authenticator: IAuthenticator | null
    more: IServiceMore | null
    actions: Array<IAction>
    reactions: Array<IAction>
}

export interface IMeta {
    authenticator: string
    applets: Array<string>
}

export interface IAuthenticator {
    name: string
    more: IServiceMore
    authorization_uri: string
    enabled: boolean
}

export interface IApplet {
    id: string
    name: string
    description: string
    public: boolean
    action: string
    active: boolean
    status: string
}

export interface IUser {
    id: string
    username: string
    email: string
}

export interface IServiceEnable {
    name: string
    enabled: boolean
}

export interface IError {
    code: string
    error: string
}

export interface IResponse<T> {
    code: number
    data: T
}
