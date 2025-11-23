//go:build e2e

package e2e

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
)

const baseURL = "http://localhost:8080"

func TestE2E_FullFlow(t *testing.T) {
	e := httpexpect.Default(t, baseURL)

	ts := time.Now().UnixNano()
	teamName := fmt.Sprintf("e2e_team_%d", ts)
	authorID := fmt.Sprintf("u_auth_%d", ts)
	rev1ID := fmt.Sprintf("u_rev1_%d", ts)
	rev2ID := fmt.Sprintf("u_rev2_%d", ts)
	rev3ID := fmt.Sprintf("u_rev3_%d", ts)
	rev4ID := fmt.Sprintf("u_rev4_%d", ts)
	prID := fmt.Sprintf("pr_%d", ts)
	prID2 := fmt.Sprintf("pr_deact_%d", ts)

	e.POST("/team/add").
		WithJSON(map[string]interface{}{
			"team_name": teamName,
			"members": []map[string]interface{}{
				{"user_id": authorID, "username": "Author", "is_active": true},
				{"user_id": rev1ID, "username": "Rev1", "is_active": true},
				{"user_id": rev2ID, "username": "Rev2", "is_active": true},
				{"user_id": rev3ID, "username": "Rev3", "is_active": true},
				{"user_id": rev4ID, "username": "Rev4", "is_active": true},
			},
		}).
		Expect().
		Status(http.StatusCreated)

	prObj := e.POST("/pullRequest/create").
		WithJSON(map[string]interface{}{
			"pull_request_id":   prID,
			"pull_request_name": "Feature Login",
			"author_id":         authorID,
		}).
		Expect().
		Status(http.StatusCreated).
		JSON().Object().
		Value("pr").Object()

	reviewers := prObj.Value("assigned_reviewers").Array()
	reviewers.Length().IsEqual(2)
	firstReviewerID := reviewers.Value(0).String().Raw()

	e.GET("/users/getReview").
		WithQuery("user_id", firstReviewerID).
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		Value("pull_requests").Array().NotEmpty()

	reassignResp := e.POST("/pullRequest/reassign").
		WithJSON(map[string]interface{}{
			"pull_request_id": prID,
			"old_user_id":     firstReviewerID,
		}).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	newReviewer := reassignResp.Value("replaced_by").String().Raw()
	reassignResp.Value("pr").Object().Value("assigned_reviewers").Array().Contains(newReviewer)

	e.POST("/pullRequest/merge").
		WithJSON(map[string]interface{}{"pull_request_id": prID}).
		Expect().Status(http.StatusOK)

	e.GET("/stats").Expect().Status(http.StatusOK)

	prObj2 := e.POST("/pullRequest/create").
		WithJSON(map[string]interface{}{
			"pull_request_id":   prID2,
			"pull_request_name": "Feature Deactivate",
			"author_id":         authorID,
		}).
		Expect().
		Status(http.StatusCreated).
		JSON().Object().
		Value("pr").Object()

	currentReviewers := prObj2.Value("assigned_reviewers").Array()
	victimID := currentReviewers.Value(0).String().Raw()

	e.POST("/users/deactivate").
		WithJSON(map[string]interface{}{
			"user_ids": []string{victimID},
		}).
		Expect().
		Status(http.StatusOK)

	prList := e.GET("/users/getReview").
		WithQuery("user_id", victimID).
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		Value("pull_requests").Array()

	for _, element := range prList.Iter() {
		element.Object().Value("pull_request_id").NotEqual(prID2)
	}
}
