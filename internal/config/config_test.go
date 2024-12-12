package config_test

import (
	"flag"
	"os"

	"github.com/RomanAgaltsev/ya_gophermart/internal/config"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*
struct{
	envName string
	envVal string
	flgName string
	flgVal string
	default string
}
*/

var _ = Describe("Config", Ordered, func() {
	var cfg *config.Config
	var err error

	defaultArgs := os.Args
	defaultCommandLine := flag.CommandLine

	BeforeAll(func() {

	})

	AfterEach(func() {
		os.Clearenv()
		os.Args = defaultArgs
		flag.CommandLine = defaultCommandLine
	})

	DescribeTable("Run address",
		func(envName, envVal, flgName, flgVal, def, expected string) {
			if envName != "" {
				t := GinkgoT()
				t.Setenv(envName, envVal)
			}

			if flgName != "" {
				flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
				os.Args = append([]string{"cmd"}, flgName, flgVal)
			}

			cfg, err = config.Get()

			Expect(err).Should(BeNil())
			Expect(cfg.RunAddress).To(Equal(expected))
		},

		EntryDescription("When env %s=%s, flag %s=%s and default=%s"),
		Entry(nil, "RUN_ADDRESS", "localhost:8081", "-a", "localhost:8091", "localhost:8080", "localhost:8081"),
		Entry(nil, "RUN_ADDRESS", "localhost:8081", "-a", "", "localhost:8080", "localhost:8081"),
		Entry(nil, "RUN_ADDRESS", "localhost:8081", "", "", "localhost:8080", "localhost:8081"),
		Entry(nil, "RUN_ADDRESS", "", "-a", "localhost:8091", "localhost:8080", "localhost:8091"),
		Entry(nil, "", "", "-a", "localhost:8091", "localhost:8080", "localhost:8091"),
		Entry(nil, "RUN_ADDRESS", "", "-a", "", "localhost:8080", ""),
		Entry(nil, "", "", "-a", "", "localhost:8080", ""),
		Entry(nil, "RUN_ADDRESS", "", "-a", "", "localhost:8080", ""),
		Entry(nil, "RUN_ADDRESS", "", "", "", "localhost:8080", "localhost:8080"),
		Entry(nil, "", "", "", "", "localhost:8080", "localhost:8080"),
		Entry(nil, "", "", "", "", "localhost:8080", "localhost:8080"),
	)
})
