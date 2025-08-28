package main

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Otp struct {
	Key       string
	CreatedAt time.Time
}

type RetentionMap map[string]Otp

func NewRetentionMap(retentionPeriod time.Duration) RetentionMap {
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	rm := make(RetentionMap)

	go rm.Retention(retentionPeriod)

	return rm
}

func (rm RetentionMap) NewOtp() Otp {
	otp := Otp{
		Key:       uuid.NewString(),
		CreatedAt: time.Now(),
	}

	rm[otp.Key] = otp
	return otp
}

func (rm RetentionMap) VerifyOtp(otp string) bool {
	if _, ok := rm[otp]; !ok {
		return false
	}

	delete(rm, otp)
	return true
}

func (rm RetentionMap) Retention(retentionPeriod time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ticker := time.NewTicker(400 * time.Millisecond)

	for {
		select {
		case <-ticker.C:
			for _, otp := range rm {
				if otp.CreatedAt.Add(retentionPeriod).Before(time.Now()) {
					delete(rm, otp.Key)
				}
			}
		case <-ctx.Done():
			return
		}

	}
}
