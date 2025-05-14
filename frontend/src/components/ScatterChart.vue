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
  potentials: number[]
  volatilities: number[]
}>()

const canvas = ref<HTMLCanvasElement>()
let chart: Chart

onMounted(() => {
  const pts = props.potentials.map((p, i) => ({ x: p, y: props.volatilities[i] }))
  chart = new Chart(canvas.value!, {
    type: 'scatter',
    data: { datasets: [{ label: 'Risk vs Reward', data: pts }] },
    options: {
      scales: {
        x: { title: { display: true, text: 'Potential (%)' } },
        y: { title: { display: true, text: 'Volatility (%)' } }
      }
    }
  })
})

watch(() => [props.potentials, props.volatilities], ([newP, newV]) => {
  chart.data.datasets[0].data = newP.map((p, i) => ({ x: p, y: newV[i] }))
  chart.update()
})
</script>
