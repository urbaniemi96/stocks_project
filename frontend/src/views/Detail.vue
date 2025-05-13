<template>
  <div class="p-8">
    <!-- Título y empresa -->
    <h1 class="text-3xl font-bold">{{ stock?.ticker }} - {{ stock?.company }}</h1>
    <p class="text-white-600 mb-4">
      {{ stock?.brokerage }} {{ stock?.action }} targets from {{ stock?.target_from }} to {{ stock?.target_to }}
    </p>

    <!-- Overview -->
    <div class="mt-6 grid grid-cols-2 gap-6">
      <!-- Estadísticas básicas -->
      <div class="p-4 bg-gray-900 rounded">
        <h2 class="font-semibold mb-2">Overview (Last {{ days }} days)</h2>
        <ul class="list-disc list-inside text-sm">
          <li>Min Close: {{ stats.min.toFixed(2) }}</li>  <!-- Calculado en script -->
          <li>Max Close: {{ stats.max.toFixed(2) }}</li>
          <li>Avg Close: {{ stats.avg.toFixed(2) }}</li>
          <li>
            Trend:
            <span :class="trend >= 0 ? 'text-green-600' : 'text-red-600'">
              {{ trend >= 0 ? '+' : '' }}{{ trend.toFixed(2) }}%
            </span>
          </li>
        </ul>
      </div>

      <!-- Close Price Chart -->
      <div class="p-4 bg-gray-900 rounded">
        <h2 class="font-semibold mb-2">Close Price</h2>
        <HistoryChart v-if="labels.length" :labels="labels" :data="closeData" />  <!-- Usamos labels y closeData computados -->
        <div v-else class="text-center text-sm text-gray-500">No data</div>
      </div>

      <!-- Volatility Chart -->
      <div class="p-4 bg-gray-900 rounded">
        <h2 class="font-semibold mb-2">Volatility (%)</h2>
        <HistoryChart
          v-if="riskReward.labels.length"
          :labels="riskReward.labels"
          :data="riskReward.volatilities"
        />
        <div v-else class="text-center text-sm text-gray-500">No data</div>
      </div>

      <!-- Potential Chart -->
      <div class="p-4 bg-gray-900 rounded">
        <h2 class="font-semibold mb-2">Potential (%)</h2>
        <HistoryChart
          v-if="riskReward.labels.length"
          :labels="riskReward.labels"
          :data="riskReward.potentials"
        />
        <div v-else class="text-center text-sm text-gray-500">No data</div>
      </div>

      <!-- Rating Distribution -->
      <div class="p-4 bg-gray-900 rounded col-span-2">
        <h2 class="font-semibold mb-2">Rating Distribution</h2>
        <ul class="grid grid-cols-3 gap-2 text-sm">
          <li v-for="(count, rating) in ratingDistribution" :key="rating">
            <span class="font-medium">{{ rating }}</span>: {{ count }}
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useStockStore } from '../stores/stocks'
import HistoryChart from '../components/HistoryChart.vue'
import type { HistoricalPoint } from '../stores/stocks'

const route = useRoute()
const ticker = route.params.ticker as string
const days = 30

const stocksStore = useStockStore()
const detail = computed(() => stocksStore.detail)

const stock = computed(() => detail.value?.stock)
const history = computed(() => detail.value?.history ?? [])
const riskReward = computed(() => detail.value?.riskReward ?? { labels: [], volatilities: [], potentials: [] })
const ratingDistribution = computed(() => detail.value?.ratingDistribution ?? {})

// Computed properties con typing explicito para evitar 'any'
const labels = computed<string[]>(() => history.value.map((h: HistoricalPoint) => h.Date.slice(0, 10)))
const closeData = computed<number[]>(() => history.value.map((h: HistoricalPoint) => h.Close))

// Estadísticas calculadas en el script
const stats = computed(() => {
  if (!history.value.length) return { min: 0, max: 0, avg: 0 }
  const closes = history.value.map((h: HistoricalPoint) => h.Close)
  const min = Math.min(...closes)
  const max = Math.max(...closes)
  const avg = closes.reduce((sum, c) => sum + c, 0) / closes.length
  return { min, max, avg }
})

// Tendencia calculada en el script
const trend = computed(() => {
  if (history.value.length < 2) return 0
  const first = history.value[0].Close
  const last = history.value[history.value.length - 1].Close
  return first ? ((last - first) / first) * 100 : 0
})

onMounted(async () => {
  try {
    await stocksStore.loadDetail(ticker, days)  // Llamado centralizado al store
  } catch (err) {
    console.error('Error loading stock detail:', err)
  }
})
</script>

<style scoped>
/* Tailwind styles */
</style>
