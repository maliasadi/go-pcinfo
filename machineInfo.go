package main

import (
	"os"       // provides a platform-independent interface to operating system functionality
	"fmt"      // implements formatted I/O with functions analogous to C's printf and scanf.
	"net"      // provides a portable interface for network I/O, including TCP/IP, UDP, domain name resolution, and Unix domain sockets.
	"flag"     // implements command-line flag parsing.
	"syscall"  // contains an interface to the low-level operating system primitives.
	"runtime"  // contains operations that interact with Go's runtime system, such as functions to control goroutines.
	"github.com/shirou/gopsutil/cpu"
)

const (
	SYSINFO_LOADS_SCALE   =  1<<16; // 2^16 
	SYSINFO_MEM_UNIT      =  1024 * 1024; // MB
)

// sys_t is go version of  "sysinfo" struct,
// for Linux versions Since 2.3.23 (i386) and  2.3.48 (all architectures).
// Reference: http://man7.org/linux/man-pages/man2/sysinfo.2.html
type sys_t struct {
	UpTime    int64      // Seconds since boot
	Loads     [3]float64  // 1, 5, and 15 minute load averages
	TotalRam  uint64     // Total usable main memory size
	FreeRam   uint64     // Available memory size
	SharedRam uint64     // Amount of shared memory
	BufferRam uint64     // Memory used by buffers
	TotalSwap uint64     // Total swap space size
	FreeSwap  uint64     // Swap space still available
	Procs     uint16     // Number of current processes
	TotalHigh uint64     // Total high memory size
	FreeHigh  uint64     // Available high memory size
}

// like Gethosname() function in unistd.h, returns the hostname.
func GetHostName() string{

	hostN, err := os.Hostname();
	if (err != nil) {
		panic(err);
	}
	
	return hostN;
}

// return GOOS and GOARCH variables as strings.
func GetOSName() (string, string) {
	return runtime.GOOS, runtime.GOARCH;
}

// like sysinfo() in sys/sysinfo.h, returns system information as sys_t type.
func GetSysInfo() sys_t {
	
	sys_obj := sys_t{};
	sys_ref := &syscall.Sysinfo_t{}
	err := syscall.Sysinfo(sys_ref)
	if  err != nil {
		panic(err);
	}
	sys_obj.UpTime = sys_ref.Uptime;
	sys_obj.Loads[0] = float64(sys_ref.Loads[0]) / SYSINFO_LOADS_SCALE;
	sys_obj.Loads[1] = float64(sys_ref.Loads[1]) / SYSINFO_LOADS_SCALE;
	sys_obj.Loads[2] = float64(sys_ref.Loads[2]) / SYSINFO_LOADS_SCALE;
	sys_obj.TotalRam = sys_ref.Totalram / SYSINFO_MEM_UNIT;
	sys_obj.FreeRam = sys_ref.Freeram / SYSINFO_MEM_UNIT;
	sys_obj.SharedRam = sys_ref.Sharedram / SYSINFO_MEM_UNIT;
	sys_obj.BufferRam = sys_ref.Bufferram / SYSINFO_MEM_UNIT;
	sys_obj.TotalSwap = sys_ref.Totalswap / SYSINFO_MEM_UNIT;	
	sys_obj.FreeSwap = sys_ref.Freeswap / SYSINFO_MEM_UNIT;
	sys_obj.Procs = sys_ref.Procs;	
	sys_obj.TotalHigh = sys_ref.Totalhigh / SYSINFO_MEM_UNIT;
	sys_obj.FreeHigh = sys_ref.Freehigh / SYSINFO_MEM_UNIT;

	return sys_obj;
}

// like sysinfo in sys/sysinfo.h, returns the cpu information in /proc/cpuinfo
// NOTE: if PRINT_FLAG != 0 then print all CPUs info. 
// Reference: https://github.com/shirou/gopsutil
func GetCPUInfo() []cpu.InfoStat {

	var cpu_info []cpu.InfoStat;
	cpu_info, err := cpu.Info()
	if err != nil {
		panic (err);
	}

	return cpu_info;
}

