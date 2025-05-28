import { defineStore } from "pinia";
// Tipos de piezas
export type TetrominoType = "I" | "J" | "L" | "O" | "S" | "T" | "Z";
// Celda puede ser vacía o contener una pieza (uso esta estructura para que sea más facil agregar propiedades futuras)
export type Cell = null | {
    tetromino: TetrominoType;
};
// Shapes estandar de Tetris (matrices 4×4 o menos)
const SHAPES: Record<TetrominoType, number[][]> = {
  I: [
    [0,0,0,0],
    [1,1,1,1],
    [0,0,0,0],
    [0,0,0,0],
  ],
  J: [
    [1,0,0],
    [1,1,1],
    [0,0,0],
  ],
  L: [
    [0,0,1],
    [1,1,1],
    [0,0,0],
  ],
  O: [
    [1,1],
    [1,1],
  ],
  S: [
    [0,1,1],
    [1,1,0],
    [0,0,0],
  ],
  T: [
    [0,1,0],
    [1,1,1],
    [0,0,0],
  ],
  Z: [
    [1,1,0],
    [0,1,1],
    [0,0,0],
  ],
};
// Matriz de celdas
type Grid = Cell[][];


export interface ActivePiece {
  shape: number[][];     // por ejemplo: [[0,1,0],[1,1,1]]
  x: number;             // posición horizontal en el grid
  y: number;             // posición vertical en el grid
  type: TetrominoType;   // para saber qué pieza es
}

// Genero grid de cells de tamaño rows x cols
function createGrid(rows: number, cols: number): Grid {
    // Creo una fila vacía
    const emptyRow: Cell[] = Array(cols).fill(null);
    // Genero un nuevo array de arrays. Esparciendo ...emptyRow para que no todas las filas apunten al mismo array de celdas emptyRow
    return Array.from({ length: rows }, () => [...emptyRow]);
}

export function createTetromino(type: TetrominoType, cols: number): ActivePiece {
  const shape = SHAPES[type];
  // Centrarlo en la parte superior: x = (cols − anchoForma) / 2
  const x = Math.floor((cols - shape[0].length) / 2);
  const y = -shape.length; // arranca “fuera” para subir efecto
  return { shape, x, y, type };
}


export function canPlace(
  shape: number[][],
  x: number,
  y: number,
  grid: Cell[][]
): boolean {
  for (let dy = 0; dy < shape.length; dy++) {
    for (let dx = 0; dx < shape[dy].length; dx++) {
      if (shape[dy][dx] === 0) continue;
      const newY = y + dy;
      const newX = x + dx;
      // Si sale del área o choca con celda ocupada → no cabe
      if (
        newX < 0 ||
        newX >= grid[0].length ||
        newY >= grid.length ||
        (newY >= 0 && grid[newY][newX] !== null)
      ) {
        return false;
      }
    }
  }
  return true;
}

function shuffle(piezas: string[]) {
    return piezas
}

export const useTetrisStore = defineStore("tetris", {
  state: () => ({
    rows: 20,
    cols: 10,
    grid: [] as Grid,
    activePiece: null as ActivePiece | null,
    queue: [] as TetrominoType[],     
    score: 0,
    status: "idle" as "idle"|"playing"|"paused"|"gameover"
  }),
  getters: {
    // Se devuelve una función porque el getter no permite recibir parámetros (pero la función que retorna si)
    isCellOccupied: (state) => (r: number, c: number) =>
      state.grid[r][c] !== null
    },
    // Equivalente legible
    /*
    isCellOccupied: function(state) {
        return function(r: number, c: number) {
            return state.grid[r][c] !== null;
        }
    }
  */
    actions: {
        initGame() {
            this.grid = createGrid(this.rows, this.cols);
            this.queue = this.generateQueue();
            this.spawnNext();
            this.score = 0;
            this.status = "playing";
            // aquí podrías arrancar un timer para caída automática
        },
        generateQueue(): TetrominoType[] {
            const all: TetrominoType[] = ["I","J","L","O","S","T","Z"];
            return all.sort(() => Math.random() - 0.5);
        },
        spawnNext() {
            if (this.queue.length === 0) this.queue = this.generateQueue();
            const next = this.queue.shift()!;
            this.activePiece = createTetromino(next, this.cols);
        },
        drop() {
            if (!this.activePiece) return;
            this.moveActive(0, 1);
        },
        moveActive(dx: number, dy: number) {
            if (!this.activePiece) return;
            const { shape, x, y } = this.activePiece;
            // Si cabe en la nueva posición, la movemos
            if (canPlace(shape, x + dx, y + dy, this.grid)) {
                this.activePiece.x += dx;
                this.activePiece.y += dy;
            } else if (dy === 1) {
                // Si no cabe **al bajar**, es momento de lockear
                this.lockPiece();
            }
        },
        lockPiece() {
    if (!this.activePiece) return;

    const { shape, x, y, type } = this.activePiece;
    const rows = this.grid.length;
    const cols = this.grid[0].length;

    // Fijo la pieza al grid
    for (let dy = 0; dy < shape.length; dy++) {
      for (let dx = 0; dx < shape[dy].length; dx++) {
        if (shape[dy][dx] === 1) {
          const newY = y + dy;
          const newX = x + dx;
          if (newY >= 0 && newY < rows && newX >= 0 && newX < cols) {
            this.grid[newY][newX] = { tetromino: type };
          }
        }
      }
    }

    // Buscar filas completas
    const fullRows: number[] = [];
    for (let r = 0; r < rows; r++) {
      if (this.grid[r].every(cell => cell !== null)) {
        fullRows.push(r);
      }
    }

    // Eliminar filas completas y contar score
    const lines = fullRows.length;
    if (lines > 0) {
      // Remover de abajo hacia arriba para no romper índices
      for (let idx = fullRows.length - 1; idx >= 0; idx--) {
        const row = fullRows[idx];
        this.grid.splice(row, 1);
        // Insertar fila vacía arriba
        this.grid.unshift(Array(cols).fill(null));
      }
    }

    // Generar siguiente pieza
    this.spawnNext();

    
  },
        moveLeft() { 
            this.moveActive(-1, 0); 
        },
        moveRight() { 
            this.moveActive(+1, 0); 
        },
        /*spawnNext() {
            if (this.queue.length < 3) {
                this.queue.push(...this.generateQueue());
            }
            const nextType = this.queue.shift()!;
            this.activePiece = createTetromino(nextType);
        },
        drop() { 
            this.moveActive(0, +1); 
        },
        rotate() { 
            if (!this.activePiece) return;
            this.activePiece.shape = rotatePiece(this.activePiece.shape);
        },
        */
        pause() { this.status = "paused"; },
        resume() { this.status = "playing"; },
        endGame() { this.status = "gameover"; }
    }
});