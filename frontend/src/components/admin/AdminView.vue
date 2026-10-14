<script setup lang="ts">
import { computed, ref } from 'vue'
import { formatRupiah } from '../../services/currency'
import {
  adjustAdminStock,
  archiveAdminProduct,
  createAdminProduct,
  getAdminProducts,
  updateAdminProduct,
} from '../../services/adminService'
import type { AdminCredentials, Product, ProductInput } from '../../types/product'

const credentials = ref<AdminCredentials | null>(null)
const username = ref('')
const password = ref('')
const products = ref<Product[]>([])
const searchQuery = ref('')
const isLoading = ref(false)
const isSaving = ref(false)
const errorMessage = ref('')
const successMessage = ref('')
const showProductForm = ref(false)
const editingProductId = ref('')
const stockProduct = ref<Product | null>(null)
const stockDelta = ref(1)
const stockReason = ref('Restock')

const emptyForm = (): ProductInput => ({
  name: '',
  category: 'Makanan',
  description: '',
  price: 10000,
  imageUrl: '',
  lowStockThreshold: 5,
  isActive: true,
})

const form = ref<ProductInput>(emptyForm())
const filteredProducts = computed(() => {
  const query = searchQuery.value.trim().toLocaleLowerCase('id-ID')
  return products.value.filter((product) =>
    `${product.name} ${product.category}`.toLocaleLowerCase('id-ID').includes(query),
  )
})
const activeCount = computed(() => products.value.filter((product) => product.isActive).length)
const lowStockCount = computed(() =>
  products.value.filter((product) => product.isActive && product.stockQuantity <= product.lowStockThreshold).length,
)
const stockTotal = computed(() =>
  products.value.filter((product) => product.isActive).reduce((total, product) => total + product.stockQuantity, 0),
)

async function signIn() {
  errorMessage.value = ''
  isLoading.value = true
  const attempt = { username: username.value.trim(), password: password.value }
  try {
    products.value = await getAdminProducts(attempt)
    credentials.value = attempt
    password.value = ''
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Gagal masuk ke admin.'
  } finally {
    isLoading.value = false
  }
}

async function refreshProducts() {
  if (!credentials.value) return
  isLoading.value = true
  errorMessage.value = ''
  try {
    products.value = await getAdminProducts(credentials.value)
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Gagal memuat produk.'
  } finally {
    isLoading.value = false
  }
}

function signOut() {
  credentials.value = null
  products.value = []
  username.value = ''
  password.value = ''
  errorMessage.value = ''
  successMessage.value = ''
}

function openNewProduct() {
  editingProductId.value = ''
  form.value = { ...emptyForm(), imageUrl: products.value[0]?.imageUrl ?? '' }
  showProductForm.value = true
  errorMessage.value = ''
}

function openEditProduct(product: Product) {
  editingProductId.value = product.id
  form.value = {
    name: product.name,
    category: product.category,
    description: product.description,
    price: product.price,
    imageUrl: product.imageUrl,
    lowStockThreshold: product.lowStockThreshold,
    isActive: product.isActive,
  }
  showProductForm.value = true
  errorMessage.value = ''
}

function closeProductForm() {
  showProductForm.value = false
  editingProductId.value = ''
}

async function saveProduct() {
  if (!credentials.value) return
  isSaving.value = true
  errorMessage.value = ''
  successMessage.value = ''
  const payload: ProductInput = {
    ...form.value,
    name: form.value.name.trim(),
    category: form.value.category.trim(),
    description: form.value.description.trim(),
    price: Number(form.value.price),
    lowStockThreshold: Number(form.value.lowStockThreshold),
  }
  try {
    if (editingProductId.value) {
      await updateAdminProduct(credentials.value, editingProductId.value, payload)
      successMessage.value = 'Perubahan menu berhasil disimpan.'
    } else {
      await createAdminProduct(credentials.value, payload)
      successMessage.value = 'Menu baru berhasil ditambahkan. Atur stok awal melalui tombol Stok.'
    }
    closeProductForm()
    await refreshProducts()
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Gagal menyimpan produk.'
  } finally {
    isSaving.value = false
  }
}

function openStockAdjustment(product: Product) {
  stockProduct.value = product
  stockDelta.value = 1
  stockReason.value = 'Restock'
  errorMessage.value = ''
}

