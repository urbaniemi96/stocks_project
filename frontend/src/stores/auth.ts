import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: {
      id: 'demo-user',
      role: 'admin' as 'admin' | 'user',
    },
  }),
  getters: {
    isAdmin: (state) => state.user.role === 'admin'
  }
})

/*
Para usar en un componente

<template>
  <button
    v-if="auth.isAdmin"
    @click="stockStore.triggerRecompute()"
  >
    Recompute recommendations
  </button>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { useStockStore } from '@/stores/stocks'

const auth = useAuthStore()
const stockStore = useStockStore()
</script>

*/