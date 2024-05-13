package staging_test

import (
	"code.cloudfoundry.org/cnbapplifecycle/pkg/staging"
	"github.com/buildpacks/lifecycle/buildpack"
	"github.com/buildpacks/lifecycle/launch"
	"github.com/buildpacks/lifecycle/platform/files"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("StagingResultFromMetadata", func() {

	It("creates an empty StagingResult if metadata is empty", func() {
		metadata := &files.BuildMetadata{}
		result := staging.StagingResultFromMetadata(metadata)

		Expect(result.LifecycleType).To(Equal("cnb"))
		Expect(result.Buildpacks).To(BeEmpty())
		Expect(result.ProcessTypes).To(BeEmpty())
	})

	It("creates StagingResult containing buildpacks", func() {
		metadata := &files.BuildMetadata{
			Buildpacks: []buildpack.GroupElement{{ID: "nodejs", Version: "1.0.0"}, {ID: "java", Version: "2.0.0"}},
		}
		result := staging.StagingResultFromMetadata(metadata)

		Expect(result.LifecycleType).To(Equal("cnb"))
		Expect(result.Buildpacks).To(ContainElements(staging.BuildpackMetadata{ID: "nodejs", Name: "nodejs@1.0.0", Version: "1.0.0"},
			staging.BuildpackMetadata{ID: "java", Name: "java@2.0.0", Version: "2.0.0"}))
		Expect(result.ProcessTypes).To(BeEmpty())
	})

	It("creates StagingResult containing process type", func() {
		metadata := &files.BuildMetadata{
			Buildpacks: []buildpack.GroupElement{{ID: "nodejs", Version: "1.0.0"}},
			Processes: []launch.Process{{Type: "web", Command: launch.NewRawCommand([]string{"start", "app"}), Args: []string{"--force"}},
				{Type: "custom", Command: launch.NewRawCommand([]string{"start.sh"}), Args: []string{"--arg1"}}},
		}
		result := staging.StagingResultFromMetadata(metadata)

		Expect(result.LifecycleType).To(Equal("cnb"))
		Expect(result.ProcessTypes).To(Equal(staging.ProcessTypes{"web": "start app --force", "custom": "start.sh --arg1"}))
	})
})
