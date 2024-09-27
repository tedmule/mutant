package mutant

import (
	"encoding/json"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"

	v1 "k8s.io/api/admission/v1"
)

var universalDeserializer = serializer.NewCodecFactory(runtime.NewScheme()).UniversalDeserializer()

// Define the variables you need in the struct
type MutantCSI struct {
	StorageClass string // Specify storageclass name to enable PVC Loadbalancing
	Annotation   string // Specify annotation to disable PVC Loadbalancing
}

// This is the function that will be called to mutate
func (nfs *MutantCSI) Mutate(request v1.AdmissionRequest, storageclass string) (v1.AdmissionResponse, error) {
	response := v1.AdmissionResponse{}

	// Default response
	response.Allowed = true
	response.UID = request.UID

	// Decode the PVC object
	var pvc corev1.PersistentVolumeClaim
	if _, _, err := universalDeserializer.Decode(request.Object.Raw, nil, &pvc); err != nil {
		return response, err
	}

	// Mutate storageClassName of PVC
	opt := JSONPatchOpt{
		Path:  "/spec/storageClassName",
		Value: storageclass,
	}

	if pvc.Spec.StorageClassName == nil {
		opt.Op = "add"
	} else {
		opt.Op = "replace"
	}

	var patch JSONPatch
	patch = append(patch, opt)

	patchBytes, _ := json.Marshal(patch)

	response.Patch = []byte(patchBytes)
	response.PatchType = func() *v1.PatchType {
		pt := v1.PatchTypeJSONPatch
		return &pt
	}()

	return response, nil
}
