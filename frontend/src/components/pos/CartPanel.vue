<script setup lang="ts">
import { computed } from 'vue'
import { formatRupiah } from '../../services/currency'
import type { CartItem } from '../../types/product'

const props = defineProps<{
  items: CartItem[]
  subtotal: number
}>()

const emit = defineEmits<{
  quantityChange: [id: string, quantity: number]
  remove: [id: string]
}>()

const totalQuantity = computed(() =>
  props.items.reduce((total, item) => total + item.quantity, 0),
)
</script>

<template>
  <aside class="cart-panel" aria-label="Keranjang belanja">
    <div class="cart-heading">
      <h2>Keranjang</h2>
      <span :key="totalQuantity" class="cart-count">{{ totalQuantity }}</span>
    </div>

    <div v-if="items.length" class="cart-items">
      <article v-for="item in items" :key="item.product.id" class="cart-item">
        <img
          class="cart-item-image"
          :src="item.product.imageUrl"
          alt=""
          width="56"
          height="56"
          loading="lazy"
          decoding="async"
        />
        <div class="cart-item-info">
          <strong>{{ item.product.name }}</strong>
          <span>{{ formatRupiah(item.product.price) }}</span>
          <div class="quantity-control" :aria-label="`Jumlah ${item.product.name}`">
            <button
              type="button"
              :aria-label="`Kurangi ${item.product.name}`"
              @click="emit('quantityChange', item.product.id, item.quantity - 1)"
            >−</button>
            <span>{{ item.quantity }}</span>
            <button
              type="button"
              :aria-label="`Tambah ${item.product.name}`"
              @click="emit('quantityChange', item.product.id, item.quantity + 1)"
            >+</button>
          </div>
        </div>
        <button
          type="button"
          class="remove-item"
          :aria-label="`Hapus ${item.product.name}`"
          @click="emit('remove', item.product.id)"
        >
          <svg aria-hidden="true" viewBox="0 0 20 20"><path d="M4 6h12m-10 0 .7 10h6.6L14 6m-6-2h4l1 2M8 9v4m4-4v4" /></svg>
        </button>
      </article>
    </div>
    <div v-else class="empty-cart">
      <svg aria-hidden="true" viewBox="0 0 48 48">
        <path d="M7 10h5l4 21h19l5-15H14" />
        <circle cx="19" cy="38" r="2" />
        <circle cx="33" cy="38" r="2" />
      </svg>
      <strong>Keranjang masih kosong</strong>
      <p>Pilih menu untuk memulai pesanan.</p>
    </div>

    <div class="cart-summary">
      <div class="summary-row">
        <span>Subtotal</span>
        <strong>{{ formatRupiah(subtotal) }}</strong>
      </div>
      <div class="summary-row total-row">
        <span>Total sementara</span>
        <strong>{{ formatRupiah(subtotal) }}</strong>
      </div>
      <button type="button" class="checkout-button" disabled>
        Checkout belum aktif
        <svg aria-hidden="true" viewBox="0 0 20 20"><path d="M3 10h13m-5-5 5 5-5 5" /></svg>
      </button>
      <p class="checkout-note">Tahap berikutnya: validasi dan pembayaran melalui API.</p>
    </div>
  </aside>
</template>
