// This file was generated by counterfeiter
package outfakes

import (
	"sync"

	go_pivnet "github.com/pivotal-cf/go-pivnet"
	"github.com/pivotal-cf/pivnet-resource/concourse"
)

type Finalizer struct {
	FinalizeStub        func(productSlug string, release go_pivnet.Release) (concourse.OutResponse, error)
	finalizeMutex       sync.RWMutex
	finalizeArgsForCall []struct {
		productSlug string
		release     go_pivnet.Release
	}
	finalizeReturns struct {
		result1 concourse.OutResponse
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Finalizer) Finalize(productSlug string, release go_pivnet.Release) (concourse.OutResponse, error) {
	fake.finalizeMutex.Lock()
	fake.finalizeArgsForCall = append(fake.finalizeArgsForCall, struct {
		productSlug string
		release     go_pivnet.Release
	}{productSlug, release})
	fake.recordInvocation("Finalize", []interface{}{productSlug, release})
	fake.finalizeMutex.Unlock()
	if fake.FinalizeStub != nil {
		return fake.FinalizeStub(productSlug, release)
	} else {
		return fake.finalizeReturns.result1, fake.finalizeReturns.result2
	}
}

func (fake *Finalizer) FinalizeCallCount() int {
	fake.finalizeMutex.RLock()
	defer fake.finalizeMutex.RUnlock()
	return len(fake.finalizeArgsForCall)
}

func (fake *Finalizer) FinalizeArgsForCall(i int) (string, go_pivnet.Release) {
	fake.finalizeMutex.RLock()
	defer fake.finalizeMutex.RUnlock()
	return fake.finalizeArgsForCall[i].productSlug, fake.finalizeArgsForCall[i].release
}

func (fake *Finalizer) FinalizeReturns(result1 concourse.OutResponse, result2 error) {
	fake.FinalizeStub = nil
	fake.finalizeReturns = struct {
		result1 concourse.OutResponse
		result2 error
	}{result1, result2}
}

func (fake *Finalizer) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.finalizeMutex.RLock()
	defer fake.finalizeMutex.RUnlock()
	return fake.invocations
}

func (fake *Finalizer) recordInvocation(key string, args []interface{}) {
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
