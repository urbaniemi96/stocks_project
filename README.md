## 📚 Secciones de la Aplicación

La aplicación cuenta con cinco secciones principales:

### 1️⃣ Dashboard de Acciones
Pantalla principal donde se listan todas las acciones disponibles.  
- Permite buscar y filtrar acciones según distintos criterios.  
- Cada acción incluye un botón **"Ver Detalle"** para acceder a su información completa.

### 2️⃣ Detalle de la Acción
Vista individual con análisis de cada acción.  
- Muestra un resumen y el historial en gráficos interactivos.  
- Permite aplicar filtros por rango de fecha.  
- Opción para solicitar la opinión de una **IA experta en finanzas**.

### 3️⃣ Acciones Recomendadas
Ranking con el **Top 20 de acciones recomendadas** para invertir.  
- Muestra el score calculado por el sistema.  
- Permite acceder al detalle completo de cada acción (sección 2).

### 4️⃣ Panel de Administración
Sección administrativa para la gestión de datos.  
- Importación desde la **API oficial del desafío**.  
- Enriquecimiento de datos desde **Yahoo Finance**.  
- Recalculo del score de las acciones recomendadas.

### 5️⃣ Suscripción Premium
Sección para obtener acceso a beneficios exclusivos.  
- Ofrece distintos niveles de suscripción.  
- Incluye un botón de pago (actualmente muestra un mensaje de agradecimiento).

---

## 🛠️ Tecnologías Utilizadas

- **Backend:** Go con Gin y GORM  
- **Frontend:** Vue.js con Pinia, TypeScript y Tailwind CSS  
- **Base de Datos:** CockroachDB  
