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
  distribution: Record<string, number>
  type: 'bar' | 'pie'
}>()

const canvas = ref<HTMLCanvasElement>()
let chart: Chart

function build() {
  const labels = Object.keys(props.distribution)
  const data = labels.map(l => props.distribution[l])
  return { labels, data }
}

onMounted(() => {
  const { labels, data } = build()
  chart = new Chart(canvas.value!, {
    type: props.type,
    data: { labels, datasets: [{ label: 'Ratings', data }] },
    options: {}
  })
})

watch(() => props.distribution, () => {
  const { labels, data } = build()
  chart.data.labels = labels
  chart.data.datasets[0].data = data
  chart.update()
})
</script>
