<template>
  <BackButton />
  <div class="p-4">
    <div class="flex justify-between items-center mb-4">
    <h1 class="text-2xl font-bold mb-4">Top 20 Recommendations</h1>
    <div v-if="isAdmin" class="flex items-center space-x-4">
        <span class="text-sm text-gray-500">Última actualización: {{ lastUpdated }}</span>
        <button
          class="px-3 py-1 bg-indigo-600 text-white rounded hover:bg-indigo-700"
          @click=""
        >
          Recalcular
        </button>
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
          'p-4 mb-2 rounded shadow flex justify-between items-center',
          index < 5 ? 'bg-green-100' : 'bg-white'
        ]"
      >
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


const router = useRouter()
const auth = useAuthStore()

const store = useStockStore()
const { recommendations } = storeToRefs(store)

const isAdmin = computed(() => auth.isAdmin)

// Deriva la última fecha de updated_at del primer elemento
const lastUpdated = computed(() => {
  if (recommendations.value.length === 0) return 'N/A'
  const dates = recommendations.value.map(r => new Date(r.UpdatedAt))
  return dates.reduce((a, b) => (a > b ? a : b)).toLocaleString()
})

onMounted(async () => {
  // Carga las recomendaciones
  await store.loadTopRecommendations()
})
</script>

<style scoped>
/* Puedes ajustar estilos adicionales si quisieras */
</style>