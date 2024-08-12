// package internal

// import (
// 	"archive/tar"
// 	"compress/gzip"

// 	"fmt"
// 	"io"
// 	"os"
// 	"os/exec"
// 	"path"
// 	"path/filepath"
// 	"runtime"
// 	"strings"
// 	"sync"

// 	"github.com/melbahja/got"

// 	"golang.org/x/net/context"
// )

// // Constants
// // const (
// // 	defaultFabricVersion = "2.5.9"
// // 	defaultCAVersion     = "1.5.12"
// // 	registry             = "docker.io/hyperledger"
// // )

// const (
// 	registry = "docker.io/hyperledger"
// )

// // Variables
// var (
// 	fabricVersion string
// 	caVersion     string
// 	components    string
// 	platform      string
// 	march         string
// )

// func RemoveDuplicates(component []string) []string {
// 	unique := make(map[string]bool)
// 	result := []string{}
// 	for _, item := range component {
// 		if !unique[item] {
// 			result = append(result, item)
// 			unique[item] = true
// 		}
// 	}
// 	return result
// }

// func checkDownloaderBin(fabricVersion, caVersion, platform string) bool {
// 	tmpDir := "tmp"
// 	binariesDir := "binaries"
// 	caArchive := fmt.Sprintf("hyperledger-fabric-ca-%s-%s.tar.gz", platform, caVersion)
// 	fabricArchive := fmt.Sprintf("hyperledger-fabric-%s-%s.tar.gz", platform, fabricVersion)
// 	fabricBinPath := filepath.Join(binariesDir, "fabric_bin", fabricVersion)
// 	caBinPath := filepath.Join(binariesDir, "ca_bin", caVersion)

// 	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
// 		return false
// 	}
// 	if _, err := os.Stat(binariesDir); os.IsNotExist(err) {
// 		return false
// 	}
// 	if _, err := os.Stat(filepath.Join(binariesDir, "fabric_bin")); os.IsNotExist(err) {
// 		return false
// 	}
// 	if _, err := os.Stat(filepath.Join(binariesDir, "fabric_bin", fabricVersion)); os.IsNotExist(err) {
// 		return false
// 	}
// 	if _, err := os.Stat(filepath.Join(binariesDir, "ca_bin")); os.IsNotExist(err) {
// 		return false
// 	}
// 	if _, err := os.Stat(filepath.Join(binariesDir, "ca_bin", caVersion)); os.IsNotExist(err) {
// 		return false
// 	}
// 	if _, err := os.Stat(filepath.Join(tmpDir, fabricArchive)); os.IsNotExist(err) {
// 		return false
// 	}
// 	if _, err := os.Stat(filepath.Join(tmpDir, caArchive)); os.IsNotExist(err) {
// 		return false
// 	}
// 	if _, err := os.Stat(filepath.Join(fabricBinPath, "bin")); os.IsNotExist(err) {
// 		return false
// 	}
// 	if _, err := os.Stat(filepath.Join(caBinPath, "bin")); os.IsNotExist(err) {
// 		return false
// 	}
// 	return true
// }

// func FabricDownload(cfgStruct *ConfigStruct) {
// 	// rempove duplicate entries from components
// 	cfgStruct.HyperledgerFabric.Components = RemoveDuplicates(cfgStruct.HyperledgerFabric.Components)

// 	for _, c := range cfgStruct.HyperledgerFabric.Components {
// 		switch c {
// 		case "all":
// 			components = "docker,binary,podman,samples"
// 		case "docker":
// 			components = "docker,"
// 		case "binary":
// 			components = "binary,"
// 		case "podman":
// 			components = "podman,"
// 		case "samples":
// 			components = "samples,"
// 		}
// 	}
// 	// remove duplicate comma
// 	components = strings.TrimRight(components, ",")
// 	fabricVersion = cfgStruct.HyperledgerFabric.Fabric
// 	caVersion = cfgStruct.HyperledgerFabric.CA

// 	osType := strings.ToLower(runtime.GOOS)
// 	arch := runtime.GOARCH
// 	if arch == "x86_64" {
// 		arch = "amd64"
// 	} else if arch == "arm64" {
// 		arch = "aarch64"
// 	}
// 	platform = fmt.Sprintf("%s-%s", osType, arch)
// 	march = runtime.GOARCH

