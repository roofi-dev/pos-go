import type { AdminCredentials, Product, ProductInput, StockAdjustment } from '../types/product'

function basicAuthorization(credentials: AdminCredentials): string {
  const bytes = new TextEncoder().encode(`${credentials.username}:${credentials.password}`)
  const binary = Array.from(bytes, (byte) => String.fromCharCode(byte)).join('')
  return `Basic ${btoa(binary)}`
}

async function adminRequest<T>(
  path: string,
  credentials: AdminCredentials,
  init: RequestInit = {},
): Promise<T> {
  const response = await fetch(path, {
    ...init,
    headers: {
      'Content-Type': 'application/json',
      ...init.headers,
      Authorization: basicAuthorization(credentials),
    },
  })

  if (!response.ok) {
    const body = await response.json().catch(() => ({})) as { error?: string }
    throw new Error(body.error || `Permintaan gagal (${response.status})`)
  }

  return response.json() as Promise<T>
}

export function getAdminProducts(credentials: AdminCredentials): Promise<Product[]> {
  return adminRequest('/api/admin/products', credentials)
}

export function createAdminProduct(credentials: AdminCredentials, product: ProductInput): Promise<Product> {
  return adminRequest('/api/admin/products', credentials, {
    method: 'POST',
    body: JSON.stringify(product),
  })
}

export function updateAdminProduct(
  credentials: AdminCredentials,
  id: string,
  product: ProductInput,
): Promise<Product> {
  return adminRequest(`/api/admin/products/${encodeURIComponent(id)}`, credentials, {
    method: 'PUT',
    body: JSON.stringify(product),
  })
}

export function archiveAdminProduct(credentials: AdminCredentials, id: string): Promise<{ status: string }> {
  return adminRequest(`/api/admin/products/${encodeURIComponent(id)}`, credentials, {
    method: 'DELETE',
  })
}

export function adjustAdminStock(
  credentials: AdminCredentials,
  id: string,
  adjustment: StockAdjustment,
): Promise<Product> {
  return adminRequest(`/api/admin/products/${encodeURIComponent(id)}/stock`, credentials, {
    method: 'POST',
    body: JSON.stringify(adjustment),
  })
}
