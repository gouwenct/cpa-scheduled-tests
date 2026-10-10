package main

/*
#include <stdint.h>
#include <stdlib.h>
#include <string.h>

typedef struct cliproxy_buffer {
    uint8_t* ptr;
    size_t len;
} cliproxy_buffer;

typedef int (*cliproxy_host_call_fn)(void* host_ctx, const char* method,
    const uint8_t* request, size_t request_len, cliproxy_buffer* response);
typedef void (*cliproxy_host_free_buffer_fn)(void* ptr, size_t len);

typedef struct cliproxy_host_api {
    uint32_t abi_version;
    void* host_ctx;
    cliproxy_host_call_fn call;
    cliproxy_host_free_buffer_fn free_buffer;
} cliproxy_host_api;

typedef int (*cliproxy_plugin_call_fn)(char* method, uint8_t* request, size_t request_len, cliproxy_buffer* response);
typedef void (*cliproxy_plugin_free_buffer_fn)(void* ptr, size_t len);
typedef void (*cliproxy_plugin_shutdown_fn)(void);

typedef struct cliproxy_plugin_api {
    uint32_t abi_version;
    cliproxy_plugin_call_fn call;
    cliproxy_plugin_free_buffer_fn free_buffer;
    cliproxy_plugin_shutdown_fn shutdown;
} cliproxy_plugin_api;

#ifdef _WIN32
#define CPA_PLUGIN_EXPORT __declspec(dllexport)
#else
#define CPA_PLUGIN_EXPORT
#endif

extern CPA_PLUGIN_EXPORT int scheduledTestsPluginCall(char* method, uint8_t* request, size_t request_len, cliproxy_buffer* response);
extern CPA_PLUGIN_EXPORT void scheduledTestsPluginFreeBuffer(void* ptr, size_t len);
extern CPA_PLUGIN_EXPORT void scheduledTestsPluginShutdown(void);

static const cliproxy_host_api* stored_host;

static inline void store_host_api(const cliproxy_host_api* host) {
    stored_host = host;
}

static inline void set_plugin_api(cliproxy_plugin_api* plugin) {
    plugin->abi_version = 1;
    plugin->call = scheduledTestsPluginCall;
    plugin->free_buffer = scheduledTestsPluginFreeBuffer;
    plugin->shutdown = scheduledTestsPluginShutdown;
}

static inline int call_host_api(const char* method, const uint8_t* request, size_t request_len, cliproxy_buffer* response) {
    if (stored_host == NULL || stored_host->call == NULL) return 1;
    return stored_host->call(stored_host->host_ctx, method, request, request_len, response);
}

static inline void free_host_buffer(void* ptr, size_t len) {
    if (stored_host != NULL && stored_host->free_buffer != NULL && ptr != NULL) {
        stored_host->free_buffer(ptr, len);
    }
}
*/
import "C"

import (
	"encoding/json"
	"fmt"
	"unsafe"
)

const (
	pluginID                      = "cpa-scheduled-tests"
	pluginVersion                 = "0.1.6"
	supportedSchemaVersion uint32 = 6
)

type lifecycleRequest struct {
	ConfigYAML    []byte `json:"config_yaml"`
	PluginDir     string `json:"plugin_dir,omitempty"`
	SchemaVersion uint32 `json:"schema_version,omitempty"`
}

type registrationResult struct {
	SchemaVersion int                    `json:"schema_version"`
	Metadata      registrationMetadata   `json:"metadata"`
	Capabilities  registrationCapability `json:"capabilities"`
}

type registrationMetadata struct {
	Name             string        `json:"Name"`
	Version          string        `json:"Version"`
	Author           string        `json:"Author"`
	GitHubRepository string        `json:"GitHubRepository,omitempty"`
	Description      string        `json:"Description"`
	ConfigFields     []configField `json:"ConfigFields,omitempty"`
}

type configField struct {
	Name         string `json:"Name"`
	Type         string `json:"Type"`
	Description  string `json:"Description"`
	DefaultValue any    `json:"DefaultValue,omitempty"`
}

type registrationCapability struct {
	ManagementAPI bool `json:"management_api"`
}

type managementRoute struct {
	Method      string `json:"Method"`
	Path        string `json:"Path"`
	Menu        string `json:"Menu,omitempty"`
	Description string `json:"Description,omitempty"`
}

