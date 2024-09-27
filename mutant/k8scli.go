package mutant

import (
	"context"
	"os"
	"path/filepath"
	"strconv"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type k8sWorker struct {
	client *kubernetes.Clientset
}

func initK8SClient(mode string) (*kubernetes.Clientset, error) {
	client := &kubernetes.Clientset{}

	if mode == "outofcluster" {
		home, err := os.UserHomeDir()
		if err != nil {
			home = "~"
		}

		kubeconfig := filepath.Join(home, ".kube", "config")
		configFromFlags, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Fatalf("Error building kubeconfig: %s\n", err.Error())
		}

		client, err = kubernetes.NewForConfig(configFromFlags)
		if err != nil {
			log.Fatalf("Error building Kubernetes client: %s\n", err.Error())
		}
		return client, nil

	}
	return client, nil
}

func (k8s *k8sWorker) listWeightedStorageClass(storageclass string) []WeightedItem {
	w := []WeightedItem{}

	// List all StorageClass
	storageClassList, err := k8s.client.StorageV1().StorageClasses().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Errorf("List StorageClass error: %s", err.Error())
		return w
	}

	for _, sc := range storageClassList.Items {
		item := WeightedItem{}
		if sc.Provisioner == storageclass {
			item.Value = sc.Name
			weight := sc.Parameters["weight"]
			if len(weight) > 0 {
				num, err := strconv.ParseInt(weight, 10, 64)
				if err != nil {
					item.Weight = 1
				} else {
					item.Weight = int(num)
				}
			} else {
				item.Weight = 1
			}
			w = append(w, item)
		}
	}
	log.Infof("%+v\n", w)
	return w
}

// func (k8s *k8sWorker) listPVC() {
// 	namespace := "default" // empty string means all namespaces
// 	labelSelector := "type=fast-storage"

// 	pvcList, err := k8s.client.CoreV1().PersistentVolumeClaims(namespace).List(context.TODO(), metav1.ListOptions{
// 		LabelSelector: labelSelector,
// 	})

// 	if err != nil {
// 		log.Fatalf("Error listing PVCs: %s", err.Error())
// 	}

// 	for _, pvc := range pvcList.Items {
// 		fmt.Printf("PVC Name: %s, StorageClass: %s\n", pvc.Name, pvc.Spec.StorageClassName)
// 	}
// }

func NewK8SWorker(config MutantConfig) (*k8sWorker, error) {
	var worker k8sWorker

	cli, err := initK8SClient(config.Mode)
	if err != nil {
		return &worker, err
	}

	worker.client = cli
	return &worker, nil
}
