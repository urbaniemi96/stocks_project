<template>
  <!-- Floating Button -->
  <button
  v-if="!isOpen"
    @click="toggleDrawer"
    class="fixed top-14 right-13 z-50 p-4 bg-purple-600 text-white rounded-full shadow-lg hover:bg-purple-700 transition"
    
  >
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
        <button @click="toggleDrawer" class="p-1 text-white">
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
import { ref, computed } from 'vue'
import dayjs from 'dayjs'
import isSameOrAfter from 'dayjs/plugin/isSameOrAfter'
import type { Stock, HistoricalPoint, RiskRewardResponse } from '../stores/stocks'

dayjs.extend(isSameOrAfter)


const props = defineProps<{
  stock: Stock
  history: HistoricalPoint[]
  riskReward: RiskRewardResponse
}>()

// Estado del drawer
const isOpen = ref(false)
const aiOpinion = ref('')

function toggleDrawer() {
  isOpen.value = !isOpen.value
  if (isOpen.value) {
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
  const prompt = `Por favor, responde solo en texto plano, sin asteriscos.
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
          { role: 'system', content: 'Eres un asistente financiero experto.' },
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
  aiOpinion.value = text.trim()
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
