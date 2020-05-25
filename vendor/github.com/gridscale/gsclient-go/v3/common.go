package gsclient

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type emptyStruct struct {
}

//retryableFunc defines a function that can be retried
type retryableFunc func() (bool, error)

//isValidUUID validates the uuid
func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

//retryWithContext reruns a function until the context is done
func retryWithContext(ctx context.Context, targetFunc retryableFunc, delay time.Duration) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			time.Sleep(delay) //delay between retries
			continueRetrying, err := targetFunc()
			if !continueRetrying {
				return err
			}
		}
	}
}

//retryWithLimitedNumOfRetries reruns a function within a number of retries
func retryWithLimitedNumOfRetries(targetFunc retryableFunc, numOfRetries int, delay time.Duration) error {
	retryNo := 0
	var err error
	var continueRetrying bool
	for retryNo <= numOfRetries {
		retryNo++
		time.Sleep(delay * time.Duration(retryNo)) //delay between retries
		continueRetrying, err = targetFunc()
		if !continueRetrying {
			return err
		}

	}
	if err != nil {
		return fmt.Errorf("Maximum number of trials has been exhausted with error: %v", err)
	}
	return errors.New("Maximum number of trials has been exhausted")
}
