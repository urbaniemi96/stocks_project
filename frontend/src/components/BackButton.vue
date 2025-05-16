<template>
  <button
    @click="goBack"
    class="fixed top-14 left-13 w-14 h-14 !bg-red-500 !text-white aspect-square  !rounded-full !shadow-lg z-50 focus:!outline-none focus:!ring-2 focus:!ring-red-400 !text-xl transform transition duration-200 ease-out hover:-translate-y-1 hover:scale-105 hover:shadow-xl active:translate-y-1 active:scale-95 active:shadow-md"
    aria-label="Volver atrás"
  >
    <font-awesome-icon :icon="['fas', 'arrow-left']" />
  </button>
</template>

<script lang="ts" setup>
import { useRouter } from 'vue-router'

const router = useRouter()
function goBack() {
  // Vuelvo atrás o al home
  if (window.history.length > 1) {
    router.back()
  } else {
    router.push('/')
  }
}
</script>