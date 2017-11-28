package registry

var defaultOs = "centos"
var defaultOsVersion = "7.3"
var defaultLanguage = "java"
var defaultLanguageVersion = "1.8"

type WorkerRegistry struct {
	// a map from manifest to workerInfo[]
	// a map from projectId to workerInfo[]  /* history affinity */
}

// Simple implementation with global var of WorkerRegistry
var GlobalWorkerRegistry WorkerRegistry = WorkerRegistry{}
