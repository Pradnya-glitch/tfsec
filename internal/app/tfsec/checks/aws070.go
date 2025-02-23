package checks

import (
	"fmt"
	"github.com/tfsec/tfsec/internal/app/tfsec/parser"
	"github.com/tfsec/tfsec/internal/app/tfsec/scanner"
)

const AWSESDomainLoggingEnabled scanner.RuleCode = "AWS070"
const AWSESDomainLoggingEnabledDescription scanner.RuleSummary = "AWS ES Domain should have logging enabled"
const AWSESDomainLoggingEnabledExplanation = `
AWS ES domain should have logging enabled by default.
`
const AWSESDomainLoggingEnabledBadExample = `
resource "aws_elasticsearch_domain" "example" {
  // other config

  // One of the log_publishing_options has to be AUDIT_LOGS
  log_publishing_options {
    cloudwatch_log_group_arn = aws_cloudwatch_log_group.example.arn
    log_type                 = "INDEX_SLOW_LOGS"
  }
}
`
const AWSESDomainLoggingEnabledGoodExample = `
resource "aws_elasticsearch_domain" "example" {
  // other config

  // At minimum we should have AUDIT_LOGS enabled
  log_publishing_options {
    cloudwatch_log_group_arn = aws_cloudwatch_log_group.example.arn
    log_type                 = "AUDIT_LOGS"
  }
}
`

func init() {
	scanner.RegisterCheck(scanner.Check{
		Code: AWSESDomainLoggingEnabled,
		Documentation: scanner.CheckDocumentation{
			Summary:     AWSESDomainLoggingEnabledDescription,
			Explanation: AWSESDomainLoggingEnabledExplanation,
			BadExample:  AWSESDomainLoggingEnabledBadExample,
			GoodExample: AWSESDomainLoggingEnabledGoodExample,
			Links: []string{
				"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/elasticsearch_domain#log_publishing_options",
			},
		},
		Provider:       scanner.AWSProvider,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"aws_elasticsearch_domain"},
		CheckFunc: func(check *scanner.Check, block *parser.Block, _ *scanner.Context) []scanner.Result {
			logPublishingOptions := block.GetBlocks("log_publishing_options")
			if len(logPublishingOptions) > 0 {
				auditLogFound := false
				for _, logPublishingOption := range logPublishingOptions {
					logType := logPublishingOption.GetAttribute("log_type")
					if logType != nil {
						if logType.Equals("AUDIT_LOGS") {
							auditLogFound = true
						}
					}
				}

				if !auditLogFound {
					return []scanner.Result{
						check.NewResult(
							fmt.Sprintf("Resource '%s' is missing 'AUDIT_LOGS` in one of the `log_publishing_options`-`log_type` attributes so audit log is not enabled", block.FullName()),
							block.Range(),
							scanner.SeverityError,
						),
					}
				}
			}

			return nil
		},
	})
}
