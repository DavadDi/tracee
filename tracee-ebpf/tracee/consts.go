package tracee

import (
	"github.com/aquasecurity/libbpfgo/helpers"
	"math"

	"github.com/aquasecurity/tracee/tracee-ebpf/external"
)

// bpfConfig is an enum that include various configurations that can be passed to bpf code
// config should match defined values in ebpf code
type bpfConfig uint32

// Max depth of each stack trace to track
// Matches 'MAX_STACK_DEPTH' in eBPF code
const maxStackDepth int = 20

// Custom KernelConfigOption's to extend kernel_config helper support
// Add here all kconfig variables used within tracee.bpf.c
const (
	CONFIG_ARCH_HAS_SYSCALL_WRAPPER helpers.KernelConfigOption = iota + helpers.CUSTOM_OPTION_START
)

const (
	configDetectOrigSyscall bpfConfig = iota + 1
	configExecEnv
	configCaptureFiles
	configExtractDynCode
	configTraceePid
	configStackAddresses
	configUIDFilter
	configMntNsFilter
	configPidNsFilter
	configUTSNsFilter
	configCommFilter
	configPidFilter
	configContFilter
	configFollowFilter
	configNewPidFilter
	configNewContFilter
	configDebugNet
	configProcTreeFilter
	configCaptureModules
	configCgroupV1
)

const (
	filterNotEqual uint32 = iota
	filterEqual
)

const (
	filterIn  uint8 = 1
	filterOut uint8 = 2
)

const (
	uidLess uint32 = iota
	uidGreater
	pidLess
	pidGreater
	mntNsLess
	mntNsGreater
	pidNsLess
	pidNsGreater
)

// Set default inequality values
// val<0 and val>math.MaxUint64 should never be used by the user as they give an empty set
const (
	LessNotSetUint    uint64 = 0
	GreaterNotSetUint uint64 = math.MaxUint64
	LessNotSetInt     int64  = math.MinInt64
	GreaterNotSetInt  int64  = math.MaxInt64
)

// an enum that specifies the index of a function to be used in a bpf tail call
// tail function indexes should match defined values in ebpf code
const (
	tailVfsWrite uint32 = iota
	tailVfsWritev
	tailSendBin
	tailSendBinTP
)

// binType is an enum that specifies the type of binary data sent in the file perf map
// binary types should match defined values in ebpf code
type binType uint8

const (
	sendVfsWrite binType = iota + 1
	sendMprotect
	sendKernelModule
)

// argType is an enum that encodes the argument types that the BPF program may write to the shared buffer
// argument types should match defined values in ebpf code
type argType uint8

const (
	noneT argType = iota
	intT
	uintT
	longT
	ulongT
	offT
	modeT
	devT
	sizeT
	pointerT
	strT
	strArrT
	sockAddrT
	bytesT
	u16T
	credT
	intArr2T
	argsArrT
)

// ProbeType is an enum that describes the mechanism used to attach the event
// Kprobes are explained here: https://github.com/iovisor/bcc/blob/master/docs/reference_guide.md#1-kprobes
// Tracepoints are explained here: https://github.com/iovisor/bcc/blob/master/docs/reference_guide.md#3-tracepoints
// Raw tracepoints are explained here: https://github.com/iovisor/bcc/blob/master/docs/reference_guide.md#7-raw-tracepoints
type probeType uint8

const (
	sysCall probeType = iota
	kprobe
	kretprobe
	tracepoint
	rawTracepoint
)

type probe struct {
	event  string
	attach probeType
	fn     string
}

// EventConfig is a struct describing an event configuration
type EventConfig struct {
	ID             int32
	ID32Bit        int32
	Name           string
	Probes         []probe
	EssentialEvent bool
	Sets           []string
}

// Non syscalls events (used by all architectures)
// events should match defined values in ebpf code
const (
	SysEnterEventID int32 = iota + 1000
	SysExitEventID
	SchedProcessForkEventID
	SchedProcessExecEventID
	SchedProcessExitEventID
	SchedSwitchEventID
	DoExitEventID
	CapCapableEventID
	VfsWriteEventID
	VfsWritevEventID
	MemProtAlertEventID
	CommitCredsEventID
	SwitchTaskNSEventID
	MagicWriteEventID
	CgroupAttachTaskEventID
	CgroupMkdirEventID
	CgroupRmdirEventID
	SecurityBprmCheckEventID
	SecurityFileOpenEventID
	SecurityInodeUnlinkEventID
	SecuritySocketCreateEventID
	SecuritySocketListenEventID
	SecuritySocketConnectEventID
	SecuritySocketAcceptEventID
	SecuritySocketBindEventID
	SecuritySbMountEventID
	SecurityBPFEventID
	SecurityBPFMapEventID
	SecurityKernelReadFileEventID
	SecurityInodeMknodEventID
	SecurityPostReadFileEventID
	SocketDupEventID
	MaxEventID
)

// Events originated from user-space
const (
	InitNamespacesEventID int32 = iota + 2000
)

const (
	NetPacket uint32 = iota
	DebugNetSecurityBind
	DebugNetUdpSendmsg
	DebugNetUdpDisconnect
	DebugNetUdpDestroySock
	DebugNetUdpV6DestroySock
	DebugNetInetSockSetState
	DebugNetTcpConnect
)

