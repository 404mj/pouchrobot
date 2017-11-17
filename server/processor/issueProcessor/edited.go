package issueProcessor

import (
	"fmt"

	"github.com/allencloud/automan/server/processor/issueProcessor/open"
	"github.com/allencloud/automan/server/utils"
	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
)

// ActToIssueEdited acts to edited issue
// This function covers the following part:
// generate labels;
// attach comments;
// assign issue to specific user;
func (fIP *IssueProcessor) ActToIssueEdited(issue *github.Issue) error {
	// generate labels
	newLabels := open.ParseToGenerateLabels(issue)
	if len(newLabels) != 0 {
		// replace the original labels for issue
		getLabels, err := fIP.Client.GetLabelsInIssue(*(issue.Number))
		if err != nil {
			return err
		}
		originalLabels := []string{}
		for _, value := range getLabels {
			originalLabels = append(originalLabels, value.GetName())
		}
		addedLabels := utils.DeltaSlice(originalLabels, newLabels)
		if err := fIP.Client.AddLabelsToIssue(*(issue.Number), addedLabels); err != nil {
			return err
		}
	}

	// attach comment
	newComment := &github.IssueComment{}

	// check if the title is too short or the body empty.
	if issue.Title == nil || len(*(issue.Title)) < 20 {
		body := fmt.Sprintf(utils.IssueTitleTooShort, *(issue.User.Login))
		newComment.Body = &body
		if err := fIP.Client.AddCommentToIssue(*(issue.Number), newComment); err != nil {
			return err
		}
		logrus.Infof("succeed in attaching TITLE TOO SHORT comment for issue %d", *(issue.Number))

		labels := []string{"status/more-info-needed"}
		fIP.Client.AddLabelsToIssue(*(issue.Number), labels)

		return nil
	}

	if issue.Body == nil || len(*(issue.Body)) < 50 {
		body := fmt.Sprintf(utils.IssueDescriptionTooShort, *(issue.User.Login))
		newComment.Body = &body
		if err := fIP.Client.AddCommentToIssue(*(issue.Number), newComment); err != nil {
			return err
		}
		logrus.Infof("secceed in attaching TITLE TOO SHORT comment for issue %d", *(issue.Number))

		labels := []string{"status/more-info-needed"}
		fIP.Client.AddLabelsToIssue(*(issue.Number), labels)

		return nil
	}

	return nil
}