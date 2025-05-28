<template>
  <div v-if="auth.isUser" class="flex justify-center">
    <button
      @click="goPay"
      class="fixed top-54 left-13 w-14 h-14 !bg-[#28A745] !text-white aspect-square  !rounded-full !shadow-lg z-50 focus:!outline-none focus:!ring-2 focus:!ring-green-400 !text-xl !transform !transition !duration-200 !z ease-out hover:-translate-y-1 hover:scale-105 hover:shadow-xl active:translate-y-1 active:scale-95 active:shadow-md"
      aria-label="Ir a Suscripción"
    >
      <font-awesome-icon :icon="['fas', 'sack-dollar']" />
    </button>
  </div>
</template>

<script lang="ts" setup>
import { useRouter } from 'vue-router'
import { useAuthStore } from '../../stores/auth'

const auth = useAuthStore()

const router = useRouter()
function goPay() {
  // Muestro vista de suscripción
  router.push('/suscription')
  
}
</script>