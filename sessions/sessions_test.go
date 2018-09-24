package sessions

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/percona/rds_exporter/client"
	"github.com/percona/rds_exporter/config"
)

func TestSession(t *testing.T) {
	cfg := &config.Config{
		Instances: []config.Instance{
			{
				Region:   "us-east-1",
				Instance: "rds-aurora1",
				// no explicit key
			},
			{
				Region:       "us-east-1",
				Instance:     "rds-aurora57",
				AWSAccessKey: os.Getenv("AWS_ACCESS_KEY_ID"),
				AWSSecretKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
			},
			{
				Region:   "us-east-1",
				Instance: "rds-mysql56",
				// no explicit key
			},
			{
				Region:   "us-west-1",
				Instance: "rds-mysql57",
				// no explicit key
			},
		},
	}

	client := client.New()
	sessions, err := New(cfg.Instances, client.HTTP(), false)
	require.NoError(t, err)

	a1s := sessions.GetSession("us-east-1", "rds-aurora1")
	a57s := sessions.GetSession("us-east-1", "rds-aurora57")
	m56s := sessions.GetSession("us-east-1", "rds-mysql56")
	m57s := sessions.GetSession("us-west-1", "rds-mysql57")

	if a1s == a57s {
		assert.Fail(t, "rds-aurora1 and rds-aurora57 should not share session (different keys)")
	}
	if a1s != m56s {
		assert.Fail(t, "rds-aurora1 and rds-mysql56 should share session")
	}
	if m57s != nil {
		assert.Fail(t, "rds-mysql57 should be skipped")
	}

	all := sessions.AllSessions()
	assert.Equal(t, all, map[*session.Session][]Instance{
		a1s: {
			{Region: "us-east-1", Instance: "rds-aurora1", ResourceID: "db-P5QCHK64NWDD5BLLBVT5NPQS2Q"},
			{Region: "us-east-1", Instance: "rds-mysql56", ResourceID: "db-J6JH3LJAWBZ6MXDDWYRG4RRJ6A"},
		},
		a57s: {
			{Region: "us-east-1", Instance: "rds-aurora57", ResourceID: "db-CDBSN4EK5SMBQCSXI4UPZVF3W4"},
		},
	})
}
