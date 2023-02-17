package process

import (
	"fmt"
	"generator/config"
	"generator/rec"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"io/fs"
	"log"
	"os"
)

// term.ReadPassword(int(syscall.Stdin))

var _ FileWriter = (*LocalWriter)(nil)
var _ FileWriter = (*RemoteWriter)(nil)

type FileWriter interface {
	Open(name string) (io.WriteCloser, error)
	MkdirAll(dir string) error
	PostRun() error
	Close() error
	Stat(name string) (fs.FileInfo, error)
}

type LocalWriter struct{}

func (LocalWriter) Open(name string) (io.WriteCloser, error) {
	return os.OpenFile(name, openMode, 0644)
}
func (LocalWriter) MkdirAll(dir string) error {
	return os.MkdirAll(dir, 0644)
}
func (LocalWriter) PostRun() error {
	return nil
}
func (LocalWriter) Close() error {
	return nil
}

func (LocalWriter) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(name)
}

type RemoteWriter struct {
	addr       string
	sshClient  *ssh.Client
	sftpClient *sftp.Client
}

func NewRemoteWriter() *RemoteWriter {
	return &RemoteWriter{}
}

func (rw *RemoteWriter) tryAuthMethod(keyStr, keyPath, password string) []ssh.AuthMethod {
	if keyStr != "" {
		signer, parseKeyErr := ssh.ParsePrivateKey([]byte(keyStr))
		if parseKeyErr != nil {
			rec.Writeln("[warning] read ssh key from", keyPath, parseKeyErr)
		} else {
			return []ssh.AuthMethod{ssh.PublicKeys(signer)}
		}
	}
	data, keyPathErr := os.ReadFile(keyPath)
	if keyPathErr != nil {
		rec.Writeln("[warning] read ssh key from", keyPath, keyPathErr)
	} else {
		signer, parseKeyErr := ssh.ParsePrivateKey(data)
		if parseKeyErr != nil {
			rec.Writeln("[warning] parse ssh key", keyPath, parseKeyErr)
		} else {
			return []ssh.AuthMethod{ssh.PublicKeys(signer)}
		}
	}
	return []ssh.AuthMethod{ssh.PasswordCallback(func() (secret string, err error) {
		return password, nil
	})}
}

func (rw *RemoteWriter) Connect(d *config.Data) error {
	var cfg ssh.ClientConfig
	cfg.Auth = rw.tryAuthMethod(d.KeyStr, d.KeyFile, d.Password)
	cfg.User = d.User
	cfg.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	addr := fmt.Sprintf("%s:%s", d.RemoteAddr, d.RemotePort)
	rec.Writeln("ssh dial addr:", addr)
	rec.Writeln("password len:", len(d.Password))
	sshClient, err := ssh.Dial("tcp", addr, &cfg)
	if err != nil {
		return err
	}
	rw.sshClient = sshClient

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return err
	}
	rw.sftpClient = sftpClient
	return nil
}

func (rw *RemoteWriter) Open(name string) (io.WriteCloser, error) {
	f, err := rw.sftpClient.OpenFile(name, openMode)
	if err != nil {
		return nil, err
	}
	rec.Writeln("[info] open file", name)
	return f, nil
}

func (rw *RemoteWriter) MkdirAll(dir string) error {
	rec.Writeln("[info] mkdir all", dir)
	return rw.sftpClient.MkdirAll(dir)
}

func (rw *RemoteWriter) PostRun() error {
	res, err := rw.Cmd(`mkhomepg -p`)
	if err != nil {
		return err
	}
	log.Println(res)
	return nil
}

func (rw *RemoteWriter) Cmd(str string) (string, error) {
	session, err := rw.sshClient.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	res, err := session.Output(str)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func (rw *RemoteWriter) Close() error {
	err1 := rw.sftpClient.Close()
	err2 := rw.sshClient.Close()
	if err1 == nil && err2 == nil {
		return nil
	}
	return fmt.Errorf("remote writer close %v %v", err1, err2)
}

func (rw *RemoteWriter) Stat(name string) (fs.FileInfo, error) {
	return rw.sftpClient.Stat(name)
}
