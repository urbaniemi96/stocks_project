<template>
  <div>
    <canvas ref="canvas" class="h-full w-full"></canvas>
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
      responsive: true,
      maintainAspectRatio: false,
      scales: {
        x: { 
          display: true, 
          title: { 
            display: false, 
            text: 'Date',
            color: 'rgb(16, 24, 40)',
          },
          ticks:  { color: 'rgb(16, 24, 40)' }
        },
        y: { 
          display: true, 
          title: { 
            display: true, 
            text: 'Price',
            color: 'rgb(16, 24, 40)'
          },
          ticks:  { color: 'rgb(16, 24, 40)' }
        }
      }
    }
  })
})

// Actualiza chart si los datos o los labels cambian
watch(
  () => props.data,
  (newData) => {
    chart.data.datasets[0].data = newData
    chart.update()
  }
)
watch(
  () => props.labels,
  (newLabel) => {
    chart.data.labels = newLabel
    chart.update()
  }
)
</script>