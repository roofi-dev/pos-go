export interface Product {
  id: string
  name: string
  category: string
  description: string
  price: number
  imageUrl: string
  stockQuantity: number
  lowStockThreshold: number
  isActive: boolean
  available: boolean
}

export type ProductInput = Pick<
  Product,
  'name' | 'category' | 'description' | 'price' | 'imageUrl' | 'lowStockThreshold' | 'isActive'
>

export interface StockAdjustment {
  delta: number
  reason: string
}

export interface AdminCredentials {
  username: string
  password: string
}

export interface CartItem {
  product: Product
  quantity: number
}
