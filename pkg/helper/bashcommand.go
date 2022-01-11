package helper

import (
	"bytes"
	"errors"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

func ExecuteCommand(cmd string) ([]byte, error) {
	log.Infof("Executing command: (%s)", cmd)
	execute := exec.Command("sh", "-c", cmd)
	out, err := execute.CombinedOutput()
	if err != nil {
		return nil, errors.New(string(out) + ". " + err.Error())
	}
	out = ReplaceNewLine(out)
	log.Infof("Done: (%s)", cmd)
	return out, nil
}

func ReplaceNewLine(output []byte) []byte {
	newLineCount := countRune(string(output),'\n')
	tabCount := countRune(string(output),'\r')
	lines := bytes.Replace(output, []byte("\n"), []byte(" "), newLineCount-1)
	lines = bytes.Replace(lines, []byte("\n"), []byte(""), 1)
	lines = bytes.Replace(lines, []byte("\r"), []byte(" "), tabCount-1)
	lines = bytes.Replace(lines, []byte("\r"), []byte(""), 1)
	return lines
}

func countRune(s string, r rune) int {
	count := 0
	for _, c := range s {
		if c == r {
			count++
		}
	}
	return count
}