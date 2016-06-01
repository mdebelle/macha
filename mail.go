package main

import (
	"io"
	"os"
	"os/exec"
)

const sendmailcmd = "/usr/sbin/sendmail"

func sendmail(email, body string) {

	msg := "From: contact@macha.fr\n"
	msg += "To: " + email + "\n"
	msg += "Subject: Testing\n\n"
	msg += body
	msg += ".\n"

	cmd := exec.Command(sendmailcmd, "-t")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	checkErr(err)

	err = cmd.Start()
	checkErr(err)

	var errs [3]error
	_, errs[0] = io.WriteString(stdin, msg)
	errs[1] = stdin.Close()
	errs[2] = cmd.Wait()
	for _, err = range errs {
		checkErr(err)
	}
}
