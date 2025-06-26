package docker

import (
	"fmt"
	"os/exec"
	"strings"
)

func CheckDockerIsAvailable() (string, error) {
	cmd := exec.Command("docker", "--version")
	output, err := cmd.CombinedOutput() // Output + error message (stderr) combined
	if err != nil {
		return "", fmt.Errorf("docker command failed: %w (maybe docker not found in your machine, make sure to have docker installer)", err)
	}

	// Output will be like: "Docker version 24.0.5, build ced0996"
	version := strings.TrimSpace(string(output))

	return version, nil
}

func DockerDownloadImages(imageName, imageTag string) error {
	cmd := exec.Command("docker", "pull", fmt.Sprintf("%s:%s", imageName, imageTag))
	_, err := cmd.CombinedOutput() // Output + error message (stderr) combined
	if err != nil {
		return fmt.Errorf("docker pull command failed: %w (maybe docker not found in your machine, make sure to have docker installer)", err)
	}

	return nil
}

func CheckDockerImageIsAvailable(imageName, imageTag string) (DockerImage, error) {
	cmd := exec.Command("docker", "images", fmt.Sprintf("%s:%s", imageName, imageTag))
	output, err := cmd.CombinedOutput() // Output + error message (stderr) combined
	if err != nil {
		return DockerImage{}, fmt.Errorf("docker command failed: %w (maybe docker not found in your machine, make sure to have docker installer)", err)
	}

	// Output will be like: "semgrep 1.0.0"
	images := strings.Split(string(output), "\n")
	if len(images) < 2 {
		return DockerImage{}, fmt.Errorf("no Docker images found, maybe Semgrep is not installed")
	}

	var image_data DockerImage
	for _, image := range images[1:] {
		image = strings.TrimSpace(image)
		if image == "" {
			continue
		}

		image_splits := strings.Fields(image)

		image_repository := image_splits[0] // e.g., "semgrep"
		image_tag := image_splits[1]        // e.g., "1.0.0"
		image_id := image_splits[2]         // e.g., "abc123def456"

		image_data = DockerImage{
			Repository: image_repository,
			Tag:        image_tag,
			ImageID:    image_id,
		}
	}

	return image_data, nil
}
