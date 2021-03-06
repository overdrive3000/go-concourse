package concourse_test

import (
	"net/http"

	"github.com/concourse/atc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("ATC Handler Pipelines", func() {
	Describe("PausePipeline", func() {
		Context("when the pipeline exists", func() {
			BeforeEach(func() {
				expectedURL := "/api/v1/pipelines/mypipeline/pause"
				atcServer.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("PUT", expectedURL),
						ghttp.RespondWithJSONEncoded(http.StatusOK, ""),
					),
				)
			})
			It("return true and no error", func() {
				found, err := client.PausePipeline("mypipeline")
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
			})
		})

		Context("when the pipeline doesn't exist", func() {
			BeforeEach(func() {
				expectedURL := "/api/v1/pipelines/mypipeline/pause"
				atcServer.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("PUT", expectedURL),
						ghttp.RespondWithJSONEncoded(http.StatusNotFound, ""),
					),
				)
			})
			It("returns false and no error", func() {
				found, err := client.PausePipeline("mypipeline")
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeFalse())
			})
		})

	})

	Describe("Pipeline", func() {
		var expectedPipeline atc.Pipeline
		pipelineName := "mypipeline"
		expectedURL := "/api/v1/pipelines/mypipeline"

		BeforeEach(func() {
			expectedPipeline = atc.Pipeline{
				Name:   "mypipeline",
				Paused: true,
				Groups: []atc.GroupConfig{
					{
						Name:      "group1",
						Jobs:      []string{"job1", "job2"},
						Resources: []string{"resource1", "resource2"},
					},
				},
			}
		})

		Context("when the pipeline is found", func() {
			BeforeEach(func() {
				atcServer.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", expectedURL),
						ghttp.RespondWithJSONEncoded(http.StatusOK, expectedPipeline),
					),
				)
			})

			It("returns the requested pipeline", func() {
				pipeline, found, err := client.Pipeline(pipelineName)
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				Expect(pipeline).To(Equal(expectedPipeline))
			})
		})

		Context("when the pipeline is not found", func() {
			BeforeEach(func() {
				atcServer.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", expectedURL),
						ghttp.RespondWith(http.StatusNotFound, ""),
					),
				)
			})

			It("returns false", func() {
				_, found, err := client.Pipeline(pipelineName)
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeFalse())
			})
		})
	})

	Describe("ListPipelines", func() {
		var expectedPipelines []atc.Pipeline

		BeforeEach(func() {
			expectedURL := "/api/v1/pipelines"

			expectedPipelines = []atc.Pipeline{
				{
					Name:   "mypipeline-1",
					Paused: true,
					Groups: []atc.GroupConfig{
						{
							Name:      "group1",
							Jobs:      []string{"job1", "job2"},
							Resources: []string{"resource1", "resource2"},
						},
					},
				},
				{
					Name:   "mypipeline-2",
					Paused: false,
					Groups: []atc.GroupConfig{
						{
							Name:      "group2",
							Jobs:      []string{"job3", "job4"},
							Resources: []string{"resource3", "resource4"},
						},
					},
				},
			}

			atcServer.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", expectedURL),
					ghttp.RespondWithJSONEncoded(http.StatusOK, expectedPipelines),
				),
			)
		})

		It("returns all the pipelines", func() {
			pipelines, err := client.ListPipelines()
			Expect(err).NotTo(HaveOccurred())
			Expect(pipelines).To(Equal(expectedPipelines))
		})
	})

	Describe("DeletePipeline", func() {
		expectedURL := "/api/v1/pipelines/mypipeline"

		Context("when the pipeline exists", func() {
			BeforeEach(func() {
				atcServer.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("DELETE", expectedURL),
						ghttp.RespondWith(http.StatusNoContent, ""),
					),
				)
			})

			It("deletes the pipeline when called", func() {
				Expect(func() {
					found, err := client.DeletePipeline("mypipeline")
					Expect(err).NotTo(HaveOccurred())
					Expect(found).To(BeTrue())
				}).To(Change(func() int {
					return len(atcServer.ReceivedRequests())
				}).By(1))
			})
		})

		Context("when the pipeline does not exist", func() {
			BeforeEach(func() {
				atcServer.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("DELETE", expectedURL),
						ghttp.RespondWith(http.StatusNotFound, ""),
					),
				)
			})

			It("returns false and no error", func() {
				found, err := client.DeletePipeline("mypipeline")
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeFalse())
			})
		})
	})

	Describe("RenamePipeline", func() {
		expectedURL := "/api/v1/pipelines/mypipeline/rename"

		Context("when the pipeline exists", func() {
			BeforeEach(func() {
				atcServer.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("PUT", expectedURL),
						ghttp.VerifyJSON(`{"name":"newpipelinename"}`),
						ghttp.RespondWith(http.StatusNoContent, ""),
					),
				)
			})

			It("renames the pipeline when called", func() {
				renamed, err := client.RenamePipeline("mypipeline", "newpipelinename")
				Expect(err).NotTo(HaveOccurred())
				Expect(renamed).To(BeTrue())
			})
		})

		Context("when the pipeline does not exist", func() {
			BeforeEach(func() {
				atcServer.AppendHandlers(
					ghttp.RespondWith(http.StatusNotFound, ""),
				)
			})

			It("returns false and no error", func() {
				renamed, err := client.RenamePipeline("mypipeline", "newpipelinename")
				Expect(err).NotTo(HaveOccurred())
				Expect(renamed).To(BeFalse())
			})
		})

		Context("when an error occurs", func() {
			BeforeEach(func() {
				atcServer.AppendHandlers(
					ghttp.RespondWith(http.StatusTeapot, ""),
				)
			})

			It("returns an error", func() {
				renamed, err := client.RenamePipeline("mypipeline", "newpipelinename")
				Expect(err).To(MatchError(ContainSubstring("418 I'm a teapot")))
				Expect(renamed).To(BeFalse())
			})
		})
	})
})