// EventsIDToEvent is list of supported events, indexed by their ID
var EventsIDToEvent = map[int32]EventConfig{
	ReadEventID:                   {ID: ReadEventID, ID32Bit: sys32read, Name: "read", Probes: []probe{{event: "read", attach: sysCall, fn: "read"}}, Sets: []string{"syscalls", "fs", "fs_read_write"}},
	WriteEventID:                  {ID: WriteEventID, ID32Bit: sys32write, Name: "write", Probes: []probe{{event: "write", attach: sysCall, fn: "write"}}, Sets: []string{"syscalls", "fs", "fs_read_write"}},
	OpenEventID:                   {ID: OpenEventID, ID32Bit: sys32open, Name: "open", Probes: []probe{{event: "open", attach: sysCall, fn: "open"}}, Sets: []string{"default", "syscalls", "fs", "fs_file_ops"}},
	CloseEventID:                  {ID: CloseEventID, ID32Bit: sys32close, Name: "close", Probes: []probe{{event: "close", attach: sysCall, fn: "close"}}, Sets: []string{"default", "syscalls", "fs", "fs_file_ops"}},
	StatEventID:                   {ID: StatEventID, ID32Bit: sys32stat, Name: "stat", Probes: []probe{{event: "newstat", attach: sysCall, fn: "newstat"}}, Sets: []string{"default", "syscalls", "fs", "fs_file_attr"}},
	FstatEventID:                  {ID: FstatEventID, ID32Bit: sys32fstat, Name: "fstat", Probes: []probe{{event: "newfstat", attach: sysCall, fn: "newfstat"}}, Sets: []string{"default", "syscalls", "fs", "fs_file_attr"}},
	LstatEventID:                  {ID: LstatEventID, ID32Bit: sys32lstat, Name: "lstat", Probes: []probe{{event: "newlstat", attach: sysCall, fn: "newlstat"}}, Sets: []string{"default", "syscalls", "fs", "fs_file_attr"}},
	PollEventID:                   {ID: PollEventID, ID32Bit: sys32poll, Name: "poll", Probes: []probe{{event: "poll", attach: sysCall, fn: "poll"}}, Sets: []string{"syscalls", "fs", "fs_mux_io"}},
	LseekEventID:                  {ID: LseekEventID, ID32Bit: sys32lseek, Name: "lseek", Probes: []probe{{event: "lseek", attach: sysCall, fn: "lseek"}}, Sets: []string{"syscalls", "fs", "fs_read_write"}},
	MmapEventID:                   {ID: MmapEventID, ID32Bit: sys32mmap, Name: "mmap", Probes: []probe{{event: "mmap", attach: sysCall, fn: "mmap"}}, Sets: []string{"syscalls", "proc", "proc_mem"}},
	MprotectEventID:               {ID: MprotectEventID, ID32Bit: sys32mprotect, Name: "mprotect", Probes: []probe{{event: "mprotect", attach: sysCall, fn: "mprotect"}}, Sets: []string{"syscalls", "proc", "proc_mem"}},
	MunmapEventID:                 {ID: MunmapEventID, ID32Bit: sys32munmap, Name: "munmap", Probes: []probe{{event: "munmap", attach: sysCall, fn: "munmap"}}, Sets: []string{"syscalls", "proc", "proc_mem"}},
	BrkEventID:                    {ID: BrkEventID, ID32Bit: sys32brk, Name: "brk", Probes: []probe{{event: "brk", attach: sysCall, fn: "brk"}}, Sets: []string{"syscalls", "proc", "proc_mem"}},
	RtSigactionEventID:            {ID: RtSigactionEventID, ID32Bit: sys32rt_sigaction, Name: "rt_sigaction", Probes: []probe{{event: "rt_sigaction", attach: sysCall, fn: "rt_sigaction"}}, Sets: []string{"syscalls", "signals"}},
	RtSigprocmaskEventID:          {ID: RtSigprocmaskEventID, ID32Bit: sys32rt_sigprocmask, Name: "rt_sigprocmask", Probes: []probe{{event: "rt_sigprocmask", attach: sysCall, fn: "rt_sigprocmask"}}, Sets: []string{"syscalls", "signals"}},
	RtSigreturnEventID:            {ID: RtSigreturnEventID, ID32Bit: sys32rt_sigreturn, Name: "rt_sigreturn", Probes: []probe{{event: "rt_sigreturn", attach: sysCall, fn: "rt_sigreturn"}}, Sets: []string{"syscalls", "signals"}},
	IoctlEventID:                  {ID: IoctlEventID, ID32Bit: sys32ioctl, Name: "ioctl", Probes: []probe{{event: "ioctl", attach: sysCall, fn: "ioctl"}}, Sets: []string{"syscalls", "fs", "fs_fd_ops"}},
	Pread64EventID:                {ID: Pread64EventID, ID32Bit: sys32pread64, Name: "pread64", Probes: []probe{{event: "pread64", attach: sysCall, fn: "pread64"}}, Sets: []string{"syscalls", "fs", "fs_read_write"}},
	Pwrite64EventID:               {ID: Pwrite64EventID, ID32Bit: sys32pwrite64, Name: "pwrite64", Probes: []probe{{event: "pwrite64", attach: sysCall, fn: "pwrite64"}}, Sets: []string{"syscalls", "fs", "fs_read_write"}},
	ReadvEventID:                  {ID: ReadvEventID, ID32Bit: sys32readv, Name: "readv", Probes: []probe{{event: "readv", attach: sysCall, fn: "readv"}}, Sets: []string{"syscalls", "fs", "fs_read_write"}},
	WritevEventID:                 {ID: WritevEventID, ID32Bit: sys32writev, Name: "writev", Probes: []probe{{event: "writev", attach: sysCall, fn: "writev"}}, Sets: []string{"syscalls", "fs", "fs_read_write"}},
	AccessEventID:                 {ID: AccessEventID, ID32Bit: sys32access, Name: "access", Probes: []probe{{event: "access", attach: sysCall, fn: "access"}}, Sets: []string{"default", "syscalls", "fs", "fs_file_attr"}},
	PipeEventID:                   {ID: PipeEventID, ID32Bit: sys32pipe, Name: "pipe", Probes: []probe{{event: "pipe", attach: sysCall, fn: "pipe"}}, Sets: []string{"syscalls", "ipc", "ipc_pipe"}},
	SelectEventID:                 {ID: SelectEventID, ID32Bit: sys32select, Name: "select", Probes: []probe{{event: "select", attach: sysCall, fn: "select"}}, Sets: []string{"syscalls", "fs", "fs_mux_io"}},
	SchedYieldEventID:             {ID: SchedYieldEventID, ID32Bit: sys32sched_yield, Name: "sched_yield", Probes: []probe{{event: "sched_yield", attach: sysCall, fn: "sched_yield"}}, Sets: []string{"syscalls", "proc", "proc_sched"}},
	MremapEventID:                 {ID: MremapEventID, ID32Bit: sys32mremap, Name: "mremap", Probes: []probe{{event: "mremap", attach: sysCall, fn: "mremap"}}, Sets: []string{"syscalls", "proc", "proc_mem"}},
	MsyncEventID:                  {ID: MsyncEventID, ID32Bit: sys32msync, Name: "msync", Probes: []probe{{event: "msync", attach: sysCall, fn: "msync"}}, Sets: []string{"syscalls", "fs", "fs_sync"}},
	MincoreEventID:                {ID: MincoreEventID, ID32Bit: sys32mincore, Name: "mincore", Probes: []probe{{event: "mincore", attach: sysCall, fn: "mincore"}}, Sets: []string{"syscalls", "proc", "proc_mem"}},
	MadviseEventID:                {ID: MadviseEventID, ID32Bit: sys32madvise, Name: "madvise", Probes: []probe{{event: "madvise", attach: sysCall, fn: "madvise"}}, Sets: []string{"syscalls", "proc", "proc_mem"}},
	ShmgetEventID:                 {ID: ShmgetEventID, ID32Bit: sys32shmget, Name: "shmget", Probes: []probe{{event: "shmget", attach: sysCall, fn: "shmget"}}, Sets: []string{"syscalls", "ipc", "ipc_shm"}},
	ShmatEventID:                  {ID: ShmatEventID, ID32Bit: sys32shmat, Name: "shmat", Probes: []probe{{event: "shmat", attach: sysCall, fn: "shmat"}}, Sets: []string{"syscalls", "ipc", "ipc_shm"}},
	ShmctlEventID:                 {ID: ShmctlEventID, ID32Bit: sys32shmctl, Name: "shmctl", Probes: []probe{{event: "shmctl", attach: sysCall, fn: "shmctl"}}, Sets: []string{"syscalls", "ipc", "ipc_shm"}},
	DupEventID:                    {ID: DupEventID, ID32Bit: sys32dup, Name: "dup", Probes: []probe{{event: "dup", attach: sysCall, fn: "dup"}}, Sets: []string{"default", "syscalls", "fs", "fs_fd_ops"}},
	Dup2EventID:                   {ID: Dup2EventID, ID32Bit: sys32dup2, Name: "dup2", Probes: []probe{{event: "dup2", attach: sysCall, fn: "dup2"}}, Sets: []string{"default", "syscalls", "fs", "fs_fd_ops"}},
	PauseEventID:                  {ID: PauseEventID, ID32Bit: sys32pause, Name: "pause", Probes: []probe{{event: "pause", attach: sysCall, fn: "pause"}}, Sets: []string{"syscalls", "signals"}},
	NanosleepEventID:              {ID: NanosleepEventID, ID32Bit: sys32nanosleep, Name: "nanosleep", Probes: []probe{{event: "nanosleep", attach: sysCall, fn: "nanosleep"}}, Sets: []string{"syscalls", "time", "time_timer"}},
	GetitimerEventID:              {ID: GetitimerEventID, ID32Bit: sys32getitimer, Name: "getitimer", Probes: []probe{{event: "getitimer", attach: sysCall, fn: "getitimer"}}, Sets: []string{"syscalls", "time", "time_timer"}},
	AlarmEventID:                  {ID: AlarmEventID, ID32Bit: sys32alarm, Name: "alarm", Probes: []probe{{event: "alarm", attach: sysCall, fn: "alarm"}}, Sets: []string{"syscalls", "time", "time_timer"}},
	SetitimerEventID:              {ID: SetitimerEventID, ID32Bit: sys32setitimer, Name: "setitimer", Probes: []probe{{event: "setitimer", attach: sysCall, fn: "setitimer"}}, Sets: []string{"syscalls", "time", "time_timer"}},
	GetpidEventID:                 {ID: GetpidEventID, ID32Bit: sys32getpid, Name: "getpid", Probes: []probe{{event: "getpid", attach: sysCall, fn: "getpid"}}, Sets: []string{"syscalls", "proc", "proc_ids"}},
	SendfileEventID:               {ID: SendfileEventID, ID32Bit: sys32sendfile, Name: "sendfile", Probes: []probe{{event: "sendfile", attach: sysCall, fn: "sendfile"}}, Sets: []string{"syscalls", "fs", "fs_read_write"}},
	SocketEventID:                 {ID: SocketEventID, ID32Bit: sys32socket, Name: "socket", Probes: []probe{{event: "socket", attach: sysCall, fn: "socket"}}, Sets: []string{"default", "syscalls", "net", "net_sock"}},
	ConnectEventID:                {ID: ConnectEventID, ID32Bit: sys32connect, Name: "connect", Probes: []probe{{event: "connect", attach: sysCall, fn: "connect"}}, Sets: []string{"default", "syscalls", "net", "net_sock"}},
	AcceptEventID:                 {ID: AcceptEventID, ID32Bit: sys32undefined, Name: "accept", Probes: []probe{{event: "accept", attach: sysCall, fn: "accept"}}, Sets: []string{"default", "syscalls", "net", "net_sock"}},
	SendtoEventID:                 {ID: SendtoEventID, ID32Bit: sys32sendto, Name: "sendto", Probes: []probe{{event: "sendto", attach: sysCall, fn: "sendto"}}, Sets: []string{"syscalls", "net", "net_snd_rcv"}},
	RecvfromEventID:               {ID: RecvfromEventID, ID32Bit: sys32recvfrom, Name: "recvfrom", Probes: []probe{{event: "recvfrom", attach: sysCall, fn: "recvfrom"}}, Sets: []string{"syscalls", "net", "net_snd_rcv"}},
	SendmsgEventID:                {ID: SendmsgEventID, ID32Bit: sys32sendmsg, Name: "sendmsg", Probes: []probe{{event: "sendmsg", attach: sysCall, fn: "sendmsg"}}, Sets: []string{"syscalls", "net", "net_snd_rcv"}},
	RecvmsgEventID:                {ID: RecvmsgEventID, ID32Bit: sys32recvmsg, Name: "recvmsg", Probes: []probe{{event: "recvmsg", attach: sysCall, fn: "recvmsg"}}, Sets: []string{"syscalls", "net", "net_snd_rcv"}},
	ShutdownEventID:               {ID: ShutdownEventID, ID32Bit: sys32shutdown, Name: "shutdown", Probes: []probe{{event: "shutdown", attach: sysCall, fn: "shutdown"}}, Sets: []string{"syscalls", "net", "net_sock"}},
	BindEventID:                   {ID: BindEventID, ID32Bit: sys32bind, Name: "bind", Probes: []probe{{event: "bind", attach: sysCall, fn: "bind"}}, Sets: []string{"default", "syscalls", "net", "net_sock"}},
	ListenEventID:                 {ID: ListenEventID, ID32Bit: sys32listen, Name: "listen", Probes: []probe{{event: "listen", attach: sysCall, fn: "listen"}}, Sets: []string{"default", "syscalls", "net", "net_sock"}},
	GetsocknameEventID:            {ID: GetsocknameEventID, ID32Bit: sys32getsockname, Name: "getsockname", Probes: []probe{{event: "getsockname", attach: sysCall, fn: "getsockname"}}, Sets: []string{"default", "syscalls", "net", "net_sock"}},
	GetpeernameEventID:            {ID: GetpeernameEventID, ID32Bit: sys32getpeername, Name: "getpeername", Probes: []probe{{event: "getpeername", attach: sysCall, fn: "getpeername"}}, Sets: []string{"syscalls", "net", "net_sock"}},
	SocketpairEventID:             {ID: SocketpairEventID, ID32Bit: sys32socketpair, Name: "socketpair", Probes: []probe{{event: "socketpair", attach: sysCall, fn: "socketpair"}}, Sets: []string{"syscalls", "net", "net_sock"}},
	SetsockoptEventID:             {ID: SetsockoptEventID, ID32Bit: sys32setsockopt, Name: "setsockopt", Probes: []probe{{event: "setsockopt", attach: sysCall, fn: "setsockopt"}}, Sets: []string{"syscalls", "net", "net_sock"}},
	GetsockoptEventID:             {ID: GetsockoptEventID, ID32Bit: sys32getsockopt, Name: "getsockopt", Probes: []probe{{event: "getsockopt", attach: sysCall, fn: "getsockopt"}}, Sets: []string{"syscalls", "net", "net_sock"}},
	CloneEventID:                  {ID: CloneEventID, ID32Bit: sys32clone, Name: "clone", Probes: []probe{{event: "clone", attach: sysCall, fn: "clone"}}, Sets: []string{"default", "syscalls", "proc", "proc_life"}},
	ForkEventID:                   {ID: ForkEventID, ID32Bit: sys32fork, Name: "fork", Probes: []probe{{event: "fork", attach: sysCall, fn: "fork"}}, Sets: []string{"default", "syscalls", "proc", "proc_life"}},
	VforkEventID:                  {ID: VforkEventID, ID32Bit: sys32vfork, Name: "vfork", Probes: []probe{{event: "vfork", attach: sysCall, fn: "vfork"}}, Sets: []string{"default", "syscalls", "proc", "proc_life"}},
	ExecveEventID:                 {ID: ExecveEventID, ID32Bit: sys32execve, Name: "execve", Probes: []probe{{event: "execve", attach: sysCall, fn: "execve"}}, Sets: []string{"default", "syscalls", "proc", "proc_life"}},
	ExitEventID:                   {ID: ExitEventID, ID32Bit: sys32exit, Name: "exit", Probes: []probe{{event: "exit", attach: sysCall, fn: "exit"}}, Sets: []string{"syscalls", "proc", "proc_life"}},
	Wait4EventID:                  {ID: Wait4EventID, ID32Bit: sys32wait4, Name: "wait4", Probes: []probe{{event: "wait4", attach: sysCall, fn: "wait4"}}, Sets: []string{"syscalls", "proc", "proc_life"}},
	KillEventID:                   {ID: KillEventID, ID32Bit: sys32kill, Name: "kill", Probes: []probe{{event: "kill", attach: sysCall, fn: "kill"}}, Sets: []string{"default", "syscalls", "signals"}},
	UnameEventID:                  {ID: UnameEventID, ID32Bit: sys32uname, Name: "uname", Probes: []probe{{event: "uname", attach: sysCall, fn: "uname"}}, Sets: []string{"syscalls", "system"}},
	SemgetEventID:                 {ID: SemgetEventID, ID32Bit: sys32semget, Name: "semget", Probes: []probe{{event: "semget", attach: sysCall, fn: "semget"}}, Sets: []string{"syscalls", "ipc", "ipc_sem"}},
	SemopEventID:                  {ID: SemopEventID, ID32Bit: sys32undefined, Name: "semop", Probes: []probe{{event: "semop", attach: sysCall, fn: "semop"}}, Sets: []string{"syscalls", "ipc", "ipc_sem"}},
	SemctlEventID:                 {ID: SemctlEventID, ID32Bit: sys32semctl, Name: "semctl", Probes: []probe{{event: "semctl", attach: sysCall, fn: "semctl"}}, Sets: []string{"syscalls", "ipc", "ipc_sem"}},
	ShmdtEventID:                  {ID: ShmdtEventID, ID32Bit: sys32shmdt, Name: "shmdt", Probes: []probe{{event: "shmdt", attach: sysCall, fn: "shmdt"}}, Sets: []string{"syscalls", "ipc", "ipc_shm"}},
	MsggetEventID:                 {ID: MsggetEventID, ID32Bit: sys32msgget, Name: "msgget", Probes: []probe{{event: "msgget", attach: sysCall, fn: "msgget"}}, Sets: []string{"syscalls", "ipc", "ipc_msgq"}},
	MsgsndEventID:                 {ID: MsgsndEventID, ID32Bit: sys32msgsnd, Name: "msgsnd", Probes: []probe{{event: "msgsnd", attach: sysCall, fn: "msgsnd"}}, Sets: []string{"syscalls", "ipc", "ipc_msgq"}},
	MsgrcvEventID:                 {ID: MsgrcvEventID, ID32Bit: sys32msgrcv, Name: "msgrcv", Probes: []probe{{event: "msgrcv", attach: sysCall, fn: "msgrcv"}}, Sets: []string{"syscalls", "ipc", "ipc_msgq"}},
	MsgctlEventID:                 {ID: MsgctlEventID, ID32Bit: sys32msgctl, Name: "msgctl", Probes: []probe{{event: "msgctl", attach: sysCall, fn: "msgctl"}}, Sets: []string{"syscalls", "ipc", "ipc_msgq"}},
	FcntlEventID:                  {ID: FcntlEventID, ID32Bit: sys32fcntl, Name: "fcntl", Probes: []probe{{event: "fcntl", attach: sysCall, fn: "fcntl"}}, Sets: []string{"syscalls", "fs", "fs_fd_ops"}},
	FlockEventID:                  {ID: FlockEventID, ID32Bit: sys32flock, Name: "flock", Probes: []probe{{event: "flock", attach: sysCall, fn: "flock"}}, Sets: []string{"syscalls", "fs", "fs_fd_ops"}},
	FsyncEventID:                  {ID: FsyncEventID, ID32Bit: sys32fsync, Name: "fsync", Probes: []probe{{event: "fsync", attach: sysCall, fn: "fsync"}}, Sets: []string{"syscalls", "fs", "fs_sync"}},
	FdatasyncEventID:              {ID: FdatasyncEventID, ID32Bit: sys32fdatasync, Name: "fdatasync", Probes: []probe{{event: "fdatasync", attach: sysCall, fn: "fdatasync"}}, Sets: []string{"syscalls", "fs", "fs_sync"}},
	TruncateEventID:               {ID: TruncateEventID, ID32Bit: sys32truncate, Name: "truncate", Probes: []probe{{event: "truncate", attach: sysCall, fn: "truncate"}}, Sets: []string{"syscalls", "fs", "fs_file_ops"}},
	FtruncateEventID:              {ID: FtruncateEventID, ID32Bit: sys32ftruncate, Name: "ftruncate", Probes: []probe{{event: "ftruncate", attach: sysCall, fn: "ftruncate"}}, Sets: []string{"syscalls", "fs", "fs_file_ops"}},
	GetdentsEventID:               {ID: GetdentsEventID, ID32Bit: sys32getdents, Name: "getdents", Probes: []probe{{event: "getdents", attach: sysCall, fn: "getdents"}}, Sets: []string{"default", "syscalls", "fs", "fs_dir_ops"}},
	GetcwdEventID:                 {ID: GetcwdEventID, ID32Bit: sys32getcwd, Name: "getcwd", Probes: []probe{{event: "getcwd", attach: sysCall, fn: "getcwd"}}, Sets: []string{"syscalls", "fs", "fs_dir_ops"}},
	ChdirEventID:                  {ID: ChdirEventID, ID32Bit: sys32chdir, Name: "chdir", Probes: []probe{{event: "chdir", attach: sysCall, fn: "chdir"}}, Sets: []string{"syscalls", "fs", "fs_dir_ops"}},
	FchdirEventID:                 {ID: FchdirEventID, ID32Bit: sys32fchdir, Name: "fchdir", Probes: []probe{{event: "fchdir", attach: sysCall, fn: "fchdir"}}, Sets: []string{"syscalls", "fs", "fs_dir_ops"}},
	RenameEventID:                 {ID: RenameEventID, ID32Bit: sys32rename, Name: "rename", Probes: []probe{{event: "rename", attach: sysCall, fn: "rename"}}, Sets: []string{"syscalls", "fs", "fs_file_ops"}},
	MkdirEventID:                  {ID: MkdirEventID, ID32Bit: sys32mkdir, Name: "mkdir", Probes: []probe{{event: "mkdir", attach: sysCall, fn: "mkdir"}}, Sets: []string{"syscalls", "fs", "fs_dir_ops"}},
	RmdirEventID:                  {ID: RmdirEventID, ID32Bit: sys32rmdir, Name: "rmdir", Probes: []probe{{event: "rmdir", attach: sysCall, fn: "rmdir"}}, Sets: []string{"syscalls", "fs", "fs_dir_ops"}},
	CreatEventID:                  {ID: CreatEventID, ID32Bit: sys32creat, Name: "creat", Probes: []probe{{event: "creat", attach: sysCall, fn: "creat"}}, Sets: []string{"default", "syscalls", "fs", "fs_file_ops"}},
	LinkEventID:                   {ID: LinkEventID, ID32Bit: sys32link, Name: "link", Probes: []probe{{event: "link", attach: sysCall, fn: "link"}}, Sets: []string{"syscalls", "fs", "fs_link_ops"}},
	UnlinkEventID:                 {ID: UnlinkEventID, ID32Bit: sys32unlink, Name: "unlink", Probes: []probe{{event: "unlink", attach: sysCall, fn: "unlink"}}, Sets: []string{"default", "syscalls", "fs", "fs_link_ops"}},
	SymlinkEventID:                {ID: SymlinkEventID, ID32Bit: sys32symlink, Name: "symlink", Probes: []probe{{event: "symlink", attach: sysCall, fn: "symlink"}}, Sets: []string{"default", "syscalls", "fs", "fs_link_ops"}},
	ReadlinkEventID:               {ID: ReadlinkEventID, ID32Bit: sys32readlink, Name: "readlink", Probes: []probe{{event: "readlink", attach: sysCall, fn: "readlink"}}, Sets: []string{"syscalls", "fs", "fs_link_ops"}},
	ChmodEventID:                  {ID: ChmodEventID, ID32Bit: sys32chmod, Name: "chmod", Probes: []probe{{event: "chmod", attach: sysCall, fn: "chmod"}}, Sets: []string{"default", "syscalls", "fs", "fs_file_attr"}},
	FchmodEventID:                 {ID: FchmodEventID, ID32Bit: sys32fchmod, Name: "fchmod", Probes: []probe{{event: "fchmod", attach: sysCall, fn: "fchmod"}}, Sets: []string{"default", "syscalls", "fs", "fs_file_attr"}},
	ChownEventID:                  {ID: ChownEventID, ID32Bit: sys32chown, Name: "chown", Probes: []probe{{event: "chown", attach: sysCall, fn: "chown"}}, Sets: []string{"default", "syscalls", "fs", "fs_file_attr"}},
	FchownEventID:                 {ID: FchownEventID, ID32Bit: sys32fchown, Name: "fchown", Probes: []probe{{event: "fchown", attach: sysCall, fn: "fchown"}}, Sets: []string{"default", "syscalls", "fs", "fs_file_attr"}},
	LchownEventID:                 {ID: LchownEventID, ID32Bit: sys32lchown, Name: "lchown", Probes: []probe{{event: "lchown", attach: sysCall, fn: "lchown"}}, Sets: []string{"default", "syscalls", "fs", "fs_file_attr"}},
	UmaskEventID:                  {ID: UmaskEventID, ID32Bit: sys32umask, Name: "umask", Probes: []probe{{event: "umask", attach: sysCall, fn: "umask"}}, Sets: []string{"syscalls", "fs", "fs_file_attr"}},
	GettimeofdayEventID:           {ID: GettimeofdayEventID, ID32Bit: sys32gettimeofday, Name: "gettimeofday", Probes: []probe{{event: "gettimeofday", attach: sysCall, fn: "gettimeofday"}}, Sets: []string{"syscalls", "time", "time_tod"}},
	GetrlimitEventID:              {ID: GetrlimitEventID, ID32Bit: sys32getrlimit, Name: "getrlimit", Probes: []probe{{event: "getrlimit", attach: sysCall, fn: "getrlimit"}}, Sets: []string{"syscalls", "proc"}},
	GetrusageEventID:              {ID: GetrusageEventID, ID32Bit: sys32getrusage, Name: "getrusage", Probes: []probe{{event: "getrusage", attach: sysCall, fn: "getrusage"}}, Sets: []string{"syscalls", "proc"}},
	SysinfoEventID:                {ID: SysinfoEventID, ID32Bit: sys32sysinfo, Name: "sysinfo", Probes: []probe{{event: "sysinfo", attach: sysCall, fn: "sysinfo"}}, Sets: []string{"syscalls", "system"}},
	TimesEventID:                  {ID: TimesEventID, ID32Bit: sys32times, Name: "times", Probes: []probe{{event: "times", attach: sysCall, fn: "times"}}, Sets: []string{"syscalls", "proc"}},
	PtraceEventID:                 {ID: PtraceEventID, ID32Bit: sys32ptrace, Name: "ptrace", Probes: []probe{{event: "ptrace", attach: sysCall, fn: "ptrace"}}, Sets: []string{"default", "syscalls", "proc"}},
	GetuidEventID:                 {ID: GetuidEventID, ID32Bit: sys32getuid, Name: "getuid", Probes: []probe{{event: "getuid", attach: sysCall, fn: "getuid"}}, Sets: []string{"syscalls", "proc", "proc_ids"}},
	SyslogEventID:                 {ID: SyslogEventID, ID32Bit: sys32syslog, Name: "syslog", Probes: []probe{{event: "syslog", attach: sysCall, fn: "syslog"}}, Sets: []string{"syscalls", "system"}},
	GetgidEventID:                 {ID: GetgidEventID, ID32Bit: sys32getgid, Name: "getgid", Probes: []probe{{event: "getgid", attach: sysCall, fn: "getgid"}}, Sets: []string{"syscalls", "proc", "proc_ids"}},
	SetuidEventID:                 {ID: SetuidEventID, ID32Bit: sys32setuid, Name: "setuid", Probes: []probe{{event: "setuid", attach: sysCall, fn: "setuid"}}, Sets: []string{"default", "syscalls", "proc", "proc_ids"}},
	SetgidEventID:                 {ID: SetgidEventID, ID32Bit: sys32setgid, Name: "setgid", Probes: []probe{{event: "setgid", attach: sysCall, fn: "setgid"}}, Sets: []string{"default", "syscalls", "proc", "proc_ids"}},
	GeteuidEventID:                {ID: GeteuidEventID, ID32Bit: sys32geteuid, Name: "geteuid", Probes: []probe{{event: "geteuid", attach: sysCall, fn: "geteuid"}}, Sets: []string{"syscalls", "proc", "proc_ids"}},
	GetegidEventID:                {ID: GetegidEventID, ID32Bit: sys32getegid, Name: "getegid", Probes: []probe{{event: "getegid", attach: sysCall, fn: "getegid"}}, Sets: []string{"syscalls", "proc", "proc_ids"}},
	SetpgidEventID:                {ID: SetpgidEventID, ID32Bit: sys32setpgid, Name: "setpgid", Probes: []probe{{event: "setpgid", attach: sysCall, fn: "setpgid"}}, Sets: []string{"syscalls", "proc", "proc_ids"}},
	GetppidEventID:                {ID: GetppidEventID, ID32Bit: sys32getppid, Name: "getppid", Probes: []probe{{event: "getppid", attach: sysCall, fn: "getppid"}}, Sets: []string{"syscalls", "proc", "proc_ids"}},
	GetpgrpEventID:                {ID: GetpgrpEventID, ID32Bit: sys32getpgrp, Name: "getpgrp", Probes: []probe{{event: "getpgrp", attach: sysCall, fn: "getpgrp"}}, Sets: []string{"syscalls", "proc", "proc_ids"}},
	SetsidEventID:                 {ID: SetsidEventID, ID32Bit: sys32setsid, Name: "setsid", Probes: []probe{{event: "setsid", attach: sysCall, fn: "setsid"}}, Sets: []string{"syscalls", "proc", "proc_ids"}},
	SetreuidEventID:               {ID: SetreuidEventID, ID32Bit: sys32setreuid, Name: "setreuid", Probes: []probe{{event: "setreuid", attach: sysCall, fn: "setreuid"}}, Sets: []string{"default", "syscalls", "proc", "proc_ids"}},
	SetregidEventID:               {ID: SetregidEventID, ID32Bit: sys32setregid, Name: "setregid", Probes: []probe{{event: "setregid", attach: sysCall, fn: "setregid"}}, Sets: []string{"default", "syscalls", "proc", "proc_ids"}},
	GetgroupsEventID:              {ID: GetgroupsEventID, ID32Bit: sys32getgroups, Name: "getgroups", Probes: []probe{{event: "getgroups", attach: sysCall, fn: "getgroups"}}, Sets: []string{"syscalls", "proc", "proc_ids"}},
	SetgroupsEventID:              {ID: SetgroupsEventID, ID32Bit: sys32setgroups, Name: "setgroups", Probes: []probe{{event: "setgroups", attach: sysCall, fn: "setgroups"}}, Sets: []string{"syscalls", "proc", "proc_ids"}},
	SetresuidEventID:              {ID: SetresuidEventID, ID32Bit: sys32setresuid, Name: "setresuid", Probes: []probe{{event: "setresuid", attach: sysCall, fn: "setresuid"}}, Sets: []string{"syscalls", "proc", "proc_ids"}},
	GetresuidEventID:              {ID: GetresuidEventID, ID32Bit: sys32getresuid, Name: "getresuid", Probes: []probe{{event: "getresuid", attach: sysCall, fn: "getresuid"}}, Sets: []string{"syscalls", "proc", "proc_ids"}},
	SetresgidEventID:              {ID: SetresgidEventID, ID32Bit: sys32setresgid, Name: "setresgid", Probes: []probe{{event: "setresgid", attach: sysCall, fn: "setresgid"}}, Sets: []string{"syscalls", "proc", "proc_ids"}},
	GetresgidEventID:              {ID: GetresgidEventID, ID32Bit: sys32getresgid, Name: "getresgid", Probes: []probe{{event: "getresgid", attach: sysCall, fn: "getresgid"}}, Sets: []string{"syscalls", "proc", "proc_ids"}},
	GetpgidEventID:                {ID: GetpgidEventID, ID32Bit: sys32getpgid, Name: "getpgid", Probes: []probe{{event: "getpgid", attach: sysCall, fn: "getpgid"}}, Sets: []string{"syscalls", "proc", "proc_ids"}},
	SetfsuidEventID:               {ID: SetfsuidEventID, ID32Bit: sys32setfsuid, Name: "setfsuid", Probes: []probe{{event: "setfsuid", attach: sysCall, fn: "setfsuid"}}, Sets: []string{"default", "syscalls", "proc", "proc_ids"}},
	SetfsgidEventID:               {ID: SetfsgidEventID, ID32Bit: sys32setfsgid, Name: "setfsgid", Probes: []probe{{event: "setfsgid", attach: sysCall, fn: "setfsgid"}}, Sets: []string{"default", "syscalls", "proc", "proc_ids"}},
	GetsidEventID:                 {ID: GetsidEventID, ID32Bit: sys32getsid, Name: "getsid", Probes: []probe{{event: "getsid", attach: sysCall, fn: "getsid"}}, Sets: []string{"syscalls", "proc", "proc_ids"}},
	CapgetEventID:                 {ID: CapgetEventID, ID32Bit: sys32capget, Name: "capget", Probes: []probe{{event: "capget", attach: sysCall, fn: "capget"}}, Sets: []string{"syscalls", "proc"}},
	CapsetEventID:                 {ID: CapsetEventID, ID32Bit: sys32capset, Name: "capset", Probes: []probe{{event: "capset", attach: sysCall, fn: "capset"}}, Sets: []string{"syscalls", "proc"}},
	RtSigpendingEventID:           {ID: RtSigpendingEventID, ID32Bit: sys32rt_sigpending, Name: "rt_sigpending", Probes: []probe{{event: "rt_sigpending", attach: sysCall, fn: "rt_sigpending"}}, Sets: []string{"syscalls", "signals"}},
	RtSigtimedwaitEventID:         {ID: RtSigtimedwaitEventID, ID32Bit: sys32rt_sigtimedwait, Name: "rt_sigtimedwait", Probes: []probe{{event: "rt_sigtimedwait", attach: sysCall, fn: "rt_sigtimedwait"}}, Sets: []string{"syscalls", "signals"}},
	RtSigqueueinfoEventID:         {ID: RtSigqueueinfoEventID, ID32Bit: sys32rt_sigqueueinfo, Name: "rt_sigqueueinfo", Probes: []probe{{event: "rt_sigqueueinfo", attach: sysCall, fn: "rt_sigqueueinfo"}}, Sets: []string{"syscalls", "signals"}},
	RtSigsuspendEventID:           {ID: RtSigsuspendEventID, ID32Bit: sys32rt_sigsuspend, Name: "rt_sigsuspend", Probes: []probe{{event: "rt_sigsuspend", attach: sysCall, fn: "rt_sigsuspend"}}, Sets: []string{"syscalls", "signals"}},
	SigaltstackEventID:            {ID: SigaltstackEventID, ID32Bit: sys32sigaltstack, Name: "sigaltstack", Probes: []probe{{event: "sigaltstack", attach: sysCall, fn: "sigaltstack"}}, Sets: []string{"syscalls", "signals"}},
	UtimeEventID:                  {ID: UtimeEventID, ID32Bit: sys32utime, Name: "utime", Probes: []probe{{event: "utime", attach: sysCall, fn: "utime"}}, Sets: []string{"syscalls", "fs", "fs_file_attr"}},
	MknodEventID:                  {ID: MknodEventID, ID32Bit: sys32mknod, Name: "mknod", Probes: []probe{{event: "mknod", attach: sysCall, fn: "mknod"}}, Sets: []string{"default", "syscalls", "fs", "fs_file_ops"}},
	UselibEventID:                 {ID: UselibEventID, ID32Bit: sys32uselib, Name: "uselib", Probes: []probe{{event: "uselib", attach: sysCall, fn: "uselib"}}, Sets: []string{"syscalls", "proc"}},
	PersonalityEventID:            {ID: PersonalityEventID, ID32Bit: sys32personality, Name: "personality", Probes: []probe{{event: "personality", attach: sysCall, fn: "personality"}}, Sets: []string{"syscalls", "system"}},
	UstatEventID:                  {ID: UstatEventID, ID32Bit: sys32ustat, Name: "ustat", Probes: []probe{{event: "ustat", attach: sysCall, fn: "ustat"}}, Sets: []string{"syscalls", "fs", "fs_info"}},
	StatfsEventID:                 {ID: StatfsEventID, ID32Bit: sys32statfs, Name: "statfs", Probes: []probe{{event: "statfs", attach: sysCall, fn: "statfs"}}, Sets: []string{"syscalls", "fs", "fs_info"}},
	FstatfsEventID:                {ID: FstatfsEventID, ID32Bit: sys32fstatfs, Name: "fstatfs", Probes: []probe{{event: "fstatfs", attach: sysCall, fn: "fstatfs"}}, Sets: []string{"syscalls", "fs", "fs_info"}},
	SysfsEventID:                  {ID: SysfsEventID, ID32Bit: sys32sysfs, Name: "sysfs", Probes: []probe{{event: "sysfs", attach: sysCall, fn: "sysfs"}}, Sets: []string{"syscalls", "fs", "fs_info"}},
	GetpriorityEventID:            {ID: GetpriorityEventID, ID32Bit: sys32getpriority, Name: "getpriority", Probes: []probe{{event: "getpriority", attach: sysCall, fn: "getpriority"}}, Sets: []string{"syscalls", "proc", "proc_sched"}},
	SetpriorityEventID:            {ID: SetpriorityEventID, ID32Bit: sys32setpriority, Name: "setpriority", Probes: []probe{{event: "setpriority", attach: sysCall, fn: "setpriority"}}, Sets: []string{"syscalls", "proc", "proc_sched"}},
	SchedSetparamEventID:          {ID: SchedSetparamEventID, ID32Bit: sys32sched_setparam, Name: "sched_setparam", Probes: []probe{{event: "sched_setparam", attach: sysCall, fn: "sched_setparam"}}, Sets: []string{"syscalls", "proc", "proc_sched"}},
	SchedGetparamEventID:          {ID: SchedGetparamEventID, ID32Bit: sys32sched_getparam, Name: "sched_getparam", Probes: []probe{{event: "sched_getparam", attach: sysCall, fn: "sched_getparam"}}, Sets: []string{"syscalls", "proc", "proc_sched"}},
	SchedSetschedulerEventID:      {ID: SchedSetschedulerEventID, ID32Bit: sys32sched_setscheduler, Name: "sched_setscheduler", Probes: []probe{{event: "sched_setscheduler", attach: sysCall, fn: "sched_setscheduler"}}, Sets: []string{"syscalls", "proc", "proc_sched"}},
	SchedGetschedulerEventID:      {ID: SchedGetschedulerEventID, ID32Bit: sys32sched_getscheduler, Name: "sched_getscheduler", Probes: []probe{{event: "sched_getscheduler", attach: sysCall, fn: "sched_getscheduler"}}, Sets: []string{"syscalls", "proc", "proc_sched"}},
	SchedGetPriorityMaxEventID:    {ID: SchedGetPriorityMaxEventID, ID32Bit: sys32sched_get_priority_max, Name: "sched_get_priority_max", Probes: []probe{{event: "sched_get_priority_max", attach: sysCall, fn: "sched_get_priority_max"}}, Sets: []string{"syscalls", "proc", "proc_sched"}},
	SchedGetPriorityMinEventID:    {ID: SchedGetPriorityMinEventID, ID32Bit: sys32sched_get_priority_min, Name: "sched_get_priority_min", Probes: []probe{{event: "sched_get_priority_min", attach: sysCall, fn: "sched_get_priority_min"}}, Sets: []string{"syscalls", "proc", "proc_sched"}},
	SchedRrGetIntervalEventID:     {ID: SchedRrGetIntervalEventID, ID32Bit: sys32sched_rr_get_interval, Name: "sched_rr_get_interval", Probes: []probe{{event: "sched_rr_get_interval", attach: sysCall, fn: "sched_rr_get_interval"}}, Sets: []string{"syscalls", "proc", "proc_sched"}},
	MlockEventID:                  {ID: MlockEventID, ID32Bit: sys32mlock, Name: "mlock", Probes: []probe{{event: "mlock", attach: sysCall, fn: "mlock"}}, Sets: []string{"syscalls", "proc", "proc_mem"}},
	MunlockEventID:                {ID: MunlockEventID, ID32Bit: sys32munlock, Name: "munlock", Probes: []probe{{event: "munlock", attach: sysCall, fn: "munlock"}}, Sets: []string{"syscalls", "proc", "proc_mem"}},
	MlockallEventID:               {ID: MlockallEventID, ID32Bit: sys32mlockall, Name: "mlockall", Probes: []probe{{event: "mlockall", attach: sysCall, fn: "mlockall"}}, Sets: []string{"syscalls", "proc", "proc_mem"}},
	MunlockallEventID:             {ID: MunlockallEventID, ID32Bit: sys32munlockall, Name: "munlockall", Probes: []probe{{event: "munlockall", attach: sysCall, fn: "munlockall"}}, Sets: []string{"syscalls", "proc", "proc_mem"}},
	VhangupEventID:                {ID: VhangupEventID, ID32Bit: sys32vhangup, Name: "vhangup", Probes: []probe{{event: "vhangup", attach: sysCall, fn: "vhangup"}}, Sets: []string{"syscalls", "system"}},
	ModifyLdtEventID:              {ID: ModifyLdtEventID, ID32Bit: sys32modify_ldt, Name: "modify_ldt", Probes: []probe{{event: "modify_ldt", attach: sysCall, fn: "modify_ldt"}}, Sets: []string{"syscalls", "proc", "proc_mem"}},
	PivotRootEventID:              {ID: PivotRootEventID, ID32Bit: sys32pivot_root, Name: "pivot_root", Probes: []probe{{event: "pivot_root", attach: sysCall, fn: "pivot_root"}}, Sets: []string{"syscalls", "fs"}},
	SysctlEventID:                 {ID: SysctlEventID, ID32Bit: sys32undefined, Name: "sysctl", Probes: []probe{{event: "sysctl", attach: sysCall, fn: "sysctl"}}, Sets: []string{"syscalls", "system"}},
	PrctlEventID:                  {ID: PrctlEventID, ID32Bit: sys32prctl, Name: "prctl", Probes: []probe{{event: "prctl", attach: sysCall, fn: "prctl"}}, Sets: []string{"default", "syscalls", "proc"}},
	ArchPrctlEventID:              {ID: ArchPrctlEventID, ID32Bit: sys32arch_prctl, Name: "arch_prctl", Probes: []probe{{event: "arch_prctl", attach: sysCall, fn: "arch_prctl"}}, Sets: []string{"syscalls", "proc"}},
	AdjtimexEventID:               {ID: AdjtimexEventID, ID32Bit: sys32adjtimex, Name: "adjtimex", Probes: []probe{{event: "adjtimex", attach: sysCall, fn: "adjtimex"}}, Sets: []string{"syscalls", "time", "time_clock"}},
	SetrlimitEventID:              {ID: SetrlimitEventID, ID32Bit: sys32setrlimit, Name: "setrlimit", Probes: []probe{{event: "setrlimit", attach: sysCall, fn: "setrlimit"}}, Sets: []string{"syscalls", "proc"}},
	ChrootEventID:                 {ID: ChrootEventID, ID32Bit: sys32chroot, Name: "chroot", Probes: []probe{{event: "chroot", attach: sysCall, fn: "chroot"}}, Sets: []string{"syscalls", "fs", "fs_dir_ops"}},
	SyncEventID:                   {ID: SyncEventID, ID32Bit: sys32sync, Name: "sync", Probes: []probe{{event: "sync", attach: sysCall, fn: "sync"}}, Sets: []string{"syscalls", "fs", "fs_sync"}},
	AcctEventID:                   {ID: AcctEventID, ID32Bit: sys32acct, Name: "acct", Probes: []probe{{event: "acct", attach: sysCall, fn: "acct"}}, Sets: []string{"syscalls", "system"}},
	SettimeofdayEventID:           {ID: SettimeofdayEventID, ID32Bit: sys32settimeofday, Name: "settimeofday", Probes: []probe{{event: "settimeofday", attach: sysCall, fn: "settimeofday"}}, Sets: []string{"syscalls", "time", "time_tod"}},
	MountEventID:                  {ID: MountEventID, ID32Bit: sys32mount, Name: "mount", Probes: []probe{{event: "mount", attach: sysCall, fn: "mount"}}, Sets: []string{"default", "syscalls", "fs"}},
	UmountEventID:                 {ID: UmountEventID, ID32Bit: sys32umount, Name: "umount", Probes: []probe{{event: "umount", attach: sysCall, fn: "umount"}}, Sets: []string{"default", "syscalls", "fs"}},
	SwaponEventID:                 {ID: SwaponEventID, ID32Bit: sys32swapon, Name: "swapon", Probes: []probe{{event: "swapon", attach: sysCall, fn: "swapon"}}, Sets: []string{"syscalls", "fs"}},
	SwapoffEventID:                {ID: SwapoffEventID, ID32Bit: sys32swapoff, Name: "swapoff", Probes: []probe{{event: "swapoff", attach: sysCall, fn: "swapoff"}}, Sets: []string{"syscalls", "fs"}},
	RebootEventID:                 {ID: RebootEventID, ID32Bit: sys32reboot, Name: "reboot", Probes: []probe{{event: "reboot", attach: sysCall, fn: "reboot"}}, Sets: []string{"syscalls", "system"}},
	SethostnameEventID:            {ID: SethostnameEventID, ID32Bit: sys32sethostname, Name: "sethostname", Probes: []probe{{event: "sethostname", attach: sysCall, fn: "sethostname"}}, Sets: []string{"syscalls", "net"}},
	SetdomainnameEventID:          {ID: SetdomainnameEventID, ID32Bit: sys32setdomainname, Name: "setdomainname", Probes: []probe{{event: "setdomainname", attach: sysCall, fn: "setdomainname"}}, Sets: []string{"syscalls", "net"}},
	IoplEventID:                   {ID: IoplEventID, ID32Bit: sys32iopl, Name: "iopl", Probes: []probe{{event: "iopl", attach: sysCall, fn: "iopl"}}, Sets: []string{"syscalls", "system"}},
	IopermEventID:                 {ID: IopermEventID, ID32Bit: sys32ioperm, Name: "ioperm", Probes: []probe{{event: "ioperm", attach: sysCall, fn: "ioperm"}}, Sets: []string{"syscalls", "system"}},
	CreateModuleEventID:           {ID: CreateModuleEventID, ID32Bit: sys32create_module, Name: "create_module", Probes: []probe{{event: "create_module", attach: sysCall, fn: "create_module"}}, Sets: []string{"syscalls", "system", "system_module"}},
	InitModuleEventID:             {ID: InitModuleEventID, ID32Bit: sys32init_module, Name: "init_module", Probes: []probe{{event: "init_module", attach: sysCall, fn: "init_module"}}, Sets: []string{"default", "syscalls", "system", "system_module"}},
	DeleteModuleEventID:           {ID: DeleteModuleEventID, ID32Bit: sys32delete_module, Name: "delete_module", Probes: []probe{{event: "delete_module", attach: sysCall, fn: "delete_module"}}, Sets: []string{"default", "syscalls", "system", "system_module"}},
	GetKernelSymsEventID:          {ID: GetKernelSymsEventID, ID32Bit: sys32get_kernel_syms, Name: "get_kernel_syms", Probes: []probe{{event: "get_kernel_syms", attach: sysCall, fn: "get_kernel_syms"}}, Sets: []string{"syscalls", "system", "system_module"}},
	QueryModuleEventID:            {ID: QueryModuleEventID, ID32Bit: sys32query_module, Name: "query_module", Probes: []probe{{event: "query_module", attach: sysCall, fn: "query_module"}}, Sets: []string{"syscalls", "system", "system_module"}},
	QuotactlEventID:               {ID: QuotactlEventID, ID32Bit: sys32quotactl, Name: "quotactl", Probes: []probe{{event: "quotactl", attach: sysCall, fn: "quotactl"}}, Sets: []string{"syscalls", "system"}},
	NfsservctlEventID:             {ID: NfsservctlEventID, ID32Bit: sys32nfsservctl, Name: "nfsservctl", Probes: []probe{{event: "nfsservctl", attach: sysCall, fn: "nfsservctl"}}, Sets: []string{"syscalls", "fs"}},
	GetpmsgEventID:                {ID: GetpmsgEventID, ID32Bit: sys32getpmsg, Name: "getpmsg", Probes: []probe{{event: "getpmsg", attach: sysCall, fn: "getpmsg"}}, Sets: []string{"syscalls"}},
	PutpmsgEventID:                {ID: PutpmsgEventID, ID32Bit: sys32putpmsg, Name: "putpmsg", Probes: []probe{{event: "putpmsg", attach: sysCall, fn: "putpmsg"}}, Sets: []string{"syscalls"}},
	AfsEventID:                    {ID: AfsEventID, ID32Bit: sys32undefined, Name: "afs", Probes: []probe{{event: "afs", attach: sysCall, fn: "afs"}}, Sets: []string{"syscalls"}},
	TuxcallEventID:                {ID: TuxcallEventID, ID32Bit: sys32undefined, Name: "tuxcall", Probes: []probe{{event: "tuxcall", attach: sysCall, fn: "tuxcall"}}, Sets: []string{"syscalls"}},
	SecurityEventID:               {ID: SecurityEventID, ID32Bit: sys32undefined, Name: "security", Probes: []probe{{event: "security", attach: sysCall, fn: "security"}}, Sets: []string{"syscalls"}},
	GettidEventID:                 {ID: GettidEventID, ID32Bit: sys32gettid, Name: "gettid", Probes: []probe{{event: "gettid", attach: sysCall, fn: "gettid"}}, Sets: []string{"syscalls", "proc", "proc_ids"}},
	ReadaheadEventID:              {ID: ReadaheadEventID, ID32Bit: sys32readahead, Name: "readahead", Probes: []probe{{event: "readahead", attach: sysCall, fn: "readahead"}}, Sets: []string{"syscalls", "fs"}},
	SetxattrEventID:               {ID: SetxattrEventID, ID32Bit: sys32setxattr, Name: "setxattr", Probes: []probe{{event: "setxattr", attach: sysCall, fn: "setxattr"}}, Sets: []string{"syscalls", "fs", "fs_file_attr"}},
	LsetxattrEventID:              {ID: LsetxattrEventID, ID32Bit: sys32lsetxattr, Name: "lsetxattr", Probes: []probe{{event: "lsetxattr", attach: sysCall, fn: "lsetxattr"}}, Sets: []string{"syscalls", "fs", "fs_file_attr"}},
	FsetxattrEventID:              {ID: FsetxattrEventID, ID32Bit: sys32fsetxattr, Name: "fsetxattr", Probes: []probe{{event: "fsetxattr", attach: sysCall, fn: "fsetxattr"}}, Sets: []string{"syscalls", "fs", "fs_file_attr"}},
	GetxattrEventID:               {ID: GetxattrEventID, ID32Bit: sys32getxattr, Name: "getxattr", Probes: []probe{{event: "getxattr", attach: sysCall, fn: "getxattr"}}, Sets: []string{"syscalls", "fs", "fs_file_attr"}},
	LgetxattrEventID:              {ID: LgetxattrEventID, ID32Bit: sys32lgetxattr, Name: "lgetxattr", Probes: []probe{{event: "lgetxattr", attach: sysCall, fn: "lgetxattr"}}, Sets: []string{"syscalls", "fs", "fs_file_attr"}},
	FgetxattrEventID:              {ID: FgetxattrEventID, ID32Bit: sys32fgetxattr, Name: "fgetxattr", Probes: []probe{{event: "fgetxattr", attach: sysCall, fn: "fgetxattr"}}, Sets: []string{"syscalls", "fs", "fs_file_attr"}},
	ListxattrEventID:              {ID: ListxattrEventID, ID32Bit: sys32listxattr, Name: "listxattr", Probes: []probe{{event: "listxattr", attach: sysCall, fn: "listxattr"}}, Sets: []string{"syscalls", "fs", "fs_file_attr"}},
	LlistxattrEventID:             {ID: LlistxattrEventID, ID32Bit: sys32llistxattr, Name: "llistxattr", Probes: []probe{{event: "llistxattr", attach: sysCall, fn: "llistxattr"}}, Sets: []string{"syscalls", "fs", "fs_file_attr"}},
	FlistxattrEventID:             {ID: FlistxattrEventID, ID32Bit: sys32flistxattr, Name: "flistxattr", Probes: []probe{{event: "flistxattr", attach: sysCall, fn: "flistxattr"}}, Sets: []string{"syscalls", "fs", "fs_file_attr"}},
	RemovexattrEventID:            {ID: RemovexattrEventID, ID32Bit: sys32removexattr, Name: "removexattr", Probes: []probe{{event: "removexattr", attach: sysCall, fn: "removexattr"}}, Sets: []string{"syscalls", "fs", "fs_file_attr"}},
	LremovexattrEventID:           {ID: LremovexattrEventID, ID32Bit: sys32lremovexattr, Name: "lremovexattr", Probes: []probe{{event: "lremovexattr", attach: sysCall, fn: "lremovexattr"}}, Sets: []string{"syscalls", "fs", "fs_file_attr"}},
	FremovexattrEventID:           {ID: FremovexattrEventID, ID32Bit: sys32fremovexattr, Name: "fremovexattr", Probes: []probe{{event: "fremovexattr", attach: sysCall, fn: "fremovexattr"}}, Sets: []string{"syscalls", "fs", "fs_file_attr"}},
	TkillEventID:                  {ID: TkillEventID, ID32Bit: sys32tkill, Name: "tkill", Probes: []probe{{event: "tkill", attach: sysCall, fn: "tkill"}}, Sets: []string{"syscalls", "signals"}},
	TimeEventID:                   {ID: TimeEventID, ID32Bit: sys32time, Name: "time", Probes: []probe{{event: "time", attach: sysCall, fn: "time"}}, Sets: []string{"syscalls", "time", "time_tod"}},
	FutexEventID:                  {ID: FutexEventID, ID32Bit: sys32futex, Name: "futex", Probes: []probe{{event: "futex", attach: sysCall, fn: "futex"}}, Sets: []string{"syscalls", "ipc", "ipc_futex"}},
	SchedSetaffinityEventID:       {ID: SchedSetaffinityEventID, ID32Bit: sys32sched_setaffinity, Name: "sched_setaffinity", Probes: []probe{{event: "sched_setaffinity", attach: sysCall, fn: "sched_setaffinity"}}, Sets: []string{"syscalls", "proc", "proc_sched"}},
	SchedGetaffinityEventID:       {ID: SchedGetaffinityEventID, ID32Bit: sys32sched_getaffinity, Name: "sched_getaffinity", Probes: []probe{{event: "sched_getaffinity", attach: sysCall, fn: "sched_getaffinity"}}, Sets: []string{"syscalls", "proc", "proc_sched"}},
	SetThreadAreaEventID:          {ID: SetThreadAreaEventID, ID32Bit: sys32set_thread_area, Name: "set_thread_area", Probes: []probe{{event: "set_thread_area", attach: sysCall, fn: "set_thread_area"}}, Sets: []string{"syscalls", "proc"}},
	IoSetupEventID:                {ID: IoSetupEventID, ID32Bit: sys32io_setup, Name: "io_setup", Probes: []probe{{event: "io_setup", attach: sysCall, fn: "io_setup"}}, Sets: []string{"syscalls", "fs", "fs_async_io"}},
	IoDestroyEventID:              {ID: IoDestroyEventID, ID32Bit: sys32io_destroy, Name: "io_destroy", Probes: []probe{{event: "io_destroy", attach: sysCall, fn: "io_destroy"}}, Sets: []string{"syscalls", "fs", "fs_async_io"}},
	IoGeteventsEventID:            {ID: IoGeteventsEventID, ID32Bit: sys32io_getevents, Name: "io_getevents", Probes: []probe{{event: "io_getevents", attach: sysCall, fn: "io_getevents"}}, Sets: []string{"syscalls", "fs", "fs_async_io"}},
	IoSubmitEventID:               {ID: IoSubmitEventID, ID32Bit: sys32io_submit, Name: "io_submit", Probes: []probe{{event: "io_submit", attach: sysCall, fn: "io_submit"}}, Sets: []string{"syscalls", "fs", "fs_async_io"}},
	IoCancelEventID:               {ID: IoCancelEventID, ID32Bit: sys32io_cancel, Name: "io_cancel", Probes: []probe{{event: "io_cancel", attach: sysCall, fn: "io_cancel"}}, Sets: []string{"syscalls", "fs", "fs_async_io"}},
	GetThreadAreaEventID:          {ID: GetThreadAreaEventID, ID32Bit: sys32get_thread_area, Name: "get_thread_area", Probes: []probe{{event: "get_thread_area", attach: sysCall, fn: "get_thread_area"}}, Sets: []string{"syscalls", "proc"}},
	LookupDcookieEventID:          {ID: LookupDcookieEventID, ID32Bit: sys32lookup_dcookie, Name: "lookup_dcookie", Probes: []probe{{event: "lookup_dcookie", attach: sysCall, fn: "lookup_dcookie"}}, Sets: []string{"syscalls", "fs", "fs_dir_ops"}},
	EpollCreateEventID:            {ID: EpollCreateEventID, ID32Bit: sys32epoll_create, Name: "epoll_create", Probes: []probe{{event: "epoll_create", attach: sysCall, fn: "epoll_create"}}, Sets: []string{"syscalls", "fs", "fs_mux_io"}},
	EpollCtlOldEventID:            {ID: EpollCtlOldEventID, ID32Bit: sys32undefined, Name: "epoll_ctl_old", Probes: []probe{{event: "epoll_ctl_old", attach: sysCall, fn: "epoll_ctl_old"}}, Sets: []string{"syscalls", "fs", "fs_mux_io"}},
	EpollWaitOldEventID:           {ID: EpollWaitOldEventID, ID32Bit: sys32undefined, Name: "epoll_wait_old", Probes: []probe{{event: "epoll_wait_old", attach: sysCall, fn: "epoll_wait_old"}}, Sets: []string{"syscalls", "fs", "fs_mux_io"}},
	RemapFilePagesEventID:         {ID: RemapFilePagesEventID, ID32Bit: sys32remap_file_pages, Name: "remap_file_pages", Probes: []probe{{event: "remap_file_pages", attach: sysCall, fn: "remap_file_pages"}}, Sets: []string{"syscalls"}},
	Getdents64EventID:             {ID: Getdents64EventID, ID32Bit: sys32getdents64, Name: "getdents64", Probes: []probe{{event: "getdents64", attach: sysCall, fn: "getdents64"}}, Sets: []string{"default", "syscalls", "fs", "fs_dir_ops"}},
	SetTidAddressEventID:          {ID: SetTidAddressEventID, ID32Bit: sys32set_tid_address, Name: "set_tid_address", Probes: []probe{{event: "set_tid_address", attach: sysCall, fn: "set_tid_address"}}, Sets: []string{"syscalls", "proc"}},
	RestartSyscallEventID:         {ID: RestartSyscallEventID, ID32Bit: sys32restart_syscall, Name: "restart_syscall", Probes: []probe{{event: "restart_syscall", attach: sysCall, fn: "restart_syscall"}}, Sets: []string{"syscalls", "signals"}},
	SemtimedopEventID:             {ID: SemtimedopEventID, ID32Bit: sys32semtimedop_time64, Name: "semtimedop", Probes: []probe{{event: "semtimedop", attach: sysCall, fn: "semtimedop"}}, Sets: []string{"syscalls", "ipc", "ipc_sem"}},
	Fadvise64EventID:              {ID: Fadvise64EventID, ID32Bit: sys32fadvise64, Name: "fadvise64", Probes: []probe{{event: "fadvise64", attach: sysCall, fn: "fadvise64"}}, Sets: []string{"syscalls", "fs"}},
	TimerCreateEventID:            {ID: TimerCreateEventID, ID32Bit: sys32timer_create, Name: "timer_create", Probes: []probe{{event: "timer_create", attach: sysCall, fn: "timer_create"}}, Sets: []string{"syscalls", "time", "time_timer"}},
	TimerSettimeEventID:           {ID: TimerSettimeEventID, ID32Bit: sys32timer_settime, Name: "timer_settime", Probes: []probe{{event: "timer_settime", attach: sysCall, fn: "timer_settime"}}, Sets: []string{"syscalls", "time", "time_timer"}},
	TimerGettimeEventID:           {ID: TimerGettimeEventID, ID32Bit: sys32timer_gettime, Name: "timer_gettime", Probes: []probe{{event: "timer_gettime", attach: sysCall, fn: "timer_gettime"}}, Sets: []string{"syscalls", "time", "time_timer"}},
	TimerGetoverrunEventID:        {ID: TimerGetoverrunEventID, ID32Bit: sys32timer_getoverrun, Name: "timer_getoverrun", Probes: []probe{{event: "timer_getoverrun", attach: sysCall, fn: "timer_getoverrun"}}, Sets: []string{"syscalls", "time", "time_timer"}},
	TimerDeleteEventID:            {ID: TimerDeleteEventID, ID32Bit: sys32timer_delete, Name: "timer_delete", Probes: []probe{{event: "timer_delete", attach: sysCall, fn: "timer_delete"}}, Sets: []string{"syscalls", "time", "time_timer"}},
	ClockSettimeEventID:           {ID: ClockSettimeEventID, ID32Bit: sys32clock_settime, Name: "clock_settime", Probes: []probe{{event: "clock_settime", attach: sysCall, fn: "clock_settime"}}, Sets: []string{"syscalls", "time", "time_clock"}},
	ClockGettimeEventID:           {ID: ClockGettimeEventID, ID32Bit: sys32clock_gettime, Name: "clock_gettime", Probes: []probe{{event: "clock_gettime", attach: sysCall, fn: "clock_gettime"}}, Sets: []string{"syscalls", "time", "time_clock"}},
	ClockGetresEventID:            {ID: ClockGetresEventID, ID32Bit: sys32clock_getres, Name: "clock_getres", Probes: []probe{{event: "clock_getres", attach: sysCall, fn: "clock_getres"}}, Sets: []string{"syscalls", "time", "time_clock"}},
	ClockNanosleepEventID:         {ID: ClockNanosleepEventID, ID32Bit: sys32clock_nanosleep, Name: "clock_nanosleep", Probes: []probe{{event: "clock_nanosleep", attach: sysCall, fn: "clock_nanosleep"}}, Sets: []string{"syscalls", "time", "time_clock"}},
	ExitGroupEventID:              {ID: ExitGroupEventID, ID32Bit: sys32exit_group, Name: "exit_group", Probes: []probe{{event: "exit_group", attach: sysCall, fn: "exit_group"}}, Sets: []string{"syscalls", "proc", "proc_life"}},
	EpollWaitEventID:              {ID: EpollWaitEventID, ID32Bit: sys32epoll_wait, Name: "epoll_wait", Probes: []probe{{event: "epoll_wait", attach: sysCall, fn: "epoll_wait"}}, Sets: []string{"syscalls", "fs", "fs_mux_io"}},
	EpollCtlEventID:               {ID: EpollCtlEventID, ID32Bit: sys32epoll_ctl, Name: "epoll_ctl", Probes: []probe{{event: "epoll_ctl", attach: sysCall, fn: "epoll_ctl"}}, Sets: []string{"syscalls", "fs", "fs_mux_io"}},
	TgkillEventID:                 {ID: TgkillEventID, ID32Bit: sys32tgkill, Name: "tgkill", Probes: []probe{{event: "tgkill", attach: sysCall, fn: "tgkill"}}, Sets: []string{"syscalls", "signals"}},
	UtimesEventID:                 {ID: UtimesEventID, ID32Bit: sys32utimes, Name: "utimes", Probes: []probe{{event: "utimes", attach: sysCall, fn: "utimes"}}, Sets: []string{"syscalls", "fs", "fs_file_attr"}},
	VserverEventID:                {ID: VserverEventID, ID32Bit: sys32vserver, Name: "vserver", Probes: []probe{{event: "vserver", attach: sysCall, fn: "vserver"}}, Sets: []string{"syscalls"}},
	MbindEventID:                  {ID: MbindEventID, ID32Bit: sys32mbind, Name: "mbind", Probes: []probe{{event: "mbind", attach: sysCall, fn: "mbind"}}, Sets: []string{"syscalls", "system", "system_numa"}},
	SetMempolicyEventID:           {ID: SetMempolicyEventID, ID32Bit: sys32set_mempolicy, Name: "set_mempolicy", Probes: []probe{{event: "set_mempolicy", attach: sysCall, fn: "set_mempolicy"}}, Sets: []string{"syscalls", "system", "system_numa"}},
	GetMempolicyEventID:           {ID: GetMempolicyEventID, ID32Bit: sys32get_mempolicy, Name: "get_mempolicy", Probes: []probe{{event: "get_mempolicy", attach: sysCall, fn: "get_mempolicy"}}, Sets: []string{"syscalls", "system", "system_numa"}},
	MqOpenEventID:                 {ID: MqOpenEventID, ID32Bit: sys32mq_open, Name: "mq_open", Probes: []probe{{event: "mq_open", attach: sysCall, fn: "mq_open"}}, Sets: []string{"syscalls", "ipc", "ipc_msgq"}},
	MqUnlinkEventID:               {ID: MqUnlinkEventID, ID32Bit: sys32mq_unlink, Name: "mq_unlink", Probes: []probe{{event: "mq_unlink", attach: sysCall, fn: "mq_unlink"}}, Sets: []string{"syscalls", "ipc", "ipc_msgq"}},
	MqTimedsendEventID:            {ID: MqTimedsendEventID, ID32Bit: sys32mq_timedsend, Name: "mq_timedsend", Probes: []probe{{event: "mq_timedsend", attach: sysCall, fn: "mq_timedsend"}}, Sets: []string{"syscalls", "ipc", "ipc_msgq"}},
	MqTimedreceiveEventID:         {ID: MqTimedreceiveEventID, ID32Bit: sys32mq_timedreceive, Name: "mq_timedreceive", Probes: []probe{{event: "mq_timedreceive", attach: sysCall, fn: "mq_timedreceive"}}, Sets: []string{"syscalls", "ipc", "ipc_msgq"}},
	MqNotifyEventID:               {ID: MqNotifyEventID, ID32Bit: sys32mq_notify, Name: "mq_notify", Probes: []probe{{event: "mq_notify", attach: sysCall, fn: "mq_notify"}}, Sets: []string{"syscalls", "ipc", "ipc_msgq"}},
	MqGetsetattrEventID:           {ID: MqGetsetattrEventID, ID32Bit: sys32mq_getsetattr, Name: "mq_getsetattr", Probes: []probe{{event: "mq_getsetattr", attach: sysCall, fn: "mq_getsetattr"}}, Sets: []string{"syscalls", "ipc", "ipc_msgq"}},
	KexecLoadEventID:              {ID: KexecLoadEventID, ID32Bit: sys32kexec_load, Name: "kexec_load", Probes: []probe{{event: "kexec_load", attach: sysCall, fn: "kexec_load"}}, Sets: []string{"syscalls", "system"}},
	WaitidEventID:                 {ID: WaitidEventID, ID32Bit: sys32waitid, Name: "waitid", Probes: []probe{{event: "waitid", attach: sysCall, fn: "waitid"}}, Sets: []string{"syscalls", "proc", "proc_life"}},
	AddKeyEventID:                 {ID: AddKeyEventID, ID32Bit: sys32add_key, Name: "add_key", Probes: []probe{{event: "add_key", attach: sysCall, fn: "add_key"}}, Sets: []string{"syscalls", "system", "system_keys"}},
	RequestKeyEventID:             {ID: RequestKeyEventID, ID32Bit: sys32request_key, Name: "request_key", Probes: []probe{{event: "request_key", attach: sysCall, fn: "request_key"}}, Sets: []string{"syscalls", "system", "system_keys"}},
	KeyctlEventID:                 {ID: KeyctlEventID, ID32Bit: sys32keyctl, Name: "keyctl", Probes: []probe{{event: "keyctl", attach: sysCall, fn: "keyctl"}}, Sets: []string{"syscalls", "system", "system_keys"}},
	IoprioSetEventID:              {ID: IoprioSetEventID, ID32Bit: sys32ioprio_set, Name: "ioprio_set", Probes: []probe{{event: "ioprio_set", attach: sysCall, fn: "ioprio_set"}}, Sets: []string{"syscalls", "proc", "proc_sched"}},
	IoprioGetEventID:              {ID: IoprioGetEventID, ID32Bit: sys32ioprio_get, Name: "ioprio_get", Probes: []probe{{event: "ioprio_get", attach: sysCall, fn: "ioprio_get"}}, Sets: []string{"syscalls", "proc", "proc_sched"}},
	InotifyInitEventID:            {ID: InotifyInitEventID, ID32Bit: sys32inotify_init, Name: "inotify_init", Probes: []probe{{event: "inotify_init", attach: sysCall, fn: "inotify_init"}}, Sets: []string{"syscalls", "fs", "fs_monitor"}},
	InotifyAddWatchEventID:        {ID: InotifyAddWatchEventID, ID32Bit: sys32inotify_add_watch, Name: "inotify_add_watch", Probes: []probe{{event: "inotify_add_watch", attach: sysCall, fn: "inotify_add_watch"}}, Sets: []string{"syscalls", "fs", "fs_monitor"}},
	InotifyRmWatchEventID:         {ID: InotifyRmWatchEventID, ID32Bit: sys32inotify_rm_watch, Name: "inotify_rm_watch", Probes: []probe{{event: "inotify_rm_watch", attach: sysCall, fn: "inotify_rm_watch"}}, Sets: []string{"syscalls", "fs", "fs_monitor"}},
	MigratePagesEventID:           {ID: MigratePagesEventID, ID32Bit: sys32migrate_pages, Name: "migrate_pages", Probes: []probe{{event: "migrate_pages", attach: sysCall, fn: "migrate_pages"}}, Sets: []string{"syscalls", "system", "system_numa"}},
	OpenatEventID:                 {ID: OpenatEventID, ID32Bit: sys32openat, Name: "openat", Probes: []probe{{event: "openat", attach: sysCall, fn: "openat"}}, Sets: []string{"default", "syscalls", "fs", "fs_file_ops"}},
	MkdiratEventID:                {ID: MkdiratEventID, ID32Bit: sys32mkdirat, Name: "mkdirat", Probes: []probe{{event: "mkdirat", attach: sysCall, fn: "mkdirat"}}, Sets: []string{"syscalls", "fs", "fs_dir_ops"}},
	MknodatEventID:                {ID: MknodatEventID, ID32Bit: sys32mknodat, Name: "mknodat", Probes: []probe{{event: "mknodat", attach: sysCall, fn: "mknodat"}}, Sets: []string{"default", "syscalls", "fs", "fs_file_ops"}},
	FchownatEventID:               {ID: FchownatEventID, ID32Bit: sys32fchownat, Name: "fchownat", Probes: []probe{{event: "fchownat", attach: sysCall, fn: "fchownat"}}, Sets: []string{"default", "syscalls", "fs", "fs_file_attr"}},
	FutimesatEventID:              {ID: FutimesatEventID, ID32Bit: sys32futimesat, Name: "futimesat", Probes: []probe{{event: "futimesat", attach: sysCall, fn: "futimesat"}}, Sets: []string{"syscalls", "fs", "fs_file_attr"}},
	NewfstatatEventID:             {ID: NewfstatatEventID, ID32Bit: sys32fstatat64, Name: "newfstatat", Probes: []probe{{event: "newfstatat", attach: sysCall, fn: "newfstatat"}}, Sets: []string{"syscalls", "fs", "fs_file_attr"}},
	UnlinkatEventID:               {ID: UnlinkatEventID, ID32Bit: sys32unlinkat, Name: "unlinkat", Probes: []probe{{event: "unlinkat", attach: sysCall, fn: "unlinkat"}}, Sets: []string{"default", "syscalls", "fs", "fs_link_ops"}},
	RenameatEventID:               {ID: RenameatEventID, ID32Bit: sys32renameat, Name: "renameat", Probes: []probe{{event: "renameat", attach: sysCall, fn: "renameat"}}, Sets: []string{"syscalls", "fs", "fs_file_ops"}},
	LinkatEventID:                 {ID: LinkatEventID, ID32Bit: sys32linkat, Name: "linkat", Probes: []probe{{event: "linkat", attach: sysCall, fn: "linkat"}}, Sets: []string{"syscalls", "fs", "fs_link_ops"}},
	SymlinkatEventID:              {ID: SymlinkatEventID, ID32Bit: sys32symlinkat, Name: "symlinkat", Probes: []probe{{event: "symlinkat", attach: sysCall, fn: "symlinkat"}}, Sets: []string{"default", "syscalls", "fs", "fs_link_ops"}},
	ReadlinkatEventID:             {ID: ReadlinkatEventID, ID32Bit: sys32readlinkat, Name: "readlinkat", Probes: []probe{{event: "readlinkat", attach: sysCall, fn: "readlinkat"}}, Sets: []string{"syscalls", "fs", "fs_link_ops"}},
	FchmodatEventID:               {ID: FchmodatEventID, ID32Bit: sys32fchmodat, Name: "fchmodat", Probes: []probe{{event: "fchmodat", attach: sysCall, fn: "fchmodat"}}, Sets: []string{"default", "syscalls", "fs", "fs_file_attr"}},
	FaccessatEventID:              {ID: FaccessatEventID, ID32Bit: sys32faccessat, Name: "faccessat", Probes: []probe{{event: "faccessat", attach: sysCall, fn: "faccessat"}}, Sets: []string{"default", "syscalls", "fs", "fs_file_attr"}},
	Pselect6EventID:               {ID: Pselect6EventID, ID32Bit: sys32pselect6, Name: "pselect6", Probes: []probe{{event: "pselect6", attach: sysCall, fn: "pselect6"}}, Sets: []string{"syscalls", "fs", "fs_mux_io"}},
	PpollEventID:                  {ID: PpollEventID, ID32Bit: sys32ppoll, Name: "ppoll", Probes: []probe{{event: "ppoll", attach: sysCall, fn: "ppoll"}}, Sets: []string{"syscalls", "fs", "fs_mux_io"}},
	UnshareEventID:                {ID: UnshareEventID, ID32Bit: sys32unshare, Name: "unshare", Probes: []probe{{event: "unshare", attach: sysCall, fn: "unshare"}}, Sets: []string{"syscalls", "proc"}},
	SetRobustListEventID:          {ID: SetRobustListEventID, ID32Bit: sys32set_robust_list, Name: "set_robust_list", Probes: []probe{{event: "set_robust_list", attach: sysCall, fn: "set_robust_list"}}, Sets: []string{"syscalls", "ipc", "ipc_futex"}},
	GetRobustListEventID:          {ID: GetRobustListEventID, ID32Bit: sys32get_robust_list, Name: "get_robust_list", Probes: []probe{{event: "get_robust_list", attach: sysCall, fn: "get_robust_list"}}, Sets: []string{"syscalls", "ipc", "ipc_futex"}},
	SpliceEventID:                 {ID: SpliceEventID, ID32Bit: sys32splice, Name: "splice", Probes: []probe{{event: "splice", attach: sysCall, fn: "splice"}}, Sets: []string{"syscalls", "ipc", "ipc_pipe"}},
	TeeEventID:                    {ID: TeeEventID, ID32Bit: sys32tee, Name: "tee", Probes: []probe{{event: "tee", attach: sysCall, fn: "tee"}}, Sets: []string{"syscalls", "ipc", "ipc_pipe"}},
	SyncFileRangeEventID:          {ID: SyncFileRangeEventID, ID32Bit: sys32sync_file_range, Name: "sync_file_range", Probes: []probe{{event: "sync_file_range", attach: sysCall, fn: "sync_file_range"}}, Sets: []string{"syscalls", "fs", "fs_sync"}},
	VmspliceEventID:               {ID: VmspliceEventID, ID32Bit: sys32vmsplice, Name: "vmsplice", Probes: []probe{{event: "vmsplice", attach: sysCall, fn: "vmsplice"}}, Sets: []string{"syscalls", "ipc", "ipc_pipe"}},
	MovePagesEventID:              {ID: MovePagesEventID, ID32Bit: sys32move_pages, Name: "move_pages", Probes: []probe{{event: "move_pages", attach: sysCall, fn: "move_pages"}}, Sets: []string{"syscalls", "system", "system_numa"}},
	UtimensatEventID:              {ID: UtimensatEventID, ID32Bit: sys32utimensat, Name: "utimensat", Probes: []probe{{event: "utimensat", attach: sysCall, fn: "utimensat"}}, Sets: []string{"syscalls", "fs", "fs_file_attr"}},
	EpollPwaitEventID:             {ID: EpollPwaitEventID, ID32Bit: sys32epoll_pwait, Name: "epoll_pwait", Probes: []probe{{event: "epoll_pwait", attach: sysCall, fn: "epoll_pwait"}}, Sets: []string{"syscalls", "fs", "fs_mux_io"}},
	SignalfdEventID:               {ID: SignalfdEventID, ID32Bit: sys32signalfd, Name: "signalfd", Probes: []probe{{event: "signalfd", attach: sysCall, fn: "signalfd"}}, Sets: []string{"syscalls", "signals"}},
	TimerfdCreateEventID:          {ID: TimerfdCreateEventID, ID32Bit: sys32timerfd_create, Name: "timerfd_create", Probes: []probe{{event: "timerfd_create", attach: sysCall, fn: "timerfd_create"}}, Sets: []string{"syscalls", "time", "time_timer"}},
	EventfdEventID:                {ID: EventfdEventID, ID32Bit: sys32eventfd, Name: "eventfd", Probes: []probe{{event: "eventfd", attach: sysCall, fn: "eventfd"}}, Sets: []string{"syscalls", "signals"}},
	FallocateEventID:              {ID: FallocateEventID, ID32Bit: sys32fallocate, Name: "fallocate", Probes: []probe{{event: "fallocate", attach: sysCall, fn: "fallocate"}}, Sets: []string{"syscalls", "fs", "fs_file_ops"}},
	TimerfdSettimeEventID:         {ID: TimerfdSettimeEventID, ID32Bit: sys32timerfd_settime, Name: "timerfd_settime", Probes: []probe{{event: "timerfd_settime", attach: sysCall, fn: "timerfd_settime"}}, Sets: []string{"syscalls", "time", "time_timer"}},
	TimerfdGettimeEventID:         {ID: TimerfdGettimeEventID, ID32Bit: sys32timerfd_gettime, Name: "timerfd_gettime", Probes: []probe{{event: "timerfd_gettime", attach: sysCall, fn: "timerfd_gettime"}}, Sets: []string{"syscalls", "time", "time_timer"}},
	Accept4EventID:                {ID: Accept4EventID, ID32Bit: sys32accept4, Name: "accept4", Probes: []probe{{event: "accept4", attach: sysCall, fn: "accept4"}}, Sets: []string{"default", "syscalls", "net", "net_sock"}},
	Signalfd4EventID:              {ID: Signalfd4EventID, ID32Bit: sys32signalfd4, Name: "signalfd4", Probes: []probe{{event: "signalfd4", attach: sysCall, fn: "signalfd4"}}, Sets: []string{"syscalls", "signals"}},
	Eventfd2EventID:               {ID: Eventfd2EventID, ID32Bit: sys32eventfd2, Name: "eventfd2", Probes: []probe{{event: "eventfd2", attach: sysCall, fn: "eventfd2"}}, Sets: []string{"syscalls", "signals"}},
	EpollCreate1EventID:           {ID: EpollCreate1EventID, ID32Bit: sys32epoll_create1, Name: "epoll_create1", Probes: []probe{{event: "epoll_create1", attach: sysCall, fn: "epoll_create1"}}, Sets: []string{"syscalls", "fs", "fs_mux_io"}},
	Dup3EventID:                   {ID: Dup3EventID, ID32Bit: sys32dup3, Name: "dup3", Probes: []probe{{event: "dup3", attach: sysCall, fn: "dup3"}}, Sets: []string{"default", "syscalls", "fs", "fs_fd_ops"}},
	Pipe2EventID:                  {ID: Pipe2EventID, ID32Bit: sys32pipe2, Name: "pipe2", Probes: []probe{{event: "pipe2", attach: sysCall, fn: "pipe2"}}, Sets: []string{"syscalls", "ipc", "ipc_pipe"}},
	InotifyInit1EventID:           {ID: InotifyInit1EventID, ID32Bit: sys32inotify_init1, Name: "inotify_init1", Probes: []probe{{event: "inotify_init1", attach: sysCall, fn: "inotify_init1"}}, Sets: []string{"syscalls", "fs", "fs_monitor"}},
	PreadvEventID:                 {ID: PreadvEventID, ID32Bit: sys32preadv, Name: "preadv", Probes: []probe{{event: "preadv", attach: sysCall, fn: "preadv"}}, Sets: []string{"syscalls", "fs", "fs_read_write"}},
	PwritevEventID:                {ID: PwritevEventID, ID32Bit: sys32pwritev, Name: "pwritev", Probes: []probe{{event: "pwritev", attach: sysCall, fn: "pwritev"}}, Sets: []string{"syscalls", "fs", "fs_read_write"}},
	RtTgsigqueueinfoEventID:       {ID: RtTgsigqueueinfoEventID, ID32Bit: sys32rt_tgsigqueueinfo, Name: "rt_tgsigqueueinfo", Probes: []probe{{event: "rt_tgsigqueueinfo", attach: sysCall, fn: "rt_tgsigqueueinfo"}}, Sets: []string{"syscalls", "signals"}},
	PerfEventOpenEventID:          {ID: PerfEventOpenEventID, ID32Bit: sys32perf_event_open, Name: "perf_event_open", Probes: []probe{{event: "perf_event_open", attach: sysCall, fn: "perf_event_open"}}, Sets: []string{"syscalls", "system"}},
	RecvmmsgEventID:               {ID: RecvmmsgEventID, ID32Bit: sys32recvmmsg, Name: "recvmmsg", Probes: []probe{{event: "recvmmsg", attach: sysCall, fn: "recvmmsg"}}, Sets: []string{"syscalls", "net", "net_snd_rcv"}},
	FanotifyInitEventID:           {ID: FanotifyInitEventID, ID32Bit: sys32fanotify_init, Name: "fanotify_init", Probes: []probe{{event: "fanotify_init", attach: sysCall, fn: "fanotify_init"}}, Sets: []string{"syscalls", "fs", "fs_monitor"}},
	FanotifyMarkEventID:           {ID: FanotifyMarkEventID, ID32Bit: sys32fanotify_mark, Name: "fanotify_mark", Probes: []probe{{event: "fanotify_mark", attach: sysCall, fn: "fanotify_mark"}}, Sets: []string{"syscalls", "fs", "fs_monitor"}},
	Prlimit64EventID:              {ID: Prlimit64EventID, ID32Bit: sys32prlimit64, Name: "prlimit64", Probes: []probe{{event: "prlimit64", attach: sysCall, fn: "prlimit64"}}, Sets: []string{"syscalls", "proc"}},
	NameToHandleAtEventID:         {ID: NameToHandleAtEventID, ID32Bit: sys32name_to_handle_at, Name: "name_to_handle_at", Probes: []probe{{event: "name_to_handle_at", attach: sysCall, fn: "name_to_handle_at"}}, Sets: []string{"syscalls", "fs", "fs_file_ops"}},
	OpenByHandleAtEventID:         {ID: OpenByHandleAtEventID, ID32Bit: sys32open_by_handle_at, Name: "open_by_handle_at", Probes: []probe{{event: "open_by_handle_at", attach: sysCall, fn: "open_by_handle_at"}}, Sets: []string{"syscalls", "fs", "fs_file_ops"}},
	ClockAdjtimeEventID:           {ID: ClockAdjtimeEventID, ID32Bit: sys32clock_adjtime, Name: "clock_adjtime", Probes: []probe{{event: "clock_adjtime", attach: sysCall, fn: "clock_adjtime"}}, Sets: []string{"syscalls", "time", "time_clock"}},
	SyncfsEventID:                 {ID: SyncfsEventID, ID32Bit: sys32syncfs, Name: "syncfs", Probes: []probe{{event: "syncfs", attach: sysCall, fn: "syncfs"}}, Sets: []string{"syscalls", "fs", "fs_sync"}},
	SendmmsgEventID:               {ID: SendmmsgEventID, ID32Bit: sys32sendmmsg, Name: "sendmmsg", Probes: []probe{{event: "sendmmsg", attach: sysCall, fn: "sendmmsg"}}, Sets: []string{"syscalls", "net", "net_snd_rcv"}},
	SetnsEventID:                  {ID: SetnsEventID, ID32Bit: sys32setns, Name: "setns", Probes: []probe{{event: "setns", attach: sysCall, fn: "setns"}}, Sets: []string{"syscalls", "proc"}},
	GetcpuEventID:                 {ID: GetcpuEventID, ID32Bit: sys32getcpu, Name: "getcpu", Probes: []probe{{event: "getcpu", attach: sysCall, fn: "getcpu"}}, Sets: []string{"syscalls", "system", "system_numa"}},
	ProcessVmReadvEventID:         {ID: ProcessVmReadvEventID, ID32Bit: sys32process_vm_readv, Name: "process_vm_readv", Probes: []probe{{event: "process_vm_readv", attach: sysCall, fn: "process_vm_readv"}}, Sets: []string{"default", "syscalls", "proc"}},
	ProcessVmWritevEventID:        {ID: ProcessVmWritevEventID, ID32Bit: sys32process_vm_writev, Name: "process_vm_writev", Probes: []probe{{event: "process_vm_writev", attach: sysCall, fn: "process_vm_writev"}}, Sets: []string{"default", "syscalls", "proc"}},
	KcmpEventID:                   {ID: KcmpEventID, ID32Bit: sys32kcmp, Name: "kcmp", Probes: []probe{{event: "kcmp", attach: sysCall, fn: "kcmp"}}, Sets: []string{"syscalls", "proc"}},
	FinitModuleEventID:            {ID: FinitModuleEventID, ID32Bit: sys32finit_module, Name: "finit_module", Probes: []probe{{event: "finit_module", attach: sysCall, fn: "finit_module"}}, Sets: []string{"default", "syscalls", "system", "system_module"}},
	SchedSetattrEventID:           {ID: SchedSetattrEventID, ID32Bit: sys32sched_setattr, Name: "sched_setattr", Probes: []probe{{event: "sched_setattr", attach: sysCall, fn: "sched_setattr"}}, Sets: []string{"syscalls", "proc", "proc_sched"}},
	SchedGetattrEventID:           {ID: SchedGetattrEventID, ID32Bit: sys32sched_getattr, Name: "sched_getattr", Probes: []probe{{event: "sched_getattr", attach: sysCall, fn: "sched_getattr"}}, Sets: []string{"syscalls", "proc", "proc_sched"}},
	Renameat2EventID:              {ID: Renameat2EventID, ID32Bit: sys32renameat2, Name: "renameat2", Probes: []probe{{event: "renameat2", attach: sysCall, fn: "renameat2"}}, Sets: []string{"syscalls", "fs", "fs_file_ops"}},
	SeccompEventID:                {ID: SeccompEventID, ID32Bit: sys32seccomp, Name: "seccomp", Probes: []probe{{event: "seccomp", attach: sysCall, fn: "seccomp"}}, Sets: []string{"syscalls", "proc"}},
	GetrandomEventID:              {ID: GetrandomEventID, ID32Bit: sys32getrandom, Name: "getrandom", Probes: []probe{{event: "getrandom", attach: sysCall, fn: "getrandom"}}, Sets: []string{"syscalls", "fs"}},
	MemfdCreateEventID:            {ID: MemfdCreateEventID, ID32Bit: sys32memfd_create, Name: "memfd_create", Probes: []probe{{event: "memfd_create", attach: sysCall, fn: "memfd_create"}}, Sets: []string{"default", "syscalls", "fs", "fs_file_ops"}},
	KexecFileLoadEventID:          {ID: KexecFileLoadEventID, ID32Bit: sys32undefined, Name: "kexec_file_load", Probes: []probe{{event: "kexec_file_load", attach: sysCall, fn: "kexec_file_load"}}, Sets: []string{"syscalls", "system"}},
	BpfEventID:                    {ID: BpfEventID, ID32Bit: sys32bpf, Name: "bpf", Probes: []probe{{event: "bpf", attach: sysCall, fn: "bpf"}}, Sets: []string{"default", "syscalls", "system"}},
	ExecveatEventID:               {ID: ExecveatEventID, ID32Bit: sys32execveat, Name: "execveat", Probes: []probe{{event: "execveat", attach: sysCall, fn: "execveat"}}, Sets: []string{"default", "syscalls", "proc", "proc_life"}},
	UserfaultfdEventID:            {ID: UserfaultfdEventID, ID32Bit: sys32userfaultfd, Name: "userfaultfd", Probes: []probe{{event: "userfaultfd", attach: sysCall, fn: "userfaultfd"}}, Sets: []string{"syscalls", "system"}},
	MembarrierEventID:             {ID: MembarrierEventID, ID32Bit: sys32membarrier, Name: "membarrier", Probes: []probe{{event: "membarrier", attach: sysCall, fn: "membarrier"}}, Sets: []string{"syscalls", "proc", "proc_mem"}},
	Mlock2EventID:                 {ID: Mlock2EventID, ID32Bit: sys32mlock2, Name: "mlock2", Probes: []probe{{event: "mlock2", attach: sysCall, fn: "mlock2"}}, Sets: []string{"syscalls", "proc", "proc_mem"}},
	CopyFileRangeEventID:          {ID: CopyFileRangeEventID, ID32Bit: sys32copy_file_range, Name: "copy_file_range", Probes: []probe{{event: "copy_file_range", attach: sysCall, fn: "copy_file_range"}}, Sets: []string{"syscalls", "fs", "fs_read_write"}},
	Preadv2EventID:                {ID: Preadv2EventID, ID32Bit: sys32preadv2, Name: "preadv2", Probes: []probe{{event: "preadv2", attach: sysCall, fn: "preadv2"}}, Sets: []string{"syscalls", "fs", "fs_read_write"}},
	Pwritev2EventID:               {ID: Pwritev2EventID, ID32Bit: sys32pwritev2, Name: "pwritev2", Probes: []probe{{event: "pwritev2", attach: sysCall, fn: "pwritev2"}}, Sets: []string{"syscalls", "fs", "fs_read_write"}},
	PkeyMprotectEventID:           {ID: PkeyMprotectEventID, ID32Bit: sys32pkey_mprotect, Name: "pkey_mprotect", Probes: []probe{{event: "pkey_mprotect", attach: sysCall, fn: "pkey_mprotect"}}, Sets: []string{"default", "syscalls", "proc", "proc_mem"}},
	PkeyAllocEventID:              {ID: PkeyAllocEventID, ID32Bit: sys32pkey_alloc, Name: "pkey_alloc", Probes: []probe{{event: "pkey_alloc", attach: sysCall, fn: "pkey_alloc"}}, Sets: []string{"syscalls", "proc", "proc_mem"}},
	PkeyFreeEventID:               {ID: PkeyFreeEventID, ID32Bit: sys32pkey_free, Name: "pkey_free", Probes: []probe{{event: "pkey_free", attach: sysCall, fn: "pkey_free"}}, Sets: []string{"syscalls", "proc", "proc_mem"}},
	StatxEventID:                  {ID: StatxEventID, ID32Bit: sys32statx, Name: "statx", Probes: []probe{{event: "statx", attach: sysCall, fn: "statx"}}, Sets: []string{"syscalls", "fs", "fs_file_attr"}},
	IoPgeteventsEventID:           {ID: IoPgeteventsEventID, ID32Bit: sys32io_pgetevents, Name: "io_pgetevents", Probes: []probe{{event: "io_pgetevents", attach: sysCall, fn: "io_pgetevents"}}, Sets: []string{"syscalls", "fs", "fs_async_io"}},
	RseqEventID:                   {ID: RseqEventID, ID32Bit: sys32rseq, Name: "rseq", Probes: []probe{{event: "rseq", attach: sysCall, fn: "rseq"}}, Sets: []string{"syscalls"}},
	PidfdSendSignalEventID:        {ID: PidfdSendSignalEventID, ID32Bit: sys32pidfd_send_signal, Name: "pidfd_send_signal", Probes: []probe{{event: "pidfd_send_signal", attach: sysCall, fn: "pidfd_send_signal"}}, Sets: []string{"syscalls", "signals"}},
	IoUringSetupEventID:           {ID: IoUringSetupEventID, ID32Bit: sys32io_uring_setup, Name: "io_uring_setup", Probes: []probe{{event: "io_uring_setup", attach: sysCall, fn: "io_uring_setup"}}, Sets: []string{"syscalls"}},
	IoUringEnterEventID:           {ID: IoUringEnterEventID, ID32Bit: sys32io_uring_enter, Name: "io_uring_enter", Probes: []probe{{event: "io_uring_enter", attach: sysCall, fn: "io_uring_enter"}}, Sets: []string{"syscalls"}},
	IoUringRegisterEventID:        {ID: IoUringRegisterEventID, ID32Bit: sys32io_uring_register, Name: "io_uring_register", Probes: []probe{{event: "io_uring_register", attach: sysCall, fn: "io_uring_register"}}, Sets: []string{"syscalls"}},
	OpenTreeEventID:               {ID: OpenTreeEventID, ID32Bit: sys32open_tree, Name: "open_tree", Probes: []probe{{event: "open_tree", attach: sysCall, fn: "open_tree"}}, Sets: []string{"syscalls"}},
	MoveMountEventID:              {ID: MoveMountEventID, ID32Bit: sys32move_mount, Name: "move_mount", Probes: []probe{{event: "move_mount", attach: sysCall, fn: "move_mount"}}, Sets: []string{"default", "syscalls", "fs"}},
	FsopenEventID:                 {ID: FsopenEventID, ID32Bit: sys32fsopen, Name: "fsopen", Probes: []probe{{event: "fsopen", attach: sysCall, fn: "fsopen"}}, Sets: []string{"syscalls", "fs"}},
	FsconfigEventID:               {ID: FsconfigEventID, ID32Bit: sys32fsconfig, Name: "fsconfig", Probes: []probe{{event: "fsconfig", attach: sysCall, fn: "fsconfig"}}, Sets: []string{"syscalls", "fs"}},
	FsmountEventID:                {ID: FsmountEventID, ID32Bit: sys32fsmount, Name: "fsmount", Probes: []probe{{event: "fsmount", attach: sysCall, fn: "fsmount"}}, Sets: []string{"syscalls", "fs"}},
	FspickEventID:                 {ID: FspickEventID, ID32Bit: sys32fspick, Name: "fspick", Probes: []probe{{event: "fspick", attach: sysCall, fn: "fspick"}}, Sets: []string{"syscalls", "fs"}},
	PidfdOpenEventID:              {ID: PidfdOpenEventID, ID32Bit: sys32pidfd_open, Name: "pidfd_open", Probes: []probe{{event: "pidfd_open", attach: sysCall, fn: "pidfd_open"}}, Sets: []string{"syscalls"}},
	Clone3EventID:                 {ID: Clone3EventID, ID32Bit: sys32clone3, Name: "clone3", Probes: []probe{{event: "clone3", attach: sysCall, fn: "clone3"}}, Sets: []string{"default", "syscalls", "proc", "proc_life"}},
	CloseRangeEventID:             {ID: CloseRangeEventID, ID32Bit: sys32close_range, Name: "close_range", Probes: []probe{{event: "close_range", attach: sysCall, fn: "close_range"}}, Sets: []string{"default", "syscalls", "fs", "fs_file_ops"}},
	Openat2EventID:                {ID: Openat2EventID, ID32Bit: sys32openat2, Name: "openat2", Probes: []probe{{event: "openat2", attach: sysCall, fn: "openat2"}}, Sets: []string{"default", "syscalls", "fs", "fs_file_ops"}},
	PidfdGetfdEventID:             {ID: PidfdGetfdEventID, ID32Bit: sys32pidfd_getfd, Name: "pidfd_getfd", Probes: []probe{{event: "pidfd_getfd", attach: sysCall, fn: "pidfd_getfd"}}, Sets: []string{"syscalls"}},
	Faccessat2EventID:             {ID: Faccessat2EventID, ID32Bit: sys32faccessat2, Name: "faccessat2", Probes: []probe{{event: "faccessat2", attach: sysCall, fn: "faccessat2"}}, Sets: []string{"default", "syscalls", "fs", "fs_file_attr"}},
	ProcessMadviseEventID:         {ID: ProcessMadviseEventID, ID32Bit: sys32process_madvise, Name: "process_madvise", Probes: []probe{{event: "process_madvise", attach: sysCall, fn: "process_madvise"}}, Sets: []string{"syscalls"}},
	EpollPwait2EventID:            {ID: EpollPwait2EventID, ID32Bit: sys32epoll_pwait2, Name: "epoll_pwait2", Probes: []probe{{event: "epoll_pwait2", attach: sysCall, fn: "epoll_pwait2"}}, Sets: []string{"syscalls", "fs", "fs_mux_io"}},
	SysEnterEventID:               {ID: SysEnterEventID, ID32Bit: sys32undefined, Name: "sys_enter", Probes: []probe{{event: "raw_syscalls:sys_enter", attach: rawTracepoint, fn: "tracepoint__raw_syscalls__sys_enter"}}, EssentialEvent: true, Sets: []string{}},
	SysExitEventID:                {ID: SysExitEventID, ID32Bit: sys32undefined, Name: "sys_exit", Probes: []probe{{event: "raw_syscalls:sys_exit", attach: rawTracepoint, fn: "tracepoint__raw_syscalls__sys_exit"}}, EssentialEvent: true, Sets: []string{}},
	SchedProcessForkEventID:       {ID: SchedProcessForkEventID, ID32Bit: sys32undefined, Name: "sched_process_fork", Probes: []probe{{event: "sched:sched_process_fork", attach: rawTracepoint, fn: "tracepoint__sched__sched_process_fork"}}, EssentialEvent: true, Sets: []string{}},
	SchedProcessExecEventID:       {ID: SchedProcessExecEventID, ID32Bit: sys32undefined, Name: "sched_process_exec", Probes: []probe{{event: "sched:sched_process_exec", attach: rawTracepoint, fn: "tracepoint__sched__sched_process_exec"}}, EssentialEvent: true, Sets: []string{"default", "proc"}},
	SchedProcessExitEventID:       {ID: SchedProcessExitEventID, ID32Bit: sys32undefined, Name: "sched_process_exit", Probes: []probe{{event: "sched:sched_process_exit", attach: rawTracepoint, fn: "tracepoint__sched__sched_process_exit"}}, EssentialEvent: true, Sets: []string{"default", "proc", "proc_life"}},
	SchedSwitchEventID:            {ID: SchedSwitchEventID, ID32Bit: sys32undefined, Name: "sched_switch", Probes: []probe{{event: "sched:sched_switch", attach: rawTracepoint, fn: "tracepoint__sched__sched_switch"}}, Sets: []string{}},
	DoExitEventID:                 {ID: DoExitEventID, ID32Bit: sys32undefined, Name: "do_exit", Probes: []probe{{event: "do_exit", attach: kprobe, fn: "trace_do_exit"}}, Sets: []string{"proc", "proc_life"}},
	CapCapableEventID:             {ID: CapCapableEventID, ID32Bit: sys32undefined, Name: "cap_capable", Probes: []probe{{event: "cap_capable", attach: kprobe, fn: "trace_cap_capable"}}, Sets: []string{"default"}},
	VfsWriteEventID:               {ID: VfsWriteEventID, ID32Bit: sys32undefined, Name: "vfs_write", Probes: []probe{{event: "vfs_write", attach: kprobe, fn: "trace_vfs_write"}, {event: "vfs_write", attach: kretprobe, fn: "trace_ret_vfs_write"}}, Sets: []string{}},
	VfsWritevEventID:              {ID: VfsWritevEventID, ID32Bit: sys32undefined, Name: "vfs_writev", Probes: []probe{{event: "vfs_writev", attach: kprobe, fn: "trace_vfs_writev"}, {event: "vfs_writev", attach: kretprobe, fn: "trace_ret_vfs_writev"}}, Sets: []string{}},
	MemProtAlertEventID:           {ID: MemProtAlertEventID, ID32Bit: sys32undefined, Name: "mem_prot_alert", Probes: []probe{{event: "security_mmap_addr", attach: kprobe, fn: "trace_mmap_alert"}, {event: "security_file_mprotect", attach: kprobe, fn: "trace_mprotect_alert"}}, Sets: []string{}},
	CommitCredsEventID:            {ID: CommitCredsEventID, ID32Bit: sys32undefined, Name: "commit_creds", Probes: []probe{{event: "commit_creds", attach: kprobe, fn: "trace_commit_creds"}}, Sets: []string{}},
	SwitchTaskNSEventID:           {ID: SwitchTaskNSEventID, ID32Bit: sys32undefined, Name: "switch_task_ns", Probes: []probe{{event: "switch_task_namespaces", attach: kprobe, fn: "trace_switch_task_namespaces"}}, Sets: []string{}},
	MagicWriteEventID:             {ID: MagicWriteEventID, ID32Bit: sys32undefined, Name: "magic_write", Probes: []probe{}, Sets: []string{}},
	CgroupAttachTaskEventID:       {ID: CgroupAttachTaskEventID, ID32Bit: sys32undefined, Name: "cgroup_attach_task", Probes: []probe{{event: "cgroup:cgroup_attach_task", attach: rawTracepoint, fn: "tracepoint__cgroup__cgroup_attach_task"}}, Sets: []string{}},
	CgroupMkdirEventID:            {ID: CgroupMkdirEventID, ID32Bit: sys32undefined, Name: "cgroup_mkdir", Probes: []probe{{event: "cgroup:cgroup_mkdir", attach: rawTracepoint, fn: "tracepoint__cgroup__cgroup_mkdir"}}, EssentialEvent: true, Sets: []string{}},
	CgroupRmdirEventID:            {ID: CgroupRmdirEventID, ID32Bit: sys32undefined, Name: "cgroup_rmdir", Probes: []probe{{event: "cgroup:cgroup_rmdir", attach: rawTracepoint, fn: "tracepoint__cgroup__cgroup_rmdir"}}, EssentialEvent: true, Sets: []string{}},
	SecurityBprmCheckEventID:      {ID: SecurityBprmCheckEventID, ID32Bit: sys32undefined, Name: "security_bprm_check", Probes: []probe{{event: "security_bprm_check", attach: kprobe, fn: "trace_security_bprm_check"}}, Sets: []string{"default", "lsm_hooks", "proc", "proc_life"}},
	SecurityFileOpenEventID:       {ID: SecurityFileOpenEventID, ID32Bit: sys32undefined, Name: "security_file_open", Probes: []probe{{event: "security_file_open", attach: kprobe, fn: "trace_security_file_open"}}, Sets: []string{"default", "lsm_hooks", "fs", "fs_file_ops"}},
	SecurityInodeUnlinkEventID:    {ID: SecurityInodeUnlinkEventID, ID32Bit: sys32undefined, Name: "security_inode_unlink", Probes: []probe{{event: "security_inode_unlink", attach: kprobe, fn: "trace_security_inode_unlink"}}, Sets: []string{"default", "lsm_hooks", "fs", "fs_file_ops"}},
	SecuritySocketCreateEventID:   {ID: SecuritySocketCreateEventID, ID32Bit: sys32undefined, Name: "security_socket_create", Probes: []probe{{event: "security_socket_create", attach: kprobe, fn: "trace_security_socket_create"}}, Sets: []string{"default", "lsm_hooks", "net", "net_sock"}},
	SecuritySocketListenEventID:   {ID: SecuritySocketListenEventID, ID32Bit: sys32undefined, Name: "security_socket_listen", Probes: []probe{{event: "security_socket_listen", attach: kprobe, fn: "trace_security_socket_listen"}}, Sets: []string{"default", "lsm_hooks", "net", "net_sock"}},
	SecuritySocketConnectEventID:  {ID: SecuritySocketConnectEventID, ID32Bit: sys32undefined, Name: "security_socket_connect", Probes: []probe{{event: "security_socket_connect", attach: kprobe, fn: "trace_security_socket_connect"}}, Sets: []string{"default", "lsm_hooks", "net", "net_sock"}},
	SecuritySocketAcceptEventID:   {ID: SecuritySocketAcceptEventID, ID32Bit: sys32undefined, Name: "security_socket_accept", Probes: []probe{{event: "security_socket_accept", attach: kprobe, fn: "trace_security_socket_accept"}}, Sets: []string{"default", "lsm_hooks", "net", "net_sock"}},
	SecuritySocketBindEventID:     {ID: SecuritySocketBindEventID, ID32Bit: sys32undefined, Name: "security_socket_bind", Probes: []probe{{event: "security_socket_bind", attach: kprobe, fn: "trace_security_socket_bind"}}, Sets: []string{"default", "lsm_hooks", "net", "net_sock"}},
	SecuritySbMountEventID:        {ID: SecuritySbMountEventID, ID32Bit: sys32undefined, Name: "security_sb_mount", Probes: []probe{{event: "security_sb_mount", attach: kprobe, fn: "trace_security_sb_mount"}}, Sets: []string{"default", "lsm_hooks", "fs"}},
	SecurityBPFEventID:            {ID: SecurityBPFEventID, ID32Bit: sys32undefined, Name: "security_bpf", Probes: []probe{{event: "security_bpf", attach: kprobe, fn: "trace_security_bpf"}}, Sets: []string{"lsm_hooks"}},
	SecurityBPFMapEventID:         {ID: SecurityBPFMapEventID, ID32Bit: sys32undefined, Name: "security_bpf_map", Probes: []probe{{event: "security_bpf_map", attach: kprobe, fn: "trace_security_bpf_map"}}, Sets: []string{"lsm_hooks"}},
	SecurityKernelReadFileEventID: {ID: SecurityKernelReadFileEventID, ID32Bit: sys32undefined, Name: "security_kernel_read_file", Probes: []probe{{event: "security_kernel_read_file", attach: kprobe, fn: "trace_security_kernel_read_file"}}, Sets: []string{"lsm_hooks"}},
	SecurityPostReadFileEventID:   {ID: SecurityPostReadFileEventID, ID32Bit: sys32undefined, Name: "security_kernel_post_read_file", Probes: []probe{{event: "security_kernel_post_read_file", attach: kprobe, fn: "trace_security_kernel_post_read_file"}}, Sets: []string{"lsm_hooks"}},
	SecurityInodeMknodEventID:     {ID: SecurityInodeMknodEventID, ID32Bit: sys32undefined, Name: "security_inode_mknod", Probes: []probe{{event: "security_inode_mknod", attach: kprobe, fn: "trace_security_inode_mknod"}}, Sets: []string{"lsm_hooks"}},
	InitNamespacesEventID:         {ID: InitNamespacesEventID, ID32Bit: sys32undefined, Name: "init_namespaces", Probes: []probe{}, Sets: []string{}},
	SocketDupEventID:              {ID: SocketDupEventID, ID32Bit: sys32undefined, Name: "socket_dup", Probes: []probe{}, Sets: []string{}},
}

