package ingestion

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"github.com/aws/smithy-go"
	"github.com/nishant-raj/multi_cloud_optimizer/apps/api/internal/store"
)

type PublicError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Hint    string `json:"hint,omitempty"`
	Details string `json:"details,omitempty"`
}

type SourceError struct {
	public PublicError
	err    error
}

func (e *SourceError) Error() string {
	if e == nil {
		return ""
	}
	if e.err == nil {
		return e.public.Message
	}
	return fmt.Sprintf("%s: %v", e.public.Message, e.err)
}

func (e *SourceError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.err
}

func (e *SourceError) PublicError() PublicError {
	if e == nil {
		return PublicError{}
	}
	return e.public
}

func AsPublicError(err error) (PublicError, bool) {
	var sourceErr *SourceError
	if errors.As(err, &sourceErr) {
		return sourceErr.PublicError(), true
	}
	return PublicError{}, false
}

type AWSCostExplorerSource struct {
	region           string
	defaultAccountID string
}

func NewAWSCostExplorerSource(region string, defaultAccountID string) *AWSCostExplorerSource {
	return &AWSCostExplorerSource{
		region:           region,
		defaultAccountID: strings.TrimSpace(defaultAccountID),
	}
}

func (s *AWSCostExplorerSource) Fetch(ctx context.Context, input FetchInput) ([]store.BillingRecord, error) {
	accountID := strings.TrimSpace(input.AccountID)
	if accountID == "" {
		accountID = s.defaultAccountID
	}

	if input.Days <= 0 {
		input.Days = 7
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(s.region))
	if err != nil {
		return nil, newSourceError(
			"aws_credentials_unavailable",
			"AWS credentials are not available to the API.",
			"Configure AWS credentials for the API process with environment variables, AWS_PROFILE, or an attached IAM role.",
			err,
		)
	}

	client := costexplorer.NewFromConfig(awsCfg)
	start, end := costExplorerWindow(input.Days)
	records := make([]store.BillingRecord, 0, input.Days*4)
	var nextPageToken *string

	for {
		request := &costexplorer.GetCostAndUsageInput{
			Granularity: types.GranularityDaily,
			Metrics:     []string{"UnblendedCost"},
			GroupBy: []types.GroupDefinition{
				{Key: stringPtr("SERVICE"), Type: types.GroupDefinitionTypeDimension},
			},
			TimePeriod: &types.DateInterval{
				Start: stringPtr(start.Format("2006-01-02")),
				End:   stringPtr(end.Format("2006-01-02")),
			},
			NextPageToken: nextPageToken,
		}

		if accountID != "" {
			request.Filter = &types.Expression{
				Dimensions: &types.DimensionValues{
					Key:          types.DimensionLinkedAccount,
					Values:       []string{accountID},
					MatchOptions: []types.MatchOption{types.MatchOptionEquals},
				},
			}
		}

		response, err := client.GetCostAndUsage(ctx, request)
		if err != nil {
			return nil, classifyAWSRequestError(err, accountID, s.region)
		}

		batch, err := costExplorerResultsToRecords(response.ResultsByTime, accountID)
		if err != nil {
			return nil, err
		}
		records = append(records, batch...)

		if response.NextPageToken == nil || strings.TrimSpace(*response.NextPageToken) == "" {
			break
		}
		nextPageToken = response.NextPageToken
	}

	return records, nil
}

// Cost Explorer expects an exclusive end date, so the window runs from the first
// included UTC day through the next UTC midnight after today.
func costExplorerWindow(days int) (time.Time, time.Time) {
	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	start := today.AddDate(0, 0, -(days - 1))
	end := today.AddDate(0, 0, 1)
	return start, end
}

// Cost Explorer returns daily buckets grouped by service. This converts each
// group into the same BillingRecord model used by the synthetic ingestion path.
func costExplorerResultsToRecords(results []types.ResultByTime, accountID string) ([]store.BillingRecord, error) {
	records := make([]store.BillingRecord, 0, len(results)*4)

	for _, result := range results {
		usageDate, err := time.Parse("2006-01-02", valueOrEmpty(result.TimePeriod.Start))
		if err != nil {
			return nil, fmt.Errorf("parse cost explorer date %q: %w", valueOrEmpty(result.TimePeriod.Start), err)
		}

		for groupIndex, group := range result.Groups {
			amountMetric, ok := group.Metrics["UnblendedCost"]
			if !ok {
				continue
			}

			amount, err := strconv.ParseFloat(valueOrEmpty(amountMetric.Amount), 64)
			if err != nil {
				return nil, fmt.Errorf("parse cost amount %q: %w", valueOrEmpty(amountMetric.Amount), err)
			}

			serviceName := "Unknown Service"
			if len(group.Keys) > 0 && strings.TrimSpace(group.Keys[0]) != "" {
				serviceName = group.Keys[0]
			}

			records = append(records, store.BillingRecord{
				ID:        fmt.Sprintf("aws-%s-%d-%d", sanitizeAWSRecordID(accountID), usageDate.Unix(), groupIndex),
				AccountID: accountID,
				Service:   serviceName,
				UsageDate: usageDate.UTC(),
				Amount:    amount,
				Currency:  defaultString(valueOrEmpty(amountMetric.Unit), "USD"),
				Source:    "aws",
				Scenario:  "live_cost_explorer",
			})
		}
	}

	return records, nil
}

