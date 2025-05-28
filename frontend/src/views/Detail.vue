<template>
  <BackButton />
  <TopButton />
  <HomeButton />
  <PayButton />
  <IAOpinion
    v-if="stock"
    :stock="stock"
    :history="history"
    :riskReward="riskReward"
  />

  <!-- Si no hay stock, muestro mensaje de error -->
  <div
    v-if="!stock && fetched"
    class="container mx-auto p-4 text-center text-xl text-red-600 dark:text-red-400"
  >
    Ticker inexistente
  </div>

  <!-- Si no hay detalle, solo muestro cabecera -->
  <div
    v-else-if="stock && !detailLoaded"
    class="container mx-auto p-4 space-y-5"
  >
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

    <!-- Overview -->
    <section class="bg-white dark:bg-gray-800 rounded-lg p-4 shadow-lg space-y-4">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-gray-100">Resumen</h2>
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 text-gray-300">
        <!-- Targets -->
        <div>
          <p class="text-m uppercase font-semibold text-gray-900 dark:text-gray-100">Rango objetivo</p>
          <p class="text-lg font-medium text-gray-900 dark:text-gray-100">
            {{ stock.target_from.toFixed(2) }} – {{ stock.target_to.toFixed(2) }}
          </p>
        </div>
        <!-- Rating -->
        <div>
          <p class="text-m uppercase font-semibold text-gray-900 dark:text-gray-100">Calificación</p>
          <p class="text-lg font-medium text-gray-900 dark:text-gray-100">
            {{ stock.rating_from }} → {{ stock.rating_to }}
          </p>
        </div>
      </div>
    </section>

    <!-- Mensaje de ausencia de detalle -->
    <div class="pt-15 font-bold text-center text-gray-600 dark:text-gray-400 italic">
      No hay datos de histórico o métricas para mostrar.
    </div>
  </div>

  <!-- Si hay stock y detalle, lo muestro una vez traído -->
  <div v-else-if="stock" class="container mx-auto p-4 space-y-5">
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
    <!-- Filtros -->
    <section class="bg-white dark:bg-gray-800 rounded-lg p-4 shadow-lg space-y-4">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-gray-100">Filtros</h2>
      <div class="flex items-end justify-center gap-4 text-gray-700 dark:text-gray-300">
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
        <button
          @click="applyFilters"
          class="mt-4 bg-blue-600 hover:bg-blue-500 text-white dark:text-gray-300 px-4 py-2 rounded-lg font-semibold"
        >
          <font-awesome-icon :icon="['fas', 'search']" />
        </button>
      </div>
    </section>

    <!-- Overview -->
    <section class="bg-white dark:bg-gray-800 rounded-lg p-4 shadow-lg space-y-4">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-gray-100">Resumen</h2>
      <div class="grid grid-cols-1 sm:grid-cols-3 lg:grid-cols-4 gap-4 text-gray-300">
        <!-- Min y Max juntos -->
        <div>
          <p class="text-m uppercase font-semibold text-gray-900 dark:text-gray-100">
            Cierre Mín / Máx
          </p>
          <p class="text-lg font-medium text-gray-900 dark:text-gray-100">
            {{ stats.min.toFixed(2) }} / {{ stats.max.toFixed(2) }}
          </p>
        </div>
        <!-- Avg Close -->
        <div>
          <p class="text-m uppercase font-semibold text-gray-900 dark:text-gray-100">Promedio Cierre</p>
          <p class="text-lg font-medium text-gray-900 dark:text-gray-100">{{ stats.avg.toFixed(2) }}</p>
        </div>
        
        <!-- Targets -->
        <div>
          <p class="text-m uppercase font-semibold text-gray-900 dark:text-gray-100">Rango objetivo</p>
          <p class="text-lg font-medium text-gray-900 dark:text-gray-100">
            {{ stock.target_from.toFixed(2) }} – {{ stock.target_to.toFixed(2) }}
          </p>
        </div>
        <!-- Rating separado en columnas -->
        <div>
          <p class="text-m uppercase font-semibold text-gray-900 dark:text-gray-100">Calificación</p>
          <p class="text-lg font-medium text-gray-900 dark:text-gray-100">
            {{ stock.rating_from }} → {{ stock.rating_to }}
          </p>
        </div>
        <!-- Trend -->
        <div class="sm:col-span-4">
          <p class="text-m uppercase font-semibold text-gray-900 dark:text-gray-100">Tendencia</p>
          <p :class="['text-lg font-medium', trend >= 0 ? 'text-green-600' : 'text-red-600']">
            <span>
              {{ trend >= 0 ? '+' : '' }}{{ trend.toFixed(2) }}%
            </span>
            <font-awesome-icon
              :icon="['fas', trend >= 0 ? 'arrow-up' : 'arrow-down']"
              :class="trend >= 0 ? 'text-green-600' : 'text-red-600'"
            />
          </p>
        </div>
      </div>
    </section>

    <!-- Candlestick Chart -->
    <Accordion title="Gráfico de velas (OHLC + Volumen)">
      <div class="w-210 h-120">
        <CandleChart v-if="history.length" :history="history" />
        <div v-else class="flex items-center justify-center h-full text-gray-500">
          No hay datos disponibles
        </div>
      </div>
    </Accordion>

    <!-- Close Price Chart -->
    <Accordion title="Precio a lo largo del tiempo">
      <div class="w-210 h-100">
        <HistoryChart v-if="labels.length" :labels="labels" :data="closeData" />
        <div v-else class="flex items-center justify-center h-full text-gray-500">No hay datos disponibles</div>
      </div>
    </Accordion>

    <!-- Volatility Chart -->
    <Accordion title="Volatilidad (%)">
      <div class="w-210 h-100">
        <HistoryChart
          v-if="riskReward.labels.length"
          :labels="riskReward.labels"
          :data="riskReward.volatilities"
        />
        <div v-else class="flex items-center justify-center h-full text-gray-500">No hay datos disponibles</div>
      </div>
    </Accordion>

    <!-- Potential Chart -->
    <Accordion title="Potencial (%)">
      <div class="w-210 h-100">
        <HistoryChart
          v-if="riskReward.labels.length"
          :labels="riskReward.labels"
          :data="riskReward.potentials"
        />
        <div v-else class="flex items-center justify-center h-full text-gray-500">No hay datos disponibles</div>
      </div>
    </Accordion>

    <!-- Risk vs Reward -->
    <Accordion title="Riesgo vs Recompensa">
      <div class="w-210 h-100">
        <ScatterChart
          v-if="riskReward.labels.length"
          :potentials="riskReward.potentials"
          :volatilities="riskReward.volatilities"
        />
        <div v-else class="flex items-center justify-center h-full text-gray-500">No hay datos disponibles</div>
      </div>
    </Accordion>
  </div>
