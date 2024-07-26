package metrics

import (
	"net/url"
	"sync/atomic"

	"github.com/canonical/lxd/shared/entity"
)

// RequestStatus represents the a completed request status category
type requestStatus string

const (
	serverError requestStatus = "server_error"
	clientError requestStatus = "client_error"
	success     requestStatus = "succeeded"
)

// PossibleRequestStatus is a slice that includes all the possible request status categories
var PossibleRequestStatus = []requestStatus{
	serverError,
	clientError,
	success,
}

func resolveStatusCode(statusCode int, asyncOperation bool) requestStatus {
	if asyncOperation {
		if statusCode == 400 {
			return serverError
		}
		return success
	}
	if statusCode > 500 && statusCode < 600 {
		return serverError
	}
	if statusCode > 400 && statusCode < 500 {
		return clientError
	}
	return success
}

func resolveEntityType(url url.URL) entity.Type {
	entityType, _, _, _, _ := entity.ParseURL(url)
	if entityType == entity.TypeContainer { // Handle containers as TypeInstance
		return entity.TypeInstance
	}
	return entityType
}

type completedMetricsLabeling struct {
	entityType entity.Type
	status     requestStatus
}

var ongoingRequests map[entity.Type]*int64
var completedRequests map[completedMetricsLabeling]*int64

// InitAPIMetrics initializes maps with initial values for the API rates metrics.
func InitAPIMetrics() {
	ongoingRequests = make(map[entity.Type]*int64)
	completedRequests = make(map[completedMetricsLabeling]*int64)

	for _, entityType := range entity.EntityTypes {
		if entityType != entity.TypeContainer { // We consider a container as an TypeInstance instead of TypeContainer
			ongoingRequests[entityType] = new(int64)
			for _, status := range PossibleRequestStatus {
				completedRequests[completedMetricsLabeling{entityType: entityType, status: status}] = new(int64)
			}
		}
	}
}

// MeasureOngoingRequest is used as a middleware before every request to keep track of ongoing requests.
func MeasureOngoingRequest(url url.URL) {
	atomic.AddInt64(ongoingRequests[resolveEntityType(url)], 1)
}

// MeasureCompletedRequest is called as a hook after each request is completed to keep track of completed requests.
func MeasureCompletedRequest(url url.URL, statusCode int, asyncOperation bool) {
	entityType := resolveEntityType(url)
	atomic.AddInt64(ongoingRequests[entityType], -1)
	atomic.AddInt64(completedRequests[completedMetricsLabeling{entityType: entityType, status: resolveStatusCode(statusCode, asyncOperation)}], 1)
}

// GetOngoingRequests gets the value for ongoing metrics filtered by entityType.
func GetOngoingRequests(entityType entity.Type) int64 {
	return *ongoingRequests[entityType]
}

// GetCompletedRequests gets the value of completed requests filtered by entityType and status.
func GetCompletedRequests(entityType entity.Type, status requestStatus) int64 {
	return *ongoingRequests[entityType]
}
