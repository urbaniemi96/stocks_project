<template>
  <BackButton />
  <TopButton />
  <HomeButton />
  <div class="p-8">
    <h1 class="text-3xl font-extrabold text-gray-900 dark:text-gray-100 mb-10">Configuración</h1>
    
    <div class="flex flex-col items-center space-y-7">
        <button class="tool-button w-136 h-22 flex items-center justify-center space-x-2" @click="refreshAll">
            <font-awesome-icon :icon="['fas', 'cloud-download-alt']" class="text-4xl"/>
            &nbsp;&nbsp;Traer de la API y cargar en la BD
        </button>
        <button class="tool-button w-136 h-22 flex items-center justify-center space-x-2" @click="enrichAll">
            <font-awesome-icon :icon="['fas', 'chart-line']" class="text-4xl"/>
            &nbsp;&nbsp;Enriquecer datos (últimos 3 meses)
        </button>
        <button class="tool-button w-136 h-22 flex items-center justify-center space-x-2" @click="recalcRecommendations">
            <font-awesome-icon :icon="['fas', 'sync']" class="text-4xl"/>
            &nbsp;&nbsp;Recalcular Reecomendaciones
        </button>
        <!-- Información de la tarea -->
        <div v-if="taskId" class="flex flex-col max-w-md mx-auto p-6 bg-white rounded-2xl shadow-lg">
          <h2 class="text-xl font-semibold mb-3 whitespace-nowrap overflow-hidden truncate">Tarea {{ taskId }}</h2>

          <!-- Indicador indeterminado -->
          <div class="relative w-full h-3 mb-4 bg-gray-200 rounded-full overflow-hidden">
            <div
              v-if="status.status === 'in-progress'"
              class="absolute inset-0 bg-indigo-600 opacity-75 animate-pulse"
            ></div>
          </div>
          <p class="text-sm text-gray-600 mb-2">
            Páginas obtenidas: <span class="font-medium">{{ status.pages_fetched }}</span>
          </p>

          <!-- Estado con badges -->
          <div class="flex self-center items-center">
            <span
              class="inline-block px-3 py-1 rounded-full text-xs font-medium"
              :class="{
                'bg-yellow-100 text-yellow-800': status.status === 'in-progress',
                'bg-green-100 text-green-800': status.status === 'done',
                'bg-red-100 text-red-800': status.status === 'error'
              }"
            >
              {{ status.status.replace('-', ' ').toUpperCase() }}
            </span>
            <p
              v-if="status.status === 'error' && status.error"
              class="ml-3 text-sm text-red-600"
            >
              {{ status.error }}
            </p>
            <svg
              v-if="status.status === 'in-progress'"
              class="ml-2 w-4 h-4 animate-spin text-blue-500"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
            >
              <circle
                class="opacity-25"
                cx="12"
                cy="12"
                r="10"
                stroke="currentColor"
                stroke-width="4"
              />
              <path
                class="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8v8H4z"
              />
            </svg>
          </div>
        </div>
    </div>

  </div>
</template>


<script lang="ts" setup>
import { useStockStore } from '../stores/stocks'
import { storeToRefs } from 'pinia'
import BackButton from '../components/BackButton.vue'
import TopButton from '../components/TopButton.vue'
import HomeButton from '../components/HomeButton.vue'

const store = useStockStore()
const { taskId, status } = storeToRefs(store)

async function refreshAll() {
  await store.fetchAndStore()
}

async function enrichAll() {
  await store.fetchAndEnrich()
}

async function recalcRecommendations() {
  await store.triggerRecalculate()
}

</script>

<style>
  .tool-button {
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 15px 25px;
    font-size: 20px;
    font-weight: bold;
    background-color: #4f46e5; /* indigo-600 */
    color: white;
    border: none;
    border-radius: 8px;
    box-shadow: 0 10px 25px rgba(0, 0, 0, 0.3);
    cursor: pointer;
    transition: transform 0.3s, box-shadow 0.3s, filter 0.3s;
  }

  .tool-button:hover {
    filter: brightness(1.12);
    transform: translateY(-5px);
    box-shadow: 0 20px 30px rgba(0, 0, 0, 0.4);
  }

  .tool-button:active {
    filter: brightness(0.95);
    transform: translateY(-2px);
    box-shadow: 0 8px 15px rgba(0, 0, 0, 0.3);
  }

  .tool-button:focus {
    outline: none;
  }
</style>
