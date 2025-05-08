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
  }),
  // Acciones para modificar estado
  actions: {
    // Disparo petición al fetch del back (para traer datos de la api y guardar en la bd)
    async fetchAndStore() {
      // dispara /fetch en el backend
      await api.get('/fetch')
    },
    // Traigo los stocks guardados en la bd y los guardo en list (estado del store)
    async loadStocks() {
      // dispara /stocks en el backend
      const res = await api.get<Stock[]>('/stocks')
      this.list = res.data 
    },
    // Llamo al back para recomendar una acción
    async computeRecommendation() {
      const res = await api.get<Stock>('/recommend')
      this.recommended = res.data
    },
  },
})