<template>
  <canvas ref="candleCanvas" class="h-full w-full"></canvas>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Chart, registerables } from 'chart.js'
import { CandlestickController, CandlestickElement } from 'chartjs-chart-financial'
import 'chartjs-adapter-luxon'

// Registra el plugin financiero
Chart.register(...registerables, CandlestickController, CandlestickElement)

interface MyFinancialDataPoint {
  x: number | Date
  o: number
  h: number
  l: number
  c: number
}

interface HistoryItem {
  Date: string
  Open: number
  High: number
  Low: number
  Close: number
}

const props = defineProps<{ history: HistoryItem[] }>()
const candleCanvas = ref<HTMLCanvasElement>()
let candleChart: Chart<'candlestick', MyFinancialDataPoint[], unknown>

onMounted(() => {
  const ohlc: MyFinancialDataPoint[] = props.history.map(h => ({
    x: new Date(h.Date).getTime(), // o new Date(h.Date).getTime()
    o: h.Open,
    h: h.High,
    l: h.Low,
    c: h.Close
  }))

  candleChart = new Chart(candleCanvas.value!, {
    type: 'candlestick',
    data: {
      datasets: [{ label: 'OHLC', data: ohlc }]
    },
    options: {
      parsing: false,
      scales: {
        x: { 
          type: 'time',
          time: {
            unit: 'day',         // unidad diaria
            tooltipFormat: 'DD LLL yyyy HH:mm' // formato en tooltip
          },
          ticks:  { color: 'rgb(16, 24, 40)' }

        },
        y: { title: { display: true, text: 'Price', color: 'rgb(16, 24, 40)' }, 
          ticks:  { color: 'rgb(16, 24, 40)' }
       }
      }
    }
  })
})
</script>
