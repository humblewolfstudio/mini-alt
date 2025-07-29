export const getLocaleDateTime = (date: string) => {
    return new Date(date).toLocaleString().replace(',', '')
}

export const getLocaleDate = (date: string) => {
    return new Date(date).toLocaleDateString()
}

export const formatSize = (bytes: number) => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};