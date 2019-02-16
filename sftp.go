package sync

import (
	"fmt"
	"github.com/kr/fs"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"net/url"
	"os"
	"strconv"
)

type Host struct {
	Location string
	Port int
	User string
	Password string
}

type SFTPTarget struct {
	ssh *ssh.Client
	sftp *sftp.Client
	h Host
}

func NewSFTPTarget(h Host) *SFTPTarget {
	return &SFTPTarget{h: h}
}

func (t *SFTPTarget) Connect() error {
	if t.h.Port == 0 {
		t.h.Port = 22
	}
	cfg := &ssh.ClientConfig{
		User:            t.h.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.Password(t.h.Password),
		},
	}
	cfg.SetDefaults()
	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%d", t.h.Location, t.h.Port), cfg)
	if err != nil {
		return fmt.Errorf("failed establishing SSH connection: %v", err)
	}
	t.ssh = client

	sftp, err := sftp.NewClient(t.ssh)
	if err != nil {
		return fmt.Errorf("failed startin SFTP session: %v", err)
	}
	t.sftp = sftp

	return nil
}

func (t *SFTPTarget) Close() error {
	t.sftp.Close()
	t.ssh.Close()
	return nil
}

func (t *SFTPTarget) Write(r io.Reader, path string) error {

	dst, err := t.sftp.Create(path)
	if err != nil {
		return fmt.Errorf("failed opening target file for writing: %v", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, r)
	if err != nil {
		return fmt.Errorf("failed writing file: %v", err)
	}
	return nil
}

func (t *SFTPTarget) Exists(path string) (bool, error) {
	_, err := t.sftp.Stat(path)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (t *SFTPTarget) CreateDirAll(dir string) error {
	// TODO: Read up on what mkdir does when dir already exists.
	return t.sftp.MkdirAll(dir)
}

func (t *SFTPTarget) Remove(path string) error {
	return t.sftp.Remove(path)
}

func (t *SFTPTarget) RemoveDir(dir string) error {
	return t.sftp.Remove(dir)
}

func (t *SFTPTarget) Walk(root string) *fs.Walker {
	return t.sftp.Walk(root)
}

func ParseTargetURI(uri string) Host {
	u, err := url.Parse(uri)
	if err != nil {
		panic(fmt.Errorf("failed parsing target URI: %v", err))
	}

	host, p, err := net.SplitHostPort(u.Host)
	if err != nil {
		panic(fmt.Errorf("failed parsing target URI: %v", err))
	}

	port := 0
	if p != "" {
		port, err = strconv.Atoi(p)
		if err != nil {
			panic(fmt.Errorf("failed parsing target port: %v", err))
		}
	}

	user := u.User.Username()
	password, _ := u.User.Password()

	return Host{
		Location: host,
		User: user,
		Password: password,
		Port: port,
	}
}
