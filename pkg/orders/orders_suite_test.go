package orders_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestOrders(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Orders Suite")
}
