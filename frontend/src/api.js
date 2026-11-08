// Prefer Vite env variable `VITE_API_BASE` if provided, otherwise use the local backend.
// Support either VITE_API_BASE or VITE_API_BASE_URL (docker-compose uses VITE_API_BASE_URL)
export const API_BASE_URL = import.meta.env.VITE_API_BASE || import.meta.env.VITE_API_BASE_URL || 'http://localhost:9000';

export function clearAuthSession() {
  localStorage.removeItem('authToken');
  localStorage.removeItem('authUser');
  localStorage.removeItem('rememberedEmail');
}

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

  const apiPath = path.startsWith('/') ? path : `/${path}`;
  if (response.status === 401 && apiPath !== '/api/login') {
    clearAuthSession();
    if (window.location.pathname !== '/login' || window.location.search !== '?sessionExpired=1') {
      window.location.assign('/login?sessionExpired=1');
    }
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