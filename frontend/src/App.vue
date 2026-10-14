<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import AdminView from './components/admin/AdminView.vue'
import CartPanel from './components/pos/CartPanel.vue'
import ProductGrid from './components/pos/ProductGrid.vue'
import { getProducts } from './services/productService'
import type { CartItem, Product } from './types/product'

const categories = ['Semua', 'Makanan', 'Minuman', 'Camilan']
const products = ref<Product[]>([])
const cartItems = ref<CartItem[]>([])
const selectedCategory = ref('Semua')
const searchQuery = ref('')
const activePage = ref<'cashier' | 'admin'>('cashier')
const isLoading = ref(true)
const loadError = ref('')
const recentlyAddedId = ref('')
const toastMessage = ref('')
let selectionTimeout: number | undefined
let toastTimeout: number | undefined

const filteredProducts = computed(() => {
  const query = searchQuery.value.trim().toLocaleLowerCase('id-ID')
  return products.value.filter((product) => {
    const matchesCategory = selectedCategory.value === 'Semua' || product.category === selectedCategory.value
    const matchesQuery = `${product.name} ${product.description}`.toLocaleLowerCase('id-ID').includes(query)
    return matchesCategory && matchesQuery
  })
})

const subtotal = computed(() =>
  cartItems.value.reduce((total, item) => total + item.product.price * item.quantity, 0),
)

function randomizeProductImages(menu: Product[]): Product[] {
  const images = menu.map((product) => product.imageUrl)
  for (let index = images.length - 1; index > 0; index -= 1) {
    const randomIndex = Math.floor(Math.random() * (index + 1))
    const currentImage = images[index]
    const randomImage = images[randomIndex]
    if (currentImage !== undefined && randomImage !== undefined) {
      images[index] = randomImage
      images[randomIndex] = currentImage
    }
  }
  return menu.map((product, index) => ({
    ...product,
    imageUrl: images[index] ?? product.imageUrl,
  }))
}

async function loadMenu() {
  isLoading.value = true
  loadError.value = ''
  try {
    products.value = randomizeProductImages(await getProducts())
  } catch {
    loadError.value = 'Menu belum bisa dimuat. Pastikan server Go berjalan di port 8080.'
  } finally {
    isLoading.value = false
  }
}

function addToCart(product: Product) {
  const existingItem = cartItems.value.find((item) => item.product.id === product.id)
  if (existingItem) {
    existingItem.quantity += 1
  } else {
    cartItems.value.push({ product, quantity: 1 })
  }

  recentlyAddedId.value = product.id
  toastMessage.value = product.name
  if (selectionTimeout !== undefined) window.clearTimeout(selectionTimeout)
  if (toastTimeout !== undefined) window.clearTimeout(toastTimeout)
  selectionTimeout = window.setTimeout(() => {
    recentlyAddedId.value = ''
  }, 850)
  toastTimeout = window.setTimeout(() => {
    toastMessage.value = ''
  }, 2200)
}

function dismissToast() {
  toastMessage.value = ''
  if (toastTimeout !== undefined) window.clearTimeout(toastTimeout)
}

function setQuantity(id: string, quantity: number) {
  if (quantity < 1) {
    removeFromCart(id)
    return
  }
  const item = cartItems.value.find((cartItem) => cartItem.product.id === id)
  if (item) item.quantity = quantity
}

function removeFromCart(id: string) {
  cartItems.value = cartItems.value.filter((item) => item.product.id !== id)
}

onMounted(loadMenu)
onUnmounted(() => {
  if (selectionTimeout !== undefined) window.clearTimeout(selectionTimeout)
  if (toastTimeout !== undefined) window.clearTimeout(toastTimeout)
})
</script>

<template>
  <div class="pos-app">
    <header class="topbar">
      <a class="brand" href="#menu" aria-label="Rasa POS, ke menu kasir">
        <span class="brand-mark" aria-hidden="true">R</span>
        <span class="brand-copy"><strong>rasa</strong><small>FOOD POS</small></span>
      </a>
      <div class="topbar-title">
        <strong>{{ activePage === 'cashier' ? 'Kasir' : 'Admin' }}</strong>
      </div>
      <nav class="topbar-nav" aria-label="Navigasi utama">
        <button
          type="button"
          :class="{ active: activePage === 'cashier' }"
          :aria-current="activePage === 'cashier' ? 'page' : undefined"
          @click="activePage = 'cashier'"
        >Kasir</button>
        <button
          type="button"
          :class="{ active: activePage === 'admin' }"
          :aria-current="activePage === 'admin' ? 'page' : undefined"
          @click="activePage = 'admin'"
        >Admin</button>
      </nav>
      <div class="connection-status" :class="{ 'is-offline': loadError }" role="status">
        <span class="status-dot"></span>
        {{ loadError ? 'API terputus' : isLoading ? 'Menghubungkan' : 'API terhubung' }}
      </div>
    </header>

    <main v-if="activePage === 'cashier'" id="menu" class="pos-layout">
      <section class="menu-section" aria-labelledby="menu-title">
        <div class="welcome-row">
          <div>
            <h1 id="menu-title">Mau pesan apa?</h1>
            <p class="welcome-copy">Pilih menu favorit, lalu masukkan ke keranjang.</p>
          </div>
          <label class="search-box">
            <svg class="search-icon" aria-hidden="true" viewBox="0 0 24 24">
              <circle cx="10.8" cy="10.8" r="6.8" />
              <path d="m16 16 4 4" />
            </svg>
            <input v-model="searchQuery" type="search" placeholder="Cari menu..." aria-label="Cari menu" />
            <kbd>/</kbd>
          </label>
        </div>

        <nav class="category-list" aria-label="Kategori menu">
          <button
            v-for="category in categories"
            :key="category"
            type="button"
            class="category-button"
            :class="{ active: selectedCategory === category }"
            :aria-pressed="selectedCategory === category"
            @click="selectedCategory = category"
          >
            {{ category }}
          </button>
        </nav>

        <p v-if="loadError" class="error-banner" role="alert">
          {{ loadError }}
          <button type="button" @click="loadMenu">Coba lagi</button>
        </p>
        <div v-else-if="isLoading" class="loading-state" role="status">Sedang memuat menu...</div>
        <ProductGrid
          v-else
          :products="filteredProducts"
          :recently-added-id="recentlyAddedId"
          @add="addToCart"
        />
      </section>

      <CartPanel
        :items="cartItems"
        :subtotal="subtotal"
        @quantity-change="setQuantity"
        @remove="removeFromCart"
      />
    </main>

    <AdminView v-else />

    <Transition name="toast">
      <div v-if="toastMessage" class="toast-notice" role="status" aria-live="polite">
        <svg class="toast-check" aria-hidden="true" viewBox="0 0 24 24"><path d="m5 12 4.5 4.5L19 7" /></svg>
        <span class="toast-copy">
          <strong>Berhasil ditambahkan</strong>
          <span>{{ toastMessage }} masuk ke keranjang</span>
        </span>
        <button type="button" aria-label="Tutup notifikasi" @click="dismissToast">
          <svg aria-hidden="true" viewBox="0 0 20 20"><path d="m5 5 10 10m0-10L5 15" /></svg>
        </button>
      </div>
    </Transition>
  </div>
</template>
