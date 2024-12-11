package config_test

import (
	"flag"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"

	"github.com/RomanAgaltsev/ya_gophermart/internal/config"
)

var _ = Describe("Config", func() {
	var cfg *config.Config
	var err error

	defaultArgs := os.Args
	defaultCommandLine := flag.CommandLine

	JustBeforeEach(func() {
		cfg, err = config.Get()
	})

	Describe("Testing run address", func() {
		When("env is set and flag is set", func() {
			BeforeEach(func() {
				t := GinkgoT()
				t.Setenv("RUN_ADDRESS", "localhost:8081")

				flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
				os.Args = append([]string{"cmd"}, "-a", "localhost:8091")
			})
			AfterEach(func() {
				t := GinkgoT()
				t.Setenv("RUN_ADDRESS", "")
				os.Args = defaultArgs
				flag.CommandLine = defaultCommandLine
			})

			It("is env used for config", func() {
				Expect(err).Should(BeNil())
				Expect(cfg.RunAddress).To(Equal("localhost:8081"))
			})
		})
		When("env is set and flag is unset", func() {
			BeforeEach(func() {
				t := GinkgoT()
				t.Setenv("RUN_ADDRESS", "localhost:8081")
			})
			AfterEach(func() {
				t := GinkgoT()
				t.Setenv("RUN_ADDRESS", "")
				os.Args = defaultArgs
				flag.CommandLine = defaultCommandLine
			})

			It("is env used fo config", func() {
				Expect(err).Should(BeNil())
				Expect(cfg.RunAddress).To(Equal("localhost:8081"))
			})
		})
		When("env is unset and flag is set", func() {
			BeforeEach(func() {
				flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
				os.Args = append([]string{"cmd"}, "-a", "localhost:8091")
			})
			AfterEach(func() {
				os.Args = defaultArgs
				flag.CommandLine = defaultCommandLine
			})

			It("is flag used for config", func() {
				Expect(err).Should(BeNil())
				Expect(cfg.RunAddress).To(Equal("localhost:8091"))
			})
		})
		When("env is unset and flag is unset", func() {
			It("is default used for config", func() {
				Expect(err).Should(BeNil())
				Expect(cfg.RunAddress).To(Equal("localhost:8080"))
			})
		})
	})

	Describe("Testing database URI", func() {
		When("env is set and flag is set", func() {
			It("is env used for config", func() {
				Expect(cfg.DatabaseURI).To(Equal(""))
			})
		})
		When("env is set and flag is unset", func() {
			It("is env used fo config", func() {
				Expect(cfg.DatabaseURI).To(Equal(""))
			})
		})
		When("env is unset and flag is set", func() {
			It("is flag used for config", func() {
				Expect(cfg.DatabaseURI).To(Equal(""))
			})
		})
		When("env is unset and flag is unset", func() {
			It("is default used for config", func() {
				Expect(cfg.DatabaseURI).To(Equal(""))
			})
		})
	})

	Describe("Testing address of accrual system", func() {
		When("env is set and flag is set", func() {
			It("is env used for config", func() {
				Expect(cfg.AccrualSystemAddress).To(Equal(""))
			})
		})
		When("env is set and flag is unset", func() {
			It("is env used fo config", func() {
				Expect(cfg.AccrualSystemAddress).To(Equal(""))
			})
		})
		When("env is unset and flag is set", func() {
			It("is flag used for config", func() {
				Expect(cfg.AccrualSystemAddress).To(Equal(""))
			})
		})
		When("env is unset and flag is unset", func() {
			It("is default used for config", func() {
				Expect(cfg.AccrualSystemAddress).To(Equal(""))
			})
		})
	})

	Describe("Testing secret key", func() {
		When("env is set and flag is set", func() {
			It("is env used for config", func() {
				Expect(cfg.SecretKey).To(Equal(""))
			})
		})
		When("env is set and flag is unset", func() {
			It("is env used fo config", func() {
				Expect(cfg.SecretKey).To(Equal(""))
			})
		})
		When("env is unset and flag is set", func() {
			It("is flag used for config", func() {
				Expect(cfg.SecretKey).To(Equal(""))
			})
		})
		When("env is unset and flag is unset", func() {
			It("is default used for config", func() {
				Expect(cfg.SecretKey).To(Equal(""))
			})
		})
	})
})
