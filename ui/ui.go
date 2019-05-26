package ui

import "context"

type UI interface {
	ReadPassword(ctx context.Context, desc string) (string, error)
}
