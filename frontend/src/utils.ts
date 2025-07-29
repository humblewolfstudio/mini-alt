export const getLocaleDateTime = (date: string) => {
    return new Date(date).toLocaleString().replace(',', '')
}

export const getLocaleDate = (date: string) => {
    return new Date(date).toLocaleDateString()
}