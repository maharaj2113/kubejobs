package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	v1 "k8s.io/api/batch/v1"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	translatorAPIURL = "https://google-translator9.p.rapidapi.com/v2"
	apiHostHeader    = "google-translator9.p.rapidapi.com"
	apiKeyHeader     = "06de4ebb25mshad1d579533fcb83p1d7bebjsne0c7e73548b9" // Replace with your actual API key
	namespace        = "default"                                            // Kubernetes namespace
)

type TranslationRequest struct {
	Q      string `json:"q"`
	Source string `json:"source"`
	Target string `json:"target"`
	Format string `json:"format"`
}

type TranslationResponse struct {
	Data struct {
		Translations []struct {
			TranslatedText string `json:"translatedText"`
		} `json:"translations"`
	} `json:"data"`
}

func translateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req TranslationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	go createKubernetesJob(req)

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Job created successfully!"))
}

func createKubernetesJob(req TranslationRequest) {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		kubeconfig = "/root/.kube/config" // Correct kubeconfig path
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("Failed to load kubeconfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}

	job := &v1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("translate-job-%d", time.Now().Unix()),
		},
		Spec: v1.JobSpec{
			Template: v12.PodTemplateSpec{
				Spec: v12.PodSpec{
					Containers: []v12.Container{
						{
							Name:  "translate-container",
							Image: "alpine/curl:latest", // Use an image that includes curl
							Command: []string{
								"sh",
								"-c",
								fmt.Sprintf(`curl -X POST %s -H "x-rapidapi-host: %s" -H "x-rapidapi-key: %s" -d '{"q":"%s","source":"%s","target":"%s","format":"%s"}' -o /dev/stdout`,
                                    translatorAPIURL, apiHostHeader, apiKeyHeader, req.Q, req.Source, req.Target, req.Format),
                            },
							Env: []v12.EnvVar{
								{
									Name:  "TRANSLATE_TEXT",
									Value: req.Q,
								},
							},
						},
					},
					RestartPolicy: v12.RestartPolicyNever,
				},
			},
		},
	}

	job, err = clientset.BatchV1().Jobs(namespace).Create(context.Background(), job, metav1.CreateOptions{})
	if err != nil {
		log.Printf("Failed to create Job: %v", err)
		return
	}

	waitForJobCompletion(clientset, job.Name)
}

func waitForJobCompletion(clientset *kubernetes.Clientset, jobName string) {
	for {
		job, err := clientset.BatchV1().Jobs(namespace).Get(context.Background(), jobName, metav1.GetOptions{})
		if err != nil {
			log.Printf("Failed to get Job status: %v", err)
			return
		}

		if job.Status.Succeeded > 0 {
			log.Printf("Job %s completed successfully", jobName)
			return
		}

		if job.Status.Failed > 0 {
			log.Printf("Job %s failed", jobName)
			return
		}

		time.Sleep(5 * time.Second)
	}
}

func main() {
	http.HandleFunc("/translate", translateHandler)
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}