<template>
  <BackButton />
  <TopButton />
  <HomeButton />
  <div class="container mx-auto p-4 space-y-6">
    <div class="flex justify-between items-center mb-4">
    <h1 class="text-3xl font-extrabold text-gray-900 dark:text-gray-100 mb-4">Top 20 Recomendaciones</h1>
    <div v-if="isAdmin" class="flex items-center space-x-4">
        <span class="text-sm text-gray-600 dark:text-gray-400">Última actualización: {{ lastUpdated }}</span>
      </div>
    </div>
    <p class="mb-4 text-gray-700">
      El <strong>Score</strong> mide el retorno ajustado por volatilidad y ponderado según el rating de la acción.
      Un score más alto indica una mejor relación riesgo-recompensa, lo cual sugiere invertir en estas acciones.<br>
      Las <span class="font-semibold text-green-800">5 mejores</span> aparecen resaltadas en verde.
    </p>
    <ul>
      <li
        v-for="(rec, index) in recommendations"
        :key="rec.Ticker"
        :class="[
          'p-4 mb-2 rounded-lg shadow flex justify-between items-center',
          index < 5
            ? 'bg-green-100 dark:bg-green-900/20'
            : 'bg-white dark:bg-gray-700'
         ]"
      >
        <font-awesome-icon
          v-if="index < 3"
          :icon="['fas', 'trophy']"
          class="text-2xl"
          :class="{
            'text-yellow-500': index === 0,
            'text-gray-500': index === 1,
            'text-yellow-700': index === 2
          }"
        />
        <font-awesome-icon
          v-else 
          :icon="['fas', 'award']"
          class="text-2xl"
          :class="{
            'text-zinc-800': index == 3,
            'text-zinc-500': index == 4,
            'text-zinc-400': index >= 5,
            }"
        />
        <div>
          <p class="font-semibold text-lg text-gray-900">
            {{ index + 1 }}. {{ rec.Ticker }}
          </p>

          <p class="text-sm text-gray-600">Score: {{ rec.Score.toFixed(4) }}</p>
        </div>
        <button class="px-2 py-1 bg-indigo-600 text-white rounded hover:bg-indigo-700"
        @click="router.push({ name: 'stock-detail', params: { ticker: rec.Ticker } })">
          Ver detalle
        </button>
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { useStockStore } from '../stores/stocks'
import { useAuthStore } from '../stores/auth'
import { storeToRefs } from 'pinia'
import { useRouter } from 'vue-router'
import BackButton from '../components/BackButton.vue'
import TopButton from '../components/TopButton.vue'
import HomeButton from '../components/HomeButton.vue'


const router = useRouter()
const auth = useAuthStore()

const store = useStockStore()
const { recommendations } = storeToRefs(store)

const isAdmin = computed(() => auth.isAdmin)

// Deriva la última fecha de updated_at del primer elemento
const lastUpdated = computed(() => {
  if (recommendations.value.length === 0) return 'N/A'
  const dates = recommendations.value.map(r => new Date(r.UpdatedAt))
  return dates.reduce((a, b) => (a > b ? a : b)).toLocaleString('es-AR', { timeZone: 'America/Argentina/Buenos_Aires', hour12: false, })
})

onMounted(async () => {
  // Carga las recomendaciones
  await store.loadTopRecommendations()
})
</script>

<style scoped>
/* Puedes ajustar estilos adicionales si quisieras */
</style>