// 	// Determine FABRIC_TAG and CA_TAG
// 	fabricTag := fabricVersion
// 	caTag := caVersion
// 	if strings.HasPrefix(fabricVersion, "1.") {
// 		fabricTag = fmt.Sprintf("%s-%s", march, fabricVersion)
// 		caTag = fmt.Sprintf("%s-%s", march, caVersion)
// 	}

// 	// Prior to fabric 2.5, use amd64 binaries on darwin-arm64
// 	if strings.HasPrefix(fabricVersion, "2.") && fabricVersion < "2.5" {
// 		platform = strings.ReplaceAll(platform, "darwin-arm64", "darwin-amd64")
// 	}

// 	binaryFile := fmt.Sprintf("hyperledger-fabric-%s-%s.tar.gz", platform, fabricVersion)
// 	caBinaryFile := fmt.Sprintf("hyperledger-fabric-ca-%s-%s.tar.gz", platform, caVersion)

// 	if checkDownloaderBin(fabricVersion, caVersion, platform) {
// 		return
// 	}

// 	ch := make(chan string)
// 	var wg sync.WaitGroup

// 	if strings.Contains(components, "s") || strings.Contains(components, "samples") {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			cloneSamplesRepo(ch)
// 		}()
// 	}

// 	if strings.Contains(components, "b") || strings.Contains(components, "binary") {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			pullBinaries(binaryFile, fabricTag, caBinaryFile, caTag, ch)
// 		}()
// 	}

// 	go func() {
// 		wg.Wait()
// 		close(ch)
// 	}()

// 	for msg := range ch {
// 		fmt.Println(msg)
// 	}

// }

// func cloneSamplesRepo(ch chan string) {
// 	if _, err := os.Stat("test-network"); err == nil {
// 		fmt.Println("==> Already in fabric-samples repo")
// 		ch <- fmt.Sprintln("==> Already in fabric-samples repo")

// 	} else if _, err := os.Stat("fabric-samples"); err == nil {
// 		fmt.Println("===> Changing directory to fabric-samples")
// 		ch <- fmt.Sprintln("===> Changing directory to fabric-samples")

// 		os.Chdir("fabric-samples")
// 	} else {
// 		fmt.Println("===> Cloning hyperledger/fabric-samples repo")
// 		ch <- fmt.Sprintln("===> Cloning hyperledger/fabric-samples repo")

// 		cmd := exec.Command("git", "clone", "-b", "main", "https://github.com/hyperledger/fabric-samples.git")
// 		cmd.Run()
// 		os.Chdir("fabric-samples")
// 	}

// 	cmd := exec.Command("git", "rev-parse", fmt.Sprintf("v%s", fabricVersion))
// 	if err := cmd.Run(); err == nil {
// 		fmt.Printf("===> Checking out v%s of hyperledger/fabric-samples\n", fabricVersion)
// 		ch <- fmt.Sprintf("===> Checking out v%s of hyperledger/fabric-samples\n", fabricVersion)

// 		cmd = exec.Command("git", "checkout", "-q", fmt.Sprintf("v%s", fabricVersion))
// 		cmd.Run()
// 	} else {
// 		fmt.Printf("fabric-samples v%s does not exist, defaulting to main. fabric-samples main branch is intended to work with recent versions of fabric.\n", fabricVersion)
// 		ch <- fmt.Sprintf("fabric-samples v%s does not exist, defaulting to main. fabric-samples main branch is intended to work with recent versions of fabric.\n", fabricVersion)

// 		cmd = exec.Command("git", "checkout", "-q", "main")
// 		cmd.Run()

// 	}
// }

// func pullBinaries(binaryFile, fabricTag, caBinaryFile, caTag string, ch chan<- string) {
// 	fmt.Printf("===> Downloading version %s platform specific fabric binaries\n", fabricTag)
// 	download(binaryFile, fmt.Sprintf("https://github.com/hyperledger/fabric/releases/download/v%s/%s", fabricVersion, binaryFile), path.Join("binaries", "fabric_bin", fabricVersion), ch)

// 	fmt.Printf("===> Downloading version %s platform specific fabric-ca-client binary\n", caTag)
// 	download(caBinaryFile, fmt.Sprintf("https://github.com/hyperledger/fabric-ca/releases/download/v%s/%s", caVersion, caBinaryFile), path.Join("binaries", "ca_bin", caVersion), ch)
// }

