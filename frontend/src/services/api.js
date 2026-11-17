const BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'



async function handleResponse(response){
    if (response.status === 204) return null
    const data = await response.json()
    if(!response.ok) throw new Error(data.error || 'Erro na requisição')
    
    return data
}

async function request(endpoint, options = {}) {
  try {
    const res = await fetch(`${BASE_URL}${endpoint}`, {
      ...options,
      signal: AbortSignal.timeout(5000)
    })
    return await handleResponse(res)
  } catch (err) {
    if (err.name === 'TimeoutError') {
      throw new Error('O servidor demorou muito para responder. Tente novamente.')
    }
    if (err.name === 'TypeError') {
      throw new Error('Não foi possível conectar ao servidor. Verifique se a API está rodando...')
    }
    throw err
  }
}

async function getAll(){
    return request('/users')
}

async function create(payload){
    return request('/users', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload)
  })
}

async function update(id, payload){
    return request(`/users/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload)
  })
}

async function remove(id){
    return request(`/users/${id}`, {
    method: 'DELETE'
  })
}

export const userService = { getAll, create, update, remove }