async function saveStockAdjustment() {
  if (!credentials.value || !stockProduct.value) return
  isSaving.value = true
  errorMessage.value = ''
  successMessage.value = ''
  try {
    await adjustAdminStock(credentials.value, stockProduct.value.id, {
      delta: Number(stockDelta.value),
      reason: stockReason.value.trim(),
    })
    stockProduct.value = null
    successMessage.value = 'Stok berhasil diperbarui.'
    await refreshProducts()
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Gagal mengubah stok.'
  } finally {
    isSaving.value = false
  }
}

async function toggleProduct(product: Product) {
  if (!credentials.value) return
  errorMessage.value = ''
  successMessage.value = ''
  try {
    await (product.isActive
      ? archiveAdminProduct(credentials.value, product.id)
      : updateAdminProduct(credentials.value, product.id, {
          name: product.name,
          category: product.category,
          description: product.description,
          price: product.price,
          imageUrl: product.imageUrl,
          lowStockThreshold: product.lowStockThreshold,
          isActive: true,
        }))
    successMessage.value = product.isActive ? 'Menu diarsipkan.' : 'Menu diaktifkan kembali.'
    await refreshProducts()
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Gagal mengubah status menu.'
  }
}

function stockStatus(product: Product): string {
  if (product.stockQuantity === 0) return 'Stok habis'
  if (product.stockQuantity <= product.lowStockThreshold) return 'Stok menipis'
  return 'Tersedia'
}
</script>

