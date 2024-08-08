package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/JonyBepary/chainkun/internal"
)

// Constants
// const (
// 	defaultFabricVersion = "2.5.9"
// 	defaultCAVersion     = "1.5.12"
// 	registry             = "docker.io/hyperledger"
// )

const (
	registry = "docker.io/hyperledger"
)

// Variables
var (
	fabricVersion string
	caVersion     string
	components    string
	platform      string
	march         string
)

func RemoveDuplicates(component []string) []string {
	unique := make(map[string]bool)
	result := []string{}
	for _, item := range component {
		if !unique[item] {
			result = append(result, item)
			unique[item] = true
		}
	}
	return result
}

func FabricDownload(cfgStruct *internal.ConfigStruct) {
	// rempove duplicate entries from components
	cfgStruct.HyperledgerFabric.Components = RemoveDuplicates(cfgStruct.HyperledgerFabric.Components)

	for _, c := range cfgStruct.HyperledgerFabric.Components {
		switch c {
		case "all":
			components = "docker,binary,podman,samples"
		case "docker":
			components = "docker,"
		case "binary":
			components = "binary,"
		case "podman":
			components = "podman,"
		case "samples":
			components = "samples,"
		}
	}
	// remove duplicate comma
	components = strings.TrimRight(components, ",")
	version := cfgStruct.HyperledgerFabric.Fabric
	caVersion := cfgStruct.HyperledgerFabric.CA

	osType := strings.ToLower(runtime.GOOS)
	arch := runtime.GOARCH
	if arch == "amd64" {
		arch = "x86_64"
	} else if arch == "arm64" {
		arch = "aarch64"
	}
	platform = fmt.Sprintf("%s-%s", osType, arch)
	march = runtime.GOARCH

	// Determine FABRIC_TAG and CA_TAG
	fabricTag := version
	caTag := caVersion
	if strings.HasPrefix(version, "1.") {
		fabricTag = fmt.Sprintf("%s-%s", march, version)
		caTag = fmt.Sprintf("%s-%s", march, caVersion)
	}

	// Prior to fabric 2.5, use amd64 binaries on darwin-arm64
	if strings.HasPrefix(version, "2.") && version < "2.5" {
		platform = strings.ReplaceAll(platform, "darwin-arm64", "darwin-amd64")
	}

	binaryFile := fmt.Sprintf("hyperledger-fabric-%s-%s.tar.gz", platform, version)
	caBinaryFile := fmt.Sprintf("hyperledger-fabric-ca-%s-%s.tar.gz", platform, caVersion)

	if strings.Contains(components, "s") || strings.Contains(components, "samples") {
		cloneSamplesRepo()
	}

	if strings.Contains(components, "b") || strings.Contains(components, "binary") {
		pullBinaries(binaryFile, fabricTag, caBinaryFile, caTag)
	}

	if strings.Contains(components, "p") || strings.Contains(components, "podman") {
		pullImages("podman", fabricTag, caTag)
	}

	if strings.Contains(components, "d") || strings.Contains(components, "docker") {
		pullImages("docker", fabricTag, caTag)
	}
}

func cloneSamplesRepo() {
	if _, err := os.Stat("test-network"); err == nil {
		fmt.Println("==> Already in fabric-samples repo")
	} else if _, err := os.Stat("fabric-samples"); err == nil {
		fmt.Println("===> Changing directory to fabric-samples")
		os.Chdir("fabric-samples")
	} else {
		fmt.Println("===> Cloning hyperledger/fabric-samples repo")
		cmd := exec.Command("git", "clone", "-b", "main", "https://github.com/hyperledger/fabric-samples.git")
		cmd.Run()
		os.Chdir("fabric-samples")
	}

	cmd := exec.Command("git", "rev-parse", fmt.Sprintf("v%s", fabricVersion))
	if err := cmd.Run(); err == nil {
		fmt.Printf("===> Checking out v%s of hyperledger/fabric-samples\n", fabricVersion)
		cmd = exec.Command("git", "checkout", "-q", fmt.Sprintf("v%s", fabricVersion))
		cmd.Run()
	} else {
		fmt.Printf("fabric-samples v%s does not exist, defaulting to main. fabric-samples main branch is intended to work with recent versions of fabric.\n", fabricVersion)
		cmd = exec.Command("git", "checkout", "-q", "main")
		cmd.Run()
	}
}

func pullBinaries(binaryFile, fabricTag, caBinaryFile, caTag string) {
	fmt.Printf("===> Downloading version %s platform specific fabric binaries\n", fabricTag)
	download(binaryFile, fmt.Sprintf("https://github.com/hyperledger/fabric/releases/download/v%s/%s", fabricVersion, binaryFile))

	fmt.Printf("===> Downloading version %s platform specific fabric-ca-client binary\n", caTag)
	download(caBinaryFile, fmt.Sprintf("https://github.com/hyperledger/fabric-ca/releases/download/v%s/%s", caVersion, caBinaryFile))
}

func download(binaryFile, url string) {
	destDir := "."
	if _, err := os.Stat("fabric-samples"); err == nil {
		destDir = "fabric-samples"
	}
	fmt.Printf("===> Downloading: %s\n", url)
	fmt.Printf("===> Will unpack to: %s\n", destDir)
	cmd := exec.Command("curl", "-L", "--retry", "5", "--retry-delay", "3", url)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("==> There was an error downloading the binary file.")
		os.Exit(22)
	}
	cmd = exec.Command("tar", "xz", "-C", destDir)
	cmd.Stdin = strings.NewReader(string(out))
	cmd.Run()
	fmt.Println("==> Done.")
}

func pullImages(containerCli, fabricTag, caTag string) {
	cmd := exec.Command("command", "-v", containerCli)
	if err := cmd.Run(); err != nil {
		fmt.Printf("=========================================================\n")
		fmt.Printf("%s not installed, bypassing download of Fabric images\n", containerCli)
		fmt.Printf("=========================================================\n")
		return
	}

	fabricImages := []string{"peer", "orderer", "ccenv"}
	if strings.HasPrefix(fabricVersion, "2.") || strings.HasPrefix(fabricVersion, "3.") {
		fabricImages = append(fabricImages, "baseos")
	}

	fmt.Printf("FABRIC_IMAGES: %v\n", fabricImages)
	fmt.Println("===> Pulling fabric Images")
	singleImagePull(containerCli, fabricTag, fabricImages)

	fmt.Println("===> Pulling fabric ca Image")
	caImages := []string{"ca"}
	singleImagePull(containerCli, caTag, caImages)

	fmt.Println("===> List out hyperledger images")
	cmd = exec.Command(containerCli, "images")
	out, _ := cmd.Output()
	fmt.Println(string(out))
}

func singleImagePull(containerCli, tag string, images []string) {
	for _, image := range images {
		fmt.Printf("====>  %s/fabric-%s:%s\n", registry, image, tag)
		cmd := exec.Command(containerCli, "pull", fmt.Sprintf("%s/fabric-%s:%s", registry, image, tag))
		cmd.Run()
		cmd = exec.Command(containerCli, "tag", fmt.Sprintf("%s/fabric-%s:%s", registry, image, tag), fmt.Sprintf("%s/fabric-%s", registry, image))
		cmd.Run()
		twoDigitTag := strings.Join(strings.Split(tag, ".")[0:2], ".")
		cmd = exec.Command(containerCli, "tag", fmt.Sprintf("%s/fabric-%s:%s", registry, image, tag), fmt.Sprintf("%s/fabric-%s:%s", registry, image, twoDigitTag))
		cmd.Run()
	}
}
