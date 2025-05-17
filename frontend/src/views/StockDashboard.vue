<template>
  <div class="container mx-auto p-8 space-y-6">
    <h1 class="text-3xl font-extrabold text-gray-900 dark:text-gray-100">Stock Dashboard</h1>
    <BackButton />
    <TopButton />
    <HomeButton />
    <table ref="stocksTable" class="min-w-full bg-white dark:bg-gray-800 dark:text-gray-200 text-sm divide-y">
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
import { useRouter } from 'vue-router'
import BackButton from '../components/BackButton.vue'
import TopButton from '../components/TopButton.vue'
import HomeButton from '../components/HomeButton.vue'

const router = useRouter()

const stocksTable = ref<HTMLTableElement>()
let dataTable: any = null

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
  // Redirigo para ver el detalle y le paso el ticker
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
