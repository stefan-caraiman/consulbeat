// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package cloudformation

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/private/protocol/query"
)

// CloudFormation provides the API operation methods for making requests to
// AWS CloudFormation. See this package's package overview docs
// for details on the service.
//
// CloudFormation methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type CloudFormation struct {
	*aws.Client
}

// Used for custom client initialization logic
var initClient func(*CloudFormation)

// Used for custom request initialization logic
var initRequest func(*CloudFormation, *aws.Request)

// Service information constants
const (
	ServiceName = "cloudformation" // Service endpoint prefix API calls made to.
	EndpointsID = ServiceName      // Service ID for Regions and Endpoints metadata.
)

// New creates a new instance of the CloudFormation client with a config.
//
// Example:
//     // Create a CloudFormation client from just a config.
//     svc := cloudformation.New(myConfig)
func New(config aws.Config) *CloudFormation {
	var signingName string
	signingRegion := config.Region

	svc := &CloudFormation{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   ServiceName,
				SigningName:   signingName,
				SigningRegion: signingRegion,
				APIVersion:    "2010-05-15",
			},
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(query.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(query.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(query.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(query.UnmarshalErrorHandler)

	// Run custom client initialization if present
	if initClient != nil {
		initClient(svc)
	}

	return svc
}

// newRequest creates a new request for a CloudFormation operation and runs any
// custom request initialization.
func (c *CloudFormation) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	// Run custom request initialization if present
	if initRequest != nil {
		initRequest(c, req)
	}

	return req
}
