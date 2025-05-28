<template>
  <BackButton />
  <HomeButton />
  <TopButton />
  <!-- Wrapper absoluto full-screen -->
  <div class="absolute inset-0 min-h-screen flex items-center justify-center p-6">
    <div class="max-w-2xl bg-white shadow-lg rounded-2xl p-8">
      <h1 class="text-4xl font-bold text-gray-800 mb-6 text-center">Suscripción Premium</h1>
      <!-- Listado con texto alineado a la izquierda -->
      <ul class="space-y-4 mb-8 text-left">
        <li v-for="(benefit, index) in benefits" :key="index" class="flex items-start">
          <svg class="flex-shrink-0 w-6 h-6 text-green-500 mr-3 mt-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
          </svg>
          <p v-html="benefit" class="text-gray-700 leading-relaxed"></p>
        </li>
      </ul>
      <button
        @click="goToThankYou"
        class="w-full py-4 bg-gradient-to-r from-indigo-500 to-teal-400 hover:from-teal-400 hover:to-indigo-500 text-white !font-bold rounded-xl shadow-lg !transform !transition hover:scale-105 !duration-300 active:!scale-95 focus:!outline-none"
      >
        Pagar
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import BackButton from '../components/navigation/BackButton.vue'
import HomeButton from '../components/navigation/HomeButton.vue'
import TopButton from '../components/navigation/TopButton.vue'

const router = useRouter()

const benefits: string[] = [
  '<b>Ranking de Recomendaciones actualizado</b> minuto a minuto para saber siempre las mejores acciones para invertir',
  '<b>Históricos actualizados</b> al instante al ver los detalles de una acción',
  '<b>IA avanzada</b> para el análisis profundo de una acción',
  '<b>Email diario</b> con las mejores acciones para invertir ese día, noticias de las empresas, blogs, opiniones de expertos, y muchas cosas más!',
  '<b>Nuestro eterno agradecimiento</b> por apoyar esta idea y permitir que este proyecto sea posible :) <3'
]

function goToThankYou() {
  router.push('/thanks')
}
</script>
