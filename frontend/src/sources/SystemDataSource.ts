import router from "../main"

export interface SystemInfoResponse {
    NumberBuckets: number
    NumberObjects: number
    Usage: number
}

export const fetchSystemInfo = async(): Promise<SystemInfoResponse | null> => {
    const res = await fetch('/api/system/info')

    if (res.status === 401) {
        router.push('/login')
        return null
    }

    if (res.ok) {
        return await res.json()
    }

    return null
}

export interface SystemSpecsResponse {
    TotalCapacity: number
    FreeCapacity: number
    UsedCapacity: number
    DrivePath: string
}

export const fetchSystemSpecs = async(): Promise<SystemSpecsResponse | null> => {
    const res = await fetch('/api/system/specs')

    if (res.status === 401) {
        router.push('/login')
        return null
    }

    if (res.ok) {
        return await res.json()
    }

    return null
}