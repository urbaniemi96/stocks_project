import { defineStore } from 'pinia' // Para definir un estado compartido
import { api } from '../services/api' // Mi api axios para peticiones

export interface HistoryFilters {
  days: number
  start_date?: string
  end_date?: string
  min_price?: number | null
  max_price?: number | null
  min_volume?: number | null
  order: 'asc' | 'desc'
}

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

// Interfaces del Detalle
export interface HistoricalPoint { 
  Date: string 
  Close: number 
  High: number 
  Low: number 
  Open: number 
  Volume: number 
}
interface RiskRewardResponse { 
  labels: string[]
  volatilities: number[]
  potentials: number[] 
}
export interface StockDetailResponse {
  stock: Stock
  history: HistoricalPoint[]
  riskReward: RiskRewardResponse
  ratingDistribution: Record<string, number>
}

// Recomendaciones
export interface Recommendation {
  Ticker: string
  Score: number
  UpdatedAt: string
}

// Defino store de pinia
export const useStockStore = defineStore('stocks', {
  // Defino estado del store
  state: () => ({
    list: [] as Stock[],
    taskId: null as string | null,
    status: { status: '', pages_fetched: 0 },
    detail: null as StockDetailResponse | null,
    recommendations: [] as Recommendation[],
  }),
  // Acciones para modificar estado
  actions: {
    // Disparo petición al fetch del back (para traer datos de la api y guardar en la bd)
    async fetchAndStore() {
      // dispara /fetch en el backend
      const res = await api.get('/admin/fetch')
      this.taskId = res.data.task_id

      // Empiezo a pollear cada 2seg hasta que status == "done" o "error"
      const interval = setInterval(async () => {
        const st = await api.get(`/admin/task/${this.taskId}`)
        this.status = st.data
        if (this.status.status !== 'in-progress') {
          clearInterval(interval)
          alert(`Descarga ${this.status.status}`)
        }
      }, 2000)
    },
    // Disparo petición para enriquecer la API con datos de Yahoo
    async fetchAndEnrich() {
      const res = await api.get('/admin/enrich')
      this.taskId = res.data.task_id

      // Empiezo a pollear cada 3seg hasta que status == "done" o "error"
      /*const interval = setInterval(async () => {
        const st = await api.get(`/admin/task/${this.taskId}`)
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

    // Acción para cargar el detalle
    async loadDetail(ticker: string, filters: HistoryFilters) {
      try {
        // Construir query params
        const params = new URLSearchParams()
        Object.entries(filters).forEach(([key, value]) => {
          if (value !== null && value !== undefined && value !== '') {
            params.append(key, String(value))
          }
        })
        const url = `/stocks/${ticker}/detail?${params.toString()}`
        const { data } = await api.get<StockDetailResponse>(url)
        this.detail = data
      } catch (error) {
        this.detail = null
        throw error
      }
    },

    // Disparo recalculo de recomendaciones
    async triggerRecalculate() {
      await api.post('/admin/recalculate')
      //MOSTRAR MENSAJE 
    },

    // Traigo las top-20 recomendaciones 
    async loadTopRecommendations() {
      const res = await api.get<Recommendation[]>('/recommendations/top20')
      this.recommendations = res.data
    },
  },
})