<template>
  <div class="border border-gray-200 dark:border-gray-700 rounded-lg mb-4">
    <button
      class="w-full flex justify-between items-center p-4 bg-indigo-600 text-white rounded shadow-lg focus:outline-none"
      @click="open = !open"
    >
      <p></p>
      <span class="text-lg font-medium">{{ title }}</span>
      <svg
        :class="{'transform rotate-180': open}"
        class="w-5 h-5 text-white transition-transform duration-200"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
        xmlns="http://www.w3.org/2000/svg"
      >
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
      </svg>
    </button>
    <div
      class="overflow-hidden transition-[max-height] duration-300 bg-white dark:bg-gray-800 rounded-lg shadow-lg"
      :class="open ? 'max-h-120' : 'max-h-0'"
    >
      <div class="p-6 flex flex-col items-center justify-center">
        <slot />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, defineProps } from 'vue';

const props = defineProps<{ title: string }>();
const open = ref(false);
</script>

<style scoped>
</style>