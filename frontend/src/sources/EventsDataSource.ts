import router from "../main"

type Event = {
    Id: number,
    Name: string,
    Description: string,
    BucketId: number,
    CreatedAt: string
}

export const fetchListEvents = async (): Promise<Event[]> => {
    const res = await fetch('/api/events')

    if (res.status === 401) {
        router.push('/login')
        return []
    }

    if (res.ok) {
        const data = await res.json()
        if (data) return data
    }

    return []
}

interface CreateEventRequest {
    name: string,
    description: string,
    bucket: number,
    endpoint: string,
    token: string,
    global: boolean
}

export const fetchCreateEvent = async (request: CreateEventRequest) => {
    const res = await fetch('/api/events', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            name: request.name,
            description: request.description,
            bucket: request.bucket,
            endpoint: request.endpoint,
            token: request.token,
            global: request.global
        })
    })

    if (res.status === 401) {
        router.push('/login')
        return
    }

    if (res.ok) {
        router.push('/events')
    } else {
        throw await res.json()
    }
}
