package release

import (
	"fmt"
	"strconv"

	pivnet "github.com/pivotal-cf/go-pivnet"
	"github.com/pivotal-cf/pivnet-resource/concourse"
	"github.com/pivotal-cf/pivnet-resource/metadata"
	"github.com/pivotal-cf/pivnet-resource/versions"
)

type ReleaseFinalizer struct {
	pivnet      updateClient
	metadata    metadata.Metadata
	params      concourse.OutParams
	sourcesDir  string
	productSlug string
}

func NewFinalizer(
	pivnetClient updateClient,
	params concourse.OutParams,
	metadata metadata.Metadata,
	sourcesDir,
	productSlug string,
) ReleaseFinalizer {
	return ReleaseFinalizer{
		pivnet:      pivnetClient,
		params:      params,
		metadata:    metadata,
		sourcesDir:  sourcesDir,
		productSlug: productSlug,
	}
}

//go:generate counterfeiter --fake-name UpdateClient . updateClient
type updateClient interface {
	UpdateRelease(productSlug string, release pivnet.Release) (pivnet.Release, error)
	ReleaseETag(productSlug string, releaseID int) (string, error)
	AddUserGroup(productSlug string, releaseID int, userGroupID int) error
	AddReleaseDependency(productSlug string, releaseID int, dependentReleaseID int) error
	GetRelease(productSlug string, releaseVersion string) (pivnet.Release, error)
}

func (rf ReleaseFinalizer) Finalize(productSlug string, release pivnet.Release) (concourse.OutResponse, error) {
	for i, d := range rf.metadata.Dependencies {
		dependentReleaseID := d.Release.ID
		if dependentReleaseID == 0 {
			if d.Release.Version == "" || d.Release.Product.Slug == "" {
				return concourse.OutResponse{}, fmt.Errorf(
					"Either ReleaseID or release version and product slug must be provided for dependency[%d]", i)
			}

			r, err := rf.pivnet.GetRelease(d.Release.Product.Slug, d.Release.Version)
			if err != nil {
				return concourse.OutResponse{}, err
			}
			dependentReleaseID = r.ID
		}

		err := rf.pivnet.AddReleaseDependency(productSlug, release.ID, dependentReleaseID)
		if err != nil {
			return concourse.OutResponse{}, err
		}
	}

	availability := rf.metadata.Release.Availability

	if availability != "Admins Only" {
		releaseUpdate := pivnet.Release{
			ID:           release.ID,
			Availability: availability,
		}

		var err error
		release, err = rf.pivnet.UpdateRelease(rf.productSlug, releaseUpdate)
		if err != nil {
			return concourse.OutResponse{}, err
		}

		if availability == "Selected User Groups Only" {
			userGroupIDs := rf.metadata.Release.UserGroupIDs

			for _, userGroupIDString := range userGroupIDs {
				userGroupID, err := strconv.Atoi(userGroupIDString)
				if err != nil {
					return concourse.OutResponse{}, err
				}

				err = rf.pivnet.AddUserGroup(rf.productSlug, release.ID, userGroupID)
				if err != nil {
					return concourse.OutResponse{}, err
				}
			}
		}
	}

	releaseETag, err := rf.pivnet.ReleaseETag(rf.productSlug, release.ID)
	if err != nil {
		return concourse.OutResponse{}, err
	}

	outputVersion, err := versions.CombineVersionAndETag(release.Version, releaseETag)
	if err != nil {
		return concourse.OutResponse{}, err // this will never return an error
	}

	metadata := []concourse.Metadata{
		{Name: "version", Value: release.Version},
		{Name: "release_type", Value: string(release.ReleaseType)},
		{Name: "release_date", Value: release.ReleaseDate},
		{Name: "description", Value: release.Description},
		{Name: "release_notes_url", Value: release.ReleaseNotesURL},
		{Name: "availability", Value: release.Availability},
		{Name: "controlled", Value: fmt.Sprintf("%t", release.Controlled)},
		{Name: "eccn", Value: release.ECCN},
		{Name: "license_exception", Value: release.LicenseException},
		{Name: "end_of_support_date", Value: release.EndOfSupportDate},
		{Name: "end_of_guidance_date", Value: release.EndOfGuidanceDate},
		{Name: "end_of_availability_date", Value: release.EndOfAvailabilityDate},
	}
	if release.EULA != nil {
		metadata = append(
			metadata,
			concourse.Metadata{Name: "eula_slug", Value: release.EULA.Slug})
	}

	return concourse.OutResponse{
		Version: concourse.Version{
			ProductVersion: outputVersion,
		},
		Metadata: metadata,
	}, nil
}