func sanitizeAWSRecordID(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "all-accounts"
	}

	return strings.NewReplacer(" ", "-", "/", "-", "_", "-").Replace(strings.ToLower(trimmed))
}

func defaultString(value string, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}

	return value
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}

func stringPtr(value string) *string {
	return &value
}

func newSourceError(code string, message string, hint string, err error) error {
	details := ""
	if err != nil {
		details = err.Error()
	}

	return &SourceError{
		public: PublicError{
			Code:    code,
			Message: message,
			Hint:    hint,
			Details: details,
		},
		err: err,
	}
}

func classifyAWSRequestError(err error, accountID string, region string) error {
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		code := strings.ToLower(strings.TrimSpace(apiErr.ErrorCode()))
		switch code {
		case "accessdeniedexception", "accessdenied", "unauthorizedoperation":
			return newSourceError(
				"aws_cost_explorer_access_denied",
				"The API credentials do not have permission to read AWS Cost Explorer data.",
				"Grant Cost Explorer read permissions such as ce:GetCostAndUsage to the IAM principal used by the API.",
				err,
			)
		case "billexpirationexception", "datanotavailableexception":
			return newSourceError(
				"aws_cost_data_unavailable",
				"AWS Cost Explorer returned no usable cost data for the requested window.",
				"Check that Cost Explorer is enabled for the account and that the selected date range contains finalized billing data.",
				err,
			)
		case "validationexception":
			hint := "Check the AWS account ID and lookback days in the request."
			if strings.TrimSpace(accountID) == "" {
				hint = "Provide an AWS account ID or set AWS_COST_EXPLORER_ACCOUNT_ID so the API can scope the Cost Explorer query."
			}
			return newSourceError(
				"aws_cost_explorer_validation_failed",
				"AWS Cost Explorer rejected the ingestion request.",
				hint,
				err,
			)
		}
	}

	lower := strings.ToLower(err.Error())
	switch {
	case strings.Contains(lower, "failed to refresh cached credentials"),
		strings.Contains(lower, "no valid credential sources"),
		strings.Contains(lower, "credential"):
		return newSourceError(
			"aws_credentials_unavailable",
			"AWS credentials are missing, expired, or unreadable by the API.",
			"Verify AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY, AWS_PROFILE, or the IAM role attached to the runtime.",
			err,
		)
	case strings.Contains(lower, "security token included in the request is invalid"),
		strings.Contains(lower, "invalidclienttokenid"),
		strings.Contains(lower, "expiredtoken"):
		return newSourceError(
			"aws_credentials_invalid",
			"The AWS credentials used by the API are invalid or expired.",
			"Refresh the credentials or session token used by the API process and try again.",
			err,
		)
	case strings.Contains(lower, "could not resolve endpoint"),
		strings.Contains(lower, "unknown endpoint"),
		strings.Contains(lower, "endpoint resolution"):
		return newSourceError(
			"aws_region_misconfigured",
			fmt.Sprintf("The configured AWS region %q could not be used for Cost Explorer.", region),
			"Set AWS_REGION to a valid region and ensure the runtime can reach AWS endpoints.",
			err,
		)
	case strings.Contains(lower, "request canceled"),
		strings.Contains(lower, "deadline exceeded"),
		strings.Contains(lower, "timeout"):
		return newSourceError(
			"aws_request_timed_out",
			"The request to AWS Cost Explorer timed out.",
			"Retry the ingestion and verify outbound network access from the API container or host.",
			err,
		)
	default:
		return newSourceError(
			"aws_cost_explorer_request_failed",
			"AWS Cost Explorer ingestion failed.",
			"Check API credentials, Cost Explorer permissions, region configuration, and network access.",
			err,
		)
	}
}
