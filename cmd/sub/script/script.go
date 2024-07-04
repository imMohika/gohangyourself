package script

import (
	"errors"
	"flag"
	"fmt"
	"github.com/pterm/pterm"
	"log/slog"
	"os"
	"regexp"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"text/template"
)

const autoRestartWindows = ":start\n{{.}}\n\necho Server restarting...\necho Press CTRL + C to stop.\ngoto :start"

const (
	autoRestartUnix = "while [ true ]; do\n    {{.}}\n    echo Server restarting...\n    echo Press CTRL + C to stop.\ndone"
)

const (
	osWindows = "windows"
	osUnix    = "unix"
)

const (
	platformPaper    = "paper"
	platformVelocity = "velocity"
)

// from https://github.com/PaperMC/docs/blob/main/src/components/StartScriptGenerator.tsx
var flagsAikar = []string{
	"-XX:+AlwaysPreTouch",
	"-XX:+DisableExplicitGC",
	"-XX:+ParallelRefProcEnabled",
	"-XX:+PerfDisableSharedMem",
	"-XX:+UnlockExperimentalVMOptions",
	"-XX:+UseG1GC",
	"-XX:G1HeapWastePercent=5",
	"-XX:G1MixedGCCountTarget=4",
	"-XX:G1MixedGCLiveThresholdPercent=90",
	"-XX:G1RSetUpdatingPauseTimePercent=5",
	"-XX:MaxGCPauseMillis=200",
	"-XX:MaxTenuringThreshold=1",
	"-XX:SurvivorRatio=32",
	"-Dusing.aikars.flags=https://mcflags.emc.gs",
	"-Daikars.new.flaautoRestartFlags=true",
}

// <12gb
var flagsLargeMem = []string{
	"-XX:G1NewSizePercent=40",
	"-XX:G1MaxNewSizePercent=50",
	"-XX:G1HeapRegionSize=16M",
	"-autoRestartFlagX:G1ReservePercent=15",
	"-XX:InitiatingHeapOccupancyPercent=20",
}

// >12gb
var flagsSmallMem = []string{
	"-XX:G1NewSizePercent=30",
	"-XX:G1MaxNewSizePercent=40",
	"-XX:G1HeapRegionSize=8M",
	"-XX:G1ReservePercent",
	"--add-modules=jdk.incubator.vec",
}

var flagsVelocity = []string{
	"-XX:+AlwaysPreTouch",
	"-XX:+ParallelRefProcEnabled",
	"-XX:+UnlockExperimentalVMOptions",
	"-XX:+UseG1GC",
	"-XX:G1HeapRegionSize=4M",
	"-XX:MaxInlineLevel=15",
}

var supportedPlatforms = []string{"paper", "velocity"}

type SubCommand struct {
}

func (s SubCommand) Handle(args []string) {
	flags := flag.NewFlagSet("script", flag.ExitOnError)
	userOsFlag := flags.String("os", "", "windows, unix")
	platformFlag := flags.String("platform", "", "paper, velocity")
	jarFlag := flags.String("jar", "", "jar name")
	ramFlag := flags.String("ram", "", "ram")
	var autoRestartFlag bool
	flags.BoolVar(&autoRestartFlag, "autorestart", false, "autorestart")
	flags.BoolVar(&autoRestartFlag, "ar", false, "autorestart")
	nameFlag := flags.String("name", "", "name")

	err := flags.Parse(args)
	if err != nil {
		slog.Error("error parsing flags", "err", err)
	}

	userOS, err := getOs(*userOsFlag)
	if err != nil {
		slog.Error("error getting os", "err", err)
	}

	jar, err := getJar(*jarFlag)
	if err != nil {
		slog.Error("error getting platform", "err", err)
	}

	platform, err := getPlatform(*platformFlag, jar)
	if err != nil {
		slog.Error("error getting platform", "err", err)
	}

	ram, err := getRAM(*ramFlag)
	if err != nil {
		slog.Error("error getting ram", "err", err)
	}

	autoRestart, err := isAutoRestart(autoRestartFlag)
	if err != nil {
		slog.Error("error getting autorestart", "err", err)
	}

	baseFlags, err := getBaseFlags(platform, ram)
	if err != nil {
		slog.Error("Error getting base flags", "err", err)
		return
	}

	name, err := getName(*nameFlag, userOS)
	if err != nil {
		slog.Error("error getting name", "err", err)
	}

	content := generateContent(platform, baseFlags, jar)

	err = output(name, content, userOS, autoRestart)
	if err != nil {
		slog.Error("failed to generate shell file", "err", err)
		return
	}

	slog.Info("generated shell file")
}

func getName(nameFlag string, userOS string) (string, error) {
	if nameFlag != "" {
		return nameFlag, nil
	}

	var defaultName string
	switch userOS {
	case osWindows:
		defaultName = "start.bat"
	case osUnix:
		defaultName = "start.sh"
	}

	name, err := pterm.DefaultInteractiveTextInput.WithDefaultValue(defaultName).Show("file name: ")
	pterm.Println()
	if err != nil {
		return "", err
	}

	return name, nil
}

