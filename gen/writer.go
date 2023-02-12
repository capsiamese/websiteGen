package main

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"mdgen/gui"
	"os"
)

// term.ReadPassword(int(syscall.Stdin))

type FileWriter interface {
	Open(name string) (io.WriteCloser, error)
	MkdirAll(dir string) error
	PostRun() error
	Close() error
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

type RemoteWriter struct {
	addr       string
	sshClient  *ssh.Client
	sftpClient *sftp.Client
}

func NewRemoteWriter() *RemoteWriter {
	return &RemoteWriter{}
}

func (rw *RemoteWriter) tryAuthMethod(keyPath, password string) []ssh.AuthMethod {
	data, keyPathErr := os.ReadFile(keyPath)
	if keyPathErr != nil {
		log.Println("[warning] read ssh key from", keyPath, keyPathErr)
	} else {
		signer, parseKeyErr := ssh.ParsePrivateKey(data)
		if parseKeyErr != nil {
			log.Println("[warning] parse ssh key", keyPath, parseKeyErr)
		} else {
			return []ssh.AuthMethod{ssh.PublicKeys(signer)}
		}
	}
	return []ssh.AuthMethod{ssh.PasswordCallback(func() (secret string, err error) {
		return password, nil
	})}
}

func (rw *RemoteWriter) Connect(d *gui.Data) error {
	var cfg ssh.ClientConfig
	cfg.Auth = rw.tryAuthMethod(d.KeyFile, d.Password)
	cfg.User = d.User
	cfg.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	addr := fmt.Sprintf("%s:%s", d.RemoteAddr, d.RemotePort)
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

type CloserLog struct {
	wc   io.WriteCloser
	name string
}

func (cl CloserLog) Write(p []byte) (int, error) {
	return cl.wc.Write(p)
}
func (cl CloserLog) Close() error {
	log.Println("[info] file", cl.name, "write done")
	return cl.wc.Close()
}

func (rw *RemoteWriter) Open(name string) (io.WriteCloser, error) {
	f, err := rw.sftpClient.OpenFile(name, openMode)
	if err != nil {
		return nil, err
	}
	log.Println("[info] open file", name)
	return CloserLog{
		wc: f, name: name,
	}, nil
}

func (rw *RemoteWriter) MkdirAll(dir string) error {
	log.Println("[info] mkdir all", dir)
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
