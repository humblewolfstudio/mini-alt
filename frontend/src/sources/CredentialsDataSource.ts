import router from "../main"

type Credential = {
    Id: number,
    AccessKey: string,
    SecretKey: string,
    ExpiresAt: string,
    Status: boolean,
    Name: string,
    Description: string,
    CreatedAt: string
}

export const fetchListCredentials = async (): Promise<Credential[]> => {
    const res = await fetch('/api/credentials')

    if(res.status === 401) {
        router.push('/login')
    }

    if (res.ok) {
        return await res.json()
    }

    return []
}

interface CreateCredentialsRequest {
    expiresAt: string,
    name: string,
    description: string
}

interface CreateCredentialsResponse {
    access_key: string,
    secret_key: string
}

export const fetchCreateCredentials = async (request: CreateCredentialsRequest): Promise<CreateCredentialsResponse> => {
    const res = await fetch('/api/credentials', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(request)
    })

    if (res.status === 401) {
        router.push('/login')
    }

    if (res.ok) {
        const data = await res.json()
        if (data) return data
    }

    return null
}

interface DeleteCredentialsRequest {
    accessKey: string
}

export const fetchDeleteCredentials = async (request: DeleteCredentialsRequest) => {
    await fetch('/api/credentials/delete', {
        method: 'POST',
        body: JSON.stringify(request)
    })
}

interface EditCredentialsRequest {
    accessKey: string,
    name: string,
    description: string,
    status: boolean,
    expiresAt: string
}

export const fetchEditCredentials = async (request: EditCredentialsRequest) => {
    await fetch('/api/credentials/edit', {
        method: 'POST',
        body: JSON.stringify(request)
    })
}
