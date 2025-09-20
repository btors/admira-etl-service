# SYSTEM_DESIGN.md

## Idempotencia & Reprocesamiento
El pipeline ETL asegura idempotencia al identificar registros únicos mediante claves naturales (e.g., `UTM` y fechas). Esto permite reprocesar datos sin duplicar resultados.

## Particionamiento & Retención
Los datos se particionan por fechas y campañas para optimizar consultas. La retención se gestiona en memoria, limitando el tamaño del dataset procesado.

## Concurrencia & Throughput
Se utilizan `goroutines` y un `worker pool` para procesar datos en paralelo, maximizando el throughput y aprovechando múltiples núcleos.

## Calidad de Datos
Se manejan UTMs ausentes con valores predeterminados (`fallbacks`) para evitar errores. Validaciones aseguran consistencia en los datos procesados.

## Observabilidad
Se implementan logs estructurados para rastrear el flujo ETL y métricas clave (e.g., tiempo de procesamiento, errores) para monitoreo en tiempo real.

## Evolución en el Ecosistema Admira
El diseño es compatible con un futuro `data lake` para almacenamiento histórico y contratos de API versionados, garantizando interoperabilidad y escalabilidad.