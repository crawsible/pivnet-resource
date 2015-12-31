package pivnet_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/pivotal-cf-experimental/pivnet-resource/logger"
	logger_fakes "github.com/pivotal-cf-experimental/pivnet-resource/logger/fakes"
	"github.com/pivotal-cf-experimental/pivnet-resource/pivnet"
)

var _ = Describe("PivnetClient - product files", func() {
	var (
		server     *ghttp.Server
		client     pivnet.Client
		token      string
		apiAddress string

		fakeLogger logger.Logger
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		apiAddress = server.URL() + apiPrefix
		token = "my-auth-token"

		fakeLogger = &logger_fakes.FakeLogger{}
		client = pivnet.NewClient(apiAddress, token, fakeLogger)
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("Get Product Files", func() {
		It("returns the product files for the given release", func() {
			response, err := json.Marshal(pivnet.ProductFiles{[]pivnet.ProductFile{
				{ID: 3, AWSObjectKey: "anything", Links: pivnet.Links{Download: map[string]string{"href": "/products/banana/releases/666/product_files/6/download"}}},
				{ID: 4, AWSObjectKey: "something", Links: pivnet.Links{Download: map[string]string{"href": "/products/banana/releases/666/product_files/8/download"}}},
			},
			})
			Expect(err).NotTo(HaveOccurred())

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", apiPrefix+"/products/banana/releases/666/product_files"),
					ghttp.RespondWith(http.StatusOK, response),
				),
			)

			release := pivnet.Release{
				Links: pivnet.Links{
					ProductFiles: map[string]string{"href": apiAddress + "/products/banana/releases/666/product_files"},
				},
			}

			product, err := client.GetProductFiles(release)
			Expect(err).NotTo(HaveOccurred())
			Expect(product.ProductFiles).To(HaveLen(2))

			Expect(product.ProductFiles[0].AWSObjectKey).To(Equal("anything"))
			Expect(product.ProductFiles[1].AWSObjectKey).To(Equal("something"))

			Expect(product.ProductFiles[0].Links.Download["href"]).To(Equal("/products/banana/releases/666/product_files/6/download"))
			Expect(product.ProductFiles[1].Links.Download["href"]).To(Equal("/products/banana/releases/666/product_files/8/download"))
		})

		Context("when the server responds with a non-2XX status code", func() {
			It("returns an error", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", apiPrefix+"/products/banana/releases/666/product_files"),
						ghttp.RespondWith(http.StatusTeapot, nil),
					),
				)
				release := pivnet.Release{
					Links: pivnet.Links{
						ProductFiles: map[string]string{"href": apiAddress + "/products/banana/releases/666/product_files"},
					},
				}

				_, err := client.GetProductFiles(release)
				Expect(err).To(MatchError(errors.New(
					"Pivnet returned status code: 418 for the request - expected 200")))
			})
		})
	})

	Describe("Create Product File", func() {
		var (
			createProductFileConfig pivnet.CreateProductFileConfig
		)

		BeforeEach(func() {
			createProductFileConfig = pivnet.CreateProductFileConfig{
				ProductName:  productName,
				Name:         "some-file-name",
				FileVersion:  "some-file-version",
				AWSObjectKey: "some-aws-object-key",
			}
		})

		Context("when the config is valid", func() {
			type requestBody struct {
				ProductFile pivnet.ProductFile `json:"product_file"`
			}

			const (
				expectedMD5      = "not-supported-yet"
				expectedFileType = "Software"
			)

			var (
				expectedRequestBody requestBody

				validResponse = `{"product_file":{"id":1234}}`
			)

			BeforeEach(func() {
				expectedRequestBody = requestBody{
					ProductFile: pivnet.ProductFile{
						FileType:     "Software",
						FileVersion:  createProductFileConfig.FileVersion,
						Name:         createProductFileConfig.Name,
						MD5:          "not-supported-yet",
						AWSObjectKey: createProductFileConfig.AWSObjectKey,
					},
				}
			})

			It("creates the release with the minimum required fields", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", apiPrefix+"/products/"+productName+"/product_files"),
						ghttp.VerifyJSONRepresenting(&expectedRequestBody),
						ghttp.RespondWith(http.StatusCreated, validResponse),
					),
				)

				release, err := client.CreateProductFile(createProductFileConfig)
				Expect(err).NotTo(HaveOccurred())
				Expect(release.ID).To(Equal(1234))
			})
		})

		Context("when the server responds with a non-201 status code", func() {
			It("returns an error", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", apiPrefix+"/products/"+productName+"/product_files"),
						ghttp.RespondWith(http.StatusTeapot, nil),
					),
				)

				_, err := client.CreateProductFile(createProductFileConfig)
				Expect(err).To(MatchError(errors.New(
					"Pivnet returned status code: 418 for the request - expected 201")))
			})
		})
	})

	Describe("Delete Product File", func() {
		var (
			id = 1234
		)

		It("deletes the product file", func() {
			response := []byte(`{"product_file":{"id":1234}}`)

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(
						"DELETE",
						fmt.Sprintf("%s/products/%s/product_files/%d", apiPrefix, productName, id)),
					ghttp.RespondWith(http.StatusOK, response),
				),
			)

			productFile, err := client.DeleteProductFile(productName, id)
			Expect(err).NotTo(HaveOccurred())

			Expect(productFile.ID).To(Equal(id))
		})

		Context("when the server responds with a non-2XX status code", func() {
			It("returns an error", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(
							"DELETE",
							fmt.Sprintf("%s/products/%s/product_files/%d", apiPrefix, productName, id)),
						ghttp.RespondWith(http.StatusTeapot, nil),
					),
				)

				_, err := client.DeleteProductFile(productName, id)
				Expect(err).To(MatchError(errors.New(
					"Pivnet returned status code: 418 for the request - expected 200")))
			})
		})
	})
})
