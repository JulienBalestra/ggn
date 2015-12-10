package env

import (
	"strings"
)

func (s *Service) Check() {
	s.log.Debug("Running check")
	s.runHook(EARLY, "service/check", "check")
	defer s.runHook(LATE, "service/check", "check")

	s.Generate()

	units, _, err := s.env.RunFleetCmdGetOutput("-strict-host-key-checking=false", "list-unit-files", "-no-legend", "-fields", "unit")
	if err != nil {
		s.log.WithError(err).Fatal("Cannot list unit files")
	}

	for _, unitName := range strings.Split(units, "\n") {
		unitInfo := strings.Split(unitName, "_")
		if len(unitInfo) != 3 {
			continue
		}
		if unitInfo[1] != s.Name {
			continue
		}
		split := strings.Split(unitInfo[2], ".")

		s.LoadUnit(split[0]).Check("service/check")
	}
}
