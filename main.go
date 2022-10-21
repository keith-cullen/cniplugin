// https://github.com/containernetworking/cni/blob/main/plugins/debug/main.go

package main

import (
	"encoding/json"
	"fmt"
	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	type100 "github.com/containernetworking/cni/pkg/types/100"
	"github.com/containernetworking/cni/pkg/version"
	bv "github.com/containernetworking/plugins/pkg/utils/buildversion"
	"io"
	"os"
)

type NetConf struct {
	types.NetConf
	CNIOutput  string `json:"cniOutput,omitempty"`
}

func main() {
	skel.PluginMain(cmdAdd, cmdCheck, cmdDel, version.All, bv.BuildString("none"))
}

func parseConf(data []byte) (*NetConf, error) {
	conf := &NetConf{}
	if err := json.Unmarshal(data, &conf); err != nil {
		return nil, fmt.Errorf("Failed to parse config")
	}
	return conf, nil
}

func outputCmdArgs(fp io.Writer, args *skel.CmdArgs) {
	fmt.Fprintf(fp, `ContainerID: %s
Netns: %s
IfName: %s
Args: %s
Path: %s
StdinData: %s
----------------------
`,
		args.ContainerID,
		args.Netns,
		args.IfName,
		args.Args,
		args.Path,
		string(args.StdinData))
}

func getResult(netConf *NetConf) *type100.Result {
	if netConf.RawPrevResult == nil {
		return &type100.Result{}
	}

	version.ParsePrevResult(&netConf.NetConf)
	result, _ := type100.NewResultFromResult(netConf.PrevResult)
	return result
}

func cmdAdd(args *skel.CmdArgs) error {
	netConf, _ := parseConf(args.StdinData)
	if netConf.CNIOutput != "" {
		fp, _ := os.OpenFile(netConf.CNIOutput, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0644)
		defer fp.Close()
		fmt.Fprintf(fp, "CmdAdd\n")
		outputCmdArgs(fp, args)
	}
	return types.PrintResult(getResult(netConf), netConf.CNIVersion)
}

func cmdDel(args *skel.CmdArgs) error {
	netConf, _ := parseConf(args.StdinData)
	if netConf.CNIOutput != "" {
		fp, _ := os.OpenFile(netConf.CNIOutput, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0644)
		defer fp.Close()
		fmt.Fprintf(fp, "CmdDel\n")
		outputCmdArgs(fp, args)
	}
	return types.PrintResult(&type100.Result{}, netConf.CNIVersion)
}

func cmdCheck(args *skel.CmdArgs) error {
	netConf, _ := parseConf(args.StdinData)
	if netConf.CNIOutput != "" {
		fp, _ := os.OpenFile(netConf.CNIOutput, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0644)
		defer fp.Close()
		fmt.Fprintf(fp, "CmdCheck\n")
		outputCmdArgs(fp, args)
	}
	return types.PrintResult(&type100.Result{}, netConf.CNIVersion)
}
