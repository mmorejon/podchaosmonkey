package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	// Define flags
	targetNamespace := flag.String("targetNamespace", "workloads", "Namespace used to remove pods.")
	excludeNamespaces := flag.String("excludeNamespaces", "kube-system", "Namespaces were pods can't be removed.")
	scheduler := flag.String("scheduler", "5s", "Scheduler to delete a random pod.")
	labelSelector := flag.String("labelSelector", "", "Label selector to filter the list of pods.")
	gracePeriod := flag.Int64("gracePeriod", int64(0), "Grace period to remove the pod.")

	flag.Parse()

	// validate target namespace
	valid := validateTargetNamespace(*targetNamespace, *excludeNamespaces)
	if valid == false {
		fmt.Println("Error: target namespace can't be found into the excluded namespaces list.")
		return
	}

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Panicln("Error creating cluster configuration.")
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panicln("Error creating a kubernetes client.")
	}

	// print init message
	fmt.Println("Starting chaos process ...")
	if len(*labelSelector) > 0 {
		fmt.Printf("Pods in the namespace %s with label %s will be removed every %s.\n\n", *targetNamespace, *labelSelector, *scheduler)
	} else {
		fmt.Printf("Pods in the namespace %s will be removed every %s.\n\n", *targetNamespace, *scheduler)
	}

	for {
		// wait until the next time to remove a random pod
		fmt.Println("Waiting for the next schedule.")

		sleepTimeDuration, _ := time.ParseDuration(*scheduler)
		time.Sleep(sleepTimeDuration)
		
		fmt.Println("It is time to remove a new pod ...")

		// get all pods in a namespace
		pods, err := clientset.CoreV1().Pods(*targetNamespace).List(context.Background(), metav1.ListOptions{
			LabelSelector: *labelSelector,
		})
		if err != nil {
			log.Panicln("Error listing pods.")
		}

		// get the number of available pods to be removed
		podsAvailable := len(pods.Items)

		if podsAvailable > 0 {
			// get random number from the list
			fmt.Printf("Number of pods available %d\n", podsAvailable)
			random := rand.Intn(podsAvailable)

			// delete pod
			podName := pods.Items[random].GetName()
			err := clientset.CoreV1().Pods(*targetNamespace).Delete(context.Background(), podName, metav1.DeleteOptions{
				GracePeriodSeconds: gracePeriod,
			})
			if err != nil {
				log.Panicln(fmt.Sprintf("Error removing the pod: %s", podName))
			} else {
				fmt.Printf("The pod %s was removed.\n\n", podName)
			}
		} else {
			fmt.Printf("\nThere no pod available to be removed in the namespace %s.\n\n", *targetNamespace)
		}
	}
}

// validateTargetNamespace checks if targetNamespace exist
// in the excluded namespaces list
func validateTargetNamespace(targetNamespace string, excludeNamespaces string) bool {

	if len(targetNamespace) == 0 {
		return false
	}
	
	nn := strings.Split(excludeNamespaces, ",")
	if contains(nn, targetNamespace) {
		return false
	} else {
		return true
	}
}

// contains checks if a string is present in a slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
