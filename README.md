# Admira ETL Backend (Go)

Servicio en Go que implementa un pipeline ETL para procesar datos de marketing digital. Este servicio cruza datos de Ads y CRM para calcular métricas clave de negocio, las expone a través de una API REST y permite exportarlas de forma segura.

---

## Características

- ETL Pipeline: Extrae datos de APIs externas, los transforma y los carga en un almacén en memoria.
- API REST: Proporciona endpoints para consultar métricas de negocio por canal y campaña.
- Exportación Segura: Exporta datos procesados a un servicio externo (`SINK_URL`) con firma digital HMAC-SHA256.
- Contenerización: Fácil de construir y ejecutar con Docker.
- Testing Unitario: Incluye pruebas para validar la lógica de transformación.

---

## Requisitos

Asegúrate de tener instalados:

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

---

## Configuración

1. **Clonar el repositorio:**
    ```bash
    git clone https://github.com/tu-usuario/tu-repositorio.git
    cd tu-repositorio
    ```

2. **Crear archivo `.env`:**
   Crea un archivo `.env` en la raíz del proyecto con las siguientes variables y añade tus URLs de Beeceptor:

    ```bash
    ADS_API_URL=<tu-url-ads>
    CRM_API_URL=<tu-url-crm>
    SINK_URL=<tu-url-sink>
    SINK_SECRET=admira_secret_example
    PORT=8080
    ```

---

## Ejecución

1. **Construir y ejecutar el servicio:**
    ```bash
    docker compose up --build
    ```

2. **Detener el servicio:**
    ```bash
    docker compose down
    ```

---

## Endpoints de la API

Base URL: `http://localhost:8080`

### 1. Ingestar y Procesar Datos
Activa el pipeline ETL.
- **POST** `/ingest/run`
- #### Sin filtro de fecha:

    ```bash
    curl -X POST http://localhost:8080/ingest/run
    ```

- #### Con filtro de fecha:

    ```bash
  
    curl -X POST "http://localhost:8080/ingest/run?since=2025-08-01"
    ```
#### Parámetros de consulta:
  since (opcional): Filtra los datos desde la fecha especificada en formato YYYY-MM-DD. Si no se proporciona, se procesarán todos los datos.
  **Response:**
    ```json
    {
      "status": "success",
      "message": "ETL pipeline executed successfully"
    }
    ```
### 2. Obtener Métricas por Canal
Consulta métricas agrupadas por canal.
- **GET** `/metrics/channel?from=YYYY-MM-DD&to=YYYY-MM-DD&channel=...`
    ```bash
    curl "http://localhost:8080/metrics/channel?channel=google_ads&from=2025-08-01&to=2025-08-31"
    ```
  **Response:**
    ```json
    [
      {
        "date": "2025-08-01",
        "channel": "google_ads",
        "clicks": 100,
        "cpc": 0.5,
        "revenue": 750.0,
        "roas": 15.0
      }
    ]
    ```

### 3. Obtener Métricas por Funnel
Consulta métricas agrupadas por campaña.
- **GET** `/metrics/funnel?from=YYYY-MM-DD&to=YYYY-MM-DD&utm_campaign=...`
    ```bash
    curl "http://localhost:8080/metrics/funnel?utm_campaign=back_to_school&from=2025-08-01&to=2025-08-31"
    ```
  **Response:**
    ```json
    [
      {
        "date": "2025-08-01",
        "campaign": "back_to_school",
        "leads": 2,
        "closed_won": 1,
        "revenue": 750.0,
        "cvr_lead_to_opp": 1.0,
        "cvr_opp_to_won": 0.5
      }
    ]
    ```

### 4. Exportar Datos
Exporta los datos procesados al servicio configurado.
- **POST** `/export/run`
    ```bash
    curl -X POST "http://localhost:8080/export/run?date=2025-08-01"
    ```
  **Response:**
    ```json
    {
      "status": "success",
      "message": "Data exported successfully"
    }
    ```

---

## Decisiones de Diseño

1. **Pipeline ETL:** Se diseñó para ser modular, permitiendo agregar nuevas fuentes de datos o transformaciones sin afectar el resto del sistema.
2. **Firma HMAC-SHA256:** Se eligió para garantizar la integridad y autenticidad de los datos exportados.
3. **Contenerización:** Se utilizó Docker para simplificar la configuración y despliegue del servicio.

---

## Limitaciones

1. **Almacenamiento en Memoria:** Actualmente, los datos procesados se almacenan en memoria, lo que puede ser un problema para grandes volúmenes de datos.
2. **Dependencia de APIs Externas:** El pipeline depende de la disponibilidad y consistencia de las APIs de Ads y CRM.
3. **Escalabilidad:** El diseño actual no está optimizado para entornos distribuidos o de alta concurrencia.

---

##  Manejo de UTMs Faltantes

El pipeline ETL está diseñado para manejar casos donde los datos de UTM (como `utm_campaign`, `utm_source` o `utm_medium`) estén ausentes o incompletos. En estos casos, se asignan valores predeterminados para garantizar que los datos puedan procesarse sin errores.

### Valores Predeterminados

- **utm_campaign:** `"unknown"`
- **utm_source:** `"unknown"`
- **utm_medium:** `"unknown"`

### Ejemplo

Si un registro de Ads tiene los siguientes valores:

```json
{
  "date": "2025-08-01",
  "campaign_id": "C-1001",
  "channel": "google_ads",
  "clicks": 100,
  "cost": 50.0,
  "utm_campaign": "",
  "utm_source": "",
  "utm_medium": ""
}
```
El pipeline asignará los valores predeterminados y generará la siguiente clave UTM: unknown|unknown|unknown
Esto asegura que los datos puedan agruparse y procesarse correctamente.
---

### Pruebas Unitarias

El proyecto incluye pruebas unitarias para validar la lógica principal del pipeline ETL. Estas pruebas aseguran que las transformaciones, exportaciones y cálculos se comporten como se espera.

#### Prueba de Exportación de Métricas

La prueba `TestExporter_ExportMetrics` valida que el componente `Exporter`:

- Genere correctamente la firma HMAC-SHA256 para los datos exportados.
- Envíe los datos al servicio de destino (`SINK_URL`) con los encabezados y formato adecuados.
- Maneje respuestas exitosas del servicio de destino.

Ejecuta esta prueba con el siguiente comando:

```bash
go test ./internal/etl -run TestExporter_ExportMetrics

```

#### Prueba de Ingesta de Datos

La prueba `TestIngestor_FetchData` valida que el componente `Ingestor`:

- Obtenga datos correctamente desde las APIs de Ads y CRM.
- Decodifique las respuestas JSON en las estructuras correspondientes.
- Maneje correctamente los casos exitosos de ingesta.

Ejecuta esta prueba con el siguiente comando:

```bash
go test ./internal/etl -run TestIngestor_FetchData
```

La prueba `TestCombineAndCalculateMetrics` valida que el componente `Transformer`:

- Combine correctamente los datos de Ads y CRM basados en las claves UTM.
- Calcule las métricas derivadas, como CPC, CPA, ROAS, y tasas de conversión (CVR).
- Maneje correctamente los casos donde los datos de entrada están incompletos o tienen valores predeterminados.

Ejecuta esta prueba con el siguiente comando:

```bash
go test ./internal/etl -run TestCombineAndCalculateMetrics
```

