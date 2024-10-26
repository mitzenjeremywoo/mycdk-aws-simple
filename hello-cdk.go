package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewHelloCdkStack(app, "HelloCdkStack-", &HelloCdkStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return &awscdk.Environment{
		Account: jsii.String(os.Getenv("ACCOUNT")),
		Region:  jsii.String("ap-southeast-2"),
	}
}

type HelloCdkStackProps struct {
	awscdk.StackProps
}

func NewHelloCdkStack(scope constructs.Construct, id string, props *HelloCdkStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	awss3.NewBucket(stack, jsii.String("kepung_bucket"), &awss3.BucketProps{})

	// Define the Lambda function resource
	myFunction := awslambda.NewFunction(stack, jsii.String("HelloWorldFunction"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_NODEJS_20_X(), // Provide any supported Node.js runtime
		Handler: jsii.String("index.handler"),
		Code: awslambda.Code_FromInline(jsii.String(`
		exports.handler = async function(event) {
		  return {
			statusCode: 200,
			body: JSON.stringify('Hello CDK yay!'),
		  };
		};
	  `)),
	})

	myFunctionUrl := myFunction.AddFunctionUrl(&awslambda.FunctionUrlOptions{
		AuthType: awslambda.FunctionUrlAuthType_NONE,
	})

	// Define a CloudFormation output for your URL
	awscdk.NewCfnOutput(stack, jsii.String("myFunctionUrlOutput"), &awscdk.CfnOutputProps{
		Value: myFunctionUrl.Url(),
	})

	return stack
}
