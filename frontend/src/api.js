// Prefer Vite env variable `VITE_API_BASE` if provided, otherwise fallback to localhost:8080
// Support either VITE_API_BASE or VITE_API_BASE_URL (docker-compose uses VITE_API_BASE_URL)
export const API_BASE_URL = import.meta.env.VITE_API_BASE || import.meta.env.VITE_API_BASE_URL || 'http://127.0.0.1:8080';

async function apiRequest(path, options = {}) {
  const headers = new Headers(options.headers || {});
  const token = localStorage.getItem('authToken');

  if (token) {
    headers.set('Authorization', `Bearer ${token}`);
  }

  if (options.body && !headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json');
  }

  let response;
  try {
    const baseUrl = API_BASE_URL.replace(/\/+$/, '');
    const apiPath = path.startsWith('/') ? path : `/${path}`;
    response = await fetch(`${baseUrl}${apiPath}`, { ...options, headers });
  } catch (error) {
    throw new Error(`Cannot connect to API at ${API_BASE_URL}: ${error.message}`);
  }

  const contentType = response.headers.get('content-type') || '';
  const data = contentType.includes('application/json')
    ? await response.json()
    : await response.text();

  if (!response.ok) {
    const message = typeof data === 'object' && data !== null
      ? data.error || data.message
      : data;
    throw new Error(message || `Request failed ${response.status}`);
  }

  return data;
}

export function apiGet(path) {
  return apiRequest(path);
}

export function apiPost(path, payload) {
  return apiRequest(path, { method: 'POST', body: JSON.stringify(payload) });
}

export function apiPut(path, payload) {
  return apiRequest(path, { method: 'PUT', body: JSON.stringify(payload) });
}

export function apiDelete(path) {
  return apiRequest(path, { method: 'DELETE' });
}