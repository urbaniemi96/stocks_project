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

    <table ref="stocksTable" class="min-w-full bg-white text-gray-800 text-sm">
      <thead>
        <tr class="bg-gray-200 font-semibold">
          <th>Ticker</th>
          <th>Company</th>
          <th>From</th>
          <th>To</th>
        </tr>
      </thead>
      <tbody></tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onBeforeUnmount, ref } from 'vue'
import $ from 'jquery'
import 'datatables.net'
import { useStockStore } from './stores/stocks'
import { storeToRefs } from 'pinia'

const store = useStockStore()
const { list: stocks, recommended, taskId, status } = storeToRefs(store)

const stocksTable = ref<HTMLTableElement>()
let dataTable: any = null

async function refreshAll() {
  await store.fetchAndStore()
  //await store.loadStocks()
}

async function getRecommendation() {
  await store.computeRecommendation()
}

onMounted(() => {
  //store.loadStocks()
  dataTable = $(stocksTable.value!).DataTable({
    processing: true,
    serverSide: true,
    ajax: { url: 'http://localhost:8080/stocks', type: 'GET' },
    columns: [
      { data: 'ticker' },
      { data: 'company' },
      { data: 'target_from' },
      { data: 'target_to' },
    ],
    order: [[0, 'asc']],
    pageLength: 10,
    lengthMenu: [[10, 20, 50], [10, 20, 50]],
  })
})

onBeforeUnmount(() => {
  if (dataTable) dataTable.destroy(true)
})

</script>

<style scoped>
/* espacio para estilos adicionales si quieres */
</style>