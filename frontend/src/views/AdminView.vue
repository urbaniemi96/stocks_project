<template>
  <div class="p-8">
    <h1 class="text-3xl font-extrabold text-gray-900 dark:text-gray-100 mb-10">Configuración</h1>
    
    <BackButton />
    <TopButton />
    <HomeButton />
    <div class="flex flex-col items-center space-y-7">
        <button class="tool-button w-136 h-22 flex items-center justify-center space-x-2" @click="refreshAll">
            <font-awesome-icon :icon="['fas', 'cloud-download-alt']" class="text-4xl"/>
            &nbsp;&nbsp;Traer de la API y cargar en la BD
        </button>
        <button class="tool-button w-136 h-22 flex items-center justify-center space-x-2" @click="enrichAll">
            <font-awesome-icon :icon="['fas', 'chart-line']" class="text-4xl"/>
            &nbsp;&nbsp;Enriquecer datos (puede demorar horas)
        </button>
        <button class="tool-button w-136 h-22 flex items-center justify-center space-x-2" @click="recalcRecommendations">
            <font-awesome-icon :icon="['fas', 'sync']" class="text-4xl"/>
            &nbsp;&nbsp;Recalcular Reecomendaciones
        </button>
        <div v-if="taskId">
            Progreso: {{ status.pages_fetched }} páginas<br/>
            Estado: {{ status.status }}
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
