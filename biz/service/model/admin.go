package model

type College struct {
	CollegeId   string
	CollegeName string
}

type Major struct {
	MajorId   string
	MajorName string
	CollegeId string
}
