package wfexec

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/stretchr/testify/require"
	"testing"
)

func gt(key string, val int) pathTester {
	return func(ctx context.Context, variables expr.Variables) (bool, error) {
		return variables.Int(key) > val, nil
	}
}

func TestJoinGateway(t *testing.T) {
	var (
		req        = require.New(t)
		p1, p2, p3 = &wfTestStep{name: "p1"}, &wfTestStep{name: "p2"}, &wfTestStep{name: "p3"}
		gw         = JoinGateway(p1, p2, p3)

		r   ExecResponse
		err error
	)

	r, err = gw.Exec(nil, &ExecRequest{Parent: p1})
	req.NoError(err)
	req.Equal(&partial{}, r)

	r, err = gw.Exec(nil, &ExecRequest{Parent: p2})
	req.NoError(err)
	req.Equal(&partial{}, r)

	r, err = gw.Exec(nil, &ExecRequest{Parent: p3})
	req.NoError(err)
	req.IsType(expr.Variables{}, r)
}

func TestForkGateway(t *testing.T) {
	var (
		req = require.New(t)
		gw  = ForkGateway()
	)

	r, err := gw.Exec(nil, nil)
	req.NoError(err)
	req.Equal(Steps{}, r)
	req.Empty(r)
}

func TestInclGateway(t *testing.T) {
	var (
		req = require.New(t)

		s1, s2, s3 = &wfTestStep{name: "s1"}, &wfTestStep{name: "s2"}, &wfTestStep{name: "s3"}
		gwp1, _    = NewGatewayPath(s1, gt("a", 10))
		gwp2, _    = NewGatewayPath(s2, gt("a", 5))
		gwp3, _    = NewGatewayPath(s3, gt("a", 0))

		gw, err = InclGateway(gwp1, gwp2, gwp3)
	)

	r, err := gw.Exec(context.Background(), &ExecRequest{Scope: expr.Variables{"a": 11}})
	req.NoError(err)
	req.Equal(Steps{s1, s2, s3}, r)

	r, err = gw.Exec(context.Background(), &ExecRequest{Scope: expr.Variables{"a": 6}})
	req.NoError(err)
	req.Equal(Steps{s2, s3}, r)

	r, err = gw.Exec(context.Background(), &ExecRequest{Scope: expr.Variables{"a": 1}})
	req.NoError(err)
	req.Equal(Steps{s3}, r)

	r, err = gw.Exec(context.Background(), &ExecRequest{Scope: expr.Variables{"a": 0}})
	req.Error(err)
	req.Nil(r)
}

func TestExclGateway(t *testing.T) {
	var (
		req = require.New(t)

		s1, s2, s3 = &wfTestStep{name: "s1"}, &wfTestStep{name: "s2"}, &wfTestStep{name: "s3"}
		gwp1, _    = NewGatewayPath(s1, gt("a", 10))
		gwp2, _    = NewGatewayPath(s2, gt("a", 5))
		gwp3, _    = NewGatewayPath(s3, nil)

		gw, err = ExclGateway(gwp1, gwp2, gwp3)
	)

	r, err := gw.Exec(context.Background(), &ExecRequest{Scope: expr.Variables{"a": 11}})
	req.NoError(err)
	req.Equal(s1, r)

	r, err = gw.Exec(context.Background(), &ExecRequest{Scope: expr.Variables{"a": 6}})
	req.NoError(err)
	req.Equal(s2, r)

	r, err = gw.Exec(context.Background(), &ExecRequest{Scope: expr.Variables{"a": 1}})
	req.NoError(err)
	req.Equal(s3, r)
}
