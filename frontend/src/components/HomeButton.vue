<template>
  <button
    @click="goHome"
    class="fixed top-34 left-13 !bg-blue-500 !text-white !rounded-full !shadow-lg z-50 !text-xl w-14 h-14 aspect-square
      transform transition duration-200 ease-out
      hover:-translate-y-1 hover:scale-105 hover:shadow-xl
      active:translate-y-1 active:scale-95 active:shadow-md
      focus:!outline-none focus:!ring-2 focus:!ring-blue-400"
    aria-label="Volver atrás"
  >
    <font-awesome-icon :icon="['fas', 'house']" />
  </button>
</template>

<script lang="ts" setup>
import { useRouter } from 'vue-router'

const router = useRouter()
function goHome() {
    router.push('/') 
}
</script>