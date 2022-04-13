package main

import	(
	"fmt"
	"context"
    "flag"
    "log"
    "os"
    "path/filepath"
    "strings"
    batchv1 "k8s.io/api/batch/v1"
    v1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    kubernetes "k8s.io/client-go/kubernetes"
    clientcmd "k8s.io/client-go/tools/clientcmd"
)

func connectToK8s() *kubernetes.Clientset {
	home, exists := os.LookupEnv("HOME")
	if !exists {
		home = "/root"
	}

	configPath := filepath.Join(home, ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		log.Panicln("Failed to create K8s config")
	}

	clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        log.Panicln("Failed to create K8s clientset")
    }

    return clientset
}

func launchK8sJob(clientset *kubernetes.Clientset, jobName *string, image *string, cmd *string) {
    jobs := clientset.BatchV1().Jobs("default")
    var backOffLimit int32 = 0

    jobSpec := &batchv1.Job{
        ObjectMeta: metav1.ObjectMeta{
            Name:      *jobName,
            Namespace: "default",
        },
        Spec: batchv1.JobSpec{
            Template: v1.PodTemplateSpec{
                Spec: v1.PodSpec{
                    Containers: []v1.Container{
                        {
                            Name:    *jobName,
                            Image:   *image,
                            Command: strings.Split(*cmd, " "),
                        },
                    },
                    RestartPolicy: v1.RestartPolicyNever,
                },
            },
            BackoffLimit: &backOffLimit,
        },
    }

    _, err := jobs.Create(context.TODO(), jobSpec, metav1.CreateOptions{})
    if err != nil {
        log.Fatalln("Failed to create K8s job.")
    }

    //print job details
    log.Println("Created K8s job successfully")
}

func main () {
	jobName := flag.String("jobname", "test-job", "the name of the Job")
	containerImage := flag.String("image", "ubuntu:latest", "the name of our container image")
	entryCommand := flag.String("command", "ls", "A command we will run inside the container we created in previous step")

	flag.Parse()

	fmt.Printf("Args : %s %s %s\n", *jobName, *containerImage, *entryCommand)

	clientset := connectToK8s()
}