<template>
  <div class="container mx-auto p-8 space-y-10">
    <header class="flex flex-col lg:flex-row lg:justify-between items-start lg:items-center gap-4">
      <div>
        <h1 class="text-4xl font-extrabold text-gray-900 dark:text-gray-100">{{ stock?.ticker }}</h1>
        <p class="mt-2 text-lg font-bold text-gray-600 dark:text-gray-300">{{ stock?.company }}</p>
      </div>
      <div class="flex gap-4">
        <div class="bg-indigo-600 text-white dark:text-gray-300 px-4 py-2 rounded-md font-semibold">
          {{ stock?.action }}
        </div>
        <div class="text-gray-700 dark:text-gray-300 self-center">
          Broker: {{ stock?.brokerage }}
        </div>
      </div>
    </header>
    <BackButton />
    <TopButton />
    <HomeButton />
    <!-- Filtros -->
    <section class="bg-white dark:bg-gray-800 rounded-lg p-6 shadow-lg space-y-4">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-gray-100">Filtros</h2>
      <div class="flex items-end justify-center gap-4 text-gray-700 dark:text-gray-300">
        <!--div>
          <label class="block text-sm">Días</label>
          <input
            type="number"
            min="1"
            v-model.number="filters.days"
            class="w-full mt-1 p-2 bg-gray-700 rounded"
          />
        </div-->
        <div>
          <label class="block text-m">Fecha inicio</label>
          <input
            type="date"
            v-model="filters.start_date"
            class="w-full mt-1 p-2 rounded border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition text-center"
          />
        </div>
        <div>
          <label class="block text-m">Fecha fin</label>
          <input
            type="date"
            v-model="filters.end_date"
            class="w-full mt-1 p-2 rounded border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition text-center"
          />
        </div>
        <!--div>
          <label class="block text-sm">Orden</label>
          <select
            v-model="filters.order"
            class="w-full mt-1 p-2 bg-gray-700 rounded"
          >
            <option value="asc">Ascendente</option>
            <option value="desc">Descendente</option>
          </select>
        </div>
        <div>
          <label class="block text-sm">Precio mínimo</label>
          <input
            type="number"
            step="0.01"
            v-model.number="filters.min_price"
            class="w-full mt-1 p-2 bg-gray-700 rounded"
          />
        </div>
        <div>
          <label class="block text-sm">Precio máximo</label>
          <input
            type="number"
            step="0.01"
            v-model.number="filters.max_price"
            class="w-full mt-1 p-2 bg-gray-700 rounded"
          />
        </div>
        <div>
          <label class="block text-sm">Volumen mínimo</label>
          <input
            type="number"
            v-model.number="filters.min_volume"
            class="w-full mt-1 p-2 bg-gray-700 rounded"
          />
        </div-->
        <button
          @click="applyFilters"
          class="mt-4 bg-blue-600 hover:bg-blue-500 text-white dark:text-gray-300 px-4 py-2 rounded-lg font-semibold"
        >
          <font-awesome-icon :icon="['fas', 'search']" />
        </button>
      </div>
    </section>

    <!-- Overview -->
    <section class="bg-white dark:bg-gray-800 rounded-lg p-6 shadow-lg space-y-4">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-gray-100">Overview (Last 30 days)</h2>
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 text-gray-300">
        <div>
          <p class="text-m uppercase font-semibold text-gray-900 dark:text-gray-100">Min Close</p>
          <p class="text-lg font-medium text-gray-900 dark:text-gray-100">{{ stats.min.toFixed(2) }}</p>
        </div>
        <div>
          <p class="text-m uppercase font-semibold text-gray-900 dark:text-gray-100">Max Close</p>
          <p class="text-lg font-medium text-gray-900 dark:text-gray-100">{{ stats.max.toFixed(2) }}</p>
        </div>
        <div>
          <p class="text-m uppercase font-semibold text-gray-900 dark:text-gray-100">Avg Close</p>
          <p class="text-lg font-medium text-gray-900 dark:text-gray-100">{{ stats.avg.toFixed(2) }}</p>
        </div>
        <div class="sm:col-span-3">
          <p class="text-m uppercase font-semibold text-gray-900 dark:text-gray-100">Trend</p>
          <p :class="['text-lg font-medium', trend >= 0 ? 'text-green-600' : 'text-red-600']">
            {{ trend >= 0 ? '+' : '' }}{{ trend.toFixed(2) }}%
          </p>
        </div>
      </div>
    </section>

    <!-- Candlestick Chart -->
    <section class="bg-white dark:bg-gray-800 rounded-lg p-6 shadow-lg space-y-4 flex flex-col items-center justify-center">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-gray-100">Candlestick Chart (OHLC + Volume)</h2>
      <div class="w-200 h-100">
        <CandleChart v-if="history.length" :history="history" />
        <div v-else class="flex items-center justify-center h-full text-gray-500">No data available</div>
      </div>
    </section>

    <!-- Close Price Chart -->
    <section class="bg-white dark:bg-gray-800 rounded-lg p-6 shadow-lg space-y-4 flex flex-col items-center justify-center">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-gray-100">Close Price Over Time</h2>
      <div class="w-200 h-100">
        <HistoryChart v-if="labels.length" :labels="labels" :data="closeData" />
        <div v-else class="flex items-center justify-center h-full text-gray-500">No data available</div>
      </div>
    </section>

    <!-- Volatility Chart -->
     <section class="bg-white dark:bg-gray-800 rounded-lg p-6 shadow-lg space-y-4 flex flex-col items-center justify-center">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-gray-100">Volatility (%)</h2>
      <div class="w-200 h-100">
        <HistoryChart
          v-if="riskReward.labels.length"
          :labels="riskReward.labels"
          :data="riskReward.volatilities"
        />
        <div v-else class="flex items-center justify-center h-full text-gray-500">No data available</div>
      </div>
    </section>

    <!-- Potential Chart -->
    <section class="bg-white dark:bg-gray-800 rounded-lg p-6 shadow-lg space-y-4 flex flex-col items-center justify-center">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-gray-100">Potential (%)</h2>
      <div class="w-200 h-100">
        <HistoryChart
          v-if="riskReward.labels.length"
          :labels="riskReward.labels"
          :data="riskReward.potentials"
        />
        <div v-else class="flex items-center justify-center h-full text-gray-500">No data available</div>
      </div>
    </section>

    <!-- Risk vs Reward -->
    <section class="bg-white dark:bg-gray-800 rounded-lg p-6 shadow-lg space-y-4 flex flex-col items-center justify-center">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-gray-100">Risk vs Reward</h2>
      <div class="w-200 h-100">
        <ScatterChart
          v-if="riskReward.labels.length"
          :potentials="riskReward.potentials"
          :volatilities="riskReward.volatilities"
        />
        <div v-else class="flex items-center justify-center h-full text-gray-500">No data available</div>
      </div>
    </section>

    <!-- Rating Distribution: Bar Chart -->
    <!--section class="bg-white dark:bg-gray-800 rounded-lg p-6 shadow-lg space-y-4 flex flex-col items-center justify-center">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-gray-100">Rating Distribution</h2>
      <div class="w-200 h-100">
        <RatingChart v-if="Object.keys(ratingDistribution).length" :distribution="ratingDistribution" type="bar" />
        <div v-else class="flex items-center justify-center h-full text-gray-500">No data available</div>
      </div>
    </section-->
  </div>