// func download(binaryFile, url, fabric_folder string, ch chan<- string) {
// 	destDir := "."
// 	if chainkunPath, ok := os.LookupEnv("CHAINKUN_ENV"); ok {
// 		destDir = path.Join(chainkunPath, "tmp")
// 	} else {
// 		destDir = path.Join(os.Getenv("HOME"), ".chainkun", "tmp")
// 	}
// 	if _, err := os.Stat("fabric-samples"); err == nil {
// 		destDir = "fabric-samples"
// 	}

// 	ch <- fmt.Sprintf("===> Downloading: %s\n", url)

// 	fmt.Printf("===> Downloading: %s\n", url)

// 	ch <- fmt.Sprintf("===> Will unpack to: %s\n", destDir)

// 	fmt.Printf("===> Will unpack to: %s\n", destDir)

// 	// Download the binary file to a temporary location

// 	tempFile := path.Join(destDir, binaryFile)

// 	ctx := context.Background()
// 	println("Downloading: ", url)
// 	g := got.NewDownload(ctx, url, tempFile)
// 	if err := g.Init(); err != nil {
// 		println(url)
// 		fmt.Println("==> There was an error downloading the binary file.")
// 		fmt.Println(err)
// 		os.Exit(22)
// 	}
// 	if err := g.Start(); err != nil {
// 		println(url)
// 		fmt.Println("==> There was an error downloading the binary file.")
// 		fmt.Println(err)
// 		os.Exit(22)
// 	}
// 	// Show progress bar
// 	go g.RunProgress(func(d *got.Download) {
// 		d.Size()
// 		d.Speed()
// 		d.AvgSpeed()
// 		d.TotalCost()
// 	})

// 	if chainkunPath, ok := os.LookupEnv("CHAINKUN_ENV"); ok {
// 		destDir = path.Join(chainkunPath, "binaries", fabric_folder)
// 	} else {
// 		destDir = path.Join(os.Getenv("HOME"), ".chainkun", "binaries", fabric_folder)
// 	}
// 	// Extract the binary file to the destination directory
// 	fmt.Printf("===> Unpacking %s\n", binaryFile)
// 	// open the tar.gz file
// 	file, err := os.Open(tempFile)
// 	if err != nil {
// 		fmt.Println("==> Error opening the tar.gz file.")
// 		fmt.Println(err)
// 		os.Exit(22)
// 	}
// 	defer file.Close()
// 	// open the gzip reader
// 	gzr, err := gzip.NewReader(file)
// 	if err != nil {
// 		os.Exit(22)
// 	}
// 	defer gzr.Close()

// 	tr := tar.NewReader(gzr)

// 	for {
// 		header, err := tr.Next()

// 		switch {

// 		// if no more files are found return
// 		case err == io.EOF:
// 			fmt.Println("==> Done.")
// 			return

// 		// return any other error
// 		case err != nil:
// 			os.Exit(22)

// 		// if the header is nil, just skip it (not sure how this happens)
// 		case header == nil:
// 			continue
// 		}

// 		// the target location where the dir/file should be created
// 		target := filepath.Join(destDir, header.Name)

// 		// the following switch could also be done using fi.Mode(), not sure if there
// 		// a benefit of using one vs. the other.
// 		// fi := header.FileInfo()

// 		// check the file type
// 		switch header.Typeflag {

// 		// if its a dir and it doesn't exist create it
// 		case tar.TypeDir:
// 			if _, err := os.Stat(target); err != nil {
// 				if err := os.MkdirAll(target, 0755); err != nil {
// 					os.Exit(22)
// 				}
// 			}

// 		// if it's a file create it
// 		case tar.TypeReg:
// 			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
// 			if err != nil {
// 				os.Exit(22)
// 			}

// 			// copy over contents
// 			if _, err := io.Copy(f, tr); err != nil {
// 				os.Exit(22)
// 			}

// 			// manually close here after each file operation; defering would cause each file close
// 			// to wait until all operations have completed.
// 			f.Close()
// 		}
// 	}

// }

package internal

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/melbahja/got"
	"golang.org/x/net/context"
)

// Variables
var (
	fabricVersion string
	caVersion     string
	components    string
	platform      string
	march         string
)

// RemoveDuplicates removes duplicate entries from a slice of strings.
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

