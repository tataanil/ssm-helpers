package invocation

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	mocks "github.com/disneystreaming/ssm-helpers/testing"
	"github.com/stretchr/testify/assert"
)

func TestRunSSMCommand(t *testing.T) {
	assert := assert.New(t)

	// Set up our mocks for testing RunSSMCommand()
	mockSvc := &mocks.MockSSMClient{}
	mockInput := &ssm.SendCommandInput{
		InstanceIds:  aws.StringSlice([]string{"i-123", "i-456"}),
		DocumentName: aws.String("AWS-RunShellScript"),
		Parameters: map[string][]*string{
			/*
				ssm.SendCommandInput objects require parameters for the DocumentName chosen

				For AWS-RunShellScript, the only required parameter is "commands",
				which is the shell command to be executed on the target. To emulate
				the original script, we also set "executionTimeout" to 10 minutes.
			*/
			"commands":         aws.StringSlice([]string{"uname -a", "hostname"}),
			"executionTimeout": aws.StringSlice([]string{"600"}),
		},
	}

	t.Run("dry run flag false", func(t *testing.T) {
		output, err := RunSSMCommand(mockSvc, mockInput, false)

		assert.NoError(err)
		assert.NotNil(output)
	})

	t.Run("dry run flag true", func(t *testing.T) {
		output, err := RunSSMCommand(mockSvc, mockInput, true)

		assert.NoError(err)
		assert.Nil(output, "Dry run flag enabled, should not have received any output")
	})
}

func TestGetResult(t *testing.T) {
	assert := assert.New(t)
	mockSvc := &mocks.MockSSMClient{}

	successCmd, badCmd, mockInstance :=
		aws.String("success-id"), aws.String("bad-id"), aws.String("i-123")

	oc := make(chan *ssm.GetCommandInvocationOutput)
	ec := make(chan error)

	t.Run("valid ID", func(t *testing.T) {
		go GetResult(mockSvc, successCmd, mockInstance, oc, ec)
		select {
		case result := <-oc:
			assert.Equal("success-id", *result.CommandId)
		case err := <-ec:
			assert.Empty(err)
		}

	})

	t.Run("invalid ID", func(t *testing.T) {
		go GetResult(mockSvc, badCmd, mockInstance, oc, ec)
		select {
		case result := <-oc:
			assert.Empty(result)
		case err := <-ec:
			assert.Error(err)
		}
	})
}
