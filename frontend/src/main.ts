import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import './style.css'

import 'datatables.net-dt/css/dataTables.dataTables.min.css'
import $ from 'jquery'
import 'datatables.net'

// Sobrescribimos las clases que DataTables aplica por defecto
$.extend($.fn.dataTable.ext.classes, {
  sWrapper:    'dataTables_wrapper dt-tailwind',
  sFilterInput:  'form-input block w-full px-3 py-2 border rounded',
  sLengthSelect: 'form-select block w-20 px-2 py-1 border rounded',
  sPaging:       'inline-flex items-center space-x-2 mt-4',
  sPageButton:   'px-3 py-1 rounded border bg-white hover:bg-gray-100',
  sPageButtonActive: 'bg-blue-500 text-white',
})

const app = createApp(App)
app.use(createPinia())
app.mount('#app')