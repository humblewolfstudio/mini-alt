import router from "../main"

type User = {
    Id: number,
    Username: string,
    Password: string,
    Token: string,
    AccessKey: string,
    ExpiresAt: string,
    CreatedAt: string
}

export const fetchListUsers = async (): Promise<User[]> => {
    const res = await fetch('/api/users/list')

    if(res.status === 401) {
        router.push('/login')
    }

    if (res.ok) {
        return await res.json()
    }

    return []
}

interface CreateUserRequest {
    username: string,
    password: string,
    expiresAt: string
}

export const fetchCreateUser = async (request: CreateUserRequest): Promise<string> => {
    const res = await fetch('/api/users/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(request)
    })

    if(res.status === 401) {
        router.push('/login')
    }

    if (res.ok) {
        router.push('/users')
    } else {
        throw await res.json()
    }

    return ""
}