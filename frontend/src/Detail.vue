<template>
  <div class="p-8">
    <!-- Título y empresa -->
    <h1 class="text-3xl font-bold">{{ stock.ticker }} - {{ stock.company }}</h1>
    <p class="text-gray-600 mb-4">
      {{ stock.brokerage }} {{ stock.action }} targets from {{ stock.target_from }} to {{ stock.target_to }}
    </p>

    <!-- Overview -->
    <div class="mt-6 grid grid-cols-2 gap-6">
      <!-- Estadísticas básicas -->
      <div class="p-4 bg-gray-900 rounded">
        <h2 class="font-semibold mb-2">Overview (Last {{ days }} days)</h2>
        <ul class="list-disc list-inside text-sm">
          <li>Min Close: {{ stats.min.toFixed(2) }}</li>
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
        <HistoryChart v-if="labels.length" :labels="labels" :data="data" />
        <div v-else class="text-center text-sm text-gray-500">No data</div>
      </div>

      <!-- Volatility Chart -->
      <div class="p-4 bg-gray-900 rounded">
        <h2 class="font-semibold mb-2">Volatility (%) </h2>
        <HistoryChart v-if="riskReward.labels.length" :labels="riskReward.labels" :data="riskReward.volatilities" />
        <div v-else class="text-center text-sm text-gray-500">No data</div>
      </div>

      <!-- Potential Chart -->
      <div class="p-4 bg-gray-900 rounded">
        <h2 class="font-semibold mb-2">Potential (%)</h2>
        <HistoryChart v-if="riskReward.labels.length" :labels="riskReward.labels" :data="riskReward.potentials" />
        <div v-else class="text-center text-sm text-gray-500">No data</div>
      </div>

      <!-- Rating Distribution -->
      <div class="p-4 bg-gray-900 rounded col-span-2">
        <h2 class="font-semibold mb-2">Rating Distribution</h2>
        <ul class="grid grid-cols-3 gap-2 text-sm">
          <li v-for="(count, rating) in sortedRatings" :key="rating">
            <span class="font-medium">{{ rating }}</span>: {{ count }}
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { api } from './services/api'
import HistoryChart from './components/HistoryChart.vue'

interface HistoricalPoint { Date: string; Close: number }
interface Stock {
  ticker: string
  company: string
  brokerage: string
  action: string
  target_from: number
  target_to: number
}
interface RiskRewardResponse { labels: string[]; volatilities: number[]; potentials: number[] }
interface StockDetailResponse {
  stock: Stock
  history: HistoricalPoint[]
  riskReward: RiskRewardResponse
  ratingDistribution: Record<string, number>
}

const route = useRoute()
const ticker = route.params.ticker as string
const days = 30

const stock = ref<Stock>({ ticker, company: '', brokerage: '', action: '', target_from: 0, target_to: 0 })
const history = ref<HistoricalPoint[]>([])
const labels = ref<string[]>([])
const data = ref<number[]>([])
const stats = ref({ min: 0, max: 0, avg: 0 })
const trend = ref(0)
const riskReward = ref<RiskRewardResponse>({ labels: [], volatilities: [], potentials: [] })
const ratingDistribution = ref<Record<string, number>>({})

const sortedRatings = computed(() => 
  Object.entries(ratingDistribution.value)
    .sort((a, b) => b[1] - a[1])
    .reduce((obj, [k, v]) => ({ ...obj, [k]: v }), {} as Record<string, number>)
)

onMounted(async () => {
  try {
    const res = await api.get<StockDetailResponse>(`/stocks/${ticker}/detail?days=${days}`)
    const { stock: st, history: hist, riskReward: rr, ratingDistribution: rd } = res.data

    stock.value = st
    history.value = hist
    riskReward.value = rr
    ratingDistribution.value = rd

    // Close data
    labels.value = hist.map(item => item.Date.slice(0, 10))
    data.value = hist.map(item => item.Close)

    if (data.value.length) {
      stats.value.min = Math.min(...data.value)
      stats.value.max = Math.max(...data.value)
      stats.value.avg = data.value.reduce((a, b) => a + b, 0) / data.value.length
      const first = data.value[0]
      const last = data.value[data.value.length - 1]
      trend.value = first ? ((last - first) / first) * 100 : 0
    }
  } catch (err) {
    console.error('Error loading stock detail:', err)
  }
})
</script>

<style scoped>
/* Tailwind styles */
</style>
