package main

import (
	"io"
	"os"
	"os/exec"
)

const sendmailcmd = "/usr/sbin/sendmail"

func sendmail(email, object, body string) {

	msg := "From: contact@macha.fr\n"
	msg += "To: " + email + "\n"
	msg += "Subject: " + object + "\n\n"
	msg += body
	msg += ".\n"

	cmd := exec.Command(sendmailcmd, "-t")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	checkErr(err, "sendmail")

	err = cmd.Start()
	checkErr(err, "sendmail")

	var errs [3]error
	_, errs[0] = io.WriteString(stdin, msg)
	errs[1] = stdin.Close()
	errs[2] = cmd.Wait()
	for _, err = range errs {
		checkErr(err, "sendmail")
	}
}
