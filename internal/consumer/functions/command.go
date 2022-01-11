package functions

import (
	"bytes"
	"errors"
	"fmt"
	"forklift/internal/domain"
	"forklift/internal/domain/entity"
	"forklift/pkg/helper"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"time"
)

type Commander struct {
	taskRepository domain.MigrationTaskRepository
}

func NewCommander(taskRepository domain.MigrationTaskRepository) *Commander {
	var commander = new(Commander)
	commander.taskRepository = taskRepository
	return commander
}

func (c *Commander) killProcessIfTaskNotExists(messageId string, out chan bool) {

	var (
		err    error
		result bool
	)

	if _, err = c.taskRepository.GetMigrationByMessageId(messageId); err != nil {
		if err.Error() == "not found" {
			result = true
		}
	}
	if result {
		out <- result
	}
}

func (c *Commander) GetMigrationTask(MessageId string) (*entity.MigrationTask, error) {

	var (
		migrationTask *entity.MigrationTask
		err           error
	)

	if migrationTask, err = c.taskRepository.GetMigrationByMessageId(MessageId); err != nil {
		return nil, err
	}

	return migrationTask, err
}

func (c *Commander) executeCommand(messageId, cmd string) ([]byte, error) {

	var (
		migrationTask  *entity.MigrationTask
		stdout, stderr bytes.Buffer
		err            error
	)

	log.Infof("Executing command: (%s)", cmd)

	migrationTask, err = c.GetMigrationTask(messageId)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("An error occurred while getting message from db. %v", err))
	}

	if migrationTask.Project != nil {
		if err = os.Setenv("OS_PROJECT_NAME", migrationTask.Project.Name); err != nil {
			return nil, errors.New(fmt.Sprintf("An error occurred while setting project name environment variable. %v", err))
		}
		if err = os.Setenv("OS_PROJECT_ID", migrationTask.Project.Id); err != nil {
			return nil, errors.New(fmt.Sprintf("An error occurred while setting project name environment variable. %v", err))
		}
	} else {
		return nil, errors.New(fmt.Sprintf("An error occurred while getting value from migrationTask.Project struct. Value may be nil"))
	}

	commands := exec.Command("sh", "-c", cmd)

	commands.Stdout = &stdout
	commands.Stderr = &stderr

	if err = commands.Start(); err != nil {
		log.WithField("error.exception", err).Fatalf("An error occurred while executing command")
	}

	done := make(chan error, 1)
	go func() {
		done <- commands.Wait()
	}()

	out := make(chan bool)
	go func() {
		for {
			<-time.After(5 * time.Second)
			if IsClosed(out) {
				return
			}
			c.killProcessIfTaskNotExists(messageId, out)
			log.Infof("Waiting %s for migration %s...", migrationTask.Stage, migrationTask.VmName)
		}
	}()

	select {
	case <-out:
		if err = commands.Process.Kill(); err != nil {
			return nil, errors.New(fmt.Sprintf("Failed to kill proccess: %s", err.Error()))
		}
		return nil, errors.New(fmt.Sprintf("Process killed"))
	case err = <-done:
		close(out)
		stderr1 := helper.ReplaceNewLine(stderr.Bytes())
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Proccess finished with error: %s. %s", string(stderr1), err.Error()))
		}
	}
	beautyOutput := helper.ReplaceNewLine(stdout.Bytes())
	log.Infof("Done: (%s)", cmd)
	return beautyOutput, nil
}

func IsClosed(ch <-chan bool) bool {
	select {
	case <-ch:
		return true
	default:
	}
	return false
}
