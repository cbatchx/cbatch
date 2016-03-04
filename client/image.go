package main

// Image representation of a container image
type Image struct {
	// Unique identifier for an image
	ID string
	// ImageName references the name that can be used to fetch it from external sources.
	ImageName string
	// ImageSource is where you want to download images from. (DockerHub, Quay, etc.)
	ImageSource string
}
