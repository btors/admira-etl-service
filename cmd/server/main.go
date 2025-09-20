package main

import (
	"log"

	// Asegúrate de que este import coincida con el nombre de tu módulo en go.mod
	"github.com/btors/admira-etl/internal/config"
	"github.com/btors/admira-etl/internal/etl"
)

func main() {
	// Paso 1: Cargar la configuración desde el archivo .env
	// Esto lee las URLs de tu API de mock (Beeceptor) y el puerto.
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("FATAL: No se pudo cargar la configuración: %v", err)
	}

	// Paso 2: Ejecutar la Ingesta (Extracción)
	// Crea una instancia del Ingestor y llama a FetchData para obtener los datos
	// de Ads y CRM de forma concurrente.
	log.Println("INFO: Iniciando la ingesta de datos...")
	ingestor := etl.NewIngestor(cfg.AdsAPIURL, cfg.CrmAPIURL)
	ads, crm, err := ingestor.FetchData()
	if err != nil {
		log.Fatalf("FATAL: La ingesta de datos falló: %v", err)
	}
	log.Printf("INFO: Ingesta exitosa. Se obtuvieron %d registros de anuncios y %d registros de CRM.", len(ads), len(crm))

	// Paso 3: Ejecutar la Transformación
	// Crea una instancia del Transformer y llama a la función para combinar los datos
	// y calcular todas las métricas de negocio.
	log.Println("INFO: Iniciando la transformación de datos...")
	transformer := etl.NewTransformer()
	enrichedData, err := transformer.CombineAndCalculateMetrics(ads, crm)
	if err != nil {
		log.Fatalf("FATAL: La transformación de datos falló: %v", err)
	}
	log.Printf("INFO: Transformación exitosa. Se generaron %d métricas enriquecidas.", len(enrichedData))

	// Paso 4: Verificar el resultado
	// Imprime el primer registro enriquecido en la consola para que podamos
	// inspeccionar los resultados y confirmar que los cálculos son correctos.
	if len(enrichedData) > 0 {
		log.Printf("DEBUG: Métrica de ejemplo: %+v\n", enrichedData[0])
	}
}
