package db_test

import (
	"github.com/concourse/concourse/atc"
	"github.com/concourse/concourse/atc/db"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PipelinePauser", func() {
	var (
		pauser               db.PipelinePauser
		twoJobPipeline       db.Pipeline
		err                  error
		twoJobPipelineConfig = atc.Config{
			Jobs: atc.JobConfigs{
				{
					Name: "job-one",
				},
				{
					Name: "job-two",
				},
			},
		}
		pipelineRef = atc.PipelineRef{Name: "twojobs-pipeline"}
	)

	BeforeEach(func() {
		pauser = db.NewPipelinePauser(dbConn, lockFactory)
	})

	Describe("PausePipelines that haven't run in more than 10 days", func() {
		Context("last run was 15 days ago", func() {
			It("should be paused", func() {
				By("creating a pipeline with two jobs")
				twoJobPipeline, _, err = defaultTeam.SavePipeline(pipelineRef, twoJobPipelineConfig, db.ConfigVersion(0), false)
				Expect(err).NotTo(HaveOccurred())
				Expect(twoJobPipeline.Paused()).To(BeFalse(), "pipeline should start unpaused")

				By("creating a job that ran 15 days ago")
				jobOne, found, err := twoJobPipeline.Job("job-one")
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				b1, err := jobOne.CreateBuild(defaultBuildCreatedBy)
				Expect(err).NotTo(HaveOccurred())
				b1.Finish(db.BuildStatusSucceeded)
				_, err = dbConn.Exec(`UPDATE jobs SET last_scheduled = NOW() - INTERVAL '15' DAY WHERE name = 'job-one'`)
				Expect(err).NotTo(HaveOccurred())

				By("creating a job that ran 20 days ago")
				jobTwo, found, err := twoJobPipeline.Job("job-two")
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				jobTwo.CreateBuild(defaultBuildCreatedBy)
				b2, err := jobTwo.CreateBuild(defaultBuildCreatedBy)
				Expect(err).NotTo(HaveOccurred())
				b2.Finish(db.BuildStatusSucceeded)
				_, err = dbConn.Exec(`UPDATE jobs SET last_scheduled = NOW() - INTERVAL '20' DAY WHERE name = 'job-two'`)
				Expect(err).NotTo(HaveOccurred())

				By("running the pipeline pauser")
				err = pauser.PausePipelines(10)
				Expect(err).NotTo(HaveOccurred())

				_, err = twoJobPipeline.Reload()
				Expect(err).To(BeNil())
				Expect(twoJobPipeline.Paused()).To(BeTrue(), "pipeline should be paused")

				pauser.PausePipelines(10)
			})
		})
		Context("last run was 1 day ago", func() {
			It("should not be paused", func() {
				By("creating a pipeline with two jobs")
				twoJobPipeline, _, err = defaultTeam.SavePipeline(pipelineRef, twoJobPipelineConfig, db.ConfigVersion(0), false)
				Expect(err).NotTo(HaveOccurred())
				Expect(twoJobPipeline.Paused()).To(BeFalse(), "pipeline should start unpaused")

				By("creating a job that ran yesterday")
				jobOne, found, err := twoJobPipeline.Job("job-one")
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				b1, err := jobOne.CreateBuild(defaultBuildCreatedBy)
				Expect(err).NotTo(HaveOccurred())
				b1.Finish(db.BuildStatusSucceeded)
				_, err = dbConn.Exec(`UPDATE jobs SET last_scheduled = NOW() - INTERVAL '1' DAY WHERE name = 'job-one'`)
				Expect(err).NotTo(HaveOccurred())

				By("creating a job that ran 11 days ago")
				jobTwo, found, err := twoJobPipeline.Job("job-two")
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				jobTwo.CreateBuild(defaultBuildCreatedBy)
				b2, err := jobTwo.CreateBuild(defaultBuildCreatedBy)
				Expect(err).NotTo(HaveOccurred())
				b2.Finish(db.BuildStatusSucceeded)
				_, err = dbConn.Exec(`UPDATE jobs SET last_scheduled = NOW() - INTERVAL '11' DAY WHERE name = 'job-two'`)
				Expect(err).NotTo(HaveOccurred())

				By("running the pipeline pauser")
				err = pauser.PausePipelines(10)
				Expect(err).NotTo(HaveOccurred())

				_, err = twoJobPipeline.Reload()
				Expect(err).To(BeNil())
				Expect(twoJobPipeline.Paused()).To(BeFalse(), "pipeline should NOT be paused")

				pauser.PausePipelines(10)
			})
		})
		Context("last run was 10 days ago", func() {
			It("should not be paused", func() {
				By("creating a pipeline with two jobs")
				twoJobPipeline, _, err = defaultTeam.SavePipeline(pipelineRef, twoJobPipelineConfig, db.ConfigVersion(0), false)
				Expect(err).NotTo(HaveOccurred())
				Expect(twoJobPipeline.Paused()).To(BeFalse(), "pipeline should start unpaused")

				By("creating a job that ran 10 days ago")
				jobOne, found, err := twoJobPipeline.Job("job-one")
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				b1, err := jobOne.CreateBuild(defaultBuildCreatedBy)
				Expect(err).NotTo(HaveOccurred())
				b1.Finish(db.BuildStatusSucceeded)
				_, err = dbConn.Exec(`UPDATE jobs SET last_scheduled = NOW() - INTERVAL '10' DAY WHERE name = 'job-one'`)
				Expect(err).NotTo(HaveOccurred())

				By("creating a job that ran 20 days ago")
				jobTwo, found, err := twoJobPipeline.Job("job-two")
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				jobTwo.CreateBuild(defaultBuildCreatedBy)
				b2, err := jobTwo.CreateBuild(defaultBuildCreatedBy)
				Expect(err).NotTo(HaveOccurred())
				b2.Finish(db.BuildStatusSucceeded)
				_, err = dbConn.Exec(`UPDATE jobs SET last_scheduled = NOW() - INTERVAL '20' DAY WHERE name = 'job-two'`)
				Expect(err).NotTo(HaveOccurred())

				By("running the pipeline pauser")
				err = pauser.PausePipelines(10)
				Expect(err).NotTo(HaveOccurred())

				_, err = twoJobPipeline.Reload()
				Expect(err).To(BeNil())
				Expect(twoJobPipeline.Paused()).To(BeFalse(), "pipeline should NOT be paused")

				pauser.PausePipelines(10)
			})
		})
	})
})
