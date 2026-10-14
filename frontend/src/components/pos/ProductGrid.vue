<script setup lang="ts">
import { formatRupiah } from '../../services/currency'
import type { Product } from '../../types/product'

defineProps<{
  products: Product[]
  recentlyAddedId: string
}>()

const emit = defineEmits<{ add: [product: Product] }>()
</script>

<template>
  <div v-if="products.length" class="product-grid">
    <button
      v-for="product in products"
      :key="product.id"
      type="button"
      class="product-card"
      :class="{ 'is-added': recentlyAddedId === product.id }"
      :disabled="!product.available"
      @click="emit('add', product)"
    >
      <span class="product-art">
        <img
          class="product-image"
          :src="product.imageUrl"
          alt=""
          width="720"
          height="520"
          loading="lazy"
          decoding="async"
        />
        <span class="product-category">{{ product.category }}</span>
      </span>
      <span class="product-info">
        <span class="product-name">{{ product.name }}</span>
        <span class="product-description">{{ product.description }}</span>
        <span class="product-price">{{ formatRupiah(product.price) }}</span>
      </span>
      <span v-if="recentlyAddedId === product.id" class="added-badge">
        <svg aria-hidden="true" viewBox="0 0 20 20"><path d="m4 10 4 4 8-8" /></svg>
        Ditambahkan
      </span>
      <span v-if="!product.available" class="sold-out">Habis</span>
    </button>
  </div>
  <div v-else class="empty-menu">Menu tidak ditemukan. Coba kata kunci lain.</div>
</template>
