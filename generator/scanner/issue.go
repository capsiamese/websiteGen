package scanner

import (
	"context"
	"fmt"
	"github.com/google/go-github/v51/github"
	"golang.org/x/oauth2"
	"net/http"
)

type IssueScanner struct {
	client *github.Client
	opt    *github.IssueListByRepoOptions
	owner  string
	repo   string
}

func NewIssueScanner(ctx context.Context, token, owner, repo string, opts ...IssueOption) *IssueScanner {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	opt := &github.IssueListByRepoOptions{}
	for _, v := range opts {
		v(opt)
	}
	return &IssueScanner{
		client: github.NewClient(oauth2.NewClient(ctx, ts)),
		opt:    opt, owner: owner, repo: repo,
	}
}

type IssueOption func(options *github.IssueListByRepoOptions)

func WithLabels(label ...string) IssueOption {
	return func(options *github.IssueListByRepoOptions) {
		options.Labels = append(options.Labels, label...)
	}
}

func WithState(state string) IssueOption {
	return func(options *github.IssueListByRepoOptions) {
		options.State = state
	}
}

func (i *IssueScanner) Scan(ctx context.Context) ([]string, error) {
	i.opt.Page = 0
	list := make([]string, 0)
	for {
		res, resp, err := i.client.Issues.ListByRepo(ctx, i.owner, i.repo, i.opt)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("response status code:%d", resp.StatusCode)
		}
		for _, issue := range res {
			list = append(list, issue.GetBody())

		}
		if resp.NextPage == 0 {
			return list, nil
		}
		i.opt.Page = resp.NextPage
	}
}

/*
TODO: 1.增加扫描文件系统和issues的通用接口
TODO: 2.支持特殊文件, 例如About, References等
TODO: 3.迁移本仓库到blog仓库中

blog仓库中目录结构
- templates
- posts
- workspace
  - generator
  - gui
- assets
*/
