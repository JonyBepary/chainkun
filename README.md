
# Chainkun

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/JonyBepary/chainkun)](https://goreportcard.com/report/github.com/JonyBepary/chainkun)
[![GoDoc](https://godoc.org/github.com/JonyBepary/chainkun?status.svg)](https://godoc.org/github.com/JonyBepary/chainkun)

Chainkun is a powerful terminal application designed to streamline the management and automation of Chaincode and Hyperledger Fabric test networks. By simplifying Chaincode development and testing, this tool helps developers focus on building and deploying smart contracts more efficiently.

## Features

- **Automated Setup**: Automatically sets up and configures Hyperledger Fabric and Chaincode environments.
- **Configuration Management**: Manages configuration files for different environments and versions.
- **Binary Downloads**: Downloads and installs required binaries for Hyperledger Fabric and CA.
- **Sample Repositories**: Clones and manages sample repositories for testing and development.
- **Cross-Platform Support**: Supports multiple operating systems and architectures.

## Installation

### Prerequisites

- Go (version 1.22.5 or later)
- Docker
- Git

### Steps

1. **Clone the Repository**:
   ```sh
   git clone https://github.com/JonyBepary/chainkun.git
   cd chainkun
   ```

2. **Build the Project**:
   ```sh
   go build -o chainkun
   ```

3. **Run the Application**:
   ```sh
   ./chainkun
   ```

## Configuration

Chainkun uses a YAML configuration file to manage settings for Hyperledger Fabric, Chaincode, and other components. An example configuration file is provided below:

```yaml
hyperledger-fabric:
  fabric: 2.5.9
  ca: 1.5.12
  components:
    - docker
    - binary
    # - samples
    # - podman

cc-tool:
  repo: https://github.com/hyperledger-labs/cc-tools-demo
  branch: main

chaincode:
  name: chainkun
  label: 1.0
  version: 0.1
  sequence: 1

network:
  ccas: false
  org_qnty: 3
  ccaas: false
  CcasTls: "1.3"
  orgs:
    - org1MSP
    - org2MSP
    - org3MSP

system:
  path: /home/jony/.chainkun
  go: "1.22.5"
```

## Usage

### Initializing Chainkun

To initialize Chainkun and set up the necessary directories and configuration files, run:

```sh
./chainkun init
```

### Downloading Binaries

To download the required Hyperledger Fabric and CA binaries, run:

```sh
./chainkun download
```

### Cloning Sample Repositories

To clone the sample repositories for testing and development, run:

```sh
./chainkun clone-samples
```

## Contributing

Contributions are welcome! Please follow these steps to contribute:

1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Make your changes and commit them with descriptive commit messages.
4. Push your branch to your fork.
5. Create a pull request to the main repository.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contact

<!-- LinkedIn, X, Mail -->
[![LinkedIn](https://img.shields.io/badge/LinkedIn-0077B5?style=for-the-badge&logo=linkedin&logoColor=white)](https://www.linkedin.com/in/sohel-ahmed-jony/)
[![X](https://img.shields.io/badge/X-000000?style=for-the-badge&logo=x&logoColor=white)](https://twitter.com/jbepary)
[![Email](https://img.shields.io/badge/Email-D14836?style=for-the-badge&logo=gmail&logoColor=white)](mailto:sohelahmedjony@gmail.com)