// checkDownloaderBin checks if the binaries are already downloaded.
func checkDownloaderBin(fabricVersion, caVersion, platform string) bool {
	tmpDir := "tmp"
	binariesDir := "binaries"
	caArchive := fmt.Sprintf("hyperledger-fabric-ca-%s-%s.tar.gz", platform, caVersion)
	fabricArchive := fmt.Sprintf("hyperledger-fabric-%s-%s.tar.gz", platform, fabricVersion)
	fabricBinPath := filepath.Join(binariesDir, "fabric_bin", fabricVersion)
	caBinPath := filepath.Join(binariesDir, "ca_bin", caVersion)

	pathsToCheck := []string{
		tmpDir,
		binariesDir,
		filepath.Join(binariesDir, "fabric_bin"),
		filepath.Join(binariesDir, "fabric_bin", fabricVersion),
		filepath.Join(binariesDir, "ca_bin"),
		filepath.Join(binariesDir, "ca_bin", caVersion),
		filepath.Join(tmpDir, fabricArchive),
		filepath.Join(tmpDir, caArchive),
		filepath.Join(fabricBinPath, "bin"),
		filepath.Join(caBinPath, "bin"),
	}

	for _, p := range pathsToCheck {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// FabricDownload downloads the required Fabric components.
func FabricDownload(cfgStruct *ConfigStruct) {
	// Remove duplicate entries from components
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
	// Remove duplicate comma
	components = strings.TrimRight(components, ",")
	fabricVersion = cfgStruct.HyperledgerFabric.Fabric
	caVersion = cfgStruct.HyperledgerFabric.CA

	osType := strings.ToLower(runtime.GOOS)
	arch := runtime.GOARCH
	if arch == "x86_64" {
		arch = "amd64"
	} else if arch == "arm64" {
		arch = "aarch64"
	}
	platform = fmt.Sprintf("%s-%s", osType, arch)
	march = runtime.GOARCH

	// Determine FABRIC_TAG and CA_TAG
	fabricTag := fabricVersion
	caTag := caVersion
	if strings.HasPrefix(fabricVersion, "1.") {
		fabricTag = fmt.Sprintf("%s-%s", march, fabricVersion)
		caTag = fmt.Sprintf("%s-%s", march, caVersion)
	}

	// Prior to fabric 2.5, use amd64 binaries on darwin-arm64
	if strings.HasPrefix(fabricVersion, "2.") && fabricVersion < "2.5" {
		platform = strings.ReplaceAll(platform, "darwin-arm64", "darwin-amd64")
	}

	binaryFile := fmt.Sprintf("hyperledger-fabric-%s-%s.tar.gz", platform, fabricVersion)
	caBinaryFile := fmt.Sprintf("hyperledger-fabric-ca-%s-%s.tar.gz", platform, caVersion)

	if checkDownloaderBin(fabricVersion, caVersion, platform) {
		return
	}

	ch := make(chan string)
	var wg sync.WaitGroup

	if strings.Contains(components, "samples") {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cloneSamplesRepo(ch)
		}()
	}

	if strings.Contains(components, "binary") {
		wg.Add(1)
		go func() {
			defer wg.Done()
			pullBinaries(binaryFile, fabricTag, caBinaryFile, caTag, ch)
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for msg := range ch {
		fmt.Println(msg)
	}
}

// cloneSamplesRepo clones the fabric-samples repository.
func cloneSamplesRepo(ch chan string) {
	if _, err := os.Stat("test-network"); err == nil {
		fmt.Println("==> Already in fabric-samples repo")
		ch <- fmt.Sprintln("==> Already in fabric-samples repo")
	} else if _, err := os.Stat("fabric-samples"); err == nil {
		fmt.Println("===> Changing directory to fabric-samples")
		ch <- fmt.Sprintln("===> Changing directory to fabric-samples")
		os.Chdir("fabric-samples")
	} else {
		fmt.Println("===> Cloning hyperledger/fabric-samples repo")
		ch <- fmt.Sprintln("===> Cloning hyperledger/fabric-samples repo")
		cmd := exec.Command("git", "clone", "-b", "main", "https://github.com/hyperledger/fabric-samples.git")
		if err := cmd.Run(); err != nil {
			fmt.Printf("Error cloning fabric-samples repo: %v\n", err)
			os.Exit(1)
		}
		os.Chdir("fabric-samples")
	}

	cmd := exec.Command("git", "rev-parse", fmt.Sprintf("v%s", fabricVersion))
	if err := cmd.Run(); err == nil {
		fmt.Printf("===> Checking out v%s of hyperledger/fabric-samples\n", fabricVersion)
		ch <- fmt.Sprintf("===> Checking out v%s of hyperledger/fabric-samples\n", fabricVersion)
		cmd = exec.Command("git", "checkout", "-q", fmt.Sprintf("v%s", fabricVersion))
		if err := cmd.Run(); err != nil {
			fmt.Printf("Error checking out v%s: %v\n", fabricVersion, err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("fabric-samples v%s does not exist, defaulting to main. fabric-samples main branch is intended to work with recent versions of fabric.\n", fabricVersion)
		ch <- fmt.Sprintf("fabric-samples v%s does not exist, defaulting to main. fabric-samples main branch is intended to work with recent versions of fabric.\n", fabricVersion)
		cmd = exec.Command("git", "checkout", "-q", "main")
		if err := cmd.Run(); err != nil {
			fmt.Printf("Error checking out main branch: %v\n", err)
			os.Exit(1)
		}
	}
}

// pullBinaries pulls the required binaries.
func pullBinaries(binaryFile, fabricTag, caBinaryFile, caTag string, ch chan<- string) {
	fmt.Printf("===> Downloading version %s platform specific fabric binaries\n", fabricTag)
	download(binaryFile, fmt.Sprintf("https://github.com/hyperledger/fabric/releases/download/v%s/%s", fabricVersion, binaryFile), path.Join("binaries", "fabric_bin", fabricVersion), ch)

	fmt.Printf("===> Downloading version %s platform specific fabric-ca-client binary\n", caTag)
	download(caBinaryFile, fmt.Sprintf("https://github.com/hyperledger/fabric-ca/releases/download/v%s/%s", caVersion, caBinaryFile), path.Join("binaries", "ca_bin", caVersion), ch)
}

// download downloads a file and extracts it.
func download(binaryFile, url, fabricFolder string, ch chan<- string) {
	destDir := "."
	if chainkunPath, ok := os.LookupEnv("CHAINKUN_ENV"); ok {
		destDir = path.Join(chainkunPath, "tmp")
	} else {
		destDir = path.Join(os.Getenv("HOME"), ".chainkun", "tmp")
	}
	if _, err := os.Stat("fabric-samples"); err == nil {
		destDir = "fabric-samples"
	}

	ch <- fmt.Sprintf("===> Downloading: %s\n", url)
	fmt.Printf("===> Downloading: %s\n", url)
	ch <- fmt.Sprintf("===> Will unpack to: %s\n", destDir)
	fmt.Printf("===> Will unpack to: %s\n", destDir)

	// Download the binary file to a temporary location
	tempFile := path.Join(destDir, binaryFile)
	ctx := context.Background()
	g := got.NewDownload(ctx, url, tempFile)
	if err := g.Init(); err != nil {
		fmt.Printf("Error initializing download: %v\n", err)
		os.Exit(22)
	}
	if err := g.Start(); err != nil {
		fmt.Printf("Error starting download: %v\n", err)
		os.Exit(22)
	}
	// Show progress bar
	go g.RunProgress(func(d *got.Download) {
		d.Size()
		d.Speed()
		d.AvgSpeed()
		d.TotalCost()
	})

	if chainkunPath, ok := os.LookupEnv("CHAINKUN_ENV"); ok {
		destDir = path.Join(chainkunPath, "binaries", fabricFolder)
	} else {
		destDir = path.Join(os.Getenv("HOME"), ".chainkun", "binaries", fabricFolder)
	}
	// Extract the binary file to the destination directory
	fmt.Printf("===> Unpacking %s\n", binaryFile)
	// open the tar.gz file
	file, err := os.Open(tempFile)
	if err != nil {
		fmt.Printf("Error opening the tar.gz file: %v\n", err)
		os.Exit(22)
	}
	defer file.Close()
	// open the gzip reader
	gzr, err := gzip.NewReader(file)
	if err != nil {
		fmt.Printf("Error creating gzip reader: %v\n", err)
		os.Exit(22)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {
		// if no more files are found return
		case err == io.EOF:
			fmt.Println("==> Done.")
			return
		// return any other error
		case err != nil:
			fmt.Printf("Error reading tar file: %v\n", err)
			os.Exit(22)
		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(destDir, header.Name)

		// check the file type
		switch header.Typeflag {
		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					fmt.Printf("Error creating directory: %v\n", err)
					os.Exit(22)
				}
			}
		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				fmt.Printf("Error creating file: %v\n", err)
				os.Exit(22)
			}
			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				fmt.Printf("Error copying file contents: %v\n", err)
				os.Exit(22)
			}
			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()
		}
	}
}
