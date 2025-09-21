# SYSTEM DESIGN

## Idempotencia & Reprocesamiento
- Cada métrica enriquecida se almacena en el repositorio en memoria bajo una clave única construida como `{fecha}-{campaignID}-{canal}`. Esta clave se genera en el método `Save` del repositorio usando `fmt.Sprintf`, por ejemplo: `2025-09-20-CAMPAIGN123-GoogleAds`.
- Si se reprocesa un registro (por ejemplo, si se vuelve a ejecutar el pipeline para la misma fecha y campaña), el registro anterior se sobrescribe automáticamente, ya que la clave es la misma. Esto garantiza idempotencia y evita duplicados.
- El método `Save` está protegido por un mutex para evitar condiciones de carrera en entornos concurrentes.

## Particionamiento & Retención
- El almacenamiento es en memoria, implementado como un `map[string]EnrichedMetric` dentro de la estructura `InMemoryRepository`.
- El acceso concurrente se gestiona con un `sync.RWMutex`, permitiendo múltiples lecturas simultáneas y escrituras exclusivas.
- El particionamiento lógico se basa en la clave de almacenamiento, permitiendo consultas eficientes por canal, campaña y rango de fechas mediante filtrado en memoria.
- La retención de datos está limitada por la vida del proceso y la memoria disponible. Al reiniciar el servicio, los datos se pierden.
- Para persistencia futura, la interfaz `MetricRepository` permite migrar a una base de datos sin cambiar la lógica de negocio.

## Concurrencia & Throughput
- La ingesta de datos de Ads y CRM se realiza concurrentemente usando goroutines y un `sync.WaitGroup` en el método `FetchData` del `Ingestor`. Esto reduce la latencia total de la ingesta.
- El repositorio usa `sync.RWMutex` para garantizar acceso seguro en operaciones concurrentes de lectura y escritura.
- El diseño permite escalar el procesamiento paralelizando la ingesta y el cálculo de métricas, aunque el almacenamiento en memoria puede ser un cuello de botella en grandes volúmenes.

## Calidad de Datos
- Los modelos (`AdPerformance`, `Opportunity`, `EnrichedMetric`) incluyen campos UTM (`utm_campaign`, `utm_source`, `utm_medium`) y validaciones lógicas en la transformación.
- El `Transformer` normaliza las claves UTM para asegurar coincidencias correctas entre Ads y CRM, y maneja la ausencia de datos con valores por defecto (por ejemplo, 0 para métricas numéricas).
- Se calculan métricas avanzadas como CPC (coste por clic), CPA (coste por adquisición), CVR (conversion rate), ROAS (return on ad spend), y ratios de conversión entre etapas del funnel.

## Observabilidad
- Se instrumentan métricas Prometheus: un contador de solicitudes (`api_requests_total`) y un histograma de duración (`api_request_duration_seconds`), ambos etiquetados por endpoint y método HTTP.
- El middleware Prometheus se aplica a cada endpoint, midiendo automáticamente cada petición.
- Se usan logs estructurados (por ejemplo, advertencias si no hay `SINK_URL` configurado en el Exporter, o errores al serializar o exportar métricas).
- Las métricas Prometheus pueden consultarse desde sistemas de monitoreo externos para análisis de performance y salud del servicio.

## Evolución en el Ecosistema Admira
- El diseño desacopla la lógica de negocio de la persistencia mediante la interfaz `MetricRepository`, permitiendo migrar a una base de datos relacional, NoSQL o data lake sin modificar el pipeline ETL.
- El sistema está preparado para exponer contratos de API versionados y para integrarse con otros sistemas del ecosistema Admira.
- El pipeline es extensible: se pueden añadir nuevos orígenes de datos (nuevos conectores de Ads o CRM), nuevos destinos (otros sinks o data lakes), y nuevas métricas calculadas simplemente extendiendo los modelos y la lógica de transformación.

## Refactorización y Principios SOLID/DRY
- Antes de hacer mas modificaciones, sería ideal refactorizar el handler siguiendo los principios SOLID y DRY. Por ejemplo, extraer un MetricsUtils para todo lo relacionado a Phrometeus, un QueryUtils que exponga funciones para validar los query parameters, y extraer en servicios la logica de negocio.