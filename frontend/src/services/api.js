const BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'



async function handleResponse(response){
    if (response.status === 204) return null
    const data = await response.json()
    if(!response.ok) throw new Error(data.error || 'Erro na requisição')
    
    return data
}

async function getAll(){
    const res = await fetch(`${BASE_URL}/users`,{
        signal: AbortSignal.timeout(5000)
    })
    
    return handleResponse(res)
}

async function create(payload){
    const res = await fetch(`${BASE_URL}/users`, {
        signal: AbortSignal.timeout(5000),
        method: 'POST',
        headers: { 'Content-Type': 'application/json'},
        body: JSON.stringify(payload)
    })
    return handleResponse(res)
}

async function update(id, payload){
    const res = await fetch(`${BASE_URL}/users/${id}`, {
        signal: AbortSignal.timeout(5000),
        method: 'PUT',
        headers: { 'Content-Type': 'application/json'},
        body: JSON.stringify(payload)
    })
    return handleResponse(res)
}

async function remove(id){
    const res = await fetch(`${BASE_URL}/users/${id}`, {
        signal: AbortSignal.timeout(5000),
        method: 'DELETE'
    })
    return handleResponse(res)
}

export const userService = { getAll, create, update, remove }