</template>

<script setup lang="ts">
import { reactive, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useStockStore } from '../stores/stocks'
import CandleChart from '../components/CandleChart.vue'
import HistoryChart from '../components/HistoryChart.vue'
import ScatterChart from '../components/ScatterChart.vue'
//import RatingChart from '../components/RatingChart.vue'
import type { HistoricalPoint, HistoryFilters } from '../stores/stocks'
import BackButton from '../components/BackButton.vue'
import TopButton from '../components/TopButton.vue'
import HomeButton from '../components/HomeButton.vue'

const route = useRoute()
const ticker = route.params.ticker as string

const stocksStore = useStockStore()
const detail = computed(() => stocksStore.detail)

// Calculamos hoy y hace un mes en formato YYYY-MM-DD
const today = new Date()
const monthAgo = new Date()
monthAgo.setMonth(monthAgo.getMonth() - 1)

// Función anonima para devolver solo la fecha de un Date
const formatISO = (d: Date) => d.toISOString().slice(0, 10)

// Estado reactivo para los filtros, con valores por defecto
const filters = reactive<HistoryFilters>({
  days: 30,
  start_date: formatISO(monthAgo),  // Fecha de hace un mes
  end_date:   formatISO(today),     // Fecha de hoy
  min_price: null,
  max_price: null,
  min_volume: null,
  order: 'asc'
})

const stock = computed(() => detail.value?.stock)
const history = computed<HistoricalPoint[]>(() => detail.value?.history ?? [])
const riskReward = computed(() => detail.value?.riskReward ?? { labels: [], volatilities: [], potentials: [] })
//const ratingDistribution = computed(() => detail.value?.ratingDistribution ?? {})

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

// Función que recarga el detalle aplicando los filtros
async function applyFilters() {
  try {
    await stocksStore.loadDetail(ticker, { ...filters })
  } catch (err) {
    console.error('Error al aplicar filtros:', err)
  }
}

onMounted(async () => {
  applyFilters()
  /*try {
    await stocksStore.loadDetail(ticker, days)  // Llamado centralizado al store
  } catch (err) {
    console.error('Error loading stock detail:', err)
  }*/
})
</script>

<style scoped>
/* Tailwind styles */
</style>
