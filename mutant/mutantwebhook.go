package mutant

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Mutant interface {
	Mutate(request v1.AdmissionRequest, storageclass string) (v1.AdmissionResponse, error)
}

type MutantWebhook struct {
	mutant Mutant
	config MutantConfig
	Server *echo.Echo
	k8s    *k8sWorker
}

// Based on:
// https://medium.com/ovni/writing-a-very-basic-kubernetes-mutating-admission-webhook-398dbbcb63ec
// https://github.com/alex-leonhardt/k8s-mutate-webhook

func (m *MutantWebhook) indexHandler(c echo.Context) error {
	log.Info("index")
	return c.String(http.StatusOK, "hello")
}

func (m *MutantWebhook) healthzHandler(c echo.Context) error {
	log.Info("healthz")
	return c.String(http.StatusOK, "healthz")
}

func (m *MutantWebhook) mutateHandler(c echo.Context) error {
	contentType := c.Request().Header.Get("Content-Type")

	if contentType != "application/json" {
		return c.String(http.StatusUnsupportedMediaType, "support application/json only")
	}

	// Set basic info(workaround, optimize later)
	reviewResponse := v1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1",
		},
	}

	// Parse the AdmissionReview from request body
	admissionReview := v1.AdmissionReview{}
	bodyBytes, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if err := json.Unmarshal(bodyBytes, &admissionReview); err != nil {
		return c.String(http.StatusBadRequest, "Invalid JSON")
	}

	uid := admissionReview.Request.UID
	log.Infof("====== Receive mutating request: %s\n", uid)
	log.Debugf("%s\n", PrettyJson(admissionReview))

	// Allowed by default
	defaultResponse := &v1.AdmissionResponse{
		UID:     uid,
		Allowed: true,
	}
	// SET DEFAULT RESPONSE(allowed)
	reviewResponse.Response = defaultResponse

	// Find annotation: mutant=no, skip mutation
	// Decode the PVC object
	var pvc corev1.PersistentVolumeClaim
	if _, _, err := universalDeserializer.Decode(admissionReview.Request.Object.Raw, nil, &pvc); err != nil {
		log.Errorf("Decode PVC in request failed: %s\n", err.Error())
		return c.JSON(http.StatusBadRequest, "Decode PVC error")
	} else {
		if val, ok := pvc.Annotations[m.config.Annotation]; ok {
			log.Infof("Found annotation [%s=%s], SKIP mutation.", m.config.Annotation, val)
			return c.JSON(http.StatusOK, reviewResponse)
		}
	}

	// Get storageclss name from MutantCSI
	sc := m.mutant.(*MutantCSI).StorageClass
	weightedStorageClasses := m.k8s.listWeightedStorageClass(sc)

	if len(weightedStorageClasses) == 0 {
		log.Warnf("Found 0 StorageClass for type %s, skip mutating.\n", sc)
		reviewResponse.Response = defaultResponse
	} else {
		selected := WeightedRandomSelect(weightedStorageClasses)
		log.Infof("SELECTED STORAGECLASS: %s\n", selected.Value)
		response, err := m.mutant.Mutate(*admissionReview.Request, selected.Value)
		if err != nil {
			log.Errorf("Mutate error: %s\n", err.Error())
			reviewResponse.Response = defaultResponse
		} else {
			// Fill Mutation
			reviewResponse.Response = &response
		}
	}

	resp, _ := json.Marshal(reviewResponse)
	log.Infof("====== Response: %s\n", string(resp))

	return c.JSON(http.StatusOK, reviewResponse)
}

func (m *MutantWebhook) debugHandler(c echo.Context) error {
	items := m.k8s.listWeightedStorageClass("nfs.csi.k8s.io")
	fmt.Println(items)
	// selected := WeightedRandomSelect(items)
	// log.Infof("Select storage class: %s\n", selected.Value)

	// fmt.Println(m.mutant.StorageClass)
	fmt.Printf("%+v\n", m.mutant)
	fmt.Println(m.mutant.(*MutantCSI).StorageClass)

	return c.String(http.StatusOK, "debug\n")
}

func NewMutantWebhook(mutant Mutant, config MutantConfig) (*MutantWebhook, error) {

	k8sWorker, err := NewK8SWorker(config)
	if err != nil {
		log.Fatal("shit")
	}

	webhook := &MutantWebhook{
		config: config,
		mutant: mutant,
		k8s:    k8sWorker,
	}

	e := echo.New()

	// e.Use(middleware.Logger())
	// log.SetFormatter(&log.JSONFormatter{})
	// log.SetReportCaller(true)
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:       true,
		LogStatus:    true,
		LogLatency:   true,
		LogUserAgent: true,
		LogRequestID: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			log.WithFields(log.Fields{
				"URI":       values.URI,
				"status":    values.Status,
				"latency":   values.Latency,
				"reqid":     values.RequestID,
				"useragent": values.UserAgent,
			}).Info("HTTP")

			return nil
		},
	}))

	e.GET("/", webhook.indexHandler)
	e.GET("/healthz", webhook.healthzHandler)
	e.POST("/mutate", webhook.mutateHandler)
	e.GET("/debug", webhook.debugHandler)

	webhook.Server = e

	return webhook, nil
}
