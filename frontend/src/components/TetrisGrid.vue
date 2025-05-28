<template>
  <div class="grid grid-container gap-0.5 "
       :style="{
         gridTemplateRows: `repeat(${rows}, 1fr)`,
         gridTemplateColumns: `repeat(${cols}, 1fr)`,
       }"
       >
    <div
      v-for="(row, rowIndex) in displayGrid "
      :key="rowIndex"
      class="contents"
    >
      <div
        v-for="(cell, colIndex) in row "
        :key="colIndex"
        class="w-6 h-6 border border-gray-800"
        :class="getCellColorClass(cell)"
      ></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onBeforeUnmount } from "vue";
import { storeToRefs } from "pinia";
import { useTetrisStore } from "../stores/tetris";
import type { Cell } from "../stores/tetris";

const tetris = useTetrisStore();
let dropTimer: number;

const { grid, rows, cols, activePiece } = storeToRefs(tetris);

// Computed que mezcla grid + pieza activa
const displayGrid = computed<Cell[][]>(() => {
  // Copia superficial de filas y columnas
  const view = grid.value.map(row => [...row]);

  const piece = activePiece.value;
  if (!piece) return view;

  for (let dy = 0; dy < piece.shape.length; dy++) {
    for (let dx = 0; dx < piece.shape[dy].length; dx++) {
      if (piece.shape[dy][dx] === 1) {
        const y = piece.y + dy;
        const x = piece.x + dx;
        // solo si está dentro del grid
        if (y >= 0 && y < rows.value && x >= 0 && x < cols.value) {
          view[y][x] = { tetromino: piece.type };
        }
      }
    }
  }
  return view;
});

// Devuelve clases de Tailwind o estilos inline según la pieza
function getCellColorClass(cell: Cell) {
  if (!cell) return "bg-gray-900"; // fondo vacío
  switch (cell.tetromino) {
    case "I": return "bg-cyan-400";
    case "J": return "bg-blue-500";
    case "L": return "bg-orange-500";
    case "O": return "bg-yellow-400";
    case "S": return "bg-green-500";
    case "T": return "bg-purple-500";
    case "Z": return "bg-red-500";
    default: return "bg-white";
  }
}
onMounted(() => {
  tetris.initGame();

  // Cada 500 ms baja la pieza
  dropTimer = window.setInterval(() => {
    if (tetris.status === "playing") {
      tetris.drop();
    }
  }, 500);

  // Listener de flechas: L y R y UP
  window.addEventListener("keydown", onKeyDown);
});

onBeforeUnmount(() => {
  clearInterval(dropTimer);
  window.removeEventListener("keydown", onKeyDown);
});

let isDropping = false;
function onKeyDown(e: KeyboardEvent) {
  if (e.key === "ArrowLeft") tetris.moveLeft();
  if (e.key === "ArrowRight") tetris.moveRight();
  if (e.key === "ArrowDown" && !isDropping) {
    //isDropping = true;
    tetris.drop();
  }
}
</script>

<style scoped>
.grid-container {
  display: grid;
  aspect-ratio: 10 / 20;
  width: 100%;
  max-width: 240px;
}
</style>
