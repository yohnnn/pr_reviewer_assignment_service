package models

import "time"

type Team struct {
	Name    string
	Members []TeamMember
}

type TeamMember struct {
	UserID   string
	UserName string
	IsActive bool
}

type User struct {
	ID       string
	Name     string
	TeamName string
	IsActive bool
}

type PullRequest struct {
	ID        string
	Name      string
	AuthorID  string
	Status    string
	CreatedAt time.Time
	MergedAt  *time.Time
	Reviewers []string
}

type PullRequestShort struct {
	ID       string
	Name     string
	AuthorID string
	Status   string
}
