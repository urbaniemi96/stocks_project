<template>
  <div>
    <canvas ref="canvas"></canvas>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { Chart, registerables } from 'chart.js'
Chart.register(...registerables)

const props = defineProps<{
  labels: string[]
  data: number[]
}>()

const canvas = ref<HTMLCanvasElement>()
let chart: Chart

onMounted(() => {
  chart = new Chart(canvas.value!, {
    type: 'line',
    data: {
      labels: props.labels,
      datasets: [{
        label: 'Close Price',
        data: props.data,
        fill: false,
        tension: 0.2
      }]
    },
    options: {
      scales: {
        x: { display: true, title: { display: true, text: 'Date' } },
        y: { display: true, title: { display: true, text: 'Price' } }
      }
    }
  })
})

// Actualiza chart si props cambian
watch(
  () => props.data,
  (newData) => {
    chart.data.datasets[0].data = newData
    chart.update()
  }
)
</script>