import router from '../main'

type Bucket = {
    Id: number,
    Name: string,
    NumberObjects: number,
    Size: number,
    CreatedAt: string
}

export const fetchListBuckets = async (): Promise<Bucket[]> => {
    const res = await fetch('/api/buckets')

    if (res.status === 401) {
        router.push('/login')
        return []
    }

    if(res.ok) {
        const data = await res.json()
        if(data) return data
    }

    return []
};

interface CreateBucketRequest {
    name: string
}

export const fetchCreateBucket = async (request: CreateBucketRequest) => {
    const res = await fetch('/api/buckets', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(request),
    });

    if (res.status === 401) {
        router.push('/login')
        return
    }

    router.push('/buckets')
}