</template>

<script setup lang="ts">
import { reactive, computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { useStockStore } from '../stores/stocks'
import CandleChart from '../components/charts/CandleChart.vue'
import HistoryChart from '../components/charts/HistoryChart.vue'
import ScatterChart from '../components/charts/ScatterChart.vue'
import type { HistoricalPoint, HistoryFilters } from '../stores/stocks'
import BackButton from '../components/navigation/BackButton.vue'
import TopButton from '../components/navigation/TopButton.vue'
import HomeButton from '../components/navigation/HomeButton.vue'
import Accordion from '../components/Accordion.vue'
import IAOpinion from '../components/IAOpinion.vue'
import PayButton from '../components/navigation/PayButton.vue'

const route = useRoute()
const ticker = route.params.ticker as string

const stocksStore = useStockStore()
const detail = computed(() => stocksStore.detail)

// Bandera para evitar parpadeo de pantalla al cargar un nuevo stock 
const fetched = ref(false)

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
// Para detectar si existe histórico de ese ticker
const detailLoaded = computed(() => {
  return Array.isArray(detail?.value?.history) && detail.value.history.length > 0
})
const riskReward = computed(() => detail.value?.riskReward ?? { labels: [], volatilities: [], potentials: [] })

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
  } finally {
    fetched.value = true   // Actualizo flag para mostrar datos
  }
}

onMounted(async () => {
  applyFilters()
})
</script>

<style scoped>
/* Tailwind styles */
</style>
