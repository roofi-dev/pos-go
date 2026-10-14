import type { Product } from '../types/product'

export async function getProducts(): Promise<Product[]> {
  const response = await fetch('/api/products')
  if (!response.ok) {
    throw new Error(`Gagal memuat menu (${response.status})`)
  }
  return response.json() as Promise<Product[]>
}
