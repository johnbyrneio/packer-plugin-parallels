// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"regexp"
)

// IfconfigIPFinder finds the host IP based on the output of `ifconfig`.
type IfconfigIPFinder struct {
	Devices []string
}

// HostIP returns the host's IP address or an error if it could not be found
// from the `ifconfig` output.
func (f *IfconfigIPFinder) HostIP() (string, error) {
	var ifconfigPath string

	// On some systems, ifconfig is in /sbin which is generally not
	// on the PATH for a standard user, so we just check that first.
	if _, err := os.Stat("/sbin/ifconfig"); err == nil {
		ifconfigPath = "/sbin/ifconfig"
	}

	if ifconfigPath == "" {
		var err error
		ifconfigPath, err = exec.LookPath("ifconfig")
		if err != nil {
			return "", err
		}
	}

	for _, device := range f.Devices {
		stdout := new(bytes.Buffer)

		cmd := exec.Command(ifconfigPath, device)
		// Force LANG=C so that the output is what we expect it to be
		// despite the locale.
		cmd.Env = append(cmd.Env, "LANG=C")
		cmd.Env = append(cmd.Env, os.Environ()...)

		cmd.Stdout = stdout
		cmd.Stderr = new(bytes.Buffer)

		if err := cmd.Run(); err == nil {
			re := regexp.MustCompile(`inet\s+(?:addr:)?(.+?)\s`)
			matches := re.FindStringSubmatch(stdout.String())
			if matches != nil {
				return matches[1], nil
			}
		}
	}

	devices_checked_list := ""
	for _, device := range f.Devices {
		devices_checked_list += device + " "
	}
	return "", errors.New("IP not found in ifconfig output, check the host_interfaces config to ensure your device is listed. Devices checked: " + devices_checked_list)
}
