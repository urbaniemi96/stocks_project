import axios from 'axios'

export const api = axios.create({
  baseURL: 'http://localhost:8080',    // ajusta si tu Go server corre en otro puerto
  timeout: 5000,
})