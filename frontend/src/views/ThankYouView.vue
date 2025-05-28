<template>
  <BackButton />
  <HomeButton />
  <TopButton />
  <!-- Wrapper absoluto full-screen -->
  <div class="absolute inset-0 min-h-screen bg-gradient-to-b from-blue-50 to-white flex items-center justify-center p-6">
    <div class="max-w-2xl bg-white shadow-lg rounded-2xl p-8 text-center">
      <h1 class="text-4xl font-bold text-gray-800 mb-6 underline decoration-3 decoration-[#28A745]">Gracias!</h1>
      <p>
        Quiero darte las <strong>gracias de corazón</strong> por considerar suscribirte y valorar el esfuerzo detrás de esta app. <strong>¡Eres increíble!</strong><br/><br/>
        De momento esto es solo una demo, así que todavía <strong>no es posible completar el pago</strong>, pero tu interés y tu voluntad significan muchísimo para mí.<br/><br/>
        ¡<strong>Mil gracias</strong> por probar <strong>Stockeando</strong>!
        
      </p>
      <!-- Badge de creador -->
      <div class="mt-8">
        <a 
        href="https://www.linkedin.com/in/emiliano-urbani/"
        target="_blank"
        rel="noopener noreferrer"
        class="inline-block bg-gradient-to-r from-indigo-500 to-teal-400 hover:from-teal-400 hover:to-indigo-500 !text-white !font-bold px-3 py-1 rounded-full !transform !transition hover:scale-105 !duration-300 active:!scale-95">
          Creador: Emiliano Urbani
        </a>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import BackButton from '../components/navigation/BackButton.vue'
import HomeButton from '../components/navigation/HomeButton.vue'
import TopButton from '../components/navigation/TopButton.vue'
</script>
