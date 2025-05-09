<template>
  <div class="p-8 space-y-6">
    <h1 class="text-2xl font-bold">Stock Dashboard</h1>
    <div class="space-x-2">
      <button @click="refreshAll" class="px-4 py-2 bg-blue-600 text-white rounded">Sync & Load</button>
      <button @click="getRecommendation" class="px-4 py-2 bg-green-600 text-white rounded">Recommend</button>
      <div v-if="taskId">
        Progreso: {{ status.pages_fetched }} páginas<br/>
        Estado: {{ status.status }}
      </div>
    </div>

    <div v-if="recommended" class="p-4 bg-green-100 rounded">
      <strong>Top Pick:</strong> {{ recommended.ticker }} 
      (Δ {{ (recommended.target_to - recommended.target_from).toFixed(2) }})
    </div>

    <table class="min-w-full bg-white text-gray-800 text-sm">
      <thead>
        <tr class="bg-gray-200 text-left font-semibold">
          <th class="p-2">Ticker</th>
          <th class="p-2">Company</th>
          <th class="p-2">From</th>
          <th class="p-2">To</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="s in stocks" :key="s.ticker" class="border-t hover:bg-gray-100">
          <td class="p-2">{{ s.ticker }}</td>
          <td class="p-2">{{ s.company }}</td>
          <td class="p-2">{{ s.target_from }}</td>
          <td class="p-2">{{ s.target_to }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useStockStore } from './stores/stocks'
import { storeToRefs } from 'pinia'

const store = useStockStore()
const { list: stocks, recommended, taskId, status } = storeToRefs(store)

async function refreshAll() {
  await store.fetchAndStore()
  //await store.loadStocks()
}

async function getRecommendation() {
  await store.computeRecommendation()
}

onMounted(() => {
  store.loadStocks()
})
</script>

<style scoped>
/* espacio para estilos adicionales si quieres */
</style>