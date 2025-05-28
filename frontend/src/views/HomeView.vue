<template>
  <div class="p-8">
    <header class="mb-12 text-center">
      <!-- Título con degradado en el texto -->
      <h1 class="text-5xl font-extrabold bg-clip-text text-transparent
                 bg-gradient-to-r from-indigo-600 via-purple-600 to-pink-500">
        ¡Bienvenido a Stockeando!
      </h1>
      <!-- Subtítulo descriptivo -->
      <p class="mt-3 text-lg text-gray-600 dark:text-gray-300">
        Gestiona tus inversiones con un solo clic
      </p>
    </header>

    <div class="flex flex-col items-center space-y-6">
      <div class="flex space-x-4 justify-center">
        <button class="!bg-[#DB162F] tool-button w-66 h-54 flex flex-col items-center justify-center" @click="goDash">
          <font-awesome-icon :icon="['fas', 'money-bills']" class="text-4xl mb-2"/>
          Dashboard
        </button>
        <button class="!bg-[#59544B] tool-button w-66 h-54 flex flex-col items-center justify-center" @click="goRecommend">
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
      <div v-if="auth.isUser" class="flex justify-center">
        <button class="!bg-[#28A745] tool-button w-136 h-22 flex items-center justify-center space-x-2" @click="goPay">
          <font-awesome-icon :icon="['fas', 'sack-dollar']" />
          &nbsp;&nbsp;Activar Suscripción
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
function goPay() {
  router.push('/suscription')
  
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
