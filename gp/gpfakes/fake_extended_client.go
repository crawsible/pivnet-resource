// This file was generated by counterfeiter
package gpfakes

import (
	"sync"

	"github.com/pivotal-cf-experimental/go-pivnet"
	"github.com/pivotal-cf-experimental/pivnet-resource/gp"
)

type FakeExtendedClient struct {
	ReleaseETagStub        func(productSlug string, releaseID int) (string, error)
	releaseETagMutex       sync.RWMutex
	releaseETagArgsForCall []struct {
		productSlug string
		releaseID   int
	}
	releaseETagReturns struct {
		result1 string
		result2 error
	}
	ProductVersionsStub        func(productSlug string, releases []pivnet.Release) ([]string, error)
	productVersionsMutex       sync.RWMutex
	productVersionsArgsForCall []struct {
		productSlug string
		releases    []pivnet.Release
	}
	productVersionsReturns struct {
		result1 []string
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeExtendedClient) ReleaseETag(productSlug string, releaseID int) (string, error) {
	fake.releaseETagMutex.Lock()
	fake.releaseETagArgsForCall = append(fake.releaseETagArgsForCall, struct {
		productSlug string
		releaseID   int
	}{productSlug, releaseID})
	fake.recordInvocation("ReleaseETag", []interface{}{productSlug, releaseID})
	fake.releaseETagMutex.Unlock()
	if fake.ReleaseETagStub != nil {
		return fake.ReleaseETagStub(productSlug, releaseID)
	} else {
		return fake.releaseETagReturns.result1, fake.releaseETagReturns.result2
	}
}

func (fake *FakeExtendedClient) ReleaseETagCallCount() int {
	fake.releaseETagMutex.RLock()
	defer fake.releaseETagMutex.RUnlock()
	return len(fake.releaseETagArgsForCall)
}

func (fake *FakeExtendedClient) ReleaseETagArgsForCall(i int) (string, int) {
	fake.releaseETagMutex.RLock()
	defer fake.releaseETagMutex.RUnlock()
	return fake.releaseETagArgsForCall[i].productSlug, fake.releaseETagArgsForCall[i].releaseID
}

func (fake *FakeExtendedClient) ReleaseETagReturns(result1 string, result2 error) {
	fake.ReleaseETagStub = nil
	fake.releaseETagReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeExtendedClient) ProductVersions(productSlug string, releases []pivnet.Release) ([]string, error) {
	var releasesCopy []pivnet.Release
	if releases != nil {
		releasesCopy = make([]pivnet.Release, len(releases))
		copy(releasesCopy, releases)
	}
	fake.productVersionsMutex.Lock()
	fake.productVersionsArgsForCall = append(fake.productVersionsArgsForCall, struct {
		productSlug string
		releases    []pivnet.Release
	}{productSlug, releasesCopy})
	fake.recordInvocation("ProductVersions", []interface{}{productSlug, releasesCopy})
	fake.productVersionsMutex.Unlock()
	if fake.ProductVersionsStub != nil {
		return fake.ProductVersionsStub(productSlug, releases)
	} else {
		return fake.productVersionsReturns.result1, fake.productVersionsReturns.result2
	}
}

func (fake *FakeExtendedClient) ProductVersionsCallCount() int {
	fake.productVersionsMutex.RLock()
	defer fake.productVersionsMutex.RUnlock()
	return len(fake.productVersionsArgsForCall)
}

func (fake *FakeExtendedClient) ProductVersionsArgsForCall(i int) (string, []pivnet.Release) {
	fake.productVersionsMutex.RLock()
	defer fake.productVersionsMutex.RUnlock()
	return fake.productVersionsArgsForCall[i].productSlug, fake.productVersionsArgsForCall[i].releases
}

func (fake *FakeExtendedClient) ProductVersionsReturns(result1 []string, result2 error) {
	fake.ProductVersionsStub = nil
	fake.productVersionsReturns = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeExtendedClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.releaseETagMutex.RLock()
	defer fake.releaseETagMutex.RUnlock()
	fake.productVersionsMutex.RLock()
	defer fake.productVersionsMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeExtendedClient) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ gp.ExtendedClient = new(FakeExtendedClient)
