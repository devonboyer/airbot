package airbot

import (
	"context"

	"github.com/devonboyer/airbot/airtable"
)

type Scope struct {
	BaseID  string
	TableID string
}

type ScopeController struct {
	Client *airtable.Client
	Scope  Scope
}

func NewScopeController(client *airtable.Client, s Scope) *ScopeController {
	return &ScopeController{Scope: s}
}

func (c *ScopeController) List(ctx context.Context, filterByForumla string, v interface{}) error {
	return c.Client.
		Base(c.Scope.BaseID).
		Table(c.Scope.TableID).
		List().
		FilterByFormula(filterByForumla).
		Do(ctx, v)
}
