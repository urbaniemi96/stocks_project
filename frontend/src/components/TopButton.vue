<template>
  <button
    v-if="showScroll"
    @click="scrollToTop"
    class="fixed bottom-14 left-13 !bg-blue-500 !text-white !rounded-full !shadow-lg z-50 !text-xl w-14 h-14 aspect-square
      transform transition duration-200 ease-out
      hover:-translate-y-1 hover:scale-105 hover:shadow-xl
      active:translate-y-1 active:scale-95 active:shadow-md
      focus:!outline-none focus:!ring-2 focus:!ring-blue-400"
    aria-label="Ir al inicio"
  >
    <font-awesome-icon :icon="['fas', 'arrow-up']" />
  </button>
</template>

<script lang="ts" setup>
import { ref, onMounted, onUnmounted } from 'vue'

const showScroll = ref(false)
// Detecto si hay scroll
const onScroll = () => {
  showScroll.value = window.scrollY > 100
}
// Vuelvo scroll arriba
function scrollToTop() {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(() => {
  window.addEventListener('scroll', onScroll, { passive: true })
})
onUnmounted(() => {
  window.removeEventListener('scroll', onScroll)
})
</script>