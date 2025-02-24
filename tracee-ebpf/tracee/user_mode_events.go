// Invoked tracee-ebpf events from user mode
// This utility can prove itself useful to generate information needed by signatures that is not provided by normal
// events in the kernel.
// Because the events in the kernel are invoked by other programs behavior, we cannot anticipate which events will be
// invoked and as a result what information will be extracted.
// This is critical because tracee-rules is independent, and doesn't have to run on the same machine as tracee-ebpf.
// This means that tracee-rules might lack basic information of the operating machine needed for some signatures.
// By creating user mode events this information could be intentionally collected and passed to tracee-ebpf afterwards.
package tracee

import (
	"github.com/aquasecurity/tracee/tracee-ebpf/external"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const InitProcNsDir = "/proc/1/ns"

// CreateInitNamespacesEvent collect the init process namespaces and create event from them.
func CreateInitNamespacesEvent() (external.Event, error) {
	initNamespacesArgs := getInitNamespaceArguments()
	initNamespacesEvent := external.Event{
		Timestamp:   int(time.Now().UnixNano()),
		ProcessName: "tracee-ebpf",
		EventID:     int(InitNamespacesEventID),
		EventName:   EventsIDToEvent[InitNamespacesEventID].Name,
		ArgsNum:     len(initNamespacesArgs),
		Args:        initNamespacesArgs,
	}
	return initNamespacesEvent, nil
}

// getInitNamespaceArguments Fetch the namespaces of the init process and parse them into event arguments.
func getInitNamespaceArguments() []external.Argument {
	initNamespaces := fetchInitNamespaces()
	initNamespacesArgs := make([]external.Argument, len(EventsIDToParams[InitNamespacesEventID]))
	for i, arg := range initNamespacesArgs {
		arg.ArgMeta = EventsIDToParams[InitNamespacesEventID][i]
		arg.Value = initNamespaces[arg.Name]
		initNamespacesArgs[i] = arg
	}
	return initNamespacesArgs
}

// fetchInitNamespaces fetch the namespaces values from the /proc/1/ns directory
func fetchInitNamespaces() map[string]uint32 {
	initNamespacesMap := make(map[string]uint32)
	namespaceValueReg := regexp.MustCompile(":[[[:digit:]]*]")
	namespacesLinks, _ := ioutil.ReadDir(InitProcNsDir)
	for _, namespaceLink := range namespacesLinks {
		linkString, _ := os.Readlink(filepath.Join(InitProcNsDir, namespaceLink.Name()))
		trim := strings.Trim(namespaceValueReg.FindString(linkString), "[]:")
		namespaceNumber, _ := strconv.ParseUint(trim, 10, 32)
		initNamespacesMap[namespaceLink.Name()] = uint32(namespaceNumber)
	}
	return initNamespacesMap
}
