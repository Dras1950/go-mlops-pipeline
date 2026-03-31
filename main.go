package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"math/rand"

	"github.com/gorilla/mux"
)

// PipelineStep represents a single step in the MLOps pipeline
type PipelineStep struct {
	Name     string
	Status   string
	Duration time.Duration
	StartedAt time.Time
	CompletedAt time.Time
}

// MLOpsPipeline represents the entire MLOps workflow
type MLOpsPipeline struct {
	ID    string
	Name  string
	Steps []PipelineStep
	Status string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewMLOpsPipeline creates a new MLOps pipeline instance
func NewMLOpsPipeline(name string) *MLOpsPipeline {
	return &MLOpsPipeline{
		ID:        fmt.Sprintf("pipeline-%d", time.Now().UnixNano()),
		Name:      name,
		Steps:     []PipelineStep{},
		Status:    "CREATED",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// AddStep adds a new step to the pipeline
func (p *MLOpsPipeline) AddStep(stepName string) {
	step := PipelineStep{
		Name:   stepName,
		Status: "PENDING",
	}
	p.Steps = append(p.Steps, step)
	p.UpdatedAt = time.Now()
}

// RunPipeline simulates the execution of the MLOps pipeline
func (p *MLOpsPipeline) RunPipeline() {
	p.Status = "RUNNING"
	p.UpdatedAt = time.Now()
	log.Printf("Running pipeline: %s (ID: %s)", p.Name, p.ID)

	for i := range p.Steps {
		p.Steps[i].Status = "RUNNING"
		p.Steps[i].StartedAt = time.Now()
		p.UpdatedAt = time.Now()
		log.Printf("  Executing step: %s", p.Steps[i].Name)
		time.Sleep(time.Duration(rand.Intn(3)+1) * time.Second) // Simulate work
		p.Steps[i].Status = "COMPLETED"
		p.Steps[i].CompletedAt = time.Now()
		p.Steps[i].Duration = p.Steps[i].CompletedAt.Sub(p.Steps[i].StartedAt)
		p.UpdatedAt = time.Now()
		log.Printf("  Step %s completed in %v", p.Steps[i].Name, p.Steps[i].Duration)
	}

	p.Status = "COMPLETED"
	p.UpdatedAt = time.Now()
	log.Printf("Pipeline %s completed.", p.Name)
}

// handleCreatePipeline creates a new pipeline
func handleCreatePipeline(w http.ResponseWriter, r *http.Request) {
	p := NewMLOpsPipeline("Data Processing and Model Training")
	p.AddStep("Data Ingestion")
	p.AddStep("Data Preprocessing")
	p.AddStep("Feature Engineering")
	p.AddStep("Model Training")
	p.AddStep("Model Evaluation")

	go p.RunPipeline() // Run pipeline in a goroutine

	fmt.Fprintf(w, "Pipeline %s created and started!\n", p.ID)
}

// handleGetPipelineStatus returns the status of a pipeline
func handleGetPipelineStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pipelineID := vars["id"]
	// In a real app, you would fetch pipeline status from a store
	fmt.Fprintf(w, "Status for pipeline %s: (simulated)\n", pipelineID)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/pipeline", handleCreatePipeline).Methods("POST")
	r.HandleFunc("/pipeline/{id}", handleGetPipelineStatus).Methods("GET")

	fmt.Println("MLOps Pipeline Service listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
