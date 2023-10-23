export const axiosLogger = request => {
    // console.log('Request: ', JSON.stringify(request, null, 2))
    return request
}
export const axiosLoggerResp = response => {
    // console.log('Response: ', JSON.stringify(response.data, null, 2))
    return response
}
export const generateQuery = (data) => {
    return Object.entries(data).map((kv) => Array.isArray(kv[1]) ? kv[1].map(v => `${encodeURIComponent(kv[0])}[]=${encodeURIComponent(v)}`).join('&') : `${encodeURIComponent(kv[0])}=${encodeURIComponent(kv[1])}`).join('&')
}

export const genUsernameHeaders = (username) => ({
    headers: {
        'X-USER-NAME': String(username), 'Content-type': 'application/json'
    }
})
