package main

// Image representation of a container image
type Image struct {
	// Unique identifier for an image
	ID string
	// ImageName references the name that can be used to fetch it from external sources. (DockerHub, Quay, etc.)
	ImageName string
	// ImageSource is where you want to download images from. http typically
	ImageSource string
	// Source or Name?
	Source bool
}
