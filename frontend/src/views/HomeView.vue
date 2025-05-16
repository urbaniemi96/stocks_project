<template>
  <div class="p-8 space-y-6">
    <h1 class="text-2xl font-bold">Stock Dashboard</h1>
    <div class="space-x-2">
      <button @click="refreshAll" class="px-4 py-2 bg-blue-600 text-white rounded">Traer de la API y cargar en la BD</button>
      <button @click="enrichAll" class="px-4 py-2 bg-blue-600 text-white rounded">Enriquecer desde Yahoo (puede demorar horas)</button>
      <button @click="showRecommendations" class="px-4 py-2 bg-green-600 text-white rounded">Mostrar Recomendaciones</button>
      <button @click="recalcRecommendations" class="px-4 py-2 bg-green-600 text-white rounded">Recalcular Reecomendaciones</button>
      <div v-if="taskId">
        Progreso: {{ status.pages_fetched }} páginas<br/>
        Estado: {{ status.status }}
      </div>
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
import { useStockStore } from '../stores/stocks'
import { storeToRefs } from 'pinia'
import { useRouter } from 'vue-router'

const router = useRouter()

const store = useStockStore()
const { taskId, status } = storeToRefs(store)

const stocksTable = ref<HTMLTableElement>()
let dataTable: any = null

async function refreshAll() {
  await store.fetchAndStore()
}

async function enrichAll() {
  await store.fetchAndEnrich()
}

async function showRecommendations() {
  router.push({ name: 'Recommendations' }) 
}

async function recalcRecommendations() {
  await store.triggerRecalculate()
}

onMounted(() => {
  dataTable = $(stocksTable.value!).DataTable({
    processing: true,
    serverSide: true,
    ajax: { url: 'http://localhost:8080/stocks', type: 'GET' },
    columns: [
      { data: 'ticker' },
      { data: 'company' },
      { data: 'target_from' },
      { data: 'target_to' },
      {
        data: null,
        orderable: false,
        searchable: false,
        defaultContent: `
          <button class="view-detail px-2 py-1 bg-indigo-600 text-white rounded">
            Ver detalle
          </button>
        `
      }
    ],
    order: [[0, 'asc']],
    pageLength: 10,
    lengthMenu: [[10, 20, 50], [10, 20, 50]],
  })
  // 2. Delegar el clic del botón para redirigir vía Vue Router
  $(stocksTable.value!).on('click', 'button.view-detail', function() {
    const rowData = dataTable.row($(this).closest('tr')).data()
    router.push({ name: 'stock-detail', params: { ticker: rowData.ticker } })
  })


})


onBeforeUnmount(() => {
  if (dataTable) dataTable.destroy(true)
})
</script>

<style scoped>
</style>
