package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

type notificationMock struct {
	message string
}

func (s *notificationMock) Send(message string) {
	s.message = message
}

func TestAlerts(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "redwood Suite")
}
