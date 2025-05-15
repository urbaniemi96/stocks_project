<template>
  <button
    @click="goBack"
    class="fixed top-14 left-13 !bg-red-500 hover:!bg-red-700 !text-white p-3 !rounded-full !shadow-lg z-50 focus:!outline-none focus:!ring-2 focus:!ring-red-400 hover:!border-red-400 !text-xl "
    aria-label="Volver atrás"
  >
    <font-awesome-icon :icon="['fas', 'arrow-left']" />
  </button>
</template>

<script lang="ts" setup>
import { useRouter } from 'vue-router'

const router = useRouter()
function goBack() {
  // Si no hay historial, podrías redirigir a una ruta por defecto:
  if (window.history.length > 1) {
    router.back()
  } else {
    router.push('/')  // o la ruta que prefieras
  }
}
</script>