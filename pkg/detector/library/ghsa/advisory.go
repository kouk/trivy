package ghsa

import (
	"strings"

	"golang.org/x/xerrors"

	"github.com/aquasecurity/trivy-db/pkg/vulnsrc/ghsa"
	"github.com/aquasecurity/trivy/pkg/detector/library/comparer"
	"github.com/aquasecurity/trivy/pkg/types"
)

// Advisory implements VulnSrc
type Advisory struct {
	vs       ghsa.VulnSrc
	comparer comparer.Comparer
}

// NewAdvisory is the factory method to return advisory
func NewAdvisory(ecosystem ghsa.Ecosystem, comparer comparer.Comparer) *Advisory {
	return &Advisory{
		vs:       ghsa.NewVulnSrc(ecosystem),
		comparer: comparer,
	}
}

// DetectVulnerabilities scans package for vulnerabilities
func (s *Advisory) DetectVulnerabilities(pkgName, pkgVer string) ([]types.DetectedVulnerability, error) {
	advisories, err := s.vs.Get(pkgName)
	if err != nil {
		return nil, xerrors.Errorf("failed to get ghsa advisories: %w", err)
	}

	var vulns []types.DetectedVulnerability
	for _, advisory := range advisories {
		if !s.comparer.IsVulnerable(pkgVer, advisory) {
			continue
		}

		vuln := types.DetectedVulnerability{
			VulnerabilityID:  advisory.VulnerabilityID,
			PkgName:          pkgName,
			InstalledVersion: pkgVer,
			FixedVersion:     strings.Join(advisory.PatchedVersions, ", "),
			DataSource:       advisory.DataSource,
		}
		vulns = append(vulns, vuln)
	}

	return vulns, nil
}