func isAutoRestart(autoRestartFlag bool) (bool, error) {
	if autoRestartFlag {
		return true, nil
	}

	// default value is false, in case of --autoRestart false don't prompt
	if isFlagPassed("autorestart") || isFlagPassed("ar") {
		return false, nil
	}

	autoRestart, err := pterm.DefaultInteractiveConfirm.WithDefaultValue(true).Show("auto restart? ")
	pterm.Println()
	if err != nil {
		return false, err
	}
	return autoRestart, nil
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func getJar(jarFlag string) (string, error) {
	if jarFlag != "" {
		if !strings.HasSuffix(jarFlag, ".jar") {
			return "", fmt.Errorf("%s has an invalid file extension", jarFlag)
		}
		return jarFlag, nil
	}

	// todo: find jars in current dir and show selector

	jar, err := pterm.DefaultInteractiveTextInput.Show("Jar File: ")
	if err != nil {
		return "", err
	}
	pterm.Println()
	return jar, nil
}

func output(name string, content string, userOS string, isAutoRestart bool) error {
	if !isAutoRestart {
		var output string

		switch userOS {
		case osWindows:
			output = fmt.Sprintf("@ECHO OFF\n%s\npause", content)
		case osUnix:
			output = fmt.Sprintf("#!/bin/sh\n%s", content)
		}

		err := os.WriteFile(name, []byte(output), 0o700)
		if err != nil {
			return err
		}

		return nil
	}

	var err error

	outputTmpl := template.New("output")
	switch userOS {
	case "windows":
		outputTmpl, err = outputTmpl.Parse(autoRestartWindows)
	case "unix":
		outputTmpl, err = outputTmpl.Parse(autoRestartUnix)
	}

	if err != nil {
		return err
	}

	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	err = outputTmpl.Execute(file, content)
	if err != nil {
		return err
	}

	return nil
}

func getBaseFlags(platform, ram string) ([]string, error) {
	base := []string{
		"-Xms" + ram,
		"-Xmx" + ram,
	}

	mb, err := memToMB(ram)
	if err != nil {
		return nil, err
	}

	switch platform {
	case platformPaper:
		base = append(base, flagsAikar...)
		if mb > 12000 {
			base = append(base, flagsLargeMem...)
		} else {
			base = append(base, flagsSmallMem...)
		}
	case platformVelocity:
		base = append(base, flagsVelocity...)
	}

	return base, nil
}

func memToMB(ram string) (int, error) {
	if strings.HasSuffix(ram, "G") {
		gb, err := strconv.Atoi(strings.TrimSuffix(ram, "G"))
		if err != nil {
			return 0, err
		}

		mb := gb * 1000
		return mb, nil
	}

	if strings.HasSuffix(ram, "M") {
		mb, err := strconv.Atoi(strings.TrimSuffix(ram, "M"))
		if err != nil {
			return 0, err
		}
		return mb, nil
	}

	return 0, errors.New("memory value should end with M or G")
}

func generateContent(platform string, baseFlags []string, jar string) string {
	content := baseFlags
	content = append(content, fmt.Sprintf("-jar %s", jar))

	if platform == platformPaper {
		content = append(content, "--nogui")
	}

	return fmt.Sprintf("java %s", strings.Join(content, " "))
}

func getRAM(ramFlag string) (string, error) {
	validString := regexp.MustCompile(`(?i)^\d+[MG]$`)
	if ramFlag != "" {
		if validString.MatchString(ramFlag) {
			return ramFlag, nil
		}
		return "", fmt.Errorf("invalid ram flag: %s", ramFlag)
	}

	ram, err := pterm.DefaultInteractiveTextInput.Show("Ram (Should end with M or G): ")
	if err != nil {
		return "", err
	}
	pterm.Println()

	if validString.MatchString(ram) {
		return ram, nil
	}
	return "", fmt.Errorf("invalid ram: %s", ram)
}

func getOs(userOsFlag string) (string, error) {
	if userOsFlag == "" {
		if runtime.GOOS == osWindows {
			return osWindows, nil
		}
		return osUnix, nil
	}

	if userOsFlag == osWindows || userOsFlag == osUnix {
		return userOsFlag, nil
	}

	return "", fmt.Errorf("%s is an unsupported os", userOsFlag)
}

func getPlatform(platformFlag string, jar string) (string, error) {
	if platformFlag != "" {
		if slices.Contains(supportedPlatforms, platformFlag) {
			return platformFlag, nil
		}
		return "", fmt.Errorf("unsupported platform: %s", platformFlag)
	}

	if strings.Contains(jar, "paper") {
		result, err := pterm.DefaultInteractiveConfirm.WithDefaultValue(true).Show(fmt.Sprintf("Detected %s, is it correct?", platformPaper))
		if err != nil {
			return "", err
		}

		pterm.Println()
		if result {
			return platformPaper, nil
		}
	}

	if strings.Contains(jar, "velocity") {
		result, err := pterm.DefaultInteractiveConfirm.Show(fmt.Sprintf("Detected %s, is it correct?", platformVelocity))
		if err != nil {
			return "", err
		}

		pterm.Println()
		if result {
			return platformVelocity, nil
		}
	}

	platform, err := pterm.DefaultInteractiveSelect.WithOptions(supportedPlatforms).Show("Select platform:")
	if err != nil {
		return "", err
	}
	return platform, nil
}
