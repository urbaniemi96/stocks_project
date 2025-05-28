<template>
  <!-- Floating Button -->
  <button
  v-if="!isOpen"
  @click="toggleDrawer"
  class="group fixed top-14 right-13 w-14 h-14 z-50 !bg-purple-600 text-white aspect-square !rounded-full !shadow-lg hover:bg-purple-700 focus:!outline-none focus:!ring-2 focus:!ring-purple-500 !text-xl !transform !transition !duration-200 ease-out hover:-translate-y-1 hover:scale-105 hover:shadow-xl active:translate-y-1 active:scale-95 active:shadow-md"
>
  <!-- Tooltip -->
  <span
    class="absolute right-full top-1/2 transform -translate-y-1/2 mr-2 whitespace-nowrap
           bg-gray-800 text-white text-sm px-2 py-1 rounded opacity-0
           pointer-events-none transition-opacity duration-200 group-hover:opacity-100"
  >
    ¿Preguntarle a una IA?
  </span>

  <!-- Icono -->
  <font-awesome-icon :icon="['fas', isOpen ? 'times' : 'robot']" />
</button>

  <!-- Slide-out Drawer -->
  <transition name="slide">
    <div
      v-if="isOpen"
      class="fixed top-0 right-0 h-full w-80 bg-white dark:bg-gray-800 shadow-xl z-40 flex flex-col"
    >
        <!-- Título -->
      <div class="flex items-center justify-between p-4 border-b dark:border-gray-700">
        <h3 class="text-lg font-semibold">Análisis de IA</h3>
        <button @click="toggleDrawer" class="p-1 text-white !bg-purple-600">
          <font-awesome-icon :icon="['fas', 'times']" />
        </button>
      </div>
      <p class="font-medium mb-2">{{ stock.company }} ({{ stock.ticker }})</p>
      <div class="flex-1 p-4  h-full">
        <textarea
          v-model="aiOpinion"
          rows="15"
          class="w-full h-full p-2 border rounded resize-none dark:bg-gray-700 dark:text-gray-100"
        />
      </div>
    </div>
  </transition>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import dayjs from 'dayjs'
import isSameOrAfter from 'dayjs/plugin/isSameOrAfter'
import type { Stock, HistoricalPoint, RiskRewardResponse } from '../stores/stocks'
import { useToast } from "vue-toastification";

dayjs.extend(isSameOrAfter)


const props = defineProps<{
  stock: Stock
  history: HistoricalPoint[]
  riskReward: RiskRewardResponse
}>()

// Toast
const toast = useToast();

// Estado del drawer
const isOpen = ref(false)
const aiOpinion = ref('')

// Bandera para detectar si ya se trajo una opinión
const opined = ref(false)

function toggleDrawer() {
  isOpen.value = !isOpen.value
  if (isOpen.value && !opined.value) {
    // Cambio mensaje del textarea
    aiOpinion.value = "Consultando con una IA experta, este proceso puede demorar minutos..."
    // Bandera para traer solo una vez la opinión
    opined.value = true
    askAI()
  }
}

async function askAI() {
  // Filtrar history últimos 7 días respecto al último punto
  const dates = props.history.map(h => dayjs(h.Date))
  const maxDate = dates.reduce((a, b) => (b.isAfter(a) ? b : a))
  const threshold = maxDate.subtract(3, 'day')
  const recentHistory = props.history.filter(h => dayjs(h.Date).isSameOrAfter(threshold))

  // Payload completo
  const payload = {
    stock: props.stock,
    history: recentHistory,
    riskReward: props.riskReward
  }

  // Prompt en texto plano
  const prompt = `Responde solo en texto plano, sin asteriscos.
Analiza esta acción y dime si vale la pena invertir. Datos:
${JSON.stringify(payload, null, 2)}`

  // Llamada al router de HF
  const res = await fetch(
    'https://router.huggingface.co/novita/v3/openai/chat/completions',
    {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${import.meta.env.VITE_HF_TOKEN}`
      },
      body: JSON.stringify({
        model: 'deepseek/deepseek-r1-turbo',
        messages: [
          { role: 'system', content: 'Eres un asistente financiero experto que recomienda si invertir en una acción.' },
          { role: 'user', content: prompt }
        ],
        temperature: 0.5,
        top_p: 0.7,
        stream: false
      })
    }
  )
  const result = await res.json()
  let text = result.choices?.[0]?.message?.content || ''
  // Remover bloques <think>
  text = text.replace(/<think>[\s\S]*?<\/think>/gi, '')
  // Elimino dobles asteriscos de negrita: **texto** → texto
  text = text.replace(/\*\*(.*?)\*\*/g, '$1')
  // Elimino asteriscos simples: *texto* → texto
  text = text.replace(/\*(.*?)\*/g, '$1')

  aiOpinion.value = text.trim()

  // Muestro mensaje de "listo" si no está abierto el slide
  if (!isOpen.value) {
    toast.success("La opinión de la IA está lista!", {
      timeout: 2000
    });
  }
}
</script>

<style scoped>
/* Transición slide desde la derecha */
.slide-enter-active, .slide-leave-active {
  transition: transform 0.3s ease;
}
.slide-enter-from {
  transform: translateX(100%);
}
.slide-enter-to {
  transform: translateX(0);
}
.slide-leave-from {
  transform: translateX(0);
}
.slide-leave-to {
  transform: translateX(100%);
}
</style>