type resourceRoute struct {
	Path        string `json:"Path"`
	Menu        string `json:"Menu,omitempty"`
	Description string `json:"Description,omitempty"`
}

type managementRegistration struct {
	Routes    []managementRoute `json:"routes,omitempty"`
	Resources []resourceRoute   `json:"resources,omitempty"`
}

type managementRequest struct {
	Method         string              `json:"Method"`
	Path           string              `json:"Path"`
	Headers        map[string][]string `json:"Headers"`
	Query          map[string][]string `json:"Query"`
	Body           []byte              `json:"Body"`
	HostCallbackID string              `json:"host_callback_id,omitempty"`
}

type managementResponse struct {
	StatusCode int                 `json:"StatusCode"`
	Headers    map[string][]string `json:"Headers,omitempty"`
	Body       []byte              `json:"Body,omitempty"`
}

func main() {}

//export cliproxy_plugin_init
func cliproxy_plugin_init(host *C.cliproxy_host_api, plugin *C.cliproxy_plugin_api) C.int {
	if host == nil || plugin == nil {
		return -1
	}
	C.store_host_api(host)
	C.set_plugin_api(plugin)
	return 0
}

//export scheduledTestsPluginCall
func scheduledTestsPluginCall(method *C.char, request *C.uint8_t, requestLen C.size_t, response *C.cliproxy_buffer) C.int {
	if response == nil {
		return -1
	}
	response.ptr = nil
	response.len = 0

	name := ""
	if method != nil {
		name = C.GoString(method)
	}
	requestBytes, ok := copyRequestBytes(request, requestLen)
	if !ok {
		return writeJSON(response, failEnvelope("invalid_request", "invalid request length"))
	}

	switch name {
	case "plugin.register", "plugin.reconfigure":
		var lifecycle lifecycleRequest
		_ = json.Unmarshal(requestBytes, &lifecycle)
		if err := ensureRuntime(lifecycle.PluginDir); err != nil {
			return writeJSON(response, failEnvelope("runtime_init_failed", err.Error()))
		}
		return writeJSON(response, okEnvelope(registration(lifecycle.SchemaVersion)))
	case "management.register":
		return writeJSON(response, okEnvelope(managementRegistrationResult()))
	case "management.handle":
		if err := ensureRuntime(""); err != nil {
			return writeJSON(response, failEnvelope("runtime_init_failed", err.Error()))
		}
		return writeJSON(response, handleManagement(requestBytes))
	case "plugin.shutdown":
		stopRuntime()
		return writeJSON(response, okEnvelope(map[string]any{"status": "stopped"}))
	default:
		return writeJSON(response, failEnvelope("unsupported_method", "unsupported plugin method"))
	}
}

//export scheduledTestsPluginFreeBuffer
func scheduledTestsPluginFreeBuffer(ptr unsafe.Pointer, length C.size_t) {
	_ = length
	if ptr != nil {
		C.free(ptr)
	}
}

//export scheduledTestsPluginShutdown
func scheduledTestsPluginShutdown() {
	stopRuntime()
}

func registration(hostSchemaVersion uint32) registrationResult {
	schemaVersion := hostSchemaVersion
	if schemaVersion == 0 {
		schemaVersion = 1
	}
	if schemaVersion > supportedSchemaVersion {
		schemaVersion = supportedSchemaVersion
	}
	return registrationResult{
		SchemaVersion: int(schemaVersion),
		Metadata: registrationMetadata{
			Name:             "CPA Scheduled 5H",
			Version:          pluginVersion,
			Author:           "gouwenct",
			GitHubRepository: "https://github.com/gouwenct/cpa-scheduled-tests",
			Description:      "Scheduled per-account Codex requests with five-hour window evidence, model selection, persistent logs, and send-all-now.",
			ConfigFields:     []configField{},
		},
		Capabilities: registrationCapability{ManagementAPI: true},
	}
}

