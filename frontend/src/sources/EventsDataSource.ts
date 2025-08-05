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
        return await res.json()
    }

    return []
}