// EventsIDToParams is list of the parameters (name and type) used by the events
var EventsIDToParams = map[int32][]external.ArgMeta{
	ReadEventID:                   {{Type: "int", Name: "fd"}, {Type: "void*", Name: "buf"}, {Type: "size_t", Name: "count"}},
	WriteEventID:                  {{Type: "int", Name: "fd"}, {Type: "void*", Name: "buf"}, {Type: "size_t", Name: "count"}},
	OpenEventID:                   {{Type: "const char*", Name: "pathname"}, {Type: "int", Name: "flags"}, {Type: "mode_t", Name: "mode"}},
	CloseEventID:                  {{Type: "int", Name: "fd"}},
	StatEventID:                   {{Type: "const char*", Name: "pathname"}, {Type: "struct stat*", Name: "statbuf"}},
	FstatEventID:                  {{Type: "int", Name: "fd"}, {Type: "struct stat*", Name: "statbuf"}},
	LstatEventID:                  {{Type: "const char*", Name: "pathname"}, {Type: "struct stat*", Name: "statbuf"}},
	PollEventID:                   {{Type: "struct pollfd*", Name: "fds"}, {Type: "unsigned int", Name: "nfds"}, {Type: "int", Name: "timeout"}},
	LseekEventID:                  {{Type: "int", Name: "fd"}, {Type: "off_t", Name: "offset"}, {Type: "unsigned int", Name: "whence"}},
	MmapEventID:                   {{Type: "void*", Name: "addr"}, {Type: "size_t", Name: "length"}, {Type: "int", Name: "prot"}, {Type: "int", Name: "flags"}, {Type: "int", Name: "fd"}, {Type: "off_t", Name: "off"}},
	MprotectEventID:               {{Type: "void*", Name: "addr"}, {Type: "size_t", Name: "len"}, {Type: "int", Name: "prot"}},
	MunmapEventID:                 {{Type: "void*", Name: "addr"}, {Type: "size_t", Name: "length"}},
	BrkEventID:                    {{Type: "void*", Name: "addr"}},
	RtSigactionEventID:            {{Type: "int", Name: "signum"}, {Type: "const struct sigaction*", Name: "act"}, {Type: "struct sigaction*", Name: "oldact"}, {Type: "size_t", Name: "sigsetsize"}},
	RtSigprocmaskEventID:          {{Type: "int", Name: "how"}, {Type: "sigset_t*", Name: "set"}, {Type: "sigset_t*", Name: "oldset"}, {Type: "size_t", Name: "sigsetsize"}},
	RtSigreturnEventID:            {},
	IoctlEventID:                  {{Type: "int", Name: "fd"}, {Type: "unsigned long", Name: "request"}, {Type: "unsigned long", Name: "arg"}},
	Pread64EventID:                {{Type: "int", Name: "fd"}, {Type: "void*", Name: "buf"}, {Type: "size_t", Name: "count"}, {Type: "off_t", Name: "offset"}},
	Pwrite64EventID:               {{Type: "int", Name: "fd"}, {Type: "const void*", Name: "buf"}, {Type: "size_t", Name: "count"}, {Type: "off_t", Name: "offset"}},
	ReadvEventID:                  {{Type: "int", Name: "fd"}, {Type: "const struct iovec*", Name: "iov"}, {Type: "int", Name: "iovcnt"}},
	WritevEventID:                 {{Type: "int", Name: "fd"}, {Type: "const struct iovec*", Name: "iov"}, {Type: "int", Name: "iovcnt"}},
	AccessEventID:                 {{Type: "const char*", Name: "pathname"}, {Type: "int", Name: "mode"}},
	PipeEventID:                   {{Type: "int[2]", Name: "pipefd"}},
	SelectEventID:                 {{Type: "int", Name: "nfds"}, {Type: "fd_set*", Name: "readfds"}, {Type: "fd_set*", Name: "writefds"}, {Type: "fd_set*", Name: "exceptfds"}, {Type: "struct timeval*", Name: "timeout"}},
	SchedYieldEventID:             {},
	MremapEventID:                 {{Type: "void*", Name: "old_address"}, {Type: "size_t", Name: "old_size"}, {Type: "size_t", Name: "new_size"}, {Type: "int", Name: "flags"}, {Type: "void*", Name: "new_address"}},
	MsyncEventID:                  {{Type: "void*", Name: "addr"}, {Type: "size_t", Name: "length"}, {Type: "int", Name: "flags"}},
	MincoreEventID:                {{Type: "void*", Name: "addr"}, {Type: "size_t", Name: "length"}, {Type: "unsigned char*", Name: "vec"}},
	MadviseEventID:                {{Type: "void*", Name: "addr"}, {Type: "size_t", Name: "length"}, {Type: "int", Name: "advice"}},
	ShmgetEventID:                 {{Type: "key_t", Name: "key"}, {Type: "size_t", Name: "size"}, {Type: "int", Name: "shmflg"}},
	ShmatEventID:                  {{Type: "int", Name: "shmid"}, {Type: "const void*", Name: "shmaddr"}, {Type: "int", Name: "shmflg"}},
	ShmctlEventID:                 {{Type: "int", Name: "shmid"}, {Type: "int", Name: "cmd"}, {Type: "struct shmid_ds*", Name: "buf"}},
	DupEventID:                    {{Type: "int", Name: "oldfd"}},
	Dup2EventID:                   {{Type: "int", Name: "oldfd"}, {Type: "int", Name: "newfd"}},
	PauseEventID:                  {},
	NanosleepEventID:              {{Type: "const struct timespec*", Name: "req"}, {Type: "struct timespec*", Name: "rem"}},
	GetitimerEventID:              {{Type: "int", Name: "which"}, {Type: "struct itimerval*", Name: "curr_value"}},
	AlarmEventID:                  {{Type: "unsigned int", Name: "seconds"}},
	SetitimerEventID:              {{Type: "int", Name: "which"}, {Type: "struct itimerval*", Name: "new_value"}, {Type: "struct itimerval*", Name: "old_value"}},
	GetpidEventID:                 {},
	SendfileEventID:               {{Type: "int", Name: "out_fd"}, {Type: "int", Name: "in_fd"}, {Type: "off_t*", Name: "offset"}, {Type: "size_t", Name: "count"}},
	SocketEventID:                 {{Type: "int", Name: "domain"}, {Type: "int", Name: "type"}, {Type: "int", Name: "protocol"}},
	ConnectEventID:                {{Type: "int", Name: "sockfd"}, {Type: "struct sockaddr*", Name: "addr"}, {Type: "int", Name: "addrlen"}},
	AcceptEventID:                 {{Type: "int", Name: "sockfd"}, {Type: "struct sockaddr*", Name: "addr"}, {Type: "int*", Name: "addrlen"}},
	SendtoEventID:                 {{Type: "int", Name: "sockfd"}, {Type: "void*", Name: "buf"}, {Type: "size_t", Name: "len"}, {Type: "int", Name: "flags"}, {Type: "struct sockaddr*", Name: "dest_addr"}, {Type: "int", Name: "addrlen"}},
	RecvfromEventID:               {{Type: "int", Name: "sockfd"}, {Type: "void*", Name: "buf"}, {Type: "size_t", Name: "len"}, {Type: "int", Name: "flags"}, {Type: "struct sockaddr*", Name: "src_addr"}, {Type: "int*", Name: "addrlen"}},
	SendmsgEventID:                {{Type: "int", Name: "sockfd"}, {Type: "struct msghdr*", Name: "msg"}, {Type: "int", Name: "flags"}},
	RecvmsgEventID:                {{Type: "int", Name: "sockfd"}, {Type: "struct msghdr*", Name: "msg"}, {Type: "int", Name: "flags"}},
	ShutdownEventID:               {{Type: "int", Name: "sockfd"}, {Type: "int", Name: "how"}},
	BindEventID:                   {{Type: "int", Name: "sockfd"}, {Type: "struct sockaddr*", Name: "addr"}, {Type: "int", Name: "addrlen"}},
	ListenEventID:                 {{Type: "int", Name: "sockfd"}, {Type: "int", Name: "backlog"}},
	GetsocknameEventID:            {{Type: "int", Name: "sockfd"}, {Type: "struct sockaddr*", Name: "addr"}, {Type: "int*", Name: "addrlen"}},
	GetpeernameEventID:            {{Type: "int", Name: "sockfd"}, {Type: "struct sockaddr*", Name: "addr"}, {Type: "int*", Name: "addrlen"}},
	SocketpairEventID:             {{Type: "int", Name: "domain"}, {Type: "int", Name: "type"}, {Type: "int", Name: "protocol"}, {Type: "int[2]", Name: "sv"}},
	SetsockoptEventID:             {{Type: "int", Name: "sockfd"}, {Type: "int", Name: "level"}, {Type: "int", Name: "optname"}, {Type: "const void*", Name: "optval"}, {Type: "int", Name: "optlen"}},
	GetsockoptEventID:             {{Type: "int", Name: "sockfd"}, {Type: "int", Name: "level"}, {Type: "int", Name: "optname"}, {Type: "void*", Name: "optval"}, {Type: "int*", Name: "optlen"}},
	CloneEventID:                  {{Type: "unsigned long", Name: "flags"}, {Type: "void*", Name: "stack"}, {Type: "int*", Name: "parent_tid"}, {Type: "int*", Name: "child_tid"}, {Type: "unsigned long", Name: "tls"}},
	ForkEventID:                   {},
	VforkEventID:                  {},
	ExecveEventID:                 {{Type: "const char*", Name: "pathname"}, {Type: "const char*const*", Name: "argv"}, {Type: "const char*const*", Name: "envp"}},
	ExitEventID:                   {{Type: "int", Name: "status"}},
	Wait4EventID:                  {{Type: "pid_t", Name: "pid"}, {Type: "int*", Name: "wstatus"}, {Type: "int", Name: "options"}, {Type: "struct rusage*", Name: "rusage"}},
	KillEventID:                   {{Type: "pid_t", Name: "pid"}, {Type: "int", Name: "sig"}},
	UnameEventID:                  {{Type: "struct utsname*", Name: "buf"}},
	SemgetEventID:                 {{Type: "key_t", Name: "key"}, {Type: "int", Name: "nsems"}, {Type: "int", Name: "semflg"}},
	SemopEventID:                  {{Type: "int", Name: "semid"}, {Type: "struct sembuf*", Name: "sops"}, {Type: "size_t", Name: "nsops"}},
	SemctlEventID:                 {{Type: "int", Name: "semid"}, {Type: "int", Name: "semnum"}, {Type: "int", Name: "cmd"}, {Type: "unsigned long", Name: "arg"}},
	ShmdtEventID:                  {{Type: "const void*", Name: "shmaddr"}},
	MsggetEventID:                 {{Type: "key_t", Name: "key"}, {Type: "int", Name: "msgflg"}},
	MsgsndEventID:                 {{Type: "int", Name: "msqid"}, {Type: "struct msgbuf*", Name: "msgp"}, {Type: "size_t", Name: "msgsz"}, {Type: "int", Name: "msgflg"}},
	MsgrcvEventID:                 {{Type: "int", Name: "msqid"}, {Type: "struct msgbuf*", Name: "msgp"}, {Type: "size_t", Name: "msgsz"}, {Type: "long", Name: "msgtyp"}, {Type: "int", Name: "msgflg"}},
	MsgctlEventID:                 {{Type: "int", Name: "msqid"}, {Type: "int", Name: "cmd"}, {Type: "struct msqid_ds*", Name: "buf"}},
	FcntlEventID:                  {{Type: "int", Name: "fd"}, {Type: "int", Name: "cmd"}, {Type: "unsigned long", Name: "arg"}},
	FlockEventID:                  {{Type: "int", Name: "fd"}, {Type: "int", Name: "operation"}},
	FsyncEventID:                  {{Type: "int", Name: "fd"}},
	FdatasyncEventID:              {{Type: "int", Name: "fd"}},
	TruncateEventID:               {{Type: "const char*", Name: "path"}, {Type: "off_t", Name: "length"}},
	FtruncateEventID:              {{Type: "int", Name: "fd"}, {Type: "off_t", Name: "length"}},
	GetdentsEventID:               {{Type: "int", Name: "fd"}, {Type: "struct linux_dirent*", Name: "dirp"}, {Type: "unsigned int", Name: "count"}},
	GetcwdEventID:                 {{Type: "char*", Name: "buf"}, {Type: "size_t", Name: "size"}},
	ChdirEventID:                  {{Type: "const char*", Name: "path"}},
	FchdirEventID:                 {{Type: "int", Name: "fd"}},
	RenameEventID:                 {{Type: "const char*", Name: "oldpath"}, {Type: "const char*", Name: "newpath"}},
	MkdirEventID:                  {{Type: "const char*", Name: "pathname"}, {Type: "mode_t", Name: "mode"}},
	RmdirEventID:                  {{Type: "const char*", Name: "pathname"}},
	CreatEventID:                  {{Type: "const char*", Name: "pathname"}, {Type: "mode_t", Name: "mode"}},
	LinkEventID:                   {{Type: "const char*", Name: "oldpath"}, {Type: "const char*", Name: "newpath"}},
	UnlinkEventID:                 {{Type: "const char*", Name: "pathname"}},
	SymlinkEventID:                {{Type: "const char*", Name: "target"}, {Type: "const char*", Name: "linkpath"}},
	ReadlinkEventID:               {{Type: "const char*", Name: "pathname"}, {Type: "char*", Name: "buf"}, {Type: "size_t", Name: "bufsiz"}},
	ChmodEventID:                  {{Type: "const char*", Name: "pathname"}, {Type: "mode_t", Name: "mode"}},
	FchmodEventID:                 {{Type: "int", Name: "fd"}, {Type: "mode_t", Name: "mode"}},
	ChownEventID:                  {{Type: "const char*", Name: "pathname"}, {Type: "uid_t", Name: "owner"}, {Type: "gid_t", Name: "group"}},
	FchownEventID:                 {{Type: "int", Name: "fd"}, {Type: "uid_t", Name: "owner"}, {Type: "gid_t", Name: "group"}},
	LchownEventID:                 {{Type: "const char*", Name: "pathname"}, {Type: "uid_t", Name: "owner"}, {Type: "gid_t", Name: "group"}},
	UmaskEventID:                  {{Type: "mode_t", Name: "mask"}},
	GettimeofdayEventID:           {{Type: "struct timeval*", Name: "tv"}, {Type: "struct timezone*", Name: "tz"}},
	GetrlimitEventID:              {{Type: "int", Name: "resource"}, {Type: "struct rlimit*", Name: "rlim"}},
	GetrusageEventID:              {{Type: "int", Name: "who"}, {Type: "struct rusage*", Name: "usage"}},
	SysinfoEventID:                {{Type: "struct sysinfo*", Name: "info"}},
	TimesEventID:                  {{Type: "struct tms*", Name: "buf"}},
	PtraceEventID:                 {{Type: "long", Name: "request"}, {Type: "pid_t", Name: "pid"}, {Type: "void*", Name: "addr"}, {Type: "void*", Name: "data"}},
	GetuidEventID:                 {},
	SyslogEventID:                 {{Type: "int", Name: "type"}, {Type: "char*", Name: "bufp"}, {Type: "int", Name: "len"}},
	GetgidEventID:                 {},
	SetuidEventID:                 {{Type: "uid_t", Name: "uid"}},
	SetgidEventID:                 {{Type: "gid_t", Name: "gid"}},
	GeteuidEventID:                {},
	GetegidEventID:                {},
	SetpgidEventID:                {{Type: "pid_t", Name: "pid"}, {Type: "pid_t", Name: "pgid"}},
	GetppidEventID:                {},
	GetpgrpEventID:                {},
	SetsidEventID:                 {},
	SetreuidEventID:               {{Type: "uid_t", Name: "ruid"}, {Type: "uid_t", Name: "euid"}},
	SetregidEventID:               {{Type: "gid_t", Name: "rgid"}, {Type: "gid_t", Name: "egid"}},
	GetgroupsEventID:              {{Type: "int", Name: "size"}, {Type: "gid_t*", Name: "list"}},
	SetgroupsEventID:              {{Type: "int", Name: "size"}, {Type: "gid_t*", Name: "list"}},
	SetresuidEventID:              {{Type: "uid_t", Name: "ruid"}, {Type: "uid_t", Name: "euid"}, {Type: "uid_t", Name: "suid"}},
	GetresuidEventID:              {{Type: "uid_t*", Name: "ruid"}, {Type: "uid_t*", Name: "euid"}, {Type: "uid_t*", Name: "suid"}},
	SetresgidEventID:              {{Type: "gid_t", Name: "rgid"}, {Type: "gid_t", Name: "egid"}, {Type: "gid_t", Name: "sgid"}},
	GetresgidEventID:              {{Type: "gid_t*", Name: "rgid"}, {Type: "gid_t*", Name: "egid"}, {Type: "gid_t*", Name: "sgid"}},
	GetpgidEventID:                {{Type: "pid_t", Name: "pid"}},
	SetfsuidEventID:               {{Type: "uid_t", Name: "fsuid"}},
	SetfsgidEventID:               {{Type: "gid_t", Name: "fsgid"}},
	GetsidEventID:                 {{Type: "pid_t", Name: "pid"}},
	CapgetEventID:                 {{Type: "cap_user_header_t", Name: "hdrp"}, {Type: "cap_user_data_t", Name: "datap"}},
	CapsetEventID:                 {{Type: "cap_user_header_t", Name: "hdrp"}, {Type: "const cap_user_data_t", Name: "datap"}},
	RtSigpendingEventID:           {{Type: "sigset_t*", Name: "set"}, {Type: "size_t", Name: "sigsetsize"}},
	RtSigtimedwaitEventID:         {{Type: "const sigset_t*", Name: "set"}, {Type: "siginfo_t*", Name: "info"}, {Type: "const struct timespec*", Name: "timeout"}, {Type: "size_t", Name: "sigsetsize"}},
	RtSigqueueinfoEventID:         {{Type: "pid_t", Name: "tgid"}, {Type: "int", Name: "sig"}, {Type: "siginfo_t*", Name: "info"}},
	RtSigsuspendEventID:           {{Type: "sigset_t*", Name: "mask"}, {Type: "size_t", Name: "sigsetsize"}},
	SigaltstackEventID:            {{Type: "const stack_t*", Name: "ss"}, {Type: "stack_t*", Name: "old_ss"}},
	UtimeEventID:                  {{Type: "const char*", Name: "filename"}, {Type: "const struct utimbuf*", Name: "times"}},
	MknodEventID:                  {{Type: "const char*", Name: "pathname"}, {Type: "mode_t", Name: "mode"}, {Type: "dev_t", Name: "dev"}},
	UselibEventID:                 {{Type: "const char*", Name: "library"}},
	PersonalityEventID:            {{Type: "unsigned long", Name: "persona"}},
	UstatEventID:                  {{Type: "dev_t", Name: "dev"}, {Type: "struct ustat*", Name: "ubuf"}},
	StatfsEventID:                 {{Type: "const char*", Name: "path"}, {Type: "struct statfs*", Name: "buf"}},
	FstatfsEventID:                {{Type: "int", Name: "fd"}, {Type: "struct statfs*", Name: "buf"}},
	SysfsEventID:                  {{Type: "int", Name: "option"}},
	GetpriorityEventID:            {{Type: "int", Name: "which"}, {Type: "int", Name: "who"}},
	SetpriorityEventID:            {{Type: "int", Name: "which"}, {Type: "int", Name: "who"}, {Type: "int", Name: "prio"}},
	SchedSetparamEventID:          {{Type: "pid_t", Name: "pid"}, {Type: "struct sched_param*", Name: "param"}},
	SchedGetparamEventID:          {{Type: "pid_t", Name: "pid"}, {Type: "struct sched_param*", Name: "param"}},
	SchedSetschedulerEventID:      {{Type: "pid_t", Name: "pid"}, {Type: "int", Name: "policy"}, {Type: "struct sched_param*", Name: "param"}},
	SchedGetschedulerEventID:      {{Type: "pid_t", Name: "pid"}},
	SchedGetPriorityMaxEventID:    {{Type: "int", Name: "policy"}},
	SchedGetPriorityMinEventID:    {{Type: "int", Name: "policy"}},
	SchedRrGetIntervalEventID:     {{Type: "pid_t", Name: "pid"}, {Type: "struct timespec*", Name: "tp"}},
	MlockEventID:                  {{Type: "const void*", Name: "addr"}, {Type: "size_t", Name: "len"}},
	MunlockEventID:                {{Type: "const void*", Name: "addr"}, {Type: "size_t", Name: "len"}},
	MlockallEventID:               {{Type: "int", Name: "flags"}},
	MunlockallEventID:             {},
	VhangupEventID:                {},
	ModifyLdtEventID:              {{Type: "int", Name: "func"}, {Type: "void*", Name: "ptr"}, {Type: "unsigned long", Name: "bytecount"}},
	PivotRootEventID:              {{Type: "const char*", Name: "new_root"}, {Type: "const char*", Name: "put_old"}},
	SysctlEventID:                 {{Type: "struct __sysctl_args*", Name: "args"}},
	PrctlEventID:                  {{Type: "int", Name: "option"}, {Type: "unsigned long", Name: "arg2"}, {Type: "unsigned long", Name: "arg3"}, {Type: "unsigned long", Name: "arg4"}, {Type: "unsigned long", Name: "arg5"}},
	ArchPrctlEventID:              {{Type: "int", Name: "option"}, {Type: "unsigned long", Name: "addr"}},
	AdjtimexEventID:               {{Type: "struct timex*", Name: "buf"}},
	SetrlimitEventID:              {{Type: "int", Name: "resource"}, {Type: "const struct rlimit*", Name: "rlim"}},
	ChrootEventID:                 {{Type: "const char*", Name: "path"}},
	SyncEventID:                   {},
	AcctEventID:                   {{Type: "const char*", Name: "filename"}},
	SettimeofdayEventID:           {{Type: "const struct timeval*", Name: "tv"}, {Type: "const struct timezone*", Name: "tz"}},
	MountEventID:                  {{Type: "const char*", Name: "source"}, {Type: "const char*", Name: "target"}, {Type: "const char*", Name: "filesystemtype"}, {Type: "unsigned long", Name: "mountflags"}, {Type: "const void*", Name: "data"}},
	UmountEventID:                 {{Type: "const char*", Name: "target"}, {Type: "int", Name: "flags"}},
	SwaponEventID:                 {{Type: "const char*", Name: "path"}, {Type: "int", Name: "swapflags"}},
	SwapoffEventID:                {{Type: "const char*", Name: "path"}},
	RebootEventID:                 {{Type: "int", Name: "magic"}, {Type: "int", Name: "magic2"}, {Type: "int", Name: "cmd"}, {Type: "void*", Name: "arg"}},
	SethostnameEventID:            {{Type: "const char*", Name: "name"}, {Type: "size_t", Name: "len"}},
	SetdomainnameEventID:          {{Type: "const char*", Name: "name"}, {Type: "size_t", Name: "len"}},
	IoplEventID:                   {{Type: "int", Name: "level"}},
	IopermEventID:                 {{Type: "unsigned long", Name: "from"}, {Type: "unsigned long", Name: "num"}, {Type: "int", Name: "turn_on"}},
	InitModuleEventID:             {{Type: "void*", Name: "module_image"}, {Type: "unsigned long", Name: "len"}, {Type: "const char*", Name: "param_values"}},
	DeleteModuleEventID:           {{Type: "const char*", Name: "name"}, {Type: "int", Name: "flags"}},
	QuotactlEventID:               {{Type: "int", Name: "cmd"}, {Type: "const char*", Name: "special"}, {Type: "int", Name: "id"}, {Type: "void*", Name: "addr"}},
	GettidEventID:                 {},
	ReadaheadEventID:              {{Type: "int", Name: "fd"}, {Type: "off_t", Name: "offset"}, {Type: "size_t", Name: "count"}},
	SetxattrEventID:               {{Type: "const char*", Name: "path"}, {Type: "const char*", Name: "name"}, {Type: "const void*", Name: "value"}, {Type: "size_t", Name: "size"}, {Type: "int", Name: "flags"}},
	LsetxattrEventID:              {{Type: "const char*", Name: "path"}, {Type: "const char*", Name: "name"}, {Type: "const void*", Name: "value"}, {Type: "size_t", Name: "size"}, {Type: "int", Name: "flags"}},
	FsetxattrEventID:              {{Type: "int", Name: "fd"}, {Type: "const char*", Name: "name"}, {Type: "const void*", Name: "value"}, {Type: "size_t", Name: "size"}, {Type: "int", Name: "flags"}},
	GetxattrEventID:               {{Type: "const char*", Name: "path"}, {Type: "const char*", Name: "name"}, {Type: "void*", Name: "value"}, {Type: "size_t", Name: "size"}},
	LgetxattrEventID:              {{Type: "const char*", Name: "path"}, {Type: "const char*", Name: "name"}, {Type: "void*", Name: "value"}, {Type: "size_t", Name: "size"}},
	FgetxattrEventID:              {{Type: "int", Name: "fd"}, {Type: "const char*", Name: "name"}, {Type: "void*", Name: "value"}, {Type: "size_t", Name: "size"}},
	ListxattrEventID:              {{Type: "const char*", Name: "path"}, {Type: "char*", Name: "list"}, {Type: "size_t", Name: "size"}},
	LlistxattrEventID:             {{Type: "const char*", Name: "path"}, {Type: "char*", Name: "list"}, {Type: "size_t", Name: "size"}},
	FlistxattrEventID:             {{Type: "int", Name: "fd"}, {Type: "char*", Name: "list"}, {Type: "size_t", Name: "size"}},
	RemovexattrEventID:            {{Type: "const char*", Name: "path"}, {Type: "const char*", Name: "name"}},
	LremovexattrEventID:           {{Type: "const char*", Name: "path"}, {Type: "const char*", Name: "name"}},
	FremovexattrEventID:           {{Type: "int", Name: "fd"}, {Type: "const char*", Name: "name"}},
	TkillEventID:                  {{Type: "int", Name: "tid"}, {Type: "int", Name: "sig"}},
	TimeEventID:                   {{Type: "time_t*", Name: "tloc"}},
	FutexEventID:                  {{Type: "int*", Name: "uaddr"}, {Type: "int", Name: "futex_op"}, {Type: "int", Name: "val"}, {Type: "const struct timespec*", Name: "timeout"}, {Type: "int*", Name: "uaddr2"}, {Type: "int", Name: "val3"}},
	SchedSetaffinityEventID:       {{Type: "pid_t", Name: "pid"}, {Type: "size_t", Name: "cpusetsize"}, {Type: "unsigned long*", Name: "mask"}},
	SchedGetaffinityEventID:       {{Type: "pid_t", Name: "pid"}, {Type: "size_t", Name: "cpusetsize"}, {Type: "unsigned long*", Name: "mask"}},
	SetThreadAreaEventID:          {{Type: "struct user_desc*", Name: "u_info"}},
	IoSetupEventID:                {{Type: "unsigned int", Name: "nr_events"}, {Type: "io_context_t*", Name: "ctx_idp"}},
	IoDestroyEventID:              {{Type: "io_context_t", Name: "ctx_id"}},
	IoGeteventsEventID:            {{Type: "io_context_t", Name: "ctx_id"}, {Type: "long", Name: "min_nr"}, {Type: "long", Name: "nr"}, {Type: "struct io_event*", Name: "events"}, {Type: "struct timespec*", Name: "timeout"}},
	IoSubmitEventID:               {{Type: "io_context_t", Name: "ctx_id"}, {Type: "long", Name: "nr"}, {Type: "struct iocb**", Name: "iocbpp"}},
	IoCancelEventID:               {{Type: "io_context_t", Name: "ctx_id"}, {Type: "struct iocb*", Name: "iocb"}, {Type: "struct io_event*", Name: "result"}},
	GetThreadAreaEventID:          {{Type: "struct user_desc*", Name: "u_info"}},
	LookupDcookieEventID:          {{Type: "u64", Name: "cookie"}, {Type: "char*", Name: "buffer"}, {Type: "size_t", Name: "len"}},
	EpollCreateEventID:            {{Type: "int", Name: "size"}},
	RemapFilePagesEventID:         {{Type: "void*", Name: "addr"}, {Type: "size_t", Name: "size"}, {Type: "int", Name: "prot"}, {Type: "size_t", Name: "pgoff"}, {Type: "int", Name: "flags"}},
	Getdents64EventID:             {{Type: "unsigned int", Name: "fd"}, {Type: "struct linux_dirent64*", Name: "dirp"}, {Type: "unsigned int", Name: "count"}},
	SetTidAddressEventID:          {{Type: "int*", Name: "tidptr"}},
	RestartSyscallEventID:         {},
	SemtimedopEventID:             {{Type: "int", Name: "semid"}, {Type: "struct sembuf*", Name: "sops"}, {Type: "size_t", Name: "nsops"}, {Type: "const struct timespec*", Name: "timeout"}},
	Fadvise64EventID:              {{Type: "int", Name: "fd"}, {Type: "off_t", Name: "offset"}, {Type: "size_t", Name: "len"}, {Type: "int", Name: "advice"}},
	TimerCreateEventID:            {{Type: "const clockid_t", Name: "clockid"}, {Type: "struct sigevent*", Name: "sevp"}, {Type: "timer_t*", Name: "timer_id"}},
	TimerSettimeEventID:           {{Type: "timer_t", Name: "timer_id"}, {Type: "int", Name: "flags"}, {Type: "const struct itimerspec*", Name: "new_value"}, {Type: "struct itimerspec*", Name: "old_value"}},
	TimerGettimeEventID:           {{Type: "timer_t", Name: "timer_id"}, {Type: "struct itimerspec*", Name: "curr_value"}},
	TimerGetoverrunEventID:        {{Type: "timer_t", Name: "timer_id"}},
	TimerDeleteEventID:            {{Type: "timer_t", Name: "timer_id"}},
	ClockSettimeEventID:           {{Type: "const clockid_t", Name: "clockid"}, {Type: "const struct timespec*", Name: "tp"}},
	ClockGettimeEventID:           {{Type: "const clockid_t", Name: "clockid"}, {Type: "struct timespec*", Name: "tp"}},
	ClockGetresEventID:            {{Type: "const clockid_t", Name: "clockid"}, {Type: "struct timespec*", Name: "res"}},
	ClockNanosleepEventID:         {{Type: "const clockid_t", Name: "clockid"}, {Type: "int", Name: "flags"}, {Type: "const struct timespec*", Name: "request"}, {Type: "struct timespec*", Name: "remain"}},
	ExitGroupEventID:              {{Type: "int", Name: "status"}},
	EpollWaitEventID:              {{Type: "int", Name: "epfd"}, {Type: "struct epoll_event*", Name: "events"}, {Type: "int", Name: "maxevents"}, {Type: "int", Name: "timeout"}},
	EpollCtlEventID:               {{Type: "int", Name: "epfd"}, {Type: "int", Name: "op"}, {Type: "int", Name: "fd"}, {Type: "struct epoll_event*", Name: "event"}},
	TgkillEventID:                 {{Type: "int", Name: "tgid"}, {Type: "int", Name: "tid"}, {Type: "int", Name: "sig"}},
	UtimesEventID:                 {{Type: "char*", Name: "filename"}, {Type: "struct timeval*", Name: "times"}},
	MbindEventID:                  {{Type: "void*", Name: "addr"}, {Type: "unsigned long", Name: "len"}, {Type: "int", Name: "mode"}, {Type: "const unsigned long*", Name: "nodemask"}, {Type: "unsigned long", Name: "maxnode"}, {Type: "unsigned int", Name: "flags"}},
	SetMempolicyEventID:           {{Type: "int", Name: "mode"}, {Type: "const unsigned long*", Name: "nodemask"}, {Type: "unsigned long", Name: "maxnode"}},
	GetMempolicyEventID:           {{Type: "int*", Name: "mode"}, {Type: "unsigned long*", Name: "nodemask"}, {Type: "unsigned long", Name: "maxnode"}, {Type: "void*", Name: "addr"}, {Type: "unsigned long", Name: "flags"}},
	MqOpenEventID:                 {{Type: "const char*", Name: "name"}, {Type: "int", Name: "oflag"}, {Type: "mode_t", Name: "mode"}, {Type: "struct mq_attr*", Name: "attr"}},
	MqUnlinkEventID:               {{Type: "const char*", Name: "name"}},
	MqTimedsendEventID:            {{Type: "mqd_t", Name: "mqdes"}, {Type: "const char*", Name: "msg_ptr"}, {Type: "size_t", Name: "msg_len"}, {Type: "unsigned int", Name: "msg_prio"}, {Type: "const struct timespec*", Name: "abs_timeout"}},
	MqTimedreceiveEventID:         {{Type: "mqd_t", Name: "mqdes"}, {Type: "char*", Name: "msg_ptr"}, {Type: "size_t", Name: "msg_len"}, {Type: "unsigned int*", Name: "msg_prio"}, {Type: "const struct timespec*", Name: "abs_timeout"}},
	MqNotifyEventID:               {{Type: "mqd_t", Name: "mqdes"}, {Type: "const struct sigevent*", Name: "sevp"}},
	MqGetsetattrEventID:           {{Type: "mqd_t", Name: "mqdes"}, {Type: "const struct mq_attr*", Name: "newattr"}, {Type: "struct mq_attr*", Name: "oldattr"}},
	KexecLoadEventID:              {{Type: "unsigned long", Name: "entry"}, {Type: "unsigned long", Name: "nr_segments"}, {Type: "struct kexec_segment*", Name: "segments"}, {Type: "unsigned long", Name: "flags"}},
	WaitidEventID:                 {{Type: "int", Name: "idtype"}, {Type: "pid_t", Name: "id"}, {Type: "struct siginfo*", Name: "infop"}, {Type: "int", Name: "options"}, {Type: "struct rusage*", Name: "rusage"}},
	AddKeyEventID:                 {{Type: "const char*", Name: "type"}, {Type: "const char*", Name: "description"}, {Type: "const void*", Name: "payload"}, {Type: "size_t", Name: "plen"}, {Type: "key_serial_t", Name: "keyring"}},
	RequestKeyEventID:             {{Type: "const char*", Name: "type"}, {Type: "const char*", Name: "description"}, {Type: "const char*", Name: "callout_info"}, {Type: "key_serial_t", Name: "dest_keyring"}},
	KeyctlEventID:                 {{Type: "int", Name: "operation"}, {Type: "unsigned long", Name: "arg2"}, {Type: "unsigned long", Name: "arg3"}, {Type: "unsigned long", Name: "arg4"}, {Type: "unsigned long", Name: "arg5"}},
	IoprioSetEventID:              {{Type: "int", Name: "which"}, {Type: "int", Name: "who"}, {Type: "int", Name: "ioprio"}},
	IoprioGetEventID:              {{Type: "int", Name: "which"}, {Type: "int", Name: "who"}},
	InotifyInitEventID:            {},
	InotifyAddWatchEventID:        {{Type: "int", Name: "fd"}, {Type: "const char*", Name: "pathname"}, {Type: "u32", Name: "mask"}},
	InotifyRmWatchEventID:         {{Type: "int", Name: "fd"}, {Type: "int", Name: "wd"}},
	MigratePagesEventID:           {{Type: "int", Name: "pid"}, {Type: "unsigned long", Name: "maxnode"}, {Type: "const unsigned long*", Name: "old_nodes"}, {Type: "const unsigned long*", Name: "new_nodes"}},
	OpenatEventID:                 {{Type: "int", Name: "dirfd"}, {Type: "const char*", Name: "pathname"}, {Type: "int", Name: "flags"}, {Type: "mode_t", Name: "mode"}},
	MkdiratEventID:                {{Type: "int", Name: "dirfd"}, {Type: "const char*", Name: "pathname"}, {Type: "mode_t", Name: "mode"}},
	MknodatEventID:                {{Type: "int", Name: "dirfd"}, {Type: "const char*", Name: "pathname"}, {Type: "mode_t", Name: "mode"}, {Type: "dev_t", Name: "dev"}},
	FchownatEventID:               {{Type: "int", Name: "dirfd"}, {Type: "const char*", Name: "pathname"}, {Type: "uid_t", Name: "owner"}, {Type: "gid_t", Name: "group"}, {Type: "int", Name: "flags"}},
	FutimesatEventID:              {{Type: "int", Name: "dirfd"}, {Type: "const char*", Name: "pathname"}, {Type: "struct timeval*", Name: "times"}},
	NewfstatatEventID:             {{Type: "int", Name: "dirfd"}, {Type: "const char*", Name: "pathname"}, {Type: "struct stat*", Name: "statbuf"}, {Type: "int", Name: "flags"}},
	UnlinkatEventID:               {{Type: "int", Name: "dirfd"}, {Type: "const char*", Name: "pathname"}, {Type: "int", Name: "flags"}},
	RenameatEventID:               {{Type: "int", Name: "olddirfd"}, {Type: "const char*", Name: "oldpath"}, {Type: "int", Name: "newdirfd"}, {Type: "const char*", Name: "newpath"}},
	LinkatEventID:                 {{Type: "int", Name: "olddirfd"}, {Type: "const char*", Name: "oldpath"}, {Type: "int", Name: "newdirfd"}, {Type: "const char*", Name: "newpath"}, {Type: "unsigned int", Name: "flags"}},
	SymlinkatEventID:              {{Type: "const char*", Name: "target"}, {Type: "int", Name: "newdirfd"}, {Type: "const char*", Name: "linkpath"}},
	ReadlinkatEventID:             {{Type: "int", Name: "dirfd"}, {Type: "const char*", Name: "pathname"}, {Type: "char*", Name: "buf"}, {Type: "int", Name: "bufsiz"}},
	FchmodatEventID:               {{Type: "int", Name: "dirfd"}, {Type: "const char*", Name: "pathname"}, {Type: "mode_t", Name: "mode"}, {Type: "int", Name: "flags"}},
	FaccessatEventID:              {{Type: "int", Name: "dirfd"}, {Type: "const char*", Name: "pathname"}, {Type: "int", Name: "mode"}, {Type: "int", Name: "flags"}},
	Pselect6EventID:               {{Type: "int", Name: "nfds"}, {Type: "fd_set*", Name: "readfds"}, {Type: "fd_set*", Name: "writefds"}, {Type: "fd_set*", Name: "exceptfds"}, {Type: "struct timespec*", Name: "timeout"}, {Type: "void*", Name: "sigmask"}},
	PpollEventID:                  {{Type: "struct pollfd*", Name: "fds"}, {Type: "unsigned int", Name: "nfds"}, {Type: "struct timespec*", Name: "tmo_p"}, {Type: "const sigset_t*", Name: "sigmask"}, {Type: "size_t", Name: "sigsetsize"}},
	UnshareEventID:                {{Type: "int", Name: "flags"}},
	SetRobustListEventID:          {{Type: "struct robust_list_head*", Name: "head"}, {Type: "size_t", Name: "len"}},
	GetRobustListEventID:          {{Type: "int", Name: "pid"}, {Type: "struct robust_list_head**", Name: "head_ptr"}, {Type: "size_t*", Name: "len_ptr"}},
	SpliceEventID:                 {{Type: "int", Name: "fd_in"}, {Type: "off_t*", Name: "off_in"}, {Type: "int", Name: "fd_out"}, {Type: "off_t*", Name: "off_out"}, {Type: "size_t", Name: "len"}, {Type: "unsigned int", Name: "flags"}},
	TeeEventID:                    {{Type: "int", Name: "fd_in"}, {Type: "int", Name: "fd_out"}, {Type: "size_t", Name: "len"}, {Type: "unsigned int", Name: "flags"}},
	SyncFileRangeEventID:          {{Type: "int", Name: "fd"}, {Type: "off_t", Name: "offset"}, {Type: "off_t", Name: "nbytes"}, {Type: "unsigned int", Name: "flags"}},
	VmspliceEventID:               {{Type: "int", Name: "fd"}, {Type: "const struct iovec*", Name: "iov"}, {Type: "unsigned long", Name: "nr_segs"}, {Type: "unsigned int", Name: "flags"}},
	MovePagesEventID:              {{Type: "int", Name: "pid"}, {Type: "unsigned long", Name: "count"}, {Type: "const void**", Name: "pages"}, {Type: "const int*", Name: "nodes"}, {Type: "int*", Name: "status"}, {Type: "int", Name: "flags"}},
	UtimensatEventID:              {{Type: "int", Name: "dirfd"}, {Type: "const char*", Name: "pathname"}, {Type: "struct timespec*", Name: "times"}, {Type: "int", Name: "flags"}},
	EpollPwaitEventID:             {{Type: "int", Name: "epfd"}, {Type: "struct epoll_event*", Name: "events"}, {Type: "int", Name: "maxevents"}, {Type: "int", Name: "timeout"}, {Type: "const sigset_t*", Name: "sigmask"}, {Type: "size_t", Name: "sigsetsize"}},
	SignalfdEventID:               {{Type: "int", Name: "fd"}, {Type: "sigset_t*", Name: "mask"}, {Type: "int", Name: "flags"}},
	TimerfdCreateEventID:          {{Type: "int", Name: "clockid"}, {Type: "int", Name: "flags"}},
	EventfdEventID:                {{Type: "unsigned int", Name: "initval"}, {Type: "int", Name: "flags"}},
	FallocateEventID:              {{Type: "int", Name: "fd"}, {Type: "int", Name: "mode"}, {Type: "off_t", Name: "offset"}, {Type: "off_t", Name: "len"}},
	TimerfdSettimeEventID:         {{Type: "int", Name: "fd"}, {Type: "int", Name: "flags"}, {Type: "const struct itimerspec*", Name: "new_value"}, {Type: "struct itimerspec*", Name: "old_value"}},
	TimerfdGettimeEventID:         {{Type: "int", Name: "fd"}, {Type: "struct itimerspec*", Name: "curr_value"}},
	Accept4EventID:                {{Type: "int", Name: "sockfd"}, {Type: "struct sockaddr*", Name: "addr"}, {Type: "int*", Name: "addrlen"}, {Type: "int", Name: "flags"}},
	Signalfd4EventID:              {{Type: "int", Name: "fd"}, {Type: "const sigset_t*", Name: "mask"}, {Type: "size_t", Name: "sizemask"}, {Type: "int", Name: "flags"}},
	Eventfd2EventID:               {{Type: "unsigned int", Name: "initval"}, {Type: "int", Name: "flags"}},
	EpollCreate1EventID:           {{Type: "int", Name: "flags"}},
	Dup3EventID:                   {{Type: "int", Name: "oldfd"}, {Type: "int", Name: "newfd"}, {Type: "int", Name: "flags"}},
	Pipe2EventID:                  {{Type: "int[2]", Name: "pipefd"}, {Type: "int", Name: "flags"}},
	InotifyInit1EventID:           {{Type: "int", Name: "flags"}},
	PreadvEventID:                 {{Type: "int", Name: "fd"}, {Type: "const struct iovec*", Name: "iov"}, {Type: "unsigned long", Name: "iovcnt"}, {Type: "unsigned long", Name: "pos_l"}, {Type: "unsigned long", Name: "pos_h"}},
	PwritevEventID:                {{Type: "int", Name: "fd"}, {Type: "const struct iovec*", Name: "iov"}, {Type: "unsigned long", Name: "iovcnt"}, {Type: "unsigned long", Name: "pos_l"}, {Type: "unsigned long", Name: "pos_h"}},
	RtTgsigqueueinfoEventID:       {{Type: "pid_t", Name: "tgid"}, {Type: "pid_t", Name: "tid"}, {Type: "int", Name: "sig"}, {Type: "siginfo_t*", Name: "info"}},
	PerfEventOpenEventID:          {{Type: "struct perf_event_attr*", Name: "attr"}, {Type: "pid_t", Name: "pid"}, {Type: "int", Name: "cpu"}, {Type: "int", Name: "group_fd"}, {Type: "unsigned long", Name: "flags"}},
	RecvmmsgEventID:               {{Type: "int", Name: "sockfd"}, {Type: "struct mmsghdr*", Name: "msgvec"}, {Type: "unsigned int", Name: "vlen"}, {Type: "int", Name: "flags"}, {Type: "struct timespec*", Name: "timeout"}},
	FanotifyInitEventID:           {{Type: "unsigned int", Name: "flags"}, {Type: "unsigned int", Name: "event_f_flags"}},
	FanotifyMarkEventID:           {{Type: "int", Name: "fanotify_fd"}, {Type: "unsigned int", Name: "flags"}, {Type: "u64", Name: "mask"}, {Type: "int", Name: "dirfd"}, {Type: "const char*", Name: "pathname"}},
	Prlimit64EventID:              {{Type: "pid_t", Name: "pid"}, {Type: "int", Name: "resource"}, {Type: "const struct rlimit64*", Name: "new_limit"}, {Type: "struct rlimit64*", Name: "old_limit"}},
	NameToHandleAtEventID:         {{Type: "int", Name: "dirfd"}, {Type: "const char*", Name: "pathname"}, {Type: "struct file_handle*", Name: "handle"}, {Type: "int*", Name: "mount_id"}, {Type: "int", Name: "flags"}},
	OpenByHandleAtEventID:         {{Type: "int", Name: "mount_fd"}, {Type: "struct file_handle*", Name: "handle"}, {Type: "int", Name: "flags"}},
	ClockAdjtimeEventID:           {{Type: "const clockid_t", Name: "clk_id"}, {Type: "struct timex*", Name: "buf"}},
	SyncfsEventID:                 {{Type: "int", Name: "fd"}},
	SendmmsgEventID:               {{Type: "int", Name: "sockfd"}, {Type: "struct mmsghdr*", Name: "msgvec"}, {Type: "unsigned int", Name: "vlen"}, {Type: "int", Name: "flags"}},
	SetnsEventID:                  {{Type: "int", Name: "fd"}, {Type: "int", Name: "nstype"}},
	GetcpuEventID:                 {{Type: "unsigned int*", Name: "cpu"}, {Type: "unsigned int*", Name: "node"}, {Type: "struct getcpu_cache*", Name: "tcache"}},
	ProcessVmReadvEventID:         {{Type: "pid_t", Name: "pid"}, {Type: "const struct iovec*", Name: "local_iov"}, {Type: "unsigned long", Name: "liovcnt"}, {Type: "const struct iovec*", Name: "remote_iov"}, {Type: "unsigned long", Name: "riovcnt"}, {Type: "unsigned long", Name: "flags"}},
	ProcessVmWritevEventID:        {{Type: "pid_t", Name: "pid"}, {Type: "const struct iovec*", Name: "local_iov"}, {Type: "unsigned long", Name: "liovcnt"}, {Type: "const struct iovec*", Name: "remote_iov"}, {Type: "unsigned long", Name: "riovcnt"}, {Type: "unsigned long", Name: "flags"}},
	KcmpEventID:                   {{Type: "pid_t", Name: "pid1"}, {Type: "pid_t", Name: "pid2"}, {Type: "int", Name: "type"}, {Type: "unsigned long", Name: "idx1"}, {Type: "unsigned long", Name: "idx2"}},
	FinitModuleEventID:            {{Type: "int", Name: "fd"}, {Type: "const char*", Name: "param_values"}, {Type: "int", Name: "flags"}},
	SchedSetattrEventID:           {{Type: "pid_t", Name: "pid"}, {Type: "struct sched_attr*", Name: "attr"}, {Type: "unsigned int", Name: "flags"}},
	SchedGetattrEventID:           {{Type: "pid_t", Name: "pid"}, {Type: "struct sched_attr*", Name: "attr"}, {Type: "unsigned int", Name: "size"}, {Type: "unsigned int", Name: "flags"}},
	Renameat2EventID:              {{Type: "int", Name: "olddirfd"}, {Type: "const char*", Name: "oldpath"}, {Type: "int", Name: "newdirfd"}, {Type: "const char*", Name: "newpath"}, {Type: "unsigned int", Name: "flags"}},
	SeccompEventID:                {{Type: "unsigned int", Name: "operation"}, {Type: "unsigned int", Name: "flags"}, {Type: "const void*", Name: "args"}},
	GetrandomEventID:              {{Type: "void*", Name: "buf"}, {Type: "size_t", Name: "buflen"}, {Type: "unsigned int", Name: "flags"}},
	MemfdCreateEventID:            {{Type: "const char*", Name: "name"}, {Type: "unsigned int", Name: "flags"}},
	KexecFileLoadEventID:          {{Type: "int", Name: "kernel_fd"}, {Type: "int", Name: "initrd_fd"}, {Type: "unsigned long", Name: "cmdline_len"}, {Type: "const char*", Name: "cmdline"}, {Type: "unsigned long", Name: "flags"}},
	BpfEventID:                    {{Type: "int", Name: "cmd"}, {Type: "union bpf_attr*", Name: "attr"}, {Type: "unsigned int", Name: "size"}},
	ExecveatEventID:               {{Type: "int", Name: "dirfd"}, {Type: "const char*", Name: "pathname"}, {Type: "const char*const*", Name: "argv"}, {Type: "const char*const*", Name: "envp"}, {Type: "int", Name: "flags"}},
	UserfaultfdEventID:            {{Type: "int", Name: "flags"}},
	MembarrierEventID:             {{Type: "int", Name: "cmd"}, {Type: "int", Name: "flags"}},
	Mlock2EventID:                 {{Type: "const void*", Name: "addr"}, {Type: "size_t", Name: "len"}, {Type: "int", Name: "flags"}},
	CopyFileRangeEventID:          {{Type: "int", Name: "fd_in"}, {Type: "off_t*", Name: "off_in"}, {Type: "int", Name: "fd_out"}, {Type: "off_t*", Name: "off_out"}, {Type: "size_t", Name: "len"}, {Type: "unsigned int", Name: "flags"}},
	Preadv2EventID:                {{Type: "int", Name: "fd"}, {Type: "const struct iovec*", Name: "iov"}, {Type: "unsigned long", Name: "iovcnt"}, {Type: "unsigned long", Name: "pos_l"}, {Type: "unsigned long", Name: "pos_h"}, {Type: "int", Name: "flags"}},
	Pwritev2EventID:               {{Type: "int", Name: "fd"}, {Type: "const struct iovec*", Name: "iov"}, {Type: "unsigned long", Name: "iovcnt"}, {Type: "unsigned long", Name: "pos_l"}, {Type: "unsigned long", Name: "pos_h"}, {Type: "int", Name: "flags"}},
	PkeyMprotectEventID:           {{Type: "void*", Name: "addr"}, {Type: "size_t", Name: "len"}, {Type: "int", Name: "prot"}, {Type: "int", Name: "pkey"}},
	PkeyAllocEventID:              {{Type: "unsigned int", Name: "flags"}, {Type: "unsigned long", Name: "access_rights"}},
	PkeyFreeEventID:               {{Type: "int", Name: "pkey"}},
	StatxEventID:                  {{Type: "int", Name: "dirfd"}, {Type: "const char*", Name: "pathname"}, {Type: "int", Name: "flags"}, {Type: "unsigned int", Name: "mask"}, {Type: "struct statx*", Name: "statxbuf"}},
	IoPgeteventsEventID:           {{Type: "aio_context_t", Name: "ctx_id"}, {Type: "long", Name: "min_nr"}, {Type: "long", Name: "nr"}, {Type: "struct io_event*", Name: "events"}, {Type: "struct timespec*", Name: "timeout"}, {Type: "const struct __aio_sigset*", Name: "usig"}},
	RseqEventID:                   {{Type: "struct rseq*", Name: "rseq"}, {Type: "u32", Name: "rseq_len"}, {Type: "int", Name: "flags"}, {Type: "u32", Name: "sig"}},
	PidfdSendSignalEventID:        {{Type: "int", Name: "pidfd"}, {Type: "int", Name: "sig"}, {Type: "siginfo_t*", Name: "info"}, {Type: "unsigned int", Name: "flags"}},
	IoUringSetupEventID:           {{Type: "unsigned int", Name: "entries"}, {Type: "struct io_uring_params*", Name: "p"}},
	IoUringEnterEventID:           {{Type: "unsigned int", Name: "fd"}, {Type: "unsigned int", Name: "to_submit"}, {Type: "unsigned int", Name: "min_complete"}, {Type: "unsigned int", Name: "flags"}, {Type: "sigset_t*", Name: "sig"}},
	IoUringRegisterEventID:        {{Type: "unsigned int", Name: "fd"}, {Type: "unsigned int", Name: "opcode"}, {Type: "void*", Name: "arg"}, {Type: "unsigned int", Name: "nr_args"}},
	OpenTreeEventID:               {{Type: "int", Name: "dfd"}, {Type: "const char*", Name: "filename"}, {Type: "unsigned int", Name: "flags"}},
	MoveMountEventID:              {{Type: "int", Name: "from_dfd"}, {Type: "const char*", Name: "from_path"}, {Type: "int", Name: "to_dfd"}, {Type: "const char*", Name: "to_path"}, {Type: "unsigned int", Name: "flags"}},
	FsopenEventID:                 {{Type: "const char*", Name: "fsname"}, {Type: "unsigned int", Name: "flags"}},
	FsconfigEventID:               {{Type: "int*", Name: "fs_fd"}, {Type: "unsigned int", Name: "cmd"}, {Type: "const char*", Name: "key"}, {Type: "const void*", Name: "value"}, {Type: "int", Name: "aux"}},
	FsmountEventID:                {{Type: "int", Name: "fsfd"}, {Type: "unsigned int", Name: "flags"}, {Type: "unsigned int", Name: "ms_flags"}},
	FspickEventID:                 {{Type: "int", Name: "dirfd"}, {Type: "const char*", Name: "pathname"}, {Type: "unsigned int", Name: "flags"}},
	PidfdOpenEventID:              {{Type: "pid_t", Name: "pid"}, {Type: "unsigned int", Name: "flags"}},
	Clone3EventID:                 {{Type: "struct clone_args*", Name: "cl_args"}, {Type: "size_t", Name: "size"}},
	CloseRangeEventID:             {{Type: "unsigned int", Name: "first"}, {Type: "unsigned int", Name: "last"}},
	Openat2EventID:                {{Type: "int", Name: "dirfd"}, {Type: "const char*", Name: "pathname"}, {Type: "struct open_how*", Name: "how"}, {Type: "size_t", Name: "size"}},
	PidfdGetfdEventID:             {{Type: "int", Name: "pidfd"}, {Type: "int", Name: "targetfd"}, {Type: "unsigned int", Name: "flags"}},
	Faccessat2EventID:             {{Type: "int", Name: "fd"}, {Type: "const char*", Name: "path"}, {Type: "int", Name: "mode"}, {Type: "int", Name: "flag"}},
	ProcessMadviseEventID:         {{Type: "int", Name: "pidfd"}, {Type: "void*", Name: "addr"}, {Type: "size_t", Name: "length"}, {Type: "int", Name: "advice"}, {Type: "unsigned long", Name: "flags"}},
	EpollPwait2EventID:            {{Type: "int", Name: "fd"}, {Type: "struct epoll_event*", Name: "events"}, {Type: "int", Name: "maxevents"}, {Type: "const struct timespec*", Name: "timeout"}, {Type: "const sigset_t*", Name: "sigset"}},
	SysEnterEventID:               {{Type: "int", Name: "syscall"}},
	SysExitEventID:                {{Type: "int", Name: "syscall"}},
	SchedProcessForkEventID:       {{Type: "int", Name: "parent_tid"}, {Type: "int", Name: "parent_ns_tid"}, {Type: "int", Name: "child_tid"}, {Type: "int", Name: "child_ns_tid"}},
	SchedProcessExecEventID:       {{Type: "const char*", Name: "cmdpath"}, {Type: "const char*", Name: "pathname"}, {Type: "const char**", Name: "argv"}, {Type: "const char**", Name: "env"}, {Type: "dev_t", Name: "dev"}, {Type: "unsigned long", Name: "inode"}, {Type: "int", Name: "invoked_from_kernel"}, {Type: "unsigned long", Name: "ctime"}},
	SchedProcessExitEventID:       {{Type: "long", Name: "exit_code"}},
	SchedSwitchEventID:            {{Type: "int", Name: "cpu"}, {Type: "int", Name: "prev_tid"}, {Type: "const char*", Name: "prev_comm"}, {Type: "int", Name: "next_tid"}, {Type: "const char*", Name: "next_comm"}},
	DoExitEventID:                 {},
	CapCapableEventID:             {{Type: "int", Name: "cap"}, {Type: "int", Name: "syscall"}},
	VfsWriteEventID:               {{Type: "const char*", Name: "pathname"}, {Type: "dev_t", Name: "dev"}, {Type: "unsigned long", Name: "inode"}, {Type: "size_t", Name: "count"}, {Type: "off_t", Name: "pos"}},
	VfsWritevEventID:              {{Type: "const char*", Name: "pathname"}, {Type: "dev_t", Name: "dev"}, {Type: "unsigned long", Name: "inode"}, {Type: "unsigned long", Name: "vlen"}, {Type: "off_t", Name: "pos"}},
	MemProtAlertEventID:           {{Type: "u32", Name: "alert"}},
	CommitCredsEventID:            {{Type: "slim_cred_t", Name: "old_cred"}, {Type: "slim_cred_t", Name: "new_cred"}, {Type: "int", Name: "syscall"}},
	SwitchTaskNSEventID:           {{Type: "pid_t", Name: "pid"}, {Type: "u32", Name: "new_mnt"}, {Type: "u32", Name: "new_pid"}, {Type: "u32", Name: "new_uts"}, {Type: "u32", Name: "new_ipc"}, {Type: "u32", Name: "new_net"}, {Type: "u32", Name: "new_cgroup"}},
	MagicWriteEventID:             {{Type: "const char*", Name: "pathname"}, {Type: "bytes", Name: "bytes"}, {Type: "dev_t", Name: "dev"}, {Type: "unsigned long", Name: "inode"}},
	CgroupAttachTaskEventID:       {{Type: "const char*", Name: "cgroup_path"}, {Type: "const char*", Name: "comm"}, {Type: "pid_t", Name: "pid"}},
	CgroupMkdirEventID:            {{Type: "u64", Name: "cgroup_id"}, {Type: "const char*", Name: "cgroup_path"}},
	CgroupRmdirEventID:            {{Type: "u64", Name: "cgroup_id"}, {Type: "const char*", Name: "cgroup_path"}},
	SecurityBprmCheckEventID:      {{Type: "const char*", Name: "pathname"}, {Type: "dev_t", Name: "dev"}, {Type: "unsigned long", Name: "inode"}},
	SecurityFileOpenEventID:       {{Type: "const char*", Name: "pathname"}, {Type: "int", Name: "flags"}, {Type: "dev_t", Name: "dev"}, {Type: "unsigned long", Name: "inode"}, {Type: "int", Name: "syscall"}},
	SecurityInodeUnlinkEventID:    {{Type: "const char*", Name: "pathname"}},
	SecuritySocketCreateEventID:   {{Type: "int", Name: "family"}, {Type: "int", Name: "type"}, {Type: "int", Name: "protocol"}, {Type: "int", Name: "kern"}},
	SecuritySocketListenEventID:   {{Type: "int", Name: "sockfd"}, {Type: "struct sockaddr*", Name: "local_addr"}, {Type: "int", Name: "backlog"}},
	SecuritySocketConnectEventID:  {{Type: "int", Name: "sockfd"}, {Type: "struct sockaddr*", Name: "remote_addr"}},
	SecuritySocketAcceptEventID:   {{Type: "int", Name: "sockfd"}, {Type: "struct sockaddr*", Name: "local_addr"}},
	SecuritySocketBindEventID:     {{Type: "int", Name: "sockfd"}, {Type: "struct sockaddr*", Name: "local_addr"}},
	SecuritySbMountEventID:        {{Type: "const char*", Name: "dev_name"}, {Type: "const char*", Name: "path"}, {Type: "const char*", Name: "type"}, {Type: "unsigned long", Name: "flags"}},
	SecurityBPFEventID:            {{Type: "int", Name: "cmd"}},
	SecurityBPFMapEventID:         {{Type: "unsigned int", Name: "map_id"}, {Type: "const char*", Name: "map_name"}},
	SecurityKernelReadFileEventID: {{Type: "const char*", Name: "pathname"}, {Type: "dev_t", Name: "dev"}, {Type: "unsigned long", Name: "inode"}, {Type: "int", Name: "type"}},
	SecurityPostReadFileEventID:   {{Type: "const char*", Name: "pathname"}, {Type: "long", Name: "size"}, {Type: "int", Name: "type"}},
	SecurityInodeMknodEventID:     {{Type: "const char*", Name: "file_name"}, {Type: "umode_t", Name: "mode"}, {Type: "dev_t", Name: "dev"}},
	InitNamespacesEventID:         {{Type: "u32", Name: "cgroup"}, {Type: "u32", Name: "ipc"}, {Type: "u32", Name: "mnt"}, {Type: "u32", Name: "net"}, {Type: "u32", Name: "pid"}, {Type: "u32", Name: "pid_for_children"}, {Type: "u32", Name: "time"}, {Type: "u32", Name: "time_for_children"}, {Type: "u32", Name: "user"}, {Type: "u32", Name: "uts"}},
	SocketDupEventID:              {{Type: "int", Name: "oldfd"}, {Type: "int", Name: "newfd"}, {Type: "struct sockaddr*", Name: "remote_addr"}},
}