func managementRegistrationResult() managementRegistration {
	routes := []managementRoute{
		{Method: "GET", Path: "/plugins/" + pluginID + "/state", Description: "Return settings, plans, accounts and jobs."},
		{Method: "POST", Path: "/plugins/" + pluginID + "/plans/save", Description: "Create or update a scheduled test plan."},
		{Method: "POST", Path: "/plugins/" + pluginID + "/plans/bulk-create", Description: "Create one scheduled test plan for every Codex account, skipping exact duplicates."},
		{Method: "POST", Path: "/plugins/" + pluginID + "/plans/delete", Description: "Delete a scheduled test plan."},
		{Method: "POST", Path: "/plugins/" + pluginID + "/settings", Description: "Update plugin runtime settings."},
		{Method: "POST", Path: "/plugins/" + pluginID + "/run", Description: "Run one plan or one account immediately."},
		{Method: "POST", Path: "/plugins/" + pluginID + "/run-all", Description: "Send one request to every Codex account, including disabled/unavailable records."},
		{Method: "GET", Path: "/plugins/" + pluginID + "/logs", Description: "Return recent persistent execution logs."},
		{Method: "POST", Path: "/plugins/" + pluginID + "/logs/clear", Description: "Clear persistent execution logs."},
	}
	resources := []resourceRoute{{Path: "/panel", Menu: "CPA Scheduled 5H", Description: "Manage scheduled account/model tests and run immediate requests."}}
	return managementRegistration{Routes: routes, Resources: resources}
}

func copyRequestBytes(request *C.uint8_t, requestLen C.size_t) ([]byte, bool) {
	length := int(requestLen)
	if length < 0 || C.size_t(length) != requestLen {
		return nil, false
	}
	if length == 0 {
		return nil, true
	}
	if request == nil {
		return nil, false
	}
	src := unsafe.Slice((*byte)(unsafe.Pointer(request)), length)
	return append([]byte(nil), src...), true
}

func callHost(method string, payload any, target any) error {
	raw, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	cMethod := C.CString(method)
	defer C.free(unsafe.Pointer(cMethod))

	var request *C.uint8_t
	var cPayload unsafe.Pointer
	if len(raw) > 0 {
		cPayload = C.CBytes(raw)
		if cPayload == nil {
			return fmt.Errorf("allocation failure")
		}
		defer C.free(cPayload)
		request = (*C.uint8_t)(cPayload)
	}

	var response C.cliproxy_buffer
	code := C.call_host_api(cMethod, request, C.size_t(len(raw)), &response)
	data := copyHostResponse(response)
	if response.ptr != nil {
		C.free_host_buffer(unsafe.Pointer(response.ptr), response.len)
	}
	if len(data) == 0 {
		return fmt.Errorf("host callback %s returned empty response (code=%d)", method, int(code))
	}

	var env struct {
		OK     bool            `json:"ok"`
		Result json.RawMessage `json:"result"`
		Error  *struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(data, &env); err != nil {
		return fmt.Errorf("decode host callback %s: %w", method, err)
	}
	if !env.OK {
		if env.Error != nil {
			return fmt.Errorf("%s: %s", env.Error.Code, env.Error.Message)
		}
		return fmt.Errorf("host callback %s failed", method)
	}
	if code != 0 {
		return fmt.Errorf("host callback %s returned code=%d", method, int(code))
	}
	if target != nil && len(env.Result) > 0 {
		if err := json.Unmarshal(env.Result, target); err != nil {
			return fmt.Errorf("decode host callback %s result: %w", method, err)
		}
	}
	return nil
}

func copyHostResponse(r C.cliproxy_buffer) []byte {
	if r.ptr == nil || r.len == 0 {
		return nil
	}
	return C.GoBytes(unsafe.Pointer(r.ptr), C.int(r.len))
}

func okEnvelope(result any) []byte {
	b, _ := json.Marshal(map[string]any{"ok": true, "result": result})
	return b
}

func failEnvelope(code, message string) []byte {
	b, _ := json.Marshal(map[string]any{
		"ok":    false,
		"error": map[string]any{"code": code, "message": message, "retryable": false},
	})
	return b
}

func writeJSON(response *C.cliproxy_buffer, data []byte) C.int {
	if len(data) == 0 {
		response.ptr = nil
		response.len = 0
		return 0
	}
	ptr := C.malloc(C.size_t(len(data)))
	if ptr == nil {
		return -1
	}
	C.memcpy(ptr, unsafe.Pointer(&data[0]), C.size_t(len(data)))
	response.ptr = (*C.uint8_t)(ptr)
	response.len = C.size_t(len(data))
	return 0
}