<template>
  <main v-if="!credentials" class="admin-login-page">
    <section class="admin-login-card" aria-labelledby="admin-login-title">
      <span class="admin-login-mark" aria-hidden="true">R</span>
      <h1 id="admin-login-title">Admin Rasa POS</h1>
      <p>Masuk untuk mengelola menu dan persediaan.</p>
      <form class="admin-login-form" @submit.prevent="signIn">
        <label>
          Username admin
          <input v-model="username" autocomplete="username" required />
        </label>
        <label>
          Password
          <input v-model="password" type="password" autocomplete="current-password" required />
        </label>
        <p v-if="errorMessage" class="admin-error" role="alert">{{ errorMessage }}</p>
        <button class="admin-primary-button" type="submit" :disabled="isLoading">
          {{ isLoading ? 'Memeriksa...' : 'Masuk ke admin' }}
        </button>
      </form>
      <p class="admin-login-note">Kredensial berasal dari environment backend dan hanya disimpan selama halaman ini terbuka.</p>
    </section>
  </main>

  <main v-else class="admin-page">
    <header class="admin-page-heading">
      <div>
        <h1>Menu & persediaan</h1>
        <p>Kelola produk yang tampil di kasir dan pantau stok per menu.</p>
      </div>
      <div class="admin-page-actions">
        <button class="admin-secondary-button" type="button" :disabled="isLoading" @click="refreshProducts">
          Muat ulang
        </button>
        <button class="admin-secondary-button" type="button" @click="signOut">Keluar</button>
        <button class="admin-primary-button" type="button" @click="openNewProduct">Tambah menu</button>
      </div>
    </header>

    <p v-if="errorMessage" class="admin-feedback admin-error" role="alert">{{ errorMessage }}</p>
    <p v-else-if="successMessage" class="admin-feedback admin-success" role="status">{{ successMessage }}</p>

    <section class="admin-stats" aria-label="Ringkasan persediaan">
      <article><span>Menu aktif</span><strong>{{ activeCount }}</strong></article>
      <article><span>Stok menipis</span><strong>{{ lowStockCount }}</strong></article>
      <article><span>Total porsi tersedia</span><strong>{{ stockTotal }}</strong></article>
    </section>

    <section class="admin-table-card" aria-labelledby="admin-products-title">
      <div class="admin-table-heading">
        <div>
          <h2 id="admin-products-title">Daftar menu</h2>
          <p>{{ filteredProducts.length }} produk ditemukan</p>
        </div>
        <label class="admin-search">
          <span class="visually-hidden">Cari produk admin</span>
          <input v-model="searchQuery" type="search" placeholder="Cari menu atau kategori" />
        </label>
      </div>
      <div v-if="isLoading" class="admin-loading" role="status">Memuat produk...</div>
      <div v-else class="admin-table-scroll">
        <table class="admin-table">
          <thead>
            <tr>
              <th scope="col">Produk</th>
              <th scope="col">Harga</th>
              <th scope="col">Stok</th>
              <th scope="col">Status</th>
              <th scope="col"><span class="visually-hidden">Aksi</span></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="product in filteredProducts" :key="product.id">
              <td>
                <div class="admin-product-cell">
                  <img :src="product.imageUrl" :alt="''" width="52" height="52" loading="lazy" />
                  <span><strong>{{ product.name }}</strong><small>{{ product.category }}</small></span>
                </div>
              </td>
              <td>{{ formatRupiah(product.price) }}</td>
              <td>
                <span class="stock-value" :class="{ 'stock-warning': product.isActive && product.stockQuantity <= product.lowStockThreshold, 'stock-empty': product.stockQuantity === 0 }">
                  {{ product.stockQuantity }} porsi
                </span>
                <small class="stock-caption">{{ stockStatus(product) }} · minimum {{ product.lowStockThreshold }}</small>
              </td>
              <td><span class="admin-status" :class="product.isActive ? 'status-active' : 'status-archived'">{{ product.isActive ? 'Aktif' : 'Diarsipkan' }}</span></td>
              <td>
                <div class="admin-row-actions">
                  <button class="admin-text-button" type="button" @click="openEditProduct(product)">Edit</button>
                  <button class="admin-text-button" type="button" @click="openStockAdjustment(product)">Stok</button>
                  <button class="admin-text-button" type="button" @click="toggleProduct(product)">{{ product.isActive ? 'Arsipkan' : 'Aktifkan' }}</button>
                </div>
              </td>
            </tr>
            <tr v-if="filteredProducts.length === 0"><td class="admin-empty" colspan="5">Belum ada produk yang cocok.</td></tr>
          </tbody>
        </table>
      </div>
    </section>

    <div v-if="showProductForm" class="admin-dialog-backdrop" @click.self="closeProductForm">
      <section class="admin-dialog" role="dialog" aria-modal="true" aria-labelledby="product-form-title">
        <header class="admin-dialog-heading">
          <div><h2 id="product-form-title">{{ editingProductId ? 'Edit menu' : 'Tambah menu' }}</h2><p>Stok dikelola terpisah agar setiap perubahan bisa diberi alasan.</p></div>
          <button class="admin-close-button" type="button" aria-label="Tutup form" @click="closeProductForm">×</button>
        </header>
        <form class="admin-product-form" @submit.prevent="saveProduct">
          <label>Nama menu<input v-model="form.name" maxlength="120" required /></label>
          <label>Kategori<input v-model="form.category" maxlength="60" required /></label>
          <label>Harga (Rp)<input v-model.number="form.price" type="number" min="1" step="1" required /></label>
          <label>Batas stok menipis<input v-model.number="form.lowStockThreshold" type="number" min="0" step="1" required /></label>
          <label class="admin-form-wide">URL foto (HTTPS)<input v-model="form.imageUrl" type="url" placeholder="https://..." /></label>
          <label class="admin-form-wide">Deskripsi<textarea v-model="form.description" maxlength="300" rows="3"></textarea></label>
          <label class="admin-checkbox"><input v-model="form.isActive" type="checkbox" /> Tampilkan di kasir</label>
          <div class="admin-form-actions">
            <button class="admin-secondary-button" type="button" @click="closeProductForm">Batal</button>
            <button class="admin-primary-button" type="submit" :disabled="isSaving">{{ isSaving ? 'Menyimpan...' : 'Simpan menu' }}</button>
          </div>
        </form>
      </section>
    </div>

    <div v-if="stockProduct" class="admin-dialog-backdrop" @click.self="stockProduct = null">
      <section class="admin-dialog admin-stock-dialog" role="dialog" aria-modal="true" aria-labelledby="stock-form-title">
        <header class="admin-dialog-heading">
          <div><h2 id="stock-form-title">Atur stok</h2><p>{{ stockProduct.name }} · stok saat ini {{ stockProduct.stockQuantity }} porsi</p></div>
          <button class="admin-close-button" type="button" aria-label="Tutup form stok" @click="stockProduct = null">×</button>
        </header>
        <form class="admin-product-form" @submit.prevent="saveStockAdjustment">
          <label>Perubahan jumlah (+ masuk, − keluar)<input v-model.number="stockDelta" type="number" step="1" required /></label>
          <label>Alasan perubahan<input v-model="stockReason" maxlength="200" required /></label>
          <div class="admin-form-actions">
            <button class="admin-secondary-button" type="button" @click="stockProduct = null">Batal</button>
            <button class="admin-primary-button" type="submit" :disabled="isSaving || stockDelta === 0">{{ isSaving ? 'Menyimpan...' : 'Simpan perubahan stok' }}</button>
          </div>
        </form>
      </section>
    </div>
  </main>
</template>
