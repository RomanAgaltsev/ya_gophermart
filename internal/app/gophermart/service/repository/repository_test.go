package repository_test

import (
	"context"

	"github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/service/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Repository", func() {
	var (
		err error

		ctx context.Context

		ctrl *gomock.Controller

		mockSQL  sqlmock.Sqlmock
		mockPool pgxpoolmock.PgxPool

		repo *repository.Repository
	)

	BeforeEach(func() {
		ctx = context.Background()

		ctrl = gomock.NewController(GinkgoT())
		Expect(ctrl).ShouldNot(BeNil())

		mockPool = pgxpoolmock.NewMockPgxPool(ctrl)

		repo, err = repository.New(mockPool)
		Expect(err).ShouldNot(HaveOccurred())
	})
})
