import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import './style.css'

import { router } from './router'

// Datatables
import 'datatables.net-dt/css/dataTables.dataTables.min.css'
import $ from 'jquery'
import 'datatables.net'

// Font Awesome
import { library } from '@fortawesome/fontawesome-svg-core'
import { fas } from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

// Vue Toastification
import Toast from "vue-toastification";
import type { PluginOptions } from 'vue-toastification'
import "vue-toastification/dist/index.css";

import { useAuthStore } from './stores/auth'



library.add(fas)

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

// Opciones de vue toastification
const options: PluginOptions = {
  timeout: 3000,
  closeOnClick: true,
  pauseOnHover: true,
  draggable: true,
};
app.use(Toast, options);

app.component('font-awesome-icon', FontAwesomeIcon)
app.use(createPinia())

// Traigo las credenciales del usuario
const auth = useAuthStore()
try {
  await auth.fetchUser()
} catch {
  // Usuario de demostración por defecto
  auth.setUser({ id: "demo-user", role: 'user' })
}


app.use(router)  
app.mount('#app')