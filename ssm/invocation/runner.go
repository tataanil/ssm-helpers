package invocation

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

// RunSSMCommand uses an SSM session, pre-defined SSM document parameters, the dry run flag, and any number of instance IDs and executes the given command
// using the AWS-RunShellScript SSM document. It returns an *ssm.SendCommandOutput object, which contains the execution ID of the command, which we use to
// check the progress/status of the invocation.
func RunSSMCommand(session ssmiface.SSMAPI, input *ssm.SendCommandInput, dryRunFlag bool) (scOutput *ssm.SendCommandOutput, err error) {
	if !dryRunFlag {
		return session.SendCommand(input)
	}

	return
}

// GetResult fetches the output of an invocation using the associated command and instance IDs
func GetResult(client ssmiface.SSMAPI, commandID *string, instanceID *string, gci chan *ssm.GetCommandInvocationOutput, ec chan error) {
	status, err := client.GetCommandInvocation(&ssm.GetCommandInvocationInput{
		CommandId:  commandID,
		InstanceId: instanceID,
	})

	switch {
	case err != nil:
		ec <- fmt.Errorf(
			`Error when calling GetCommandInvocation API with args:\n
			CommandId: %v\n
			InstanceId: %v\n%v`,
			*commandID, *instanceID, err)
	case status != nil:
		gci <- status
	}

}
