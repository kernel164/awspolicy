package engine

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/kernel164/awspolicy/internal/capture"
	"github.com/kernel164/awspolicy/internal/capture/terraformcli"
	"github.com/kernel164/awspolicy/internal/out"
	"github.com/kernel164/awspolicy/internal/out/ssm"
	"github.com/kernel164/awspolicy/internal/util"
)

type CaptureEngine struct {
	out    out.Output
	stdout out.Output
}

func New() *CaptureEngine {
	stdout := out.New()
	return &CaptureEngine{
		out:    determineOutput(util.Getenv("AWSPOLICY_OUTPUT", "-"), stdout),
		stdout: stdout,
	}
}

func (e *CaptureEngine) Run() error {
	policyCaptureType := determinePolicyCaptureType()
	switch policyCaptureType {
	case "terraform-cli":
		tfcli := terraformcli.New()
		data, err := tfcli.Run(&capture.Params{Args: os.Args[1:]})
		if err != nil {
			return err
		}
		e.output(data)
	default:
		fmt.Printf("%s - Policy Capture Method Not yet implemented.\n", policyCaptureType)
	}
	return nil
}

func (e *CaptureEngine) output(data *capture.Results) {
	hasNew := false
	for range data.Actions {
		hasNew = true
		break
	}
	if !hasNew {
		// Nothing to do.
		return
	}
	out := e.out
	// get
	doc, err := out.Get()
	if err != nil {
		fmt.Println(err)
		doc, _ = e.stdout.Get()
	}

	// update
	stmt := doc.Statement[0]
	for _, action := range stmt.Action {
		data.Actions[action] = struct{}{}
	}

	// conv set to sorted array

	actions := []string{}
	for key := range data.Actions {
		actions = append(actions, key)
	}
	sort.Strings(actions)
	stmt.Action = actions

	// set output.
	err = out.Set(doc)
	if err != nil {
		fmt.Println(err)
	}
	if err != nil || out != e.stdout {
		e.stdout.Set(doc)
	}
}

func determineOutput(path string, def out.Output) out.Output {
	switch {
	case strings.HasPrefix(path, "ssm:"):
		return ssm.New(path)
	}
	return def
}

func determinePolicyCaptureType() string {
	captureType := ""
	basename := path.Base(os.Args[0])
	if strings.HasPrefix(basename, "terraform") {
		captureType = "terraform-cli"
	}
	return util.Getenv("AWSPOLICY_CAPTURE_TYPE", captureType)
}
