

<template>
  <div class="p-8">
    <h1 class="text-3xl font-extrabold text-gray-900 dark:text-gray-100 mb-10">Menú</h1>

    <div class="flex flex-col items-center space-y-6">
      <div class="flex space-x-4 justify-center">
        <button style="background-color: #DB162F;" class="tool-button w-66 h-54 flex flex-col items-center justify-center" @click="goDash">
          <font-awesome-icon :icon="['fas', 'money-bills']" class="text-4xl mb-2"/>
          Dashboard
        </button>
        <button style="background-color: #59544B;" class="tool-button w-66 h-54 flex flex-col items-center justify-center" @click="goRecommend">
          <font-awesome-icon :icon="['fas', 'money-bill-trend-up']" class="text-4xl mb-2"/>
          Recomendaciones
        </button>
      </div>

      <div v-if="auth.isAdmin" class="flex justify-center">
        <button class="tool-button w-136 h-22 flex items-center justify-center space-x-2" @click="goConfig">
          <font-awesome-icon :icon="['fas', 'gear']" class="text-4xl"/>
          &nbsp;&nbsp;Panel Administración
        </button>
      </div>
    </div>
  </div>
</template>


<script lang="ts" setup>
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()

const router = useRouter()
function goDash() {
  router.push('/stocks')
}

function goRecommend() {
  router.push('/recommendations')
}

function goConfig() {
  router.push('/configuration')
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
