package cid

import (
	"context"
	"errors"
	"os"
	"strconv"

	"github.com/appleboy/easyssh-proxy"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SSHClient struct {
	sshConf   *easyssh.MakeConfig
	scp       *sftp.Client
	session   *ssh.Session
	client    *ssh.Client
	connected bool
}

func NewClient(ctx context.Context, cid *CloudInitDriveProviderModel) (*SSHClient, diag.Diagnostics) {

	c := &SSHClient{
		sshConf: &easyssh.MakeConfig{
			User:       cid.SSH.User.ValueString(),
			Server:     cid.SSH.Host.ValueString(),
			Port:       strconv.FormatInt(cid.SSH.Port.ValueInt64(), 10),
			Password:   os.Getenv("SSH_PASSWORD"),
			Key:        cid.SSH.PrivateKey.ValueString(),
			Passphrase: os.Getenv("SSH_PRIVATE_KEY_PASSPHRASE"),
		},
	}

	// "github.com/melbahja/goph"
	// "github.com/xanzy/ssh-agent"

	// 	if b := d.Get("become"); b != nil {
	// 		m := b.(map[string]interface{})
	// 		user, ok := m["user"].(string)
	// 		if !ok {
	// 			return nil, diag.Diagnostics{diag.Diagnostic{
	// 				Severity: diag.Error,
	// 				Summary:  "The username for privilege escalation was not set.",
	// 			}}
	// 		}
	// 		method, ok := m["method"].(string)
	// 		if !ok {
	// 			return nil, diag.Diagnostics{diag.Diagnostic{
	// 				Severity: diag.Error,
	// 				Summary:  "The method for privilege escalation is not set.",
	// 			}}
	// 		}
	// 		// Can only be set to an environment variable!
	// 		pass := os.Getenv("CID_BECOME_PASSWORD")
	// 		if method == "su" && pass == "" {
	// 			return nil, diag.Diagnostics{diag.Diagnostic{
	// 				Severity: diag.Error,
	// 				Summary:  "The password for the 'su' command is not set.",
	// 			}}
	// 		}
	// 	}
	return c, nil
}

func (c *SSHClient) Connect() diag.Diagnostic {
	var err error
	c.session, c.client, err = c.sshConf.Connect()
	if err != nil {
		return diag.NewErrorDiagnostic("SSH connection error", err.Error())
	}

	c.scp, err = sftp.NewClient(c.client)
	if err != nil {
		return diag.NewErrorDiagnostic("SFTP session error", err.Error())
	}

	c.connected = true

	return nil
}

func (c *SSHClient) Close() {
	if c.scp != nil {
		c.scp.Close()
	}
	if c.session != nil {
		c.session.Close()
	}
	if c.client != nil {
		c.client.Close()
	}
	c.connected = false
}

func (c *SSHClient) RemoteWrite(remoteFile string) (*sftp.File, error) {
	if !c.connected {
		return nil, errors.New("error writing the file to the remote host: no active SSH connection found")
	}

	dst, err := c.scp.Create(remoteFile)
	if err != nil {
		return nil, err
	}

	return dst, nil
}

func (c *SSHClient) RemoteRead(remoteFile string) (*sftp.File, error) {
	if !c.connected {
		return nil, errors.New("error reading the file on the remote host: no active SSH connection found")
	}

	src, err := c.scp.Open(remoteFile)
	if err != nil {
		return nil, err
	}

	return src, nil
}
