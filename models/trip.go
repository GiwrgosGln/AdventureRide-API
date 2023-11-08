package models

type Waypoint struct {
	Name        string `json:"name"`
	Location    string `json:"location"`
	Description string `json:"description"`
}

type Trip struct {
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	CreatorID     string     `json:"creatorid"`
	StartLocation string     `json:"startlocation"`
	EndLocation   string     `json:"endlocation"`
	Waypoints     []Waypoint `json:"waypoints"`
}