// like ifconfig, returns the network interfaces of /sys/class/net
// Reference:  https://golang.org/pkg/net/
func GetNetworkInterface() []net.Interface{

	var ifaces []net.Interface;
	ifaces, err := net.Interfaces();
	if err != nil {
		panic(err);
	}
	
	return ifaces;
}

/////////////////////////////////////////////////////////////////////////////

func PrintHostName() {	
	fmt.Printf("hostname:\t%v\n", GetHostName());
}

func PrintOSName() {
	fmt.Printf("OS name:\t%v\n", runtime.GOOS);
	fmt.Printf("OS arch:\t%v\n", runtime.GOARCH);
}

func PrintSysInfo() {
	sys_obj := GetSysInfo();
	fmt.Printf("sys uptime:\t%d\n", sys_obj.UpTime);
	fmt.Printf("sys load avg:\t%2.2f, %2.2f, %2.2f\n",
		sys_obj.Loads[0],sys_obj.Loads[1],sys_obj.Loads[2]);
	fmt.Printf("sys totalRam:\t%d MB\n", sys_obj.TotalRam);
	fmt.Printf("sys freeRam:\t%d MB\n", sys_obj.TotalRam);
	fmt.Printf("sys sharedRam:\t%d MB\n", sys_obj.SharedRam);
	fmt.Printf("sys bufferRam:\t%d MB\n", sys_obj.BufferRam);
	fmt.Printf("sys totalSwap:\t%d MB\n", sys_obj.TotalSwap);
	fmt.Printf("sys freeSwap:\t%d MB\n", sys_obj.FreeSwap);
	fmt.Printf("sys totalHigh:\t%d MB\n", sys_obj.TotalHigh);
	fmt.Printf("sys freeHigh:\t%d MB\n", sys_obj.FreeHigh);
	fmt.Printf("sys procs:\t%d\n", sys_obj.Procs);
}

func PrintCPUInfo() {
	cpu_info := GetCPUInfo();	
	for i := 0; i < len(cpu_info); i++ {
		fmt.Printf("cpuID:\t%d\n", cpu_info[i].CPU);
		fmt.Printf("--vendorID:\t%s\n",cpu_info[i].VendorID);
		fmt.Printf("--family:\t%s\n", cpu_info[i].Family);
		fmt.Printf("--model:\t%s\n", cpu_info[i].Model);
		fmt.Printf("--stepping:\t%d\n", cpu_info[i].Stepping);
		fmt.Printf("--physicalID:\t%s\n", cpu_info[i].PhysicalID);
		fmt.Printf("--coreID:\t%s\n", cpu_info[i].CoreID);
		fmt.Printf("--cores:\t%d\n", cpu_info[i].Cores);
		fmt.Printf("--modelName:\t%s\n", cpu_info[i].ModelName);
		fmt.Printf("--MHz:\t\t%g\n", cpu_info[i].Mhz);
		fmt.Printf("--facheSize:\t%d\n", cpu_info[i].CacheSize);
		fmt.Printf("--flags:\t%v\n", cpu_info[i].Flags);
		fmt.Printf("--Microcode:\t%s\n", cpu_info[i].Microcode);
	}
}

func PrintNetworkInterface() {
	ifaces := GetNetworkInterface();
	fmt.Printf("interfaces:\t%v\n", ifaces[0]);
	for i := 1; i < len(ifaces); i++ {
		fmt.Printf("\t\t%v\n", ifaces[i]);
	}
}


func main () {

	var hn = flag.Bool("hostname", false, "if 'hostname=true' => print hostname");
	var OS = flag.Bool("os", false, "if 'os=true' => print operating system info");
	var cpu = flag.Bool("cpu", false, "if 'cpu=true' => print cpu info");
	var netw = flag.Bool("network", false,"if 'network=true' => print network interfaces info");

	flag.Parse()
	fmt.Println (len(os.Args));
	if len(os.Args) == 1 {
		PrintHostName();
		PrintNetworkInterface();
		PrintOSName();
		PrintSysInfo();
		PrintCPUInfo();
		return;
	}

	if *hn { PrintHostName(); }
	if *OS { PrintSysInfo(); PrintCPUInfo(); }
	if *cpu { PrintCPUInfo(); }
	if *netw { PrintNetworkInterface(); }
}
