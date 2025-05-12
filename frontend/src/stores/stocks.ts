import { defineStore } from 'pinia' // Para definir un estado compartido
import { api } from '../services/api' // Mi api axios para peticiones

// Defino interfaz stock de typescript
export interface Stock {
  ticker: string
  company: string
  target_from: number
  target_to: number
  action: string
  brokerage: string
  rating_from: string
  rating_to: string
  time: string
}

// Defino store de pinia
export const useStockStore = defineStore('stocks', {
  // Defino estado del store
  state: () => ({
    list: [] as Stock[],
    recommended: null as Stock | null,
    taskId: null as string | null,
    status: { status: '', pages_fetched: 0 },
  }),
  // Acciones para modificar estado
  actions: {
    // Disparo petición al fetch del back (para traer datos de la api y guardar en la bd)
    async fetchAndStore() {
      // dispara /fetch en el backend
      const res = await api.get('/fetch')
      this.taskId = res.data.task_id

      // Empiezo a pollear cada 2seg hasta que status == "done" o "error"
      const interval = setInterval(async () => {
        const st = await api.get(`/task/${this.taskId}`)
        this.status = st.data
        if (this.status.status !== 'in-progress') {
          clearInterval(interval)
          alert(`Descarga ${this.status.status}`)
        }
      }, 2000)
    },
    // Disparo petición para enriquecer la API con datos de Yahoo
    async fetchAndEnrich() {
      const res = await api.get('/enrich')
      this.taskId = res.data.task_id

      // Empiezo a pollear cada 3seg hasta que status == "done" o "error"
      /*const interval = setInterval(async () => {
        const st = await api.get(`/task/${this.taskId}`)
        this.status = st.data
        if (this.status.status !== 'in-progress') {
          clearInterval(interval)
          alert(`Enrich ${this.status.status}`)
        }
      }, 3000)*/
    },
    // Traigo los stocks guardados en la bd y los guardo en list (estado del store)
    async loadStocks() {
      // dispara /stocks en el backend
      const res = await api.get<Stock[]>('/stocks')
      this.list = res.data 
    },
    // Llamo al back para recomendar una acción
    async computeRecommendation() {
      //const res = await api.get<Stock>('/recommend')
      //this.recommended = res.data
    },